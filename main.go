package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nurbol/cinema/internal/db"
	"github.com/nurbol/cinema/internal/routes"
)

func main() {
	db.InitDB()
	r := gin.Default()

	routes.SetupRoutes(r, db.DB)

	r.Run(":8080")
}
