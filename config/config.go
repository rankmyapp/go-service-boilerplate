package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	Log       LogConfig
	Auth      AuthConfig
	Databases map[string]DatabaseConfig
}

type ServerConfig struct {
	Port               int
	CORSAllowedOrigins []string
}

type LogConfig struct {
	Level     string
	Format    string
	AddSource bool
}

type AuthConfig struct {
	Enabled     bool
	JWTSecret   string
	JWTIssuer   string
	JWTAudience string
}

type DatabaseConfig struct {
	Kind     string
	URI      string
	Database string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	corsOrigins := splitCSVEnv(getEnv("CORS_ALLOWED_ORIGINS", ""))

	logLevel := strings.ToLower(strings.TrimSpace(getEnv("LOG_LEVEL", "info")))
	switch logLevel {
	case "debug", "info", "warn", "error":
	default:
		return nil, fmt.Errorf("invalid LOG_LEVEL: %q (valid: debug, info, warn, error)", logLevel)
	}

	logFormat := strings.ToLower(strings.TrimSpace(getEnv("LOG_FORMAT", "json")))
	switch logFormat {
	case "json", "text":
	default:
		return nil, fmt.Errorf("invalid LOG_FORMAT: %q (valid: json, text)", logFormat)
	}

	logAddSource, err := strconv.ParseBool(getEnv("LOG_ADD_SOURCE", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid LOG_ADD_SOURCE: %w", err)
	}

	authEnabled, err := strconv.ParseBool(getEnv("AUTH_ENABLED", "false"))
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_ENABLED: %w", err)
	}
	authSecret := strings.TrimSpace(getEnv("AUTH_JWT_SECRET", ""))
	if authEnabled && authSecret == "" {
		return nil, fmt.Errorf("AUTH_JWT_SECRET is required when AUTH_ENABLED=true")
	}

	dbNames := splitCSVEnv(getEnv("DB_CONNECTIONS", "primary"))
	if len(dbNames) == 0 {
		return nil, fmt.Errorf("DB_CONNECTIONS must include at least one name")
	}
	databases := make(map[string]DatabaseConfig, len(dbNames))
	for _, name := range dbNames {
		key := normalizeEnvKey(name)
		databases[name] = DatabaseConfig{
			Kind:     getEnv(fmt.Sprintf("DB_%s_KIND", key), "mongodb"),
			URI:      getEnv(fmt.Sprintf("DB_%s_URI", key), "mongodb://localhost:27017"),
			Database: getEnv(fmt.Sprintf("DB_%s_DATABASE", key), "appdb"),
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:               port,
			CORSAllowedOrigins: corsOrigins,
		},
		Log: LogConfig{
			Level:     logLevel,
			Format:    logFormat,
			AddSource: logAddSource,
		},
		Auth: AuthConfig{
			Enabled:     authEnabled,
			JWTSecret:   authSecret,
			JWTIssuer:   strings.TrimSpace(getEnv("AUTH_JWT_ISSUER", "")),
			JWTAudience: strings.TrimSpace(getEnv("AUTH_JWT_AUDIENCE", "")),
		},
		Databases: databases,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func splitCSVEnv(raw string) []string {
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		item := strings.TrimSpace(part)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func normalizeEnvKey(name string) string {
	clean := strings.TrimSpace(name)
	clean = strings.ReplaceAll(clean, "-", "_")
	clean = strings.ReplaceAll(clean, ".", "_")
	clean = strings.ReplaceAll(clean, " ", "_")
	return strings.ToUpper(clean)
}
