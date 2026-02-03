package model

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Profile struct {
	ID          string    `gorm:"primaryKey;type:char(36)" json:"id"`
	UserID      string    `gorm:"not null;type:char(36)" json:"user_id"`
	BussinesID  int       `gorm:"autoIncrement;unique;not null" json:"bussines_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Slug        string    `gorm:"unique;not null" json:"slug"`
	Address     *string   `gorm:"type:varchar(255)" json:"address"`
	Phone       *string   `gorm:"type:varchar(20)" json:"phone"`
	Image       *string   `gorm:"type:varchar(255)" json:"image"`
	Status      string    `gorm:"default:'active'" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName menentukan nama tabel secara eksplisit
func (Profile) TableName() string {
	return "profiles"
}

// BeforeCreate akan otomatis generate UUID jika ID kosong, generate slug dari name, dan set default status
func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	// Generate slug dari name
	if p.Name != "" {
		p.Slug = GenerateSlug(p.Name) // Gunakan fungsi diekspor
	}

	// Set default status jika kosong
	if p.Status == "" {
		p.Status = "active"
	}
	return
}

// GenerateSlug membuat slug dari name: lowercase, replace spasi dengan "-", remove karakter aneh (diekspor dengan huruf besar)
func GenerateSlug(name string) string {
	// Lowercase
	slug := strings.ToLower(name)
	// Replace spasi dengan "-"
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove karakter non-alphanumeric kecuali "-"
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "")
	// Remove multiple "-" berturut-turut
	reg2 := regexp.MustCompile(`-+`)
	slug = reg2.ReplaceAllString(slug, "-")
	// Trim "-" di awal/akhir
	slug = strings.Trim(slug, "-")
	return slug
}
