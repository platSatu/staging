package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Voucher struct {
	ID          string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string     `gorm:"not null;type:char(36)" json:"user_id"`
	KodeVoucher string     `gorm:"unique;not null" json:"kode_voucher"`
	Status      string     `gorm:"default:'active'" json:"status"`
	ValidFrom   *time.Time `gorm:"type:date" json:"valid_from"`
	ValidUntil  *time.Time `gorm:"type:date" json:"valid_until"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (Voucher) TableName() string {
	return "vouchers"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong, dan set default status
func (v *Voucher) BeforeCreate(tx *gorm.DB) (err error) {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if v.Status == "" {
		v.Status = "active"
	}
	return
}
