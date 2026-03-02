package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	"github.com/user/gin-microservice-boilerplate/pkg/logging"
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
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	logger, err := logging.New(logging.Config{
		Level:     cfg.Log.Level,
		Format:    cfg.Log.Format,
		AddSource: cfg.Log.AddSource,
	})
	if err != nil {
		slog.Error("failed to initialize logger", "error", err)
		os.Exit(1)
	}
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mgr := db.NewConnectionManager()
	mgr.RegisterProvider("mongodb", dbMongo.Registration())

	for name, dbCfg := range cfg.Databases {
		params := map[string]string{
			"uri":      dbCfg.URI,
			"database": dbCfg.Database,
		}
		if err := mgr.Connect(ctx, name, dbCfg.Kind, params); err != nil {
			logger.Error("failed to connect database", "name", name, "error", err)
			os.Exit(1)
		}
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()
		if err := mgr.CloseAll(closeCtx); err != nil {
			logger.Error("failed to close database connections", "error", err)
		}
	}()

	primaryDB, err := db.GetTyped[*mongoDriver.Database](mgr, "primary")
	if err != nil {
		logger.Error("primary database not configured", "error", err)
		os.Exit(1)
	}

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

	router := web.NewRouterWithLogger(logger, cfg.Server.CORSAllowedOrigins)
	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	exportHandler.RegisterRoutes(api)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down - waiting for active requests")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}
}
