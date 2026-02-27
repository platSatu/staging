package model

import (
	"time"

	"gorm.io/gorm"
)

type ScSubjectTypeGeneralSubject struct {
	ID               string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID           string    `gorm:"type:char(36);not null" json:"user_id"`      // FK ke user
	GpaWeight        *float64  `gorm:"type:decimal(5,2)" json:"gpa_weight"`        // Nullable
	MinPassingScore  *float64  `gorm:"type:decimal(5,2)" json:"min_passing_score"` // Nullable
	StartPage        *int      `gorm:"type:int" json:"start_page"`                 // Nullable
	EndingPace       *int      `gorm:"type:int" json:"ending_pace"`                // Nullable
	SubjectName      *string   `gorm:"type:varchar(255)" json:"subject_name"`      // Nullable
	ProductReference *string   `gorm:"type:varchar(255)" json:"product_reference"` // Nullable
	StartingPace     *int      `gorm:"type:int" json:"starting_pace"`              // Nullable
	Unit             *string   `gorm:"type:varchar(50)" json:"unit"`               // Nullable
	IsAssignable     *bool     `gorm:"type:tinyint(1)" json:"is_assignable"`       // Nullable
	IsPace           *bool     `gorm:"type:tinyint(1)" json:"is_pace"`             // Nullable
	PacesNumber      *int      `gorm:"type:int" json:"paces_number"`               // Nullable
	SubjectType      *string   `gorm:"type:char(36)" json:"subject_type"`          // Nullable, mungkin FK
	EndPage          *int      `gorm:"type:int" json:"end_page"`                   // Nullable
	PerCredit        *float64  `gorm:"type:decimal(5,2)" json:"per_credit"`        // Nullable
	TotalPages       *int      `gorm:"type:int" json:"total_pages"`                // Nullable
	CourseName       *string   `gorm:"type:varchar(255)" json:"course_name"`       // Nullable
	PrevPage         *int      `gorm:"type:int" json:"prev_page"`                  // Nullable
	NextPage         *int      `gorm:"type:int" json:"next_page"`                  // Nullable
	Units            *int      `gorm:"type:int" json:"units"`                      // Nullable
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScSubjectTypeGeneralSubject) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScSubjectTypeGeneralSubject) TableName() string {
	return "sc_subject_type_general_subject"
}
