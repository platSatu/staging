package model

import (
	"time"

	"gorm.io/gorm"
)

type ScLearningCenter struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"type:char(36);not null" json:"user_id"`               // FK ke user
	Name        *string   `gorm:"type:varchar(255)" json:"name"`                       // Nullable
	GroupIDN    *string   `gorm:"column:group_idn;type:varchar(255)" json:"group_idn"` // Nullable - Eksplisit column group_idn   // Nullable
	Principal   *string   `gorm:"type:varchar(255)" json:"principal"`                  // Nullable
	HomeTeacher *string   `gorm:"type:varchar(255)" json:"home_teacher"`               // Nullable
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScLearningCenter) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScLearningCenter) TableName() string {
	return "sc_learning_center"
}
