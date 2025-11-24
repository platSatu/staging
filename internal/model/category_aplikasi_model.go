package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryAplikasi struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    string    `gorm:"not null;type:char(36)" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Status    string    `gorm:"default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (CategoryAplikasi) TableName() string {
	return "category_aplikasi"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (c *CategoryAplikasi) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if c.Status == "" {
		c.Status = "active"
	}
	return
}
