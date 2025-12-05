package model

import (
	"time"
)

type PaymentPenaltySettings struct {
	ID               string    `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID           string    `json:"user_id" gorm:"type:char(36);not null"`
	CategoryID       string    `json:"category_id" gorm:"type:char(36);not null"`
	PenaltyType      string    `json:"penalty_type" gorm:"type:enum('percent','flat');not null"`
	PercentValue     *float64  `json:"percent_value" gorm:"type:decimal(5,2)"`
	FlatValue        *float64  `json:"flat_value" gorm:"type:decimal(12,2)"`
	MaxPenaltyAmount float64   `json:"max_penalty_amount" gorm:"type:decimal(12,2);not null"`
	ApplyOn          string    `json:"apply_on" gorm:"type:enum('invoice','installment','both');not null"`
	Active           string    `json:"active" gorm:"type:enum('active','inactive');not null;default:'active'"` // Diubah ke string untuk enum
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
