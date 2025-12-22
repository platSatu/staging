package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketEvent struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID       string    `gorm:"type:char(36)" json:"user_id"`
	ParentID     string    `gorm:"type:char(36)" json:"parent_id"`
	Nama         string    `json:"nama"`
	Keterangan   string    `gorm:"type:text" json:"keterangan"`
	Slug         string    `json:"slug"`
	TanggalEvent time.Time `gorm:"type:date" json:"tanggal_event"`
	Alamat       string    `gorm:"type:text" json:"alamat"`
	Status       string    `gorm:"default:'active'" json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UnmarshalJSON untuk parsing tanggal "YYYY-MM-DD"
func (t *TicketEvent) UnmarshalJSON(data []byte) error {
	type Alias TicketEvent
	aux := &struct {
		TanggalEvent string `json:"tanggal_event"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// Parse tanggal dari string jika ada
	if aux.TanggalEvent != "" {
		parsedTime, err := time.Parse("2006-01-02", aux.TanggalEvent)
		if err != nil {
			return err
		}
		t.TanggalEvent = parsedTime
	}

	return nil
}

// BeforeCreate untuk auto-generate UUID dan set default status
func (t *TicketEvent) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String() // Generate UUID di Go langsung
	}
	if t.Status == "" {
		t.Status = "active"
	}
	return
}

func (TicketEvent) TableName() string {
	return "ticket_event"
}
