package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService struct {
	DB *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{DB: db}
}

// CREATE
func (s *ProfileService) CreateProfile(profile *model.Profile) error {
	if profile.ID == "" {
		profile.ID = uuid.New().String()
	}

	if profile.Name == "" {
		return fmt.Errorf("name is required")
	}

	if profile.Status == "" {
		profile.Status = "active"
	}

	// Generate bussines_id unik (increment dari max di database)
	var maxBussinesID int
	err := s.DB.Model(&model.Profile{}).Select("COALESCE(MAX(bussines_id), 0)").Scan(&maxBussinesID).Error
	if err != nil {
		return fmt.Errorf("failed to generate bussines_id: %v", err)
	}
	profile.BussinesID = maxBussinesID + 1

	// Validasi user_id ada di tabel users (opsional, jika ingin foreign key constraint)
	var user model.User
	if err := s.DB.First(&user, "id = ?", profile.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	return s.DB.Create(profile).Error
}

// READ ALL
func (s *ProfileService) GetAllProfiles() ([]model.Profile, error) {
	var profiles []model.Profile
	result := s.DB.Find(&profiles)
	return profiles, result.Error
}

// READ BY ID
func (s *ProfileService) GetProfileByID(id string) (*model.Profile, error) {
	var profile model.Profile
	result := s.DB.First(&profile, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, result.Error
}

// READ BY USER ID
func (s *ProfileService) GetProfilesByUserID(userID string) ([]model.Profile, error) {
	var profiles []model.Profile
	result := s.DB.Where("user_id = ?", userID).Find(&profiles)
	return profiles, result.Error
}

// UPDATE
func (s *ProfileService) UpdateProfile(profile *model.Profile) error {
	if profile.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldProfile model.Profile
	if err := s.DB.First(&oldProfile, "id = ?", profile.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("profile not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if profile.Name != "" && profile.Name != oldProfile.Name {
		updateData["name"] = profile.Name
		// Regenerate slug jika name berubah
		updateData["slug"] = generateSlug(profile.Name)
	}

	if profile.Description != "" && profile.Description != oldProfile.Description {
		updateData["description"] = profile.Description
	}

	if profile.Address != nil && (oldProfile.Address == nil || *profile.Address != *oldProfile.Address) {
		updateData["address"] = profile.Address
	}

	if profile.Phone != nil && (oldProfile.Phone == nil || *profile.Phone != *oldProfile.Phone) {
		updateData["phone"] = profile.Phone
	}

	if profile.Image != nil && (oldProfile.Image == nil || *profile.Image != *oldProfile.Image) {
		updateData["image"] = profile.Image
	}

	if profile.Status != "" && profile.Status != oldProfile.Status {
		updateData["status"] = profile.Status
	}

	if profile.UserID != "" && profile.UserID != oldProfile.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", profile.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = profile.UserID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Profile{}).Where("id = ?", profile.ID).Updates(updateData).Error
}

// DELETE
func (s *ProfileService) DeleteProfile(id string) error {
	return s.DB.Delete(&model.Profile{}, "id = ?", id).Error
}

// generateSlug helper (duplikat dari model, untuk konsistensi)
func generateSlug(name string) string {
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = reg.ReplaceAllString(slug, "")
	reg2 := regexp.MustCompile(`-+`)
	slug = reg2.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	return slug
}
