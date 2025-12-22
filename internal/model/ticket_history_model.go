package model

import (
	"time"

	"gorm.io/gorm"
)

type TicketHistory struct {
	ID              string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID          string    `gorm:"type:char(36)" json:"user_id"`
	ParentID        string    `gorm:"type:char(36)" json:"parent_id"`
	Qrcode          string    `gorm:"type:text" json:"qrcode"`
	ScannedAt       time.Time `json:"scanned_at"`
	ScannedByUser   string    `gorm:"type:char(36)" json:"scanned_by_user"`
	ScannedByDevice string    `json:"scanned_by_device"`
	IPAddress       string    `json:"ip_address"`
	Browser         string    `json:"browser"`
	CreatedAt       time.Time `json:"created_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketHistory) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (TicketHistory) TableName() string {
	return "ticket_history"
}