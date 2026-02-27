package model

import (
	"time"

	"gorm.io/gorm"
)

type ScSubjectListSemesterYearly struct {
	ID         string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID     string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	StartLevel string    `gorm:"type:varchar(255)" json:"start_level"`
	EndLevel   string    `gorm:"type:varchar(255)" json:"end_level"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScSubjectListSemesterYearly) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScSubjectListSemesterYearly) TableName() string {
	return "sc_subject_list_semester_yearly"
}