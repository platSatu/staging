package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventUser struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_email_user" json:"email"`
	Phone     string    `gorm:"type:varchar(50)" json:"phone,omitempty"`
	UserID    string    `gorm:"type:char(36);uniqueIndex:idx_email_user" json:"user_id"`
	Status    string    `gorm:"type:enum('pending','verified','active','inactive');default:'pending'" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (e *EventUser) BeforeCreate(tx *gorm.DB) (err error) {
	e.ID = uuid.New().String()
	return
}
