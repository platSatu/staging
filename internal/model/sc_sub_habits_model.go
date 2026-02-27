package model

import (
	"time"

	"gorm.io/gorm"
)

type ScSubHabits struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	HabitsID    string    `gorm:"type:char(36)" json:"habits_id"` // FK ke sc_habits, tipe char(36)
	Subject     *string   `gorm:"type:varchar(255);null" json:"subject,omitempty"`
	Description *string   `gorm:"type:text;null" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScSubHabits) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScSubHabits) TableName() string {
	return "sc_sub_habits"
}