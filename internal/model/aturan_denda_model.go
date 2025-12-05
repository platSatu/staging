package model

import (
	"time"

	"gorm.io/gorm"
)

type AturanDenda struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	TipeDenda   string    `gorm:"not null" json:"tipe_denda"`   // 'harian_persentase', 'harian_flat', 'flat_cap'
	Persentase  *float64  `json:"persentase,omitempty"`
	JumlahFlat  *float64  `json:"jumlah_flat,omitempty"`
	CapMaksimal *float64  `json:"cap_maksimal,omitempty"`
	Catatan     *string   `json:"catatan,omitempty"`
	Status      string    `gorm:"default:'active'" json:"status"` // 'active' atau 'inactive'
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (a *AturanDenda) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if a.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (AturanDenda) TableName() string {
	return "aturan_denda"
}
