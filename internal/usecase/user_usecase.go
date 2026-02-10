package usecase

import (
	"context"

	"github.com/user/gin-microservice-boilerplate/internal/repository"
	"github.com/user/gin-microservice-boilerplate/models"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, user *models.User) (string, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: repo}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *models.User) (string, error) {
	return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	return u.userRepo.GetAll(ctx)
}

func (u *userUsecase) UpdateUser(ctx context.Context, user *models.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) DeleteUser(ctx context.Context, id string) error {
	return u.userRepo.Delete(ctx, id)
}
