package model

import (
	"time"

	"gorm.io/gorm"
)

type TicketRegister struct {
	ID                 string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID             string     `gorm:"type:char(36)" json:"user_id"`
	ParentID           string     `gorm:"type:char(36)" json:"parent_id"`
	FullName           string     `json:"full_name"`
	Email              string     `json:"email"`
	Handphone          *string    `json:"handphone,omitempty"`
	Status             string     `gorm:"default:'active'" json:"status"`
	TanggalRegister    time.Time  `json:"tanggal_register"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	PurchaseToken      string     `gorm:"column:purchase_token;type:varchar(255)" json:"purchase_token"`
	TokenExpiredAt     *time.Time `gorm:"column:token_expired_at" json:"token_expired_at,omitempty"`
	PurchaseLinkSentAt *time.Time `gorm:"column:purchase_link_sent_at" json:"purchase_link_sent_at,omitempty"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketRegister) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TicketRegister) TableName() string {
	return "ticket_register"
}
