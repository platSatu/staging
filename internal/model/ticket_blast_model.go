package model

import (
	"time"

	"gorm.io/gorm"
)

type TicketBlast struct {
	ID               string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           string    `gorm:"type:char(36)" json:"user_id"`
	ParentID         string    `gorm:"type:char(36)" json:"parent_id"`
	TicketTemplateID string    `gorm:"type:char(36)" json:"ticket_template_id"`
	Nama             string    `json:"nama"`
	Subject          string    `json:"subject"`
	CcOrBcc          string    `gorm:"type:enum('cc','bcc')" json:"cc_or_bcc"`
	StatusPengiriman string    `json:"status_pengiriman"`
	Status           string    `gorm:"default:'active'" json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketBlast) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TicketBlast) TableName() string {
	return "ticket_blast"
}