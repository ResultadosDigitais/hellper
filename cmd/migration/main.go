package main

import (
	"fmt"
	"hellper/internal/config/database"
	"hellper/internal/model"
	"os"
)

// env stores the dependencies of migration app
type env struct {
	db *database.DB
}

func main() {
	dsn := os.Getenv("HELLPER_DSN")
	db, err := database.NewConnectionWithDSN(dsn)
	if err != nil {
		panic(err)
	}

	env := &env{db: db}

	env.migration()
}

// migration runs the GORM Migrator
// IMPORTANT: Remember to use this function to change anything on the database
func (env *env) migration() {
	// Table renaming from `incident` to `incidents`
	err := env.db.Migrator().RenameTable("incident", &model.Incident{})
	if err != nil {
		panic(err)
	}

	// Columns renaming
	err = env.db.Migrator().RenameColumn(&model.Incident{}, "start_ts", "started_at")
	if err != nil {
		panic(err)
	}
	err = env.db.Migrator().RenameColumn(&model.Incident{}, "end_ts", "ended_at")
	if err != nil {
		panic(err)
	}
	err = env.db.Migrator().RenameColumn(&model.Incident{}, "identification_ts", "identified_at")
	if err != nil {
		panic(err)
	}

	// Change table to use GORM struct
	err = env.db.AutoMigrate(&model.Incident{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Incident Table created using GORM")
}
