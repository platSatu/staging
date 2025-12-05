package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deposit struct {
	ID               string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           string    `gorm:"type:char(36)" json:"user_id"`
	Debit            float64   `gorm:"type:decimal(10,2)" json:"debit"`
	Kredit           float64   `gorm:"type:decimal(10,2)" json:"kredit"`
	Saldo            float64   `gorm:"type:decimal(10,2)" json:"saldo"`
	TransaksiTanggal time.Time `gorm:"type:datetime" json:"transaksi_tanggal"`
	NoInvoice        string    `gorm:"type:varchar(255);unique" json:"no_invoice"`
	OrderID          string    `gorm:"type:varchar(255);unique" json:"order_id"`
	TransaksiStatus  string    `gorm:"type:enum('paid','pending','deposit','expired')" json:"transaksi_status"`
	TransaksiMethod  string    `gorm:"type:varchar(255)" json:"transaksi_method"`
	Keterangan       string    `gorm:"type:text" json:"keterangan"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (Deposit) TableName() string {
	return "deposits"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (d *Deposit) BeforeCreate(tx *gorm.DB) (err error) {
	if d.ID == "" {
		d.ID = uuid.New().String()
	}
	return
}
