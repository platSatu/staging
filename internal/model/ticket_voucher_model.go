package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type TicketVoucher struct {
	ID               string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           string    `gorm:"type:char(36)" json:"user_id"`
	ParentID         string    `gorm:"type:char(36)" json:"parent_id"`
	TicketKategoryID string    `gorm:"type:char(36)" json:"ticket_kategory_id"`
	Nama             string    `json:"nama"`
	KodeVoucher      string    `json:"kode_voucher"`
	TanggalMulai     time.Time `json:"tanggal_mulai"`
	TanggalExpired   time.Time `json:"tanggal_expired"`
	HargaFlat        float64   `json:"harga_flat"`
	Persentase       float64   `json:"persentase"`
	Status           string    `gorm:"default:'active'" json:"status"`
	Quota            int       `json:"quota"`
	Terpakai         int       `json:"terpakai"`
	Sisa             int       `json:"sisa"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Custom UnmarshalJSON untuk TanggalMulai dan TanggalExpired agar menerima date string atau RFC3339
func (t *TicketVoucher) UnmarshalJSON(data []byte) error {
	type Alias TicketVoucher
	aux := &struct {
		TanggalMulai   string `json:"tanggal_mulai"`
		TanggalExpired string `json:"tanggal_expired"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	parseDateTime := func(value string) (time.Time, error) {
		// coba parse RFC3339 / ISO8601
		if parsed, err := time.Parse(time.RFC3339, value); err == nil {
			return parsed, nil
		}
		// fallback parse YYYY-MM-DD
		return time.Parse("2006-01-02", value)
	}

	if aux.TanggalMulai != "" {
		parsed, err := parseDateTime(aux.TanggalMulai)
		if err != nil {
			return err
		}
		t.TanggalMulai = parsed
	}

	if aux.TanggalExpired != "" {
		parsed, err := parseDateTime(aux.TanggalExpired)
		if err != nil {
			return err
		}
		t.TanggalExpired = parsed
	}

	return nil
}

// BeforeCreate otomatis generate UUID jika ID kosong
func (t *TicketVoucher) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	if t.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}

func (TicketVoucher) TableName() string {
	return "ticket_voucher"
}
