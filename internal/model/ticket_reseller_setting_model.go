package model

import (
	"gorm.io/gorm"
)

type TicketResellerSetting struct {
	ID                 string  `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID             string  `gorm:"type:char(36)" json:"user_id"`
	ParentID           string  `gorm:"type:char(36)" json:"parent_id"`
	EventID            string  `gorm:"type:char(36)" json:"event_id"`
	Slug               string  `json:"slug"`
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	IDReseller         string  `json:"id_reseller"`
	Status             string  `gorm:"default:'active'" json:"status"`
	TicketCategoryID   string  `gorm:"type:char(36)" json:"ticket_category_id"`
	MethodPembayaran   string  `json:"method_pembayaran"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketResellerSetting) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TicketResellerSetting) TableName() string {
	return "ticket_reseller_setting"
}