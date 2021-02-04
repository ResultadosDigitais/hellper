package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB wraps the gorm DB struct
type DB struct {
	*gorm.DB
}

// NewConnectionWithDSN opens a new connection with the database using the GORM package
func NewConnectionWithDSN(dsn string) (*DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}
