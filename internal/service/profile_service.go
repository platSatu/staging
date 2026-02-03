package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProfileService struct {
	DB *gorm.DB
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{DB: db}
}

// CREATE - Profil dibuat untuk user yang login
func (s *ProfileService) CreateProfile(profile *model.Profile, userID string) error {
	if profile.ID == "" {
		profile.ID = uuid.New().String()
	}

	if profile.Name == "" {
		return fmt.Errorf("name is required")
	}

	if profile.Status == "" {
		profile.Status = "active"
	}

	// Set UserID dari parameter (dari token)
	profile.UserID = userID

	// Validasi user_id ada di tabel users
	var user model.User
	if err := s.DB.First(&user, "id = ?", profile.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// BussinesID sekarang auto-increment, tidak perlu generate manual
	return s.DB.Create(profile).Error
}

// READ ALL - Ubah jadi GetProfilesByUserID (user hanya lihat profilnya sendiri)
func (s *ProfileService) GetAllProfiles(userID string) ([]model.Profile, error) {
	// Alias untuk GetProfilesByUserID, karena semua user hanya lihat profil mereka sendiri
	return s.GetProfilesByUserID(userID)
}

// READ BY ID - User hanya bisa lihat profilnya sendiri
func (s *ProfileService) GetProfileByID(id string, userID string) (*model.Profile, error) {
	var profile model.Profile
	result := s.DB.Where("id = ? AND user_id = ?", id, userID).First(&profile)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &profile, result.Error
}

// READ BY USER ID - Sudah filter berdasarkan userID
func (s *ProfileService) GetProfilesByUserID(userID string) ([]model.Profile, error) {
	var profiles []model.Profile
	result := s.DB.Where("user_id = ?", userID).Find(&profiles)
	return profiles, result.Error
}

// CHECK USER HAS PROFILE - Mengecek apakah user sudah punya profile
func (s *ProfileService) CheckUserHasProfile(userID string) (bool, error) {
	var count int64
	err := s.DB.Model(&model.Profile{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check user profile: %v", err)
	}
	return count > 0, nil
}

// UPDATE - User hanya bisa update profilnya sendiri
func (s *ProfileService) UpdateProfile(profile *model.Profile, userID string) error {
	if profile.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldProfile model.Profile
	if err := s.DB.Where("id = ? AND user_id = ?", profile.ID, userID).First(&oldProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("profile not found or access denied")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if profile.Name != "" && profile.Name != oldProfile.Name {
		updateData["name"] = profile.Name
		updateData["slug"] = model.GenerateSlug(profile.Name)
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

	// User tidak bisa ubah user_id (hanya milik mereka sendiri)
	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Profile{}).Where("id = ? AND user_id = ?", profile.ID, userID).Updates(updateData).Error
}

// DELETE - User hanya bisa hapus profilnya sendiri
func (s *ProfileService) DeleteProfile(id string, userID string) error {
	return s.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Profile{}).Error
}
