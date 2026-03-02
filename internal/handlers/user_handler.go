package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/user/gin-microservice-boilerplate/internal/middleware"
	"github.com/user/gin-microservice-boilerplate/internal/usecase"
	"github.com/user/gin-microservice-boilerplate/models"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

// RegisterRoutes attaches user routes to the given router group.
func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup, perms UserRoutePermissions) {
	users := rg.Group("/users")
	{
		users.POST("", middleware.RequirePermissions(perms.Create...), h.CreateUser)
		users.GET("", middleware.RequirePermissions(perms.List...), h.GetAllUsers)
		users.GET("/:id", middleware.RequirePermissions(perms.Get...), h.GetUserByID)
		users.PUT("/:id", middleware.RequirePermissions(perms.Update...), h.UpdateUser)
		users.DELETE("/:id", middleware.RequirePermissions(perms.Delete...), h.DeleteUser)
	}
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with name and email
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      models.User  true  "User object"
// @Security     BearerAuth
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.usecase.CreateUser(c.Request.Context(), &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// GetAllUsers godoc
// @Summary      List all users
// @Description  Get all users
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.User
// @Failure      500  {object}  map[string]string
// @Router       /users [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.usecase.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Get a single user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  models.User
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.usecase.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Update user data by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path      string       true  "User ID"
// @Param        user  body      models.User  true  "User object"
// @Security     BearerAuth
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.ID = id
	if err := h.usecase.UpdateUser(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Security     BearerAuth
// @Success      200  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.DeleteUser(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
