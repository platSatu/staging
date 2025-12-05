package model

import (
	"time"
)

type PaymentInvoices struct {
	ID         string    `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID     string    `json:"user_id" gorm:"type:char(36);not null"`
	ParentID   string    `json:"parent_id" gorm:"type:char(36);not null"`
	CategoryID string    `json:"category_id" gorm:"type:char(36);not null"`
	Amount     float64   `json:"amount" gorm:"type:decimal(12,2);not null"`
	DueDate    time.Time `json:"due_date" gorm:"type:date;not null"`
	Status     string    `json:"status" gorm:"type:enum('unpaid','partial','paid');not null"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
