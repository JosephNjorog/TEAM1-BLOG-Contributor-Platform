package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"team1blog/backend/internal/admin"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/avalanche"
	"team1blog/backend/internal/banners"
	"team1blog/backend/internal/cloudinary"
	"team1blog/backend/internal/config"
	"team1blog/backend/internal/db"
	"team1blog/backend/internal/email"
	"team1blog/backend/internal/notifications"
	"team1blog/backend/internal/payments"
	"team1blog/backend/internal/profile"
	"team1blog/backend/internal/reviews"
	"team1blog/backend/internal/substack"
	"team1blog/backend/internal/uploads"
	"team1blog/backend/internal/users"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer pool.Close()

	if err := db.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	auditLogger := audit.NewLogger(pool)
	mailer := email.NewSender(cfg.ResendAPIKey, cfg.ResendFromAddr, cfg.MockEmail)

	usersRepo := users.NewRepository(pool)
	authRepo := auth.NewRepository(pool)
	tokenIssuer := auth.NewTokenIssuer(cfg.JWTSecret, cfg.AccessTokenTTL)
	authService := auth.NewService(usersRepo, authRepo, tokenIssuer, auditLogger, mailer, cfg.InviteTTL, cfg.RefreshTokenTTL, cfg.FrontendURL, cfg.AdminAppURL)
	authHandler := auth.NewHandler(authService)

	profileHandler := profile.NewHandler(usersRepo)

	notificationsRepo := notifications.NewRepository(pool)
	notificationsHandler := notifications.NewHandler(notificationsRepo)

	articlesRepo := articles.NewRepository(pool)
	articlesService := articles.NewService(articlesRepo, usersRepo, notificationsRepo, auditLogger, mailer, cfg.FrontendURL)
	articlesHandler := articles.NewHandler(articlesService)

	uploader := cloudinary.NewUploader(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret, cfg.PublicAPIURL, cfg.MockImages, "uploads")

	reviewsRepo := reviews.NewRepository(pool)
	reviewsService := reviews.NewService(reviewsRepo, articlesRepo, usersRepo, notificationsRepo, auditLogger, mailer, cfg.FrontendURL)
	reviewsHandler := reviews.NewHandler(reviewsService, articlesService)

	bannersRepo := banners.NewRepository(pool)
	bannersService := banners.NewService(bannersRepo, articlesRepo, usersRepo, notificationsRepo, uploader, auditLogger, mailer, cfg.FrontendURL)
	bannersHandler := banners.NewHandler(bannersService, articlesService)

	uploadsHandler := uploads.NewHandler(articlesRepo, uploader)

	avalancheSender, err := avalanche.NewSender(cfg.AvalancheRPCURL, cfg.AvalancheTreasuryKey, cfg.AvalancheUSDCContract, cfg.AvalancheChainID, cfg.MockPayments)
	if err != nil {
		log.Fatalf("avalanche sender setup failed: %v", err)
	}
	paymentsRepo := payments.NewRepository(pool)
	paymentsService := payments.NewService(paymentsRepo, articlesRepo, usersRepo, notificationsRepo, avalancheSender, auditLogger, mailer, cfg.FrontendURL, cfg.MockPayments)
	paymentsHandler := payments.NewHandler(paymentsService)

	adminRepo := admin.NewRepository(pool)
	adminService := admin.NewService(adminRepo, usersRepo, auditLogger)
	adminHandler := admin.NewHandler(adminService)

	substackFetcher := substack.NewFetcher(cfg.SubstackPublicationURL, cfg.MockSubstack)
	substackRepo := substack.NewRepository(pool)
	substackService := substack.NewService(substackRepo, usersRepo, substackFetcher)
	substackHandler := substack.NewHandler(substackService)

	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Serves files written by the mock (local-disk) Cloudinary uploader so
	// banner/inline image URLs resolve to something real before Cloudinary
	// credentials exist.
	if cfg.MockImages {
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	}

	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", auth.Routes(authHandler, tokenIssuer))
		api.Mount("/me", profile.Routes(profileHandler, tokenIssuer))
		api.Mount("/notifications", notifications.Routes(notificationsHandler, tokenIssuer))
		api.Mount("/articles", articles.Routes(articlesHandler, tokenIssuer))
		api.Mount("/reviews", reviews.Routes(reviewsHandler, tokenIssuer))
		api.Mount("/banners", banners.Routes(bannersHandler, tokenIssuer))
		api.Mount("/uploads", uploads.Routes(uploadsHandler, tokenIssuer))
		api.Mount("/payments", payments.Routes(paymentsHandler, tokenIssuer))
		api.Mount("/admin", admin.Routes(adminHandler, tokenIssuer))
		api.Mount("/sync/substack", substack.Routes(substackHandler, tokenIssuer))
	})

	startSubstackScheduler(ctx, substackService)

	logIntegrationModes(cfg)

	addr := ":" + cfg.Port
	log.Printf("team1blog api listening on %s (env=%s)", addr, cfg.Env)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

// startSubstackScheduler runs an initial sync shortly after boot, then
// repeats every 6 hours per the PRD, independent of the manual admin-
// triggered sync endpoint.
func startSubstackScheduler(ctx context.Context, service *substack.Service) {
	go func() {
		time.Sleep(10 * time.Second)
		runSubstackSync(ctx, service)

		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				runSubstackSync(ctx, service)
			}
		}
	}()
}

func runSubstackSync(ctx context.Context, service *substack.Service) {
	count, err := service.Sync(ctx)
	if err != nil {
		log.Printf("substack sync failed: %v", err)
		return
	}
	log.Printf("substack sync complete: %d posts", count)
}

func logIntegrationModes(cfg *config.Config) {
	mode := func(mock bool) string {
		if mock {
			return "MOCK"
		}
		return "LIVE"
	}
	log.Printf("integrations: email=%s images=%s payments=%s substack=%s",
		mode(cfg.MockEmail), mode(cfg.MockImages), mode(cfg.MockPayments), mode(cfg.MockSubstack))
}
