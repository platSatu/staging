package model

import (
	"time"

	"gorm.io/gorm"
)

type ScSubjectTypeGeneral struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36);not null" json:"user_id"` // FK ke user
	SubjectName *string   `gorm:"type:varchar(255)" json:"subject_name"` // Nullable
	Curriculum  *string   `gorm:"type:varchar(255)" json:"curriculum"`   // Nullable
	Group       *string   `gorm:"type:varchar(255)" json:"group"`        // Nullable
	Status      *string   `gorm:"type:varchar(255)" json:"status"`       // Nullable
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScSubjectTypeGeneral) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScSubjectTypeGeneral) TableName() string {
	return "sc_subject_type_general"
}
