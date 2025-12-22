package model

import (
	"time"

	"gorm.io/gorm"
)

type TicketQrcode struct {
	ID                string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID            string    `gorm:"type:char(36)" json:"user_id"`
	ParentID          string    `gorm:"type:char(36)" json:"parent_id"`
	ResellerID        string    `gorm:"type:char(36)" json:"reseller_id"`
	TicketEventID     string    `gorm:"type:char(36)" json:"ticket_event_id"`
	TicketKategoriID  string    `gorm:"type:char(36)" json:"ticket_kategori_id"`
	TicketJenisID     string    `gorm:"type:char(36)" json:"ticket_jenis_id"`
	TicketVoucherID   string    `gorm:"type:char(36)" json:"ticket_voucher_id"`
	Quantity          int       `json:"quantity"`
	Price             float64   `json:"price"`
	OrderID           string    `json:"order_id"`
	KodeBooking       string    `json:"kode_booking"`
	Qrcode            string    `gorm:"type:text" json:"qrcode"`
	DirectoryQrcode   string    `json:"directory_qrcode"`
	PaymentStatus     string    `json:"payment_status"`
	PaymentMethod     string    `json:"payment_method"`
	PaymentDate       time.Time `json:"payment_date"`
	PaymentFeeID      string    `gorm:"type:char(36)" json:"payment_fee_id"`
	IsScanned         bool      `json:"is_scanned"`
	DateScanned       time.Time `json:"date_scanned"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *TicketQrcode) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (TicketQrcode) TableName() string {
	return "ticket_qrcode"
}