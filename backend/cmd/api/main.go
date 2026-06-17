package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/config"
	"team1blog/backend/internal/db"
	"team1blog/backend/internal/email"
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
	authService := auth.NewService(usersRepo, authRepo, tokenIssuer, auditLogger, mailer, cfg.InviteTTL, cfg.RefreshTokenTTL, frontendAppURL(cfg))
	authHandler := auth.NewHandler(authService)

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

	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/auth", auth.Routes(authHandler, tokenIssuer))
	})

	logIntegrationModes(cfg)

	addr := ":" + cfg.Port
	log.Printf("team1blog api listening on %s (env=%s)", addr, cfg.Env)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func frontendAppURL(cfg *config.Config) string {
	if len(cfg.CORSOrigins) > 0 {
		return cfg.CORSOrigins[0]
	}
	return "http://localhost:5173"
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
