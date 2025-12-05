package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Packages struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID       string    `gorm:"not null;type:char(36)" json:"user_id"`
	PackagesID   string    `gorm:"type:char(36)" json:"packages_id"`
	Name         string    `gorm:"not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	Price        float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	DurationDays int       `gorm:"not null" json:"duration_days"`
	Status       string    `gorm:"default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (Packages) TableName() string {
	return "packages"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (p *Packages) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if p.Status == "" {
		p.Status = "active"
	}
	return
}
