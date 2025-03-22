package store

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           int64     `gorm:"primaryKey"` // Telegram user ID
	Email        string    `gorm:"size:255"`
	AccessToken  string    `gorm:"type:text"`
	RefreshToken string    `gorm:"type:text"`
	TokenExpiry  time.Time
	CreatedAt    time.Time
}

// GetOrCreateUser finds a user by Telegram ID or creates a new one
func (s *Store) GetOrCreateUser(telegramID int64) (*User, error) {
	var user User
	tx := s.DB.First(&user, telegramID)
	if tx.Error == nil {
		return &user, nil
	}
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		user.ID = telegramID
		if err := s.DB.Create(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, tx.Error
}
