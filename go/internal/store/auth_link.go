package store

import (
	"time"
)

// AuthLink represents a one-time OAuth token request
type AuthLink struct {
	ID         uint      `gorm:"primaryKey"`
	Token      string    `gorm:"uniqueIndex;not null"`
	TelegramID int64     `gorm:"not null"`
	ExpiresAt  time.Time `gorm:"not null"`
	Used       bool      `gorm:"default:false"`
	CreatedAt  time.Time
}

// CreateAuthLink inserts a new auth token for a given Telegram user
func (s *Store) CreateAuthLink(token string, telegramID int64, ttl time.Duration) (*AuthLink, error) {
	authLink := &AuthLink{
		Token:      token,
		TelegramID: telegramID,
		ExpiresAt:  time.Now().Add(ttl),
	}
	if err := s.DB.Create(authLink).Error; err != nil {
		return nil, err
	}
	return authLink, nil
}

// GetValidAuthLink returns an unused, non-expired auth link
func (s *Store) GetValidAuthLink(token string) (*AuthLink, error) {
	var link AuthLink
	err := s.DB.Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).First(&link).Error
	if err != nil {
		return nil, err
	}
	return &link, nil
}

// MarkAuthLinkUsed marks the link as used (after successful auth)
func (s *Store) MarkAuthLinkUsed(token string) error {
	return s.DB.Model(&AuthLink{}).
		Where("token = ?", token).
		Update("used", true).Error
}
