package model

import (
	"time"

	"gorm.io/gorm"
)

type ScTeacher struct {
	ID        string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID    string    `gorm:"type:char(36);not null" json:"user_id"` // FK ke user
	Name      *string   `gorm:"type:varchar(255)" json:"name"`         // Nullable
	Email     *string   `gorm:"type:varchar(255)" json:"email"`        // Nullable
	Phone     *string   `gorm:"type:varchar(255)" json:"phone"`        // Nullable
	Mobile    *string   `gorm:"type:varchar(255)" json:"mobile"`       // Nullable
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScTeacher) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScTeacher) TableName() string {
	return "sc_teacher"
}
