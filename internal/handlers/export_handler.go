package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/user/gin-microservice-boilerplate/internal/usecase"
	"github.com/user/gin-microservice-boilerplate/models"
)

type ExportHandler struct {
	usecase usecase.ExportUsecase
}

func NewExportHandler(uc usecase.ExportUsecase) *ExportHandler {
	return &ExportHandler{usecase: uc}
}

// RegisterRoutes attaches export routes to the given router group.
func (h *ExportHandler) RegisterRoutes(rg *gin.RouterGroup) {
	exports := rg.Group("/exports")
	{
		exports.POST("", h.CreateExport)
	}
}

// CreateExport godoc
// @Summary      Request an export
// @Description  Generates a file synchronously or queues an asynchronous export job
// @Tags         exports
// @Accept       json
// @Produce      text/csv
// @Produce      image/jpeg
// @Produce      application/json
// @Param        request  body      models.ExportRequest  true  "Export request"
// @Success      200      {file}    file
// @Success      202      {object}  models.ExportResult
// @Failure      400      {object}  map[string]string
// @Failure      501      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /api/v1/exports [post]
func (h *ExportHandler) CreateExport(c *gin.Context) {
	var req models.ExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.usecase.RequestExport(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrInvalidExportRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, usecase.ErrUnsupportedExportStrategy), errors.Is(err, usecase.ErrAsyncModeNotImplemented):
			c.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if result.Mode == models.ExportModeAsync {
		c.JSON(http.StatusAccepted, result)
		return
	}

	if result.File == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "export file not generated"})
		return
	}

	fileName := result.File.FileName
	if fileName == "" {
		fileName = defaultExportFileName(req.Format)
	}

	contentType := result.File.ContentType
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	c.Data(http.StatusOK, contentType, result.File.Data)
}

func defaultExportFileName(format models.ExportFormat) string {
	return fmt.Sprintf("export.%s", format)
}
