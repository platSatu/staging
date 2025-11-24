package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JenisTiket struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    string    `gorm:"not null;type:char(36)" json:"user_id"`
	EventID   string    `gorm:"not null;type:char(36)" json:"event_id"`
	Name      string    `gorm:"not null" json:"name"`
	Stok      int       `gorm:"not null" json:"stok"`
	Terjual   int       `gorm:"default:0" json:"terjual"`
	Sisa      int       `gorm:"not null" json:"sisa"`
	Harga     float64   `gorm:"type:decimal(10,2);not null" json:"harga"`
	Status    string    `gorm:"default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (JenisTiket) TableName() string {
	return "jenis_tiket"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (jt *JenisTiket) BeforeCreate(tx *gorm.DB) (err error) {
	if jt.ID == "" {
		jt.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if jt.Status == "" {
		jt.Status = "active"
	}

	// Set default terjual jika kosong
	if jt.Terjual == 0 {
		jt.Terjual = 0
	}
	return
}
