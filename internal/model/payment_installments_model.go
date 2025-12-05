package model

import (
	"time"
)

type PaymentInstallments struct {
	ID                string     `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID            string     `json:"user_id" gorm:"type:char(36);not null"`
	ParentID          string     `json:"parent_id" gorm:"type:char(36);not null"`
	FormID            string     `json:"form_id" gorm:"type:char(36);not null"`
	InstallmentNumber int        `json:"installment_number" gorm:"not null"`
	Amount            float64    `json:"amount" gorm:"type:decimal(12,2);not null"`
	DueDate           time.Time  `json:"due_date" gorm:"type:date;not null"`
	Status            string     `json:"status" gorm:"type:enum('unpaid','paid','late');default:'unpaid';not null"`
	PenaltyAmount     float64    `json:"penalty_amount" gorm:"type:decimal(12,2);default:0.00;not null"`
	PaidAt            *time.Time `json:"paid_at" gorm:"type:datetime"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PaymentInstallments) TableName() string {
	return "payment_installment"
}
