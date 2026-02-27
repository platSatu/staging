package model

import (
	"time"

	"gorm.io/gorm"
)

type ScAcademicProjectionDetail struct {
	ID              string     `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID          string     `gorm:"type:char(36);not null" json:"user_id"`                                // FK ke user, tipe char(36)
	SubjectTypeID   *string    `gorm:"type:char(36)" json:"subject_type_id"`                                 // Nullable
	No              *int       `gorm:"type:int" json:"no"`                                                   // Nullable
	SubjectNameID   *string    `gorm:"type:char(36)" json:"subject_name_id"`                                 // Nullable
	Status          *string    `gorm:"type:enum('stocked','issued','completed','carry over')" json:"status"` // Nullable enum
	IssueDate       *time.Time `gorm:"type:date" json:"issue_date"`                                          // Nullable DATE
	PtDate          *time.Time `gorm:"type:date" json:"pt_date"`                                             // Nullable DATE
	PtScore         *float64   `gorm:"type:decimal(10,2)" json:"pt_score"`                                   // Nullable DECIMAL
	AlphabetPtScore *string    `gorm:"type:char(36)" json:"alphabet_pt_score"`                               // Nullable
	EndDate         *time.Time `gorm:"type:date" json:"end_date"`                                            // Nullable DATE
	Paces           *int       `gorm:"type:int" json:"paces"`                                                // Nullable
	PrevPace        *int       `gorm:"type:int" json:"prev_pace"`                                            // Nullable
	NextPace        *int       `gorm:"type:int" json:"next_pace"`                                            // Nullable
	Pages           *int       `gorm:"type:int" json:"pages"`                                                // Nullable
	OrderListID     *string    `gorm:"type:varchar(255)" json:"order_list_id"`                               // Nullable
	Order           *bool      `gorm:"type:boolean" json:"order"`                                            // Nullable
	ProductID       *string    `gorm:"type:char(36)" json:"product_id"`                                      // Nullable
	AcademicYearID  *string    `gorm:"type:char(36)" json:"academic_year_id"`                                // Nullable (FK ke sc_academic_projection jika diperlukan)
	IsProcessed     *bool      `gorm:"type:boolean" json:"is_processed"`                                     // Nullable
	OrderNote       *string    `gorm:"type:varchar(255)" json:"order_note"`                                  // Nullable
	SubjectID       *string    `gorm:"type:char(36)" json:"subject_id"`                                      // Nullable
	AssignmentID    *string    `gorm:"type:char(36)" json:"assignment_id"`                                   // Nullable
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScAcademicProjectionDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScAcademicProjectionDetail) TableName() string {
	return "sc_academic_projection_detail"
}
