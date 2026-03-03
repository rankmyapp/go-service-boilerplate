package middleware

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	claimPermissions = "permissions"
)

// RequirePermissions checks if the token includes at least one of the required permissions.
// If no permissions are required, the request is allowed.
func RequirePermissions(required ...int) gin.HandlerFunc {
	requiredSet := normalizeRequiredPermissions(required)

	return func(c *gin.Context) {
		if len(requiredSet) == 0 {
			c.Next()
			return
		}

		permissions, ok := permissionsFromContext(c)
		if !ok {
			abortUnauthorized(c, "invalid_token_claims")
			return
		}

		if !hasAnyPermission(permissions, requiredSet) {
			abortUnauthorized(c, "insufficient_permissions")
			return
		}

		c.Next()
	}
}

func normalizeRequiredPermissions(required []int) map[int]struct{} {
	if len(required) == 0 {
		return map[int]struct{}{}
	}

	set := make(map[int]struct{}, len(required))
	for _, id := range required {
		if id > 0 {
			set[id] = struct{}{}
		}
	}
	return set
}

func permissionsFromContext(c *gin.Context) ([]int, bool) {
	rawClaims, ok := c.Get(ContextKeyTokenClaims)
	if !ok || rawClaims == nil {
		return nil, false
	}

	claims, ok := rawClaims.(jwt.MapClaims)
	if !ok {
		if mapClaims, ok := rawClaims.(map[string]interface{}); ok {
			claims = jwt.MapClaims(mapClaims)
		} else {
			return nil, false
		}
	}

	rawPermissions, ok := claims[claimPermissions]
	if !ok {
		return nil, false
	}

	return parsePermissions(rawPermissions)
}

func parsePermissions(value interface{}) ([]int, bool) {
	switch v := value.(type) {
	case []interface{}:
		return parseInterfaceSlice(v), true
	case []int:
		return append([]int(nil), v...), true
	case []int64:
		out := make([]int, 0, len(v))
		for _, item := range v {
			out = append(out, int(item))
		}
		return out, true
	case []float64:
		out := make([]int, 0, len(v))
		for _, item := range v {
			if int64(item) == int64(int(item)) {
				out = append(out, int(item))
			}
		}
		return out, true
	case []string:
		return parseStringSlice(v), true
	case string:
		return parseCommaSeparated(v), true
	default:
		return nil, false
	}
}

func parseInterfaceSlice(values []interface{}) []int {
	out := make([]int, 0, len(values))
	for _, item := range values {
		if id, ok := parsePermissionID(item); ok {
			out = append(out, id)
		}
	}
	return out
}

func parseStringSlice(values []string) []int {
	out := make([]int, 0, len(values))
	for _, item := range values {
		if id, ok := parsePermissionID(item); ok {
			out = append(out, id)
		}
	}
	return out
}

func parseCommaSeparated(raw string) []int {
	parts := strings.Split(raw, ",")
	out := make([]int, 0, len(parts))
	for _, part := range parts {
		if id, ok := parsePermissionID(strings.TrimSpace(part)); ok {
			out = append(out, id)
		}
	}
	return out
}

func parsePermissionID(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, v > 0
	case int32:
		return int(v), v > 0
	case int64:
		return int(v), v > 0
	case float64:
		if int64(v) != int64(int(v)) {
			return 0, false
		}
		return int(v), v > 0
	case float32:
		if int64(v) != int64(int(v)) {
			return 0, false
		}
		return int(v), v > 0
	case json.Number:
		parsed, err := strconv.Atoi(v.String())
		if err != nil || parsed <= 0 {
			return 0, false
		}
		return parsed, true
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(v))
		if err != nil || parsed <= 0 {
			return 0, false
		}
		return parsed, true
	default:
		return 0, false
	}
}

func hasAnyPermission(permissions []int, required map[int]struct{}) bool {
	if len(required) == 0 {
		return true
	}

	for _, permission := range permissions {
		if _, ok := required[permission]; ok {
			return true
		}
	}
	return false
}
