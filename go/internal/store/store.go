package store

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Store struct {
	DB *gorm.DB
}

func NewStore(dsn string) (*Store, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) InitSchema() error {
	err := s.DB.AutoMigrate(&User{})
	if err != nil {
		log.Println("[ERROR] AutoMigrate failed:", err)
	}
	return err
}