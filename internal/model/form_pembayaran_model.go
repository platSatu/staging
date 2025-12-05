package model

import (
	"time"

	"gorm.io/gorm"
)

type FormPembayaran struct {
	ID             string  `gorm:"primaryKey;type:char(36)" json:"id"`
	KategoriID     string  `gorm:"not null" json:"kategori_id"` // FK ke kategori_pembayaran
	DendaID        string  `gorm:"not null" json:"denda_id"`    // FK ke aturan_denda
	NamaForm       string  `gorm:"not null" json:"nama_form"`
	Jumlah         float64 `gorm:"not null" json:"jumlah"`
	BolehCicilan   bool    `gorm:"default:false" json:"boleh_cicilan"`
	JumlahCicilan  *int    `json:"jumlah_cicilan,omitempty"`
	Status         string  `gorm:"default:'active'" json:"status"` // Kolom status active
	UserID         string  `gorm:"type:char(36)" json:"user_id"`   // Kolom user_id dengan char(36)
	TanggalMulai   string  `gorm:"type:date" json:"tanggal_mulai"`
	TanggalSelesai *string `gorm:"type:date" json:"tanggal_selesai,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relasi (opsional, untuk preload jika diperlukan)
	Kategori KategoriPembayaran `gorm:"foreignKey:KategoriID" json:"kategori,omitempty"`
	Denda    AturanDenda        `gorm:"foreignKey:DendaID" json:"denda,omitempty"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (f *FormPembayaran) BeforeCreate(tx *gorm.DB) (err error) {
	if f.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if f.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (FormPembayaran) TableName() string {
	return "form_pembayaran"
}
