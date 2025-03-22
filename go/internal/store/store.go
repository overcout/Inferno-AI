package store

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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
	err := s.DB.AutoMigrate(&User{}, &AuthLink{})
	if err != nil {
		log.Println("[ERROR] AutoMigrate failed:", err)
	}
	return err
}
