package model

import (
	"time"
)

type PaymentForm struct {
	ID                string    `json:"id" gorm:"type:char(36);primaryKey;not null"`
	UserID            string    `json:"user_id" gorm:"type:char(36);not null"`
	CategoryID        string    `json:"category_id" gorm:"type:char(36);not null"`
	Name              string    `json:"name" gorm:"type:varchar(255);not null"`
	Description       *string   `json:"description" gorm:"type:text"`
	BaseAmount        float64   `json:"base_amount" gorm:"type:decimal(12,2);not null"`
	EnableInstallment bool      `json:"enable_installment" gorm:"type:boolean;default:false;not null"`
	EnablePenalty     bool      `json:"enable_penalty" gorm:"type:boolean;default:false;not null"`
	Status            string    `json:"status" gorm:"type:enum('active','inactive');default:'active';not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PaymentForm) TableName() string {
	return "payment_form"
}
