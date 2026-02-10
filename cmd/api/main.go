package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/gin-microservice-boilerplate/config"
	_ "github.com/user/gin-microservice-boilerplate/docs"
	"github.com/user/gin-microservice-boilerplate/internal/handlers"
	repoMongo "github.com/user/gin-microservice-boilerplate/internal/repository/mongo"
	"github.com/user/gin-microservice-boilerplate/internal/usecase"
	"github.com/user/gin-microservice-boilerplate/models"
	"github.com/user/gin-microservice-boilerplate/pkg/db"
	dbMongo "github.com/user/gin-microservice-boilerplate/pkg/db/mongo"
	csvExport "github.com/user/gin-microservice-boilerplate/pkg/export/csv"
	jpegExport "github.com/user/gin-microservice-boilerplate/pkg/export/jpeg"
	pdfExport "github.com/user/gin-microservice-boilerplate/pkg/export/pdf"
	"github.com/user/gin-microservice-boilerplate/pkg/web"

	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

// @title           Gin Microservice Boilerplate API
// @version         1.0
// @description     A clean architecture boilerplate with Gin and MongoDB
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx := context.Background()

	mgr := db.NewConnectionManager()
	mgr.RegisterProvider("mongodb", dbMongo.Registration())

	for name, dbCfg := range cfg.Databases {
		params := map[string]string{
			"uri":      dbCfg.URI,
			"database": dbCfg.Database,
		}
		if err := mgr.Connect(ctx, name, dbCfg.Kind, params); err != nil {
			log.Fatalf("failed to connect database %q: %v", name, err)
		}
	}
	defer mgr.CloseAll(ctx)

	primaryConn, err := mgr.Get("primary")
	if err != nil {
		log.Fatalf("primary database not configured: %v", err)
	}
	primaryDB := primaryConn.(*mongoDriver.Database)

	userRepo := repoMongo.NewUserRepo(primaryDB)
	userUC := usecase.NewUserUsecase(userRepo)
	userHandler := handlers.NewUserHandler(userUC)

	exportStrategies := map[usecase.ExportStrategyKey]usecase.ExportStrategy{
		usecase.NewExportStrategyKey(models.ExportFormatCSV, models.ExportSourceChart):  csvExport.NewChartStrategy(),
		usecase.NewExportStrategyKey(models.ExportFormatCSV, models.ExportSourceTable):  csvExport.NewTableStrategy(),
		usecase.NewExportStrategyKey(models.ExportFormatJPEG, models.ExportSourceChart): jpegExport.NewChartStrategy(),
		usecase.NewExportStrategyKey(models.ExportFormatPDF, models.ExportSourceTable):  pdfExport.NewTableStrategy(),
	}
	exportUC := usecase.NewExportUsecase(exportStrategies, nil)
	exportHandler := handlers.NewExportHandler(exportUC)

	router := web.NewRouter()
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	exportHandler.RegisterRoutes(api)

	go func() {
		if err := web.Start(router, cfg.Server.Port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down...")
}
