package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/user/gin-microservice-boilerplate/internal/repository"
	"github.com/user/gin-microservice-boilerplate/models"
)

var (
	ErrInvalidExportRequest      = errors.New("invalid export request")
	ErrUnsupportedExportStrategy = errors.New("unsupported export strategy")
	ErrAsyncModeNotImplemented   = errors.New("async export mode not implemented")
)

type ExportStrategy interface {
	Generate(ctx context.Context, req models.ExportRequest) (*models.ExportFile, error)
}

type ExportStrategyKey struct {
	Format     models.ExportFormat
	SourceType models.ExportSourceType
}

func NewExportStrategyKey(format models.ExportFormat, sourceType models.ExportSourceType) ExportStrategyKey {
	return ExportStrategyKey{
		Format:     format,
		SourceType: sourceType,
	}
}

type ExportUsecase interface {
	RequestExport(ctx context.Context, req models.ExportRequest) (*models.ExportResult, error)
}

type exportUsecase struct {
	strategies map[ExportStrategyKey]ExportStrategy
	jobRepo    repository.ExportJobRepository
}

func NewExportUsecase(strategies map[ExportStrategyKey]ExportStrategy, jobRepo repository.ExportJobRepository) ExportUsecase {
	if strategies == nil {
		strategies = make(map[ExportStrategyKey]ExportStrategy)
	}

	return &exportUsecase{
		strategies: strategies,
		jobRepo:    jobRepo,
	}
}

func (u *exportUsecase) RequestExport(ctx context.Context, req models.ExportRequest) (*models.ExportResult, error) {
	req = normalizeExportRequest(req)

	if err := validateExportRequest(req); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidExportRequest, err)
	}

	if req.Mode == models.ExportModeAsync {
		return u.queueExport(ctx, req)
	}

	strategy, ok := u.strategies[NewExportStrategyKey(req.Format, req.SourceType)]
	if !ok {
		return nil, ErrUnsupportedExportStrategy
	}

	file, err := strategy.Generate(ctx, req)
	if err != nil {
		return nil, err
	}

	return &models.ExportResult{
		Mode:   models.ExportModeSync,
		Status: models.ExportJobStatusCompleted,
		File:   file,
	}, nil
}

func (u *exportUsecase) queueExport(ctx context.Context, req models.ExportRequest) (*models.ExportResult, error) {
	if u.jobRepo == nil {
		return nil, ErrAsyncModeNotImplemented
	}

	job := &models.ExportJob{
		Format:      req.Format,
		SourceType:  req.SourceType,
		Status:      models.ExportJobStatusQueued,
		RequestedAt: time.Now().UTC(),
	}

	jobID, err := u.jobRepo.Create(ctx, job)
	if err != nil {
		return nil, err
	}

	return &models.ExportResult{
		Mode:   models.ExportModeAsync,
		Status: models.ExportJobStatusQueued,
		JobID:  jobID,
	}, nil
}

func normalizeExportRequest(req models.ExportRequest) models.ExportRequest {
	req.Format = models.ExportFormat(strings.ToLower(string(req.Format)))
	req.SourceType = models.ExportSourceType(strings.ToLower(string(req.SourceType)))
	req.Mode = models.ExportMode(strings.ToLower(string(req.Mode)))

	req.FileName = strings.TrimSpace(req.FileName)
	req.Locale = strings.TrimSpace(req.Locale)
	req.Timezone = strings.TrimSpace(req.Timezone)
	req.Delimiter = strings.TrimSpace(req.Delimiter)

	if req.Mode == "" {
		req.Mode = models.ExportModeSync
	}
	if req.Delimiter == "" {
		req.Delimiter = ","
	}

	return req
}

func validateExportRequest(req models.ExportRequest) error {
	if req.Payload == nil {
		return errors.New("payload is required")
	}

	switch req.Format {
	case models.ExportFormatCSV, models.ExportFormatJPEG, models.ExportFormatPDF:
	default:
		return fmt.Errorf("invalid format %q", req.Format)
	}

	switch req.SourceType {
	case models.ExportSourceChart, models.ExportSourceTable:
	default:
		return fmt.Errorf("invalid source_type %q", req.SourceType)
	}

	switch req.Mode {
	case models.ExportModeSync, models.ExportModeAsync:
	default:
		return fmt.Errorf("invalid mode %q", req.Mode)
	}

	return nil
}
