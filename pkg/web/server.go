package web

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/user/gin-microservice-boilerplate/docs"
)

// NewRouter creates and configures the Gin engine with Swagger.
func NewRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}

// Start runs the server on the given port.
func Start(r *gin.Engine, port int) error {
	addr := fmt.Sprintf(":%d", port)
	return r.Run(addr)
}
