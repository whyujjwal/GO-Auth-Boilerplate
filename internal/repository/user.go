package repository

import (
	"auth/internal/models"
	"auth/pkg/database"
	"auth/pkg/utils"
	"fmt"
)

type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	// GetByID(id uint) (*models.User, error)
	// GetByEmail(email string) (*models.User, error)
	// Update(user *models.User) (*models.User, error)
}

type userRepository struct {
	db database.Database
}

func NewUserRepository(db database.Database) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) (*models.User, error) {
	var err error
	user.PasswordHash, err = utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}

	result := r.db.Create(user)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to create user: %v", result.Error)
	}

	return user, nil
}
