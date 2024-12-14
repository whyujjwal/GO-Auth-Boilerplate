package routes

import (
	"auth/internal/api/handlers"
	"auth/internal/repository"
	"auth/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func Setup(db database.Database, logger zerolog.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userRepo, logger)

	// Public routes
	v1 := router.Group("/aspi/v1")
	{
		v1.POST("/users", userHandler.Create)
		// v1.POST("/login", userHandler.Login)
	}

	return router
}
