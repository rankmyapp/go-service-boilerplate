package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/user/gin-microservice-boilerplate/models"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	user := &models.User{Name: "Alice", Email: "alice@test.com"}
	mockRepo.On("Create", mock.Anything, user).Return("abc123", nil)

	id, err := uc.CreateUser(context.Background(), user)
	assert.NoError(t, err)
	assert.Equal(t, "abc123", id)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Error(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	user := &models.User{Name: "Alice", Email: "alice@test.com"}
	mockRepo.On("Create", mock.Anything, user).Return("", errors.New("db error"))

	id, err := uc.CreateUser(context.Background(), user)
	assert.Error(t, err)
	assert.Empty(t, id)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	expected := &models.User{ID: "abc123", Name: "Alice", Email: "alice@test.com"}
	mockRepo.On("GetByID", mock.Anything, "abc123").Return(expected, nil)

	user, err := uc.GetUserByID(context.Background(), "abc123")
	assert.NoError(t, err)
	assert.Equal(t, "Alice", user.Name)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	mockRepo.On("GetByID", mock.Anything, "notfound").Return(nil, errors.New("not found"))

	user, err := uc.GetUserByID(context.Background(), "notfound")
	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	expected := []*models.User{
		{ID: "1", Name: "Alice"},
		{ID: "2", Name: "Bob"},
	}
	mockRepo.On("GetAll", mock.Anything).Return(expected, nil)

	users, err := uc.GetAllUsers(context.Background())
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	user := &models.User{ID: "abc123", Name: "Alice Updated"}
	mockRepo.On("Update", mock.Anything, user).Return(nil)

	err := uc.UpdateUser(context.Background(), user)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	uc := NewUserUsecase(mockRepo)

	mockRepo.On("Delete", mock.Anything, "abc123").Return(nil)

	err := uc.DeleteUser(context.Background(), "abc123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
