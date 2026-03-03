package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestRequirePermissions_AllowsWhenNoRequired(t *testing.T) {
	router := setupPermissionRouter(nil)
	token := signPermissionToken(t, testSecret, map[string]interface{}{
		"userId": 10,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func TestRequirePermissions_MissingPermissionsClaim(t *testing.T) {
	router := setupPermissionRouter([]int{339})
	token := signPermissionToken(t, testSecret, map[string]interface{}{
		"userId": 10,
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"invalid_token_claims"}`, rec.Body.String())
}

func TestRequirePermissions_DeniesWhenMissingRequired(t *testing.T) {
	router := setupPermissionRouter([]int{339})
	token := signPermissionToken(t, testSecret, map[string]interface{}{
		"userId":      10,
		"exp":         time.Now().Add(time.Hour).Unix(),
		"permissions": []int{200, 201, 202},
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"insufficient_permissions"}`, rec.Body.String())
}

func TestRequirePermissions_AllowsWhenHasRequired(t *testing.T) {
	router := setupPermissionRouter([]int{339})
	token := signPermissionToken(t, testSecret, map[string]interface{}{
		"userId":      10,
		"exp":         time.Now().Add(time.Hour).Unix(),
		"permissions": []int{339, 200},
	})

	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

func setupPermissionRouter(required []int) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	api.Use(JWTAuth(testSecret))
	api.Use(RequirePermissions(required...))
	api.GET("/protected", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}

func signPermissionToken(t *testing.T, secret string, claims map[string]interface{}) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return signed
}
