package config

import (
	"os"
	"strings"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	AppEnv             string
	DatabaseURL        string
	EnableMetricWorker bool
	AllowedOrigins     []string

	JWTSecret           string
	GithubWebhookSecret string

	GithubClientID     string
	GithubClientSecret string

	GithubToken string
	GithubPAT   string
	GistID      string
	GithubOwner string
	GithubRepo  string
}

func Load() *Config {
	_ = godotenv.Load()

	logger := middleware.LogGet()

	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "")
	if dbURL == "" {
		logger.Warn().Msg("DATABASE_URL is not set in environment variables")
	}

	originsRaw := getEnv("ALLOWED_ORIGINS", "http://localhost:5173")
	origins := strings.Split(originsRaw, ",")

	return &Config{
		Port:               port,
		AppEnv:             getEnv("APP_ENV", "development"),
		DatabaseURL:        dbURL,
		EnableMetricWorker: getEnvAsBool("ENABLE_METRICS_TOGGLE", false),
		AllowedOrigins:     origins,

		JWTSecret:           getEnv("JWT_SECRET", "default_dev_secret"),
		GithubWebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),

		GithubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),

		GithubToken: getEnv("GITHUB_TOKEN", ""),
		GithubPAT:   getEnv("GITHUB_PAT", ""),
		GistID:      getEnv("GIST_ID", ""),
		GithubOwner: getEnv("GITHUB_OWNER", ""),
		GithubRepo:  getEnv("GITHUB_REPO", ""),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return defaultValue
	}

	return strings.ToLower(value) == "true" || value == "1"
}
