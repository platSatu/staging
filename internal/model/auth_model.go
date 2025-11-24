package model

import (
	"time"
)

type RefreshToken struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    string    `gorm:"type:char(36);not null" json:"user_id"`
	Token     string    `gorm:"type:text;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Revoked   bool      `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time `json:"created_at"`
}
