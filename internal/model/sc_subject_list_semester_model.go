package model

import (
	"time"

	"gorm.io/gorm"
)

type ScSubjectListSemester struct {
	ID         string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID     string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	No         string    `gorm:"type:varchar(255)" json:"no"`
	Subject    string    `gorm:"type:varchar(255)" json:"subject"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScSubjectListSemester) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScSubjectListSemester) TableName() string {
	return "sc_subject_list_semester"
}