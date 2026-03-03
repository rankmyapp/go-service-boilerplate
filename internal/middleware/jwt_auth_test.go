package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const testSecret = "ranksecret"

func TestJWTAuth_MissingAuthorizationHeader(t *testing.T) {
	router := setupProtectedRouter()

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"invalid_authorization_header"}`, rec.Body.String())
}

func TestJWTAuth_InvalidSignature(t *testing.T) {
	router := setupProtectedRouter()
	token := signToken(t, "other-secret", map[string]interface{}{
		"userId": 10,
		"name":   "User",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"invalid_token"}`, rec.Body.String())
}

func TestJWTAuth_ExpiredToken(t *testing.T) {
	router := setupProtectedRouter()
	token := signToken(t, testSecret, map[string]interface{}{
		"userId": 10,
		"name":   "User",
		"exp":    time.Now().Add(-time.Minute).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"token_expired"}`, rec.Body.String())
}

func TestJWTAuth_MissingUserIDClaim(t *testing.T) {
	router := setupProtectedRouter()
	token := signToken(t, testSecret, map[string]interface{}{
		"name": "User",
		"exp":  time.Now().Add(time.Hour).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"invalid_token_claims"}`, rec.Body.String())
}

func TestJWTAuth_ValidToken(t *testing.T) {
	router := setupProtectedRouter()
	token := signToken(t, testSecret, map[string]interface{}{
		"userId": 42,
		"name":   "Lucas",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"user_id":42}`, rec.Body.String())
}

func TestJWTAuth_ValidCompressedToken(t *testing.T) {
	router := setupProtectedRouter()
	token := signToken(t, testSecret, map[string]interface{}{
		"userId": 99,
		"name":   "Compressed",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	compressed := compressLikeAuthService(t, token)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+compressed)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"user_id":99}`, rec.Body.String())
}

func TestJWTAuth_ValidURLEncodedCompressedToken(t *testing.T) {
	router := setupProtectedRouter()
	token := signToken(t, testSecret, map[string]interface{}{
		"userId": 77,
		"name":   "Encoded",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	compressed := compressLikeAuthService(t, token)
	encoded := url.PathEscape(compressed)

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+encoded)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"user_id":77}`, rec.Body.String())
}

func TestJWTAuth_ValidCookieToken(t *testing.T) {
	router := setupProtectedRouter(WithTokenCookieName("authService_production_token"))
	token := signToken(t, testSecret, map[string]interface{}{
		"userId": 55,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.AddCookie(&http.Cookie{Name: "authService_production_token", Value: token})
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"user_id":55}`, rec.Body.String())
}

func setupProtectedRouter(opts ...JWTAuthOption) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(JWTAuth(testSecret, opts...))
	api.GET("/protected", func(c *gin.Context) {
		userID := c.GetInt(ContextKeyUserID)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})
	return r
}

func signToken(t *testing.T, secret string, claims map[string]interface{}) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}

func compressLikeAuthService(t *testing.T, token string) string {
	t.Helper()

	jsonPayload, err := json.Marshal(token)
	if err != nil {
		t.Fatalf("failed to marshal token: %v", err)
	}

	var compressed bytes.Buffer
	writer := gzip.NewWriter(&compressed)
	if _, err := writer.Write(jsonPayload); err != nil {
		t.Fatalf("failed to gzip token: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}

	return base64.StdEncoding.EncodeToString(compressed.Bytes())
}
