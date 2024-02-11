package storage

import (
	"file-chunker/internal/config"
	"file-chunker/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

func New() *Storage {
	conf := config.GetConfig()
	db, err := gorm.Open(sqlite.Open(conf.Database.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &Storage{
		DB: db,
	}
}

func (s Storage) MustMigrate() {
	err := s.DB.AutoMigrate(&models.File{})
	if err != nil {
		panic(err)
	}
	err = s.DB.AutoMigrate(&models.Chunks{})
	if err != nil {
		panic(err)
	}
}
