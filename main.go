package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nurbol/cinema/internal/models"
	"github.com/nurbol/cinema/internal/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func main() {
	dsn := "host=localhost user=user password=password dbname=postgres port=5433 sslmode=disable TimeZone=Asia/Almaty"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	// ✅ Кестені автоматты түрде қосу
	err = db.AutoMigrate(&models.Cinema{})
	if err != nil {
		log.Fatal("Failed to migrate the database: ", err)
	}

	r := gin.Default()
	routes.SetupCinemaRoutes(r, db)
	r.Run(":8080")
}
