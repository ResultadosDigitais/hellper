package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB .
var DB *gorm.DB

// Init .
func Init() {
	var err error
	dsn := os.Getenv("HELLPER_DSN")
	// dsn := "postgres://hellper_dev:hellper_dev@localhost:5432/hellper_dev"

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
