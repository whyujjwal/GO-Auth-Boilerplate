package handlers

import (
	"auth/internal/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserCreator interface {
	Create(user *models.User) (*models.User, error)
}

type UserHandler struct {
	userRepo UserCreator
	logger   zerolog.Logger
}

func NewUserHandler(ur UserCreator, l zerolog.Logger) *UserHandler {
	return &UserHandler{
		userRepo: ur,
		logger:   l,
	}
}
func username(email string) string {
	parts := strings.Split(email, "@")
	return parts[0] + "." + parts[1]
}
func (h *UserHandler) Create(c *gin.Context) {
	type CreateUserRequest struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	user := models.User{
		Username: username(req.Email),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	createdUser, err := h.userRepo.Create(&user)
	if err != nil {
		h.logger.Error().Err(err).Msg("Failed to create user")
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(201, createdUser)
}
