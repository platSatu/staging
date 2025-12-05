package model

import (
	"time"
)

type PaymentPayments struct {
	ID            string    `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID        string    `json:"user_id" gorm:"type:char(36);not null"`
	ParentID      string    `json:"parent_id" gorm:"type:char(36);not null"`
	InvoiceID     string    `json:"invoice_id" gorm:"type:char(36);not null"`
	InstallmentID string    `json:"installment_id" gorm:"type:char(36);not null"`
	Amount        float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	PaymentMethod string    `json:"payment_method" gorm:"type:varchar(50);not null"`
	PaymentDate   time.Time `json:"payment_date" gorm:"type:datetime;not null"`
	Status        string    `json:"status" gorm:"type:enum('pending','success','failed');not null"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
