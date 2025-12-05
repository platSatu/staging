package model

import (
	"time"
)

type PaymentPenalty struct {
	ID           string    `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID       string    `json:"user_id" gorm:"type:char(36);not null"`
	CategoryID   string    `json:"category_id" gorm:"type:char(36);not null"`
	Name         string    `json:"name" gorm:"type:varchar(255);not null"`
	Description  *string   `json:"description" gorm:"type:text"`
	PenaltyType  string    `json:"penalty_type" gorm:"type:enum('flat','percent');not null"`
	FlatValue    *float64  `json:"flat_value" gorm:"type:decimal(12,2)"`
	PercentValue *float64  `json:"percent_value" gorm:"type:decimal(5,2)"`
	MaxPenalty   float64   `json:"max_penalty" gorm:"type:decimal(12,2);not null"`
	ApplyOn      string    `json:"apply_on" gorm:"type:enum('invoice','installment','both');not null"`
	Active       string    `json:"active" gorm:"type:enum('active','inactive');default:'active';not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PaymentPenalty) TableName() string {
	return "payment_penalty"
}
