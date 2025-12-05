package model

import (
	"time"
)

type PaymentUser struct {
	ID        string    `json:"id" gorm:"primaryKey;type:char(36)"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	ParentID  string    `json:"parent_id" gorm:"type:char(36);not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// Tambahkan method TableName untuk mengarahkan ke tabel "users"
func (PaymentUser) TableName() string {
	return "users"
}
