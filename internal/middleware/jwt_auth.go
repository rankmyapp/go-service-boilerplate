package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ContextKeyTokenClaims = "auth_token_claims"
	ContextKeyUserID      = "auth_user_id"
)

type jwtAuthConfig struct {
	tokenCookieName string
}

type JWTAuthOption func(*jwtAuthConfig)

// WithTokenCookieName enables cookie-based token extraction.
// If name is empty, cookie fallback is disabled.
func WithTokenCookieName(name string) JWTAuthOption {
	return func(cfg *jwtAuthConfig) {
		cfg.tokenCookieName = strings.TrimSpace(name)
	}
}

// JWTAuth validates Bearer tokens signed by an auth service using HS256.
func JWTAuth(secret string, opts ...JWTAuthOption) gin.HandlerFunc {
	secretBytes := []byte(strings.TrimSpace(secret))

	cfg := jwtAuthConfig{}
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}

	return func(c *gin.Context) {
		tokenString := ""

		// 1) try Authorization Bearer
		if h := c.GetHeader("Authorization"); h != "" {
			t, err := extractBearerToken(h)
			if err == nil {
				tokenString = t
			}
		}

		// 2) fallback to HttpOnly cookie (if enabled)
		if tokenString == "" && cfg.tokenCookieName != "" {
			if ck, err := c.Cookie(cfg.tokenCookieName); err == nil {
				tokenString = ck
			}
		}

		if tokenString == "" {
			abortUnauthorized(c, "invalid_authorization_header")
			return
		}

		tokenString = normalizeTokenString(tokenString)
		tokenString = unwrapCompressedToken(tokenString)

		claims, err := parseTokenClaims(tokenString, secretBytes)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				abortUnauthorized(c, "token_expired")
				return
			}
			abortUnauthorized(c, "invalid_token")
			return
		}

		userID, ok := parseUserID(claims["userId"])
		if !ok || userID <= 0 {
			abortUnauthorized(c, "invalid_token_claims")
			return
		}

		c.Set(ContextKeyTokenClaims, claims)
		c.Set(ContextKeyUserID, userID)
		c.Next()
	}
}

func normalizeTokenString(token string) string {
	decoded, err := url.PathUnescape(strings.TrimSpace(token))
	if err != nil {
		return token
	}
	decoded = strings.TrimSpace(decoded)
	if decoded == "" {
		return token
	}
	return decoded
}

func extractBearerToken(header string) (string, error) {
	parts := strings.SplitN(strings.TrimSpace(header), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid authorization header")
	}

	token := strings.TrimSpace(parts[1])
	if token == "" {
		return "", errors.New("empty bearer token")
	}
	return token, nil
}

func parseTokenClaims(tokenString string, secret []byte) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return nil, errors.New("missing exp claim")
	}

	return claims, nil
}

func parseUserID(value interface{}) (int, bool) {
	switch v := value.(type) {
	case float64:
		parsed := int(v)
		return parsed, float64(parsed) == v
	case int:
		return v, true
	case int32:
		return int(v), true
	case int64:
		return int(v), true
	case json.Number:
		parsed, err := strconv.Atoi(v.String())
		if err != nil {
			return 0, false
		}
		return parsed, true
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func unwrapCompressedToken(token string) string {
	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return token
	}

	reader, err := gzip.NewReader(bytes.NewReader(decoded))
	if err != nil {
		return token
	}
	defer reader.Close()

	payload, err := io.ReadAll(reader)
	if err != nil {
		return token
	}

	var parsed string
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return token
	}

	parsed = strings.TrimSpace(parsed)
	if parsed == "" {
		return token
	}
	return parsed
}

func abortUnauthorized(c *gin.Context, code string) {
	c.Header("WWW-Authenticate", "Bearer")
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": code})
}
