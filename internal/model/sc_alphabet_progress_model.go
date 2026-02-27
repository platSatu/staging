package model

import (
	"time"

	"gorm.io/gorm"
)

type ScAlphabetProgress struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	Level       string    `gorm:"type:varchar(255)" json:"level"`
	LVertical   string    `gorm:"type:varchar(255)" json:"l_vertical"`
	LHorizontal string    `gorm:"type:varchar(255)" json:"l_horizontal"`
	Score       string    `gorm:"type:varchar(255)" json:"score"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScAlphabetProgress) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScAlphabetProgress) TableName() string {
	return "sc_alphabet_progress"
}