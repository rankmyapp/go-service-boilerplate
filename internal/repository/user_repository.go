package repository

import (
	"context"

	"github.com/user/gin-microservice-boilerplate/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (string, error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}
