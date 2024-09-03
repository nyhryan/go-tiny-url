package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type URLRecords struct {
	gorm.Model
	LongURL    string
	TinyURL    string
	ClickCount uint
}

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db/url.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database...")
	}
	db.AutoMigrate(&URLRecords{})

	return db
}
