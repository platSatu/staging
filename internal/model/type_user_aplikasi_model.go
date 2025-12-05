package model

import (
	"time"

	"gorm.io/gorm"
)

type TypeUserAplikasi struct {
	ID         string `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID     string `gorm:"type:char(36)" json:"user_id"`
	ParentID   string `gorm:"type:char(36)" json:"parent_id"`
	AplikasiID string `gorm:"type:char(36)" json:"aplikasi_id"`
	Status     string `gorm:"default:'active'" json:"status"` // Kolom status active atau inactive

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi (opsional, untuk preload jika diperlukan)
	// Tambahkan relasi jika ada, misalnya ke User atau Aplikasi
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TypeUserAplikasi) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TypeUserAplikasi) TableName() string {
	return "type_user_aplikasi"
}
