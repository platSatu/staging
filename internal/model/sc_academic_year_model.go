package model

import (
	"time"

	"gorm.io/gorm"
)

type ScAcademicYear struct {
	ID           string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID       string     `gorm:"type:char(36);not null" json:"user_id"`                      // FK ke user
	Name         *string    `gorm:"type:varchar(255)" json:"name"`                              // Nullable - Tambahan kolom name
	Status       string     `gorm:"type:enum('open','active','closed');not null" json:"status"` // Enum
	BeginDate    *time.Time `gorm:"type:date" json:"begin_date"`                                // Nullable DATE
	EndDate      *time.Time `gorm:"type:date" json:"end_date"`                                  // Nullable DATE
	AcademicYear *string    `gorm:"type:varchar(255)" json:"academic_year"`                     // Nullable
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScAcademicYear) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScAcademicYear) TableName() string {
	return "sc_academic_year"
}
