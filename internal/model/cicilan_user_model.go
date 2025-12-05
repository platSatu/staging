package model

import (
	"time"

	"gorm.io/gorm"
)

type CicilanUser struct {
	ID     string `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID string `gorm:"type:char(36)" json:"user_id"`

	KewajibanID   string     `gorm:"not null" json:"kewajiban_id"`   // FK ke kewajiban_user
	ParentID      string     `gorm:"type:char(36)" json:"parent_id"` // FK ke cicilan_user (untuk hierarki cicilan)
	JumlahCicilan float64    `gorm:"not null" json:"jumlah_cicilan"`
	JatuhTempo    time.Time  `gorm:"not null" json:"jatuh_tempo"`
	Status        string     `gorm:"not null;default:'belum'" json:"status"` // 'belum', 'telat', 'lunas'
	Denda         float64    `gorm:"default:0" json:"denda"`
	TanggalBayar  *time.Time `json:"tanggal_bayar,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`

	// Relasi (opsional, untuk preload jika diperlukan)
	User      User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Kewajiban KewajibanUser `gorm:"foreignKey:KewajibanID" json:"kewajiban,omitempty"`
	Parent    *CicilanUser  `gorm:"foreignKey:ParentID" json:"parent,omitempty"` // Relasi ke parent cicilan (opsional)
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (c *CicilanUser) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if c.Status == "" {
		tx.Statement.SetColumn("Status", "belum")
	}
	return
}

func (CicilanUser) TableName() string {
	return "cicilan_user"
}
