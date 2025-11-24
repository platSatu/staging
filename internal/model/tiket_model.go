package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tiket struct {
	ID              string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID          *string    `gorm:"type:char(36)" json:"user_id"`
	Name            string     `gorm:"not null" json:"name"`
	Email           string     `gorm:"not null" json:"email"`
	Phone           string     `gorm:"not null" json:"phone"`
	OrderID         string     `gorm:"not null" json:"order_id"`
	EventID         string     `gorm:"not null;type:char(36)" json:"event_id"`
	JenisTicketID   string     `gorm:"not null;type:char(36)" json:"jenis_ticket_id"`
	Quantity        int        `gorm:"not null" json:"quantity"`
	Price           float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Subtotal        float64    `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Total           float64    `gorm:"type:decimal(10,2);not null" json:"total"`
	Qrcode          string     `gorm:"not null" json:"qrcode"`
	QrcodeDirectory string     `gorm:"not null" json:"qrcode_directory"`
	PaymentStatus   string     `gorm:"not null" json:"payment_status"`
	PaymentMethod   string     `gorm:"not null" json:"payment_method"`
	PaymentDate     *time.Time `json:"payment_date"`
	IsScanned       bool       `gorm:"default:false" json:"is_scanned"`
	ScanDate        *time.Time `json:"scan_date"`
	KodeBooking     string     `gorm:"not null" json:"kode_booking"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (Tiket) TableName() string {
	return "ticket"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (t *Tiket) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	// Set default is_scanned jika kosong
	if !t.IsScanned {
		t.IsScanned = false
	}
	return
}
