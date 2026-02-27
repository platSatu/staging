package model

import (
	"time"

	"gorm.io/gorm"
)

type ScTemplateHabits struct {
	ID            string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID        string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	Level         string    `gorm:"type:varchar(255)" json:"level"`
	SubjectID     string    `gorm:"type:char(36)" json:"subject_id"`
	SubSubjectID  string    `gorm:"type:char(36)" json:"sub_subject_id"`
	Grade         string    `gorm:"type:char(36)" json:"grade"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScTemplateHabits) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScTemplateHabits) TableName() string {
	return "sc_template_habits"
}