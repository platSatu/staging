package model

import (
	"time"

	"gorm.io/gorm"
)

type KategoriPembayaran struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID       string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	NamaKategori string    `gorm:"not null" json:"nama_kategori"`
	Deskripsi    *string   `json:"deskripsi,omitempty"`
	Status       string    `gorm:"default:'active'" json:"status"` // Mengganti Aktif dengan Status (string)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (k *KategoriPembayaran) BeforeCreate(tx *gorm.DB) (err error) {
	if k.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if k.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (KategoriPembayaran) TableName() string {
	return "kategori_pembayaran"
}
