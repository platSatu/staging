package model

import (
	"time"

	"gorm.io/gorm"
)

type ScGrade struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    string    `gorm:"type:char(36)" json:"user_id"` // FK ke user, tipe char(36)
	GradeName string    `gorm:"type:varchar(255)" json:"grade_name"`
	MinScore  string    `gorm:"type:varchar(255)" json:"min_score"`
	MaxScore  string    `gorm:"type:varchar(255)" json:"max_score"`
	Status    string    `gorm:"type:varchar(255)" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScGrade) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScGrade) TableName() string {
	return "sc_grade"
}