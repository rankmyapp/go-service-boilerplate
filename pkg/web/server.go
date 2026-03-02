package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/time/rate"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/user/gin-microservice-boilerplate/docs"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// NewRouter creates and configures the Gin engine with Swagger.
func NewRouter() *gin.Engine {
	return NewRouterWithLogger(slog.Default(), nil)
}

// NewRouterWithLogger creates and configures the Gin engine with Swagger and structured request logging.
func NewRouterWithLogger(logger *slog.Logger, allowedOrigins []string) *gin.Engine {
	if logger == nil {
		logger = slog.Default()
	}

	r := gin.New()
	r.Use(requestIDMiddleware())
	r.Use(rateLimitMiddleware(10, 20))
	r.Use(requestLoggerMiddleware(logger))
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error("panic recovered",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"client_ip", c.ClientIP(),
			"error", fmt.Sprint(recovered),
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
	}))
	r.Use(corsMiddleware(allowedOrigins))

	healthHandler := func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) }
	r.GET("/health", healthHandler)
	r.GET("/api/v1/health", healthHandler)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.OPTIONS("/*any", func(c *gin.Context) { c.Status(204) })
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"error":  "not_found",
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"host":   c.Request.Host,
		})
	})
	return r
}

// Start runs the server on the given port.
func Start(r *gin.Engine, port int) error {
	addr := fmt.Sprintf(":%d", port)
	return r.Run(addr)
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	const (
		allowMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
		allowHeaders = "Origin,Content-Type,Accept,Authorization,X-Requested-With,X-Request-ID,X-XSRF-Token"
		exposeHeader = "Content-Disposition,Content-Type,X-Request-ID,Content-Length"
		maxAge       = "43200"
	)

	originSet := make(map[string]struct{}, len(allowedOrigins))
	for _, orig := range allowedOrigins {
		originSet[strings.ToLower(strings.TrimSpace(orig))] = struct{}{}
	}

	_, wildcardAllowed := originSet["*"]
	allowCredentials := !wildcardAllowed && len(originSet) > 0

	return func(c *gin.Context) {
		origin := strings.TrimSpace(c.GetHeader("Origin"))
		allowed := false

		if origin != "" {
			_, ok := originSet[strings.ToLower(origin)]
			allowed = wildcardAllowed || ok
		}

		if allowed {
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", allowMethods)
			c.Header("Access-Control-Expose-Headers", exposeHeader)
			c.Header("Access-Control-Max-Age", maxAge)
			if allowCredentials {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			requestedHeaders := strings.TrimSpace(c.GetHeader("Access-Control-Request-Headers"))
			if requestedHeaders != "" {
				c.Header("Access-Control-Allow-Headers", requestedHeaders)
			} else {
				c.Header("Access-Control-Allow-Headers", allowHeaders)
			}
		}

		if c.Request.Method == http.MethodOptions {
			if allowed {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.AbortWithStatus(http.StatusForbidden)
			}
			return
		}
		c.Next()
	}
}

func requestLoggerMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path = path + "?" + c.Request.URL.RawQuery
		}

		status := c.Writer.Status()
		attrs := []any{
			"request_id", c.GetString(string(RequestIDKey)),
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		}

		if ua := strings.TrimSpace(c.Request.UserAgent()); ua != "" {
			attrs = append(attrs, "user_agent", ua)
		}
		if len(c.Errors) > 0 {
			attrs = append(attrs, "errors", c.Errors.String())
		}

		switch {
		case status >= http.StatusInternalServerError:
			logger.Error("http request", attrs...)
		case status >= http.StatusBadRequest:
			logger.Warn("http request", attrs...)
		default:
			logger.Info("http request", attrs...)
		}
	}
}

func rateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		key := c.ClientIP()
		val, _ := limiters.LoadOrStore(key, rate.NewLimiter(rate.Limit(rps), burst))
		limiter := val.(*rate.Limiter)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limits exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set(string(RequestIDKey), requestID)
		c.Header("X-Request-ID", requestID)

		ctx := context.WithValue(c.Request.Context(), RequestIDKey, requestID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
