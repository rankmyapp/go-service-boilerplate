package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/user/gin-microservice-boilerplate/models"
)

type MockExportStrategy struct {
	mock.Mock
}

func (m *MockExportStrategy) Generate(ctx context.Context, req models.ExportRequest) (*models.ExportFile, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExportFile), args.Error(1)
}

type MockExportJobRepository struct {
	mock.Mock
}

func (m *MockExportJobRepository) Create(ctx context.Context, job *models.ExportJob) (string, error) {
	args := m.Called(ctx, job)
	return args.String(0), args.Error(1)
}

func (m *MockExportJobRepository) GetByID(ctx context.Context, id string) (*models.ExportJob, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.ExportJob), args.Error(1)
}

func (m *MockExportJobRepository) UpdateStatus(ctx context.Context, id string, status models.ExportJobStatus, location string, errMessage string) error {
	args := m.Called(ctx, id, status, location, errMessage)
	return args.Error(0)
}

func validExportRequest() models.ExportRequest {
	return models.ExportRequest{
		Format:     models.ExportFormatCSV,
		SourceType: models.ExportSourceTable,
		Payload:    []byte(`{"columns":["name"],"rows":[["alice"]]}`),
	}
}

func TestRequestExport_SyncSuccess(t *testing.T) {
	strategy := new(MockExportStrategy)
	uc := NewExportUsecase(
		map[ExportStrategyKey]ExportStrategy{
			NewExportStrategyKey(models.ExportFormatCSV, models.ExportSourceTable): strategy,
		},
		nil,
	)

	strategy.On("Generate", mock.Anything, mock.Anything).Return(&models.ExportFile{
		FileName:    "report.csv",
		ContentType: "text/csv",
		Data:        []byte("name\nalice\n"),
	}, nil)

	result, err := uc.RequestExport(context.Background(), validExportRequest())
	assert.NoError(t, err)
	assert.Equal(t, models.ExportModeSync, result.Mode)
	assert.Equal(t, models.ExportJobStatusCompleted, result.Status)
	assert.NotNil(t, result.File)
	strategy.AssertExpectations(t)
}

func TestRequestExport_UnsupportedStrategy(t *testing.T) {
	uc := NewExportUsecase(nil, nil)

	_, err := uc.RequestExport(context.Background(), validExportRequest())
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrUnsupportedExportStrategy))
}

func TestRequestExport_AsyncWithoutRepo(t *testing.T) {
	uc := NewExportUsecase(nil, nil)
	req := validExportRequest()
	req.Mode = models.ExportModeAsync

	_, err := uc.RequestExport(context.Background(), req)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrAsyncModeNotImplemented))
}

func TestRequestExport_AsyncWithRepo(t *testing.T) {
	jobRepo := new(MockExportJobRepository)
	uc := NewExportUsecase(nil, jobRepo)
	req := validExportRequest()
	req.Mode = models.ExportModeAsync

	jobRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.ExportJob")).Return("job-123", nil)

	result, err := uc.RequestExport(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, models.ExportModeAsync, result.Mode)
	assert.Equal(t, models.ExportJobStatusQueued, result.Status)
	assert.Equal(t, "job-123", result.JobID)
	jobRepo.AssertExpectations(t)
}

func TestRequestExport_InvalidRequest(t *testing.T) {
	uc := NewExportUsecase(nil, nil)
	req := validExportRequest()
	req.Format = "xlsx"

	_, err := uc.RequestExport(context.Background(), req)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidExportRequest))
}
