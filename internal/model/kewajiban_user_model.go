package model

import (
	"time"

	"gorm.io/gorm"
)

type KewajibanUser struct {
	ID                string  `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID            string  `gorm:"not null" json:"user_id"`        // FK ke user
	FormID            string  `gorm:"not null" json:"form_id"`        // FK ke form_pembayaran
	KategoriID        string  `gorm:"not null" json:"kategori_id"`    // FK ke kategori_pembayaran
	DendaID           string  `gorm:"not null" json:"denda_id"`       // FK ke aturan_denda
	ParentID          string  `gorm:"type:char(36)" json:"parent_id"` // Kolom parent_id dengan char(36)
	JumlahTotal       float64 `gorm:"not null" json:"jumlah_total"`
	JumlahSisa        float64 `gorm:"not null" json:"jumlah_sisa"`
	Status            string  `gorm:"not null;default:'active'" json:"status"` // 'active' atau 'lunas'
	TanggalMulai      string  `gorm:"type:date;not null" json:"tanggal_mulai"`
	TanggalJatuhTempo *string `gorm:"type:date" json:"tanggal_jatuh_tempo,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi (opsional, untuk preload jika diperlukan)
	User     User               `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Form     FormPembayaran     `gorm:"foreignKey:FormID" json:"form,omitempty"`
	Kategori KategoriPembayaran `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Denda    AturanDenda        `gorm:"foreignKey:DendaID" json:"denda,omitempty"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (k *KewajibanUser) BeforeCreate(tx *gorm.DB) (err error) {
	if k.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if k.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (KewajibanUser) TableName() string {
	return "kewajiban_user"
}
