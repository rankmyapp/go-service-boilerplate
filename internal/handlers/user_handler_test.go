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

	"github.com/user/gin-microservice-boilerplate/models"
)

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) CreateUser(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func (m *MockUserUsecase) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserUsecase) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserUsecase) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserUsecase) DeleteUser(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter(handler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api/v1")
	handler.RegisterRoutes(api)
	return r
}

func TestCreateUser_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	user := models.User{Name: "Alice", Email: "alice@test.com"}
	body, _ := json.Marshal(user)

	mockUC.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return("id123", nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "id123", resp["id"])
	mockUC.AssertExpectations(t)
}

func TestCreateUser_BadRequest(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_InternalError(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	user := models.User{Name: "Alice", Email: "alice@test.com"}
	body, _ := json.Marshal(user)

	mockUC.On("CreateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return("", errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockUC.AssertExpectations(t)
}

func TestGetAllUsers_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	users := []*models.User{
		{ID: "1", Name: "Alice", Email: "alice@test.com"},
		{ID: "2", Name: "Bob", Email: "bob@test.com"},
	}
	mockUC.On("GetAllUsers", mock.Anything).Return(users, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp []models.User
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	mockUC.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	user := &models.User{ID: "abc123", Name: "Alice", Email: "alice@test.com"}
	mockUC.On("GetUserByID", mock.Anything, "abc123").Return(user, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/abc123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.User
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "Alice", resp.Name)
	mockUC.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	mockUC.On("GetUserByID", mock.Anything, "notfound").Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/users/notfound", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockUC.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	user := models.User{Name: "Alice Updated", Email: "alice@test.com"}
	body, _ := json.Marshal(user)

	mockUC.On("UpdateUser", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/users/abc123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	mockUC := new(MockUserUsecase)
	handler := NewUserHandler(mockUC)
	router := setupRouter(handler)

	mockUC.On("DeleteUser", mock.Anything, "abc123").Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/users/abc123", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUC.AssertExpectations(t)
}
