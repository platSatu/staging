package model

import (
	"time"

	"gorm.io/gorm"
)

type ScPembelianBuku struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36);not null" json:"user_id"`
	ParentID    *string   `gorm:"type:char(36)" json:"parent_id"`
	Subject     *string   `gorm:"type:varchar(255)" json:"subject"`
	PacesNumber *string   `gorm:"type:varchar(255)" json:"paces_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScPembelianBuku) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScPembelianBuku) TableName() string {
	return "sc_pembelian_book"
}
