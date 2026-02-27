package model

import (
	"time"

	"gorm.io/gorm"
)

type ScHabits struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	Subject     string    `gorm:"not null" json:"subject"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScHabits) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScHabits) TableName() string {
	return "sc_habits"
}