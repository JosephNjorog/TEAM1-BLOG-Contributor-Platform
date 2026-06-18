package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	PublicAPIURL    string
	Env             string
	DatabaseURL     string
	JWTSecret       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	InviteTTL       time.Duration
	CORSOrigins     []string
	FrontendURL     string
	AdminAppURL     string

	// Resend
	ResendAPIKey   string
	ResendFromAddr string
	MockEmail      bool

	// Cloudinary
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	MockImages          bool

	// Avalanche / Core wallet payments
	AvalancheRPCURL       string
	AvalancheChainID      int64
	AvalancheTreasuryKey  string
	AvalancheUSDCContract string
	MockPayments          bool

	// Substack
	SubstackPublicationURL string
	MockSubstack           bool
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		PublicAPIURL:    getEnv("PUBLIC_API_URL", "http://localhost:8080"),
		Env:             getEnv("ENV", "development"),
		DatabaseURL:     getEnv("DATABASE_URL", "postgres://team1:team1@localhost:5433/team1blog?sslmode=disable"),
		JWTSecret:       getEnv("JWT_SECRET", "dev-insecure-secret-change-me"),
		AccessTokenTTL:  durationEnv("ACCESS_TOKEN_TTL", 24*time.Hour),
		RefreshTokenTTL: durationEnv("REFRESH_TOKEN_TTL", 7*24*time.Hour),
		InviteTTL:       durationEnv("INVITE_TTL", 72*time.Hour),
		CORSOrigins:     splitCSV(getEnv("CORS_ORIGINS", "http://localhost:5173,http://localhost:5174")),
		FrontendURL:     getEnv("FRONTEND_URL", "http://localhost:5173"),
		AdminAppURL:     getEnv("ADMIN_APP_URL", "http://localhost:5174"),

		ResendAPIKey:   os.Getenv("RESEND_API_KEY"),
		ResendFromAddr: getEnv("RESEND_FROM_ADDR", "noreply@team1.blog"),

		CloudinaryCloudName: os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret: os.Getenv("CLOUDINARY_API_SECRET"),

		AvalancheRPCURL:       getEnv("AVALANCHE_RPC_URL", "https://api.avax-test.network/ext/bc/C/rpc"),
		AvalancheChainID:      int64Env("AVALANCHE_CHAIN_ID", 43113), // Fuji testnet; 43114 is C-Chain mainnet
		AvalancheTreasuryKey:  os.Getenv("AVALANCHE_TREASURY_PRIVATE_KEY"),
		AvalancheUSDCContract: getEnv("AVALANCHE_USDC_CONTRACT_ADDRESS", ""),

		SubstackPublicationURL: os.Getenv("SUBSTACK_PUBLICATION_URL"),
	}

	cfg.MockEmail = cfg.ResendAPIKey == ""
	cfg.MockImages = cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == ""
	cfg.MockPayments = cfg.AvalancheTreasuryKey == ""
	cfg.MockSubstack = cfg.SubstackPublicationURL == ""

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}

func int64Env(key string, fallback int64) int64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(v string) []string {
	var out []string
	cur := ""
	for _, r := range v {
		if r == ',' {
			if cur != "" {
				out = append(out, cur)
			}
			cur = ""
			continue
		}
		cur += string(r)
	}
	if cur != "" {
		out = append(out, cur)
	}
	return out
}
