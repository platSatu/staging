package model

import (
	"time"
)

type PaymentInvoice struct {
	ID                string     `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID            string     `json:"user_id" gorm:"type:char(36);not null"`
	ParentID          string     `json:"parent_id" gorm:"type:char(36);not null"`
	CategoryID        string     `json:"category_id" gorm:"type:char(36);not null"`
	FormID            string     `json:"form_id" gorm:"type:char(36);not null"`
	PenaltyID         *string    `json:"penalty_id" gorm:"type:char(36)"`
	InstallmentID     *string    `json:"installment_id" gorm:"type:char(36)"`
	Name              string     `json:"name" gorm:"type:varchar(255);not null"`
	Description       *string    `json:"description" gorm:"type:text"`
	Amount            float64    `json:"amount" gorm:"type:decimal(12,2);not null"`
	AmountPaid        float64    `json:"amount_paid" gorm:"type:decimal(12,2);not null;default:0.00"`
	OutstandingAmount float64    `json:"outstanding_amount" gorm:"type:decimal(12,2);not null;default:0.00"`
	DueDate           time.Time  `json:"due_date" gorm:"type:date;not null"`
	EnableInstallment bool       `json:"enable_installment" gorm:"type:boolean;not null;default:false"`
	EnablePenalty     bool       `json:"enable_penalty" gorm:"type:boolean;not null;default:false"`
	PaymentStatus     string     `json:"payment_status" gorm:"type:enum('unpaid','partial','paid');not null;default:'unpaid'"`
	PaymentMethod     *string    `json:"payment_method" gorm:"type:varchar(50)"`
	PaymentDate       *time.Time `json:"payment_date" gorm:"type:datetime"`
	OrderID           *string    `json:"order_id" gorm:"type:varchar(100)"`
	Notes             *string    `json:"notes" gorm:"type:text"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PaymentInvoice) TableName() string {
	return "payment_invoice"
}
