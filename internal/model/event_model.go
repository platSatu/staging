package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"not null;type:char(36)" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Date        time.Time `gorm:"not null" json:"date"`
	Location    string    `gorm:"not null" json:"location"`
	Status      string    `gorm:"default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (Event) TableName() string {
	return "events"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (e *Event) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if e.Status == "" {
		e.Status = "active"
	}
	return
}
