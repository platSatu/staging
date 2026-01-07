package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              string     `gorm:"primaryKey;type:char(36)" json:"id"`
	Username        string     `gorm:"not null" json:"username"`          // boleh duplikat
	Email           string     `gorm:"unique;not null" json:"email"`      // unik
	Password        string     `gorm:"column:password;not null" json:"-"` // hash password, kolom DB "password", tidak tampil di JSON
	FullName        string     `json:"full_name"`
	Status          string     `gorm:"default:'active'" json:"status"`
	Role            string     `gorm:"default:'user'" json:"role"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
	LastLogout      *time.Time `json:"last_logout,omitempty"`
	ParentID        *string    `json:"parent_id,omitempty"`
	KodeReferal     *string    `json:"kode_referal,omitempty"` // boleh duplikat
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	RememberToken   *string    `json:"-"`
	Saldo           *float64   `json:"saldo,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// BeforeCreate akan otomatis generate UUID jika ID kosong
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		tx.Statement.SetColumn("ID", gorm.Expr("UUID()"))
	}

	// Set default status jika kosong
	if u.Status == "" {
		tx.Statement.SetColumn("Status", "active")
	}
	return
}
