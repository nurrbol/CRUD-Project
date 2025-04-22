package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nurbol/cinema/internal/auth"
	"github.com/nurbol/cinema/internal/delivery"
	"github.com/nurbol/cinema/internal/middleware"
	"github.com/nurbol/cinema/internal/repository"
	"github.com/nurbol/cinema/internal/services"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	setupAuthRoutes(r)
	setupCinemaRoutes(r, db)
}

func setupAuthRoutes(r *gin.Engine) {
	authRoutes := r.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", auth.Register)
		authRoutes.POST("/login", auth.Login)
		authRoutes.GET("/me", middleware.AuthMiddleware(), auth.Me)
	}
}

func setupCinemaRoutes(r *gin.Engine, db *gorm.DB) {
	cinemaRepo := repository.NewCinemaRepository(db)
	cinemaService := services.NewCinemaService(cinemaRepo)
	cinemaHandler := delivery.NewCinemaHandler(cinemaService)

	protected := r.Group("/api/v1/cinemas", middleware.AuthMiddleware())
	{
		protected.GET("/", cinemaHandler.GetAllCinemas)
		protected.GET("/:id", cinemaHandler.GetCinema)
		protected.POST("/", cinemaHandler.CreateCinema)
		protected.PUT("/:id", cinemaHandler.UpdateCinema)
		protected.DELETE("/:id", cinemaHandler.DeleteCinema)
	}
}
