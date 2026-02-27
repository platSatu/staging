package model

import (
	"time"

	"gorm.io/gorm"
)

type ScAcademicProjection struct {
	ID           string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID       string    `gorm:"type:char(36);not null" json:"user_id"`             // FK ke user, tipe char(36)
	StudentID    *string   `gorm:"type:char(36)" json:"student_id"`                   // Nullable
	AcademicYear *string   `gorm:"type:char(36)" json:"academic_year"`                // Nullable
	Semester     *string   `gorm:"type:char(36)" json:"semester"`                     // Nullable
	QuarterID    *string   `gorm:"type:char(36)" json:"quarter_id"`                   // Nullable
	IsSplitIso   *bool     `gorm:"type:boolean" json:"is_split_iso"`                  // Nullable boolean
	Level        *string   `gorm:"type:varchar(255)" json:"level"`                    // Nullable
	LcID         *string   `gorm:"type:char(36)" json:"lc_id"`                        // Nullable
	TotalSchool  *int      `gorm:"type:int" json:"total_school"`                      // Nullable
	TotalPages   *int      `gorm:"type:int" json:"total_pages"`                       // Nullable (perbaiki nama dari total_Pages)
	Status       *string   `gorm:"type:enum('on_progress','complete')" json:"status"` // Nullable enum
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScAcademicProjection) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScAcademicProjection) TableName() string {
	return "sc_academic_projection"
}
