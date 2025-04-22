package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/nurbol/cinema/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var DB *gorm.DB

func InitDB() {
	_ = godotenv.Load()

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslmode := "disable"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPass, dbHost, dbPort, dbName, sslmode)
	fmt.Println("Connecting to:", dsn)

	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("sql.Open error:", err)
	}

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		log.Fatal("migratepg error:", err)
	}

	cwd, _ := os.Getwd()
	migrationsPath := filepath.Join(cwd, "internal", "db", "migrations")
	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatal("migrate.New error:", err)
	}

	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal("migration up error:", err)
	}

	// ✔️ GORM тек DSN арқылы ашылады
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("gorm.Open error:", err)
	}

	DB = gormDB

	// 👇 AutoMigrate осы жерде
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("AutoMigrate error:", err)
	}
}
