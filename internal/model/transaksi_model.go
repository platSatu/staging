package model

import (
	"time"

	"gorm.io/gorm"
)

type Transaksi struct {
	ID               string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           string    `gorm:"not null" json:"user_id"`        // FK ke user
	KewajibanID      string    `gorm:"not null" json:"kewajiban_id"`   // FK ke kewajiban_user
	ParentID         string    `gorm:"type:char(36)" json:"parent_id"` // Kolom parent_id dengan char(36)
	TipeTransaksi    string    `gorm:"not null" json:"tipe_transaksi"` // 'pembayaran', 'denda', 'penyesuaian'
	Jumlah           float64   `gorm:"not null" json:"jumlah"`
	Tanggal          time.Time `gorm:"not null" json:"tanggal"`
	MetodePembayaran *string   `json:"metode_pembayaran,omitempty"`
	ReferenceID      *string   `json:"reference_id,omitempty"`
	StatusGateway    *string   `json:"status_gateway,omitempty"` // 'pending', 'success', 'failed'
	Catatan          *string   `json:"catatan,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`

	// Relasi (opsional, untuk preload jika diperlukan)
	User      User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Kewajiban KewajibanUser `gorm:"foreignKey:KewajibanID" json:"kewajiban,omitempty"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *Transaksi) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (Transaksi) TableName() string {
	return "transaksi"
}
