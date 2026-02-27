package model

import (
	"time"

	"gorm.io/gorm"
)

type ScStudent struct {
	ID            string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID        string    `gorm:"type:char(36);not null" json:"user_id"`
	Name          *string   `gorm:"type:varchar(255)" json:"name"`           // FK ke user, tipe char(36)
	Address       *string   `gorm:"type:varchar(500)" json:"address"`        // Nullable
	Tin           *string   `gorm:"type:varchar(255)" json:"tin"`            // Nullable
	Tags          *string   `gorm:"type:varchar(255)" json:"tags"`           // Nullable
	StudentType   *string   `gorm:"type:varchar(255)" json:"student_type"`   // Nullable
	LcID          *string   `gorm:"type:char(36)" json:"lc_id"`              // Nullable
	LevelID       *string   `gorm:"type:char(36)" json:"level_id"`           // Nullable
	Phone         *string   `gorm:"type:varchar(255)" json:"phone"`          // Nullable
	Mobile        *string   `gorm:"type:varchar(255)" json:"mobile"`         // Nullable
	Email         *string   `gorm:"type:varchar(255)" json:"email"`          // Nullable
	Language      *string   `gorm:"type:varchar(255)" json:"language"`       // Nullable
	StudentStatus *string   `gorm:"type:varchar(255)" json:"student_status"` // Nullable
	PartnerType   *string   `gorm:"type:varchar(255)" json:"partner_type"`   // Nullable
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (s *ScStudent) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}
	return
}

func (ScStudent) TableName() string {
	return "sc_student"
}
