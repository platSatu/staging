package model

import (
	"time"

	"gorm.io/gorm"
)

type TicketMethodPembayaran struct {
	ID         string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID     string    `gorm:"type:char(36)" json:"user_id"`
	ParentID   string    `gorm:"type:char(36)" json:"parent_id"`
	KategoriID string    `gorm:"type:char(36)" json:"kategori_id"`
	Nama       string    `json:"nama"`
	Status     string    `gorm:"default:'active'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketMethodPembayaran) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TicketMethodPembayaran) TableName() string {
	return "ticket_method_pembayaran"
}