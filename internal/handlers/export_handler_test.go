package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/user/gin-microservice-boilerplate/internal/usecase"
	"github.com/user/gin-microservice-boilerplate/models"
)

type MockExportUsecase struct {
	mock.Mock
}

func (m *MockExportUsecase) RequestExport(ctx context.Context, req models.ExportRequest) (*models.ExportResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExportResult), args.Error(1)
}

func setupExportRouter(handler *ExportHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api, ExportRoutePermissions{})
	return r
}

func TestCreateExport_SyncSuccess(t *testing.T) {
	mockUC := new(MockExportUsecase)
	handler := NewExportHandler(mockUC)
	router := setupExportRouter(handler)

	body := map[string]interface{}{
		"format":      "csv",
		"source_type": "table",
		"payload":     map[string]interface{}{"columns": []string{"name"}, "rows": [][]string{{"alice"}}},
	}
	data, _ := json.Marshal(body)

	mockUC.On("RequestExport", mock.Anything, mock.AnythingOfType("models.ExportRequest")).Return(&models.ExportResult{
		Mode:   models.ExportModeSync,
		Status: models.ExportJobStatusCompleted,
		File: &models.ExportFile{
			FileName:    "report.csv",
			ContentType: "text/csv",
			Data:        []byte("name\nalice\n"),
		},
	}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/exports", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment;")
	assert.Equal(t, "name\nalice\n", w.Body.String())
	mockUC.AssertExpectations(t)
}

func TestCreateExport_BadRequest(t *testing.T) {
	mockUC := new(MockExportUsecase)
	handler := NewExportHandler(mockUC)
	router := setupExportRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/exports", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateExport_NotImplemented(t *testing.T) {
	mockUC := new(MockExportUsecase)
	handler := NewExportHandler(mockUC)
	router := setupExportRouter(handler)

	body := map[string]interface{}{
		"format":      "pdf",
		"source_type": "table",
		"payload":     map[string]interface{}{"columns": []string{"name"}, "rows": [][]string{{"alice"}}},
	}
	data, _ := json.Marshal(body)

	mockUC.On("RequestExport", mock.Anything, mock.AnythingOfType("models.ExportRequest")).Return(nil, usecase.ErrUnsupportedExportStrategy)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/exports", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotImplemented, w.Code)
	mockUC.AssertExpectations(t)
}

func TestCreateExport_InternalError(t *testing.T) {
	mockUC := new(MockExportUsecase)
	handler := NewExportHandler(mockUC)
	router := setupExportRouter(handler)

	body := map[string]interface{}{
		"format":      "csv",
		"source_type": "table",
		"payload":     map[string]interface{}{"columns": []string{"name"}, "rows": [][]string{{"alice"}}},
	}
	data, _ := json.Marshal(body)

	mockUC.On("RequestExport", mock.Anything, mock.AnythingOfType("models.ExportRequest")).Return(nil, errors.New("unexpected failure"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/exports", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}
