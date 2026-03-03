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
	Databases map[string]DatabaseConfig
	Auth      AuthConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Kind     string
	URI      string
	Database string
}

type AuthConfig struct {
	Enabled         bool
	JWTSecret       string
	TokenCookieName string
	Permissions     AuthPermissionsConfig
}

type AuthPermissionsConfig struct {
	Default []int
	Users   UserPermissionsConfig
	Exports ExportPermissionsConfig
}

type UserPermissionsConfig struct {
	Create []int
	List   []int
	Get    []int
	Update []int
	Delete []int
}

type ExportPermissionsConfig struct {
	Create []int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port, err := strconv.Atoi(getEnvMany([]string{"SERVER_PORT", "PORT"}, "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid SERVER_PORT: %w", err)
	}

	authEnabled, err := parseBoolEnv("AUTH_ENABLED", false)
	if err != nil {
		return nil, fmt.Errorf("invalid AUTH_ENABLED: %w", err)
	}

	authCfg := AuthConfig{
		Enabled: authEnabled,
	}

	if authEnabled {
		jwtSecret := strings.TrimSpace(getEnvMany([]string{"JWT_SECRET", "API_SECRET"}, ""))
		if jwtSecret == "" {
			return nil, fmt.Errorf("missing JWT secret: set JWT_SECRET or API_SECRET")
		}

		tokenCookieName := "authService_production_token"
		if raw, ok := os.LookupEnv("AUTH_TOKEN_COOKIE_NAME"); ok {
			tokenCookieName = strings.TrimSpace(raw)
		}

		requiredPermissions, err := parseIntListEnv("AUTH_REQUIRED_PERMISSIONS")
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_REQUIRED_PERMISSIONS: %w", err)
		}

		usersCreate, err := parseIntListEnvWithFallback("AUTH_PERMISSIONS_USERS_CREATE", requiredPermissions)
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_PERMISSIONS_USERS_CREATE: %w", err)
		}
		usersList, err := parseIntListEnvWithFallback("AUTH_PERMISSIONS_USERS_LIST", requiredPermissions)
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_PERMISSIONS_USERS_LIST: %w", err)
		}
		usersGet, err := parseIntListEnvWithFallback("AUTH_PERMISSIONS_USERS_GET", requiredPermissions)
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_PERMISSIONS_USERS_GET: %w", err)
		}
		usersUpdate, err := parseIntListEnvWithFallback("AUTH_PERMISSIONS_USERS_UPDATE", requiredPermissions)
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_PERMISSIONS_USERS_UPDATE: %w", err)
		}
		usersDelete, err := parseIntListEnvWithFallback("AUTH_PERMISSIONS_USERS_DELETE", requiredPermissions)
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_PERMISSIONS_USERS_DELETE: %w", err)
		}
		exportsCreate, err := parseIntListEnvWithFallback("AUTH_PERMISSIONS_EXPORTS_CREATE", requiredPermissions)
		if err != nil {
			return nil, fmt.Errorf("invalid AUTH_PERMISSIONS_EXPORTS_CREATE: %w", err)
		}

		authCfg.JWTSecret = jwtSecret
		authCfg.TokenCookieName = tokenCookieName
		authCfg.Permissions = AuthPermissionsConfig{
			Default: requiredPermissions,
			Users: UserPermissionsConfig{
				Create: usersCreate,
				List:   usersList,
				Get:    usersGet,
				Update: usersUpdate,
				Delete: usersDelete,
			},
			Exports: ExportPermissionsConfig{
				Create: exportsCreate,
			},
		}
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: port,
		},
		Databases: map[string]DatabaseConfig{
			"primary": {
				Kind:     getEnv("DB_PRIMARY_KIND", "mongodb"),
				URI:      getEnv("DB_PRIMARY_URI", "mongodb://localhost:27017"),
				Database: getEnv("DB_PRIMARY_DATABASE", "appdb"),
			},
		},
		Auth: authCfg,
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvMany(keys []string, fallback string) string {
	for _, key := range keys {
		if val := os.Getenv(key); val != "" {
			return val
		}
	}
	return fallback
}

func parseBoolEnv(key string, fallback bool) (bool, error) {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return fallback, nil
	}

	switch strings.ToLower(raw) {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("must be true/false/1/0")
	}
}

func parseIntListEnv(key string) ([]int, error) {
	raw := strings.TrimSpace(getEnv(key, ""))
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		value := strings.TrimSpace(part)
		if value == "" {
			continue
		}
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return nil, err
		}
		if parsed > 0 {
			out = append(out, parsed)
		}
	}
	return out, nil
}

func parseIntListEnvWithFallback(key string, fallback []int) ([]int, error) {
	values, err := parseIntListEnv(key)
	if err != nil {
		return nil, err
	}
	if len(values) == 0 {
		return fallback, nil
	}
	return values, nil
}
