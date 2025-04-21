package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nurbol/cinema/internal/delivery"
	"github.com/nurbol/cinema/internal/repository"
	"github.com/nurbol/cinema/internal/services"
	"gorm.io/gorm"
)

func SetupCinemaRoutes(r *gin.Engine, db *gorm.DB) {
	cinemaRepo := repository.NewCinemaRepository(db)
	cinemaService := services.NewCinemaService(cinemaRepo)
	cinemaHandler := delivery.NewCinemaHandler(cinemaService)

	cinemaRoutes := r.Group("/api/v1/cinemas")
	{
		cinemaRoutes.GET("/", cinemaHandler.GetAllCinemas)
		cinemaRoutes.GET("/:id", cinemaHandler.GetCinema)
		cinemaRoutes.POST("/", cinemaHandler.CreateCinema)
		cinemaRoutes.PUT("/:id", cinemaHandler.UpdateCinema)
		cinemaRoutes.DELETE("/:id", cinemaHandler.DeleteCinema)
	}
}
