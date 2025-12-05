package model

import (
	"time"
)

type PaymentCategory struct {
	ID               string     `json:"id" gorm:"primaryKey;type:char(36)"`
	UserID           string     `json:"user_id" gorm:"type:char(36);not null"`
	Name             string     `json:"name" gorm:"type:varchar(255);not null"`
	Description      *string    `json:"description" gorm:"type:text"`
	AllowInstallment bool       `json:"allow_installment" gorm:"type:boolean;not null;default:false"`
	AllowPenalty     bool       `json:"allow_penalty" gorm:"type:boolean;not null;default:false"`
	Status           string     `json:"status" gorm:"type:enum('active','inactive');not null;default:'active'"`
	DueDateDefault   *time.Time `json:"due_date_default" gorm:"type:date"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

func (PaymentCategory) TableName() string {
	return "payment_category"
}
