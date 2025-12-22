package model

import (
	"time"

	"gorm.io/gorm"
)

type TicketKategori struct {
	ID              string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID          string    `gorm:"type:char(36)" json:"user_id"`
	ParentID        string    `gorm:"type:char(36)" json:"parent_id"`
	EventID         string    `gorm:"type:char(36)" json:"event_id"`
	FeeID           string    `gorm:"type:char(36)" json:"fee_id"`
	JenisQuantityID string    `gorm:"type:char(36)" json:"jenis_quantity_id"`
	Nama            string    `json:"nama"`
	Image           string    `json:"image"` // Kolom baru untuk gambar (URL atau path)
	StokAwal        int       `json:"stok_awal"`
	Terjual         int       `json:"terjual"`
	Sisa            int       `json:"sisa"`
	Harga           float64   `json:"harga"`
	MinQuantity     int       `json:"min_quantity"` // Kolom baru, kosongkan saja
	MaxQuantity     int       `json:"max_quantity"` // Kolom baru, kosongkan saja
	Status          string    `gorm:"default:'active'" json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketKategori) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TicketKategori) TableName() string {
	return "ticket_kategori"
}
