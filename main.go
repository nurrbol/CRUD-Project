package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nurbol/cinema/internal/db"
	"github.com/nurbol/cinema/internal/routes"
)

func main() {
	db.InitDB() // ✅ DB және миграция

	r := gin.Default()
	routes.SetupCinemaRoutes(r, db.DB) // db.DB - GORM instance
	r.Run(":8080")
}
