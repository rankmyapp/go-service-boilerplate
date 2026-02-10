package repository

import (
	"context"

	"github.com/user/gin-microservice-boilerplate/models"
)

type ExportJobRepository interface {
	Create(ctx context.Context, job *models.ExportJob) (string, error)
	GetByID(ctx context.Context, id string) (*models.ExportJob, error)
	UpdateStatus(ctx context.Context, id string, status models.ExportJobStatus, location string, errMessage string) error
}
