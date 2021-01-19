package main

import (
	"fmt"
	"hellper/internal/config/database"
	"hellper/internal/model"
)

func init() {
	database.Init()
}

func main() {
	migration()
}

func migration() {
	// Table renaming from `incident` to `incidents`
	err := database.DB.Migrator().RenameTable("incident", &model.Incident{})
	if err != nil {
		panic(err)
	}

	// Columns renaming
	err = database.DB.Migrator().RenameColumn(&model.Incident{}, "start_ts", "started_at")
	if err != nil {
		panic(err)
	}
	err = database.DB.Migrator().RenameColumn(&model.Incident{}, "end_ts", "ended_at")
	if err != nil {
		panic(err)
	}
	err = database.DB.Migrator().RenameColumn(&model.Incident{}, "identification_ts", "identified_at")
	if err != nil {
		panic(err)
	}

	// Change table to use GORM struct
	err = database.DB.AutoMigrate(&model.Incident{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Incident Table created using GORM")
}
