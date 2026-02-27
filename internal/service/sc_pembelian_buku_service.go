package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScPembelianBukuService struct {
	DB *gorm.DB
}

func NewScPembelianBukuService(db *gorm.DB) *ScPembelianBukuService {
	return &ScPembelianBukuService{DB: db}
}

// CREATE
func (s *ScPembelianBukuService) CreateScPembelianBuku(pembelianBuku *model.ScPembelianBuku) error {
	if pembelianBuku.ID == "" {
		pembelianBuku.ID = uuid.New().String()
	}

	if pembelianBuku.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", pembelianBuku.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(pembelianBuku).Error
}

// READ ALL
func (s *ScPembelianBukuService) GetAllScPembelianBuku() ([]model.ScPembelianBuku, error) {
	var pembelianBukus []model.ScPembelianBuku
	result := s.DB.Find(&pembelianBukus)
	return pembelianBukus, result.Error
}

// READ BY ID
func (s *ScPembelianBukuService) GetScPembelianBukuByID(id string) (*model.ScPembelianBuku, error) {
	var pembelianBuku model.ScPembelianBuku
	result := s.DB.First(&pembelianBuku, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &pembelianBuku, result.Error
}

// UPDATE
func (s *ScPembelianBukuService) UpdateScPembelianBuku(pembelianBuku *model.ScPembelianBuku) error {
	if pembelianBuku.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldPembelianBuku model.ScPembelianBuku
	if err := s.DB.First(&oldPembelianBuku, "id = ?", pembelianBuku.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_pembelian_buku not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if pembelianBuku.UserID != "" && pembelianBuku.UserID != oldPembelianBuku.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", pembelianBuku.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = pembelianBuku.UserID
	}

	if pembelianBuku.ParentID != nil && (oldPembelianBuku.ParentID == nil || *pembelianBuku.ParentID != *oldPembelianBuku.ParentID) {
		updateData["parent_id"] = pembelianBuku.ParentID
	}

	if pembelianBuku.Subject != nil && (oldPembelianBuku.Subject == nil || *pembelianBuku.Subject != *oldPembelianBuku.Subject) {
		updateData["subject"] = pembelianBuku.Subject
	}

	if pembelianBuku.PacesNumber != nil && (oldPembelianBuku.PacesNumber == nil || *pembelianBuku.PacesNumber != *oldPembelianBuku.PacesNumber) {
		updateData["paces_number"] = pembelianBuku.PacesNumber
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScPembelianBuku{}).Where("id = ?", pembelianBuku.ID).Updates(updateData).Error
}

// DELETE
func (s *ScPembelianBukuService) DeleteScPembelianBuku(id string) error {
	return s.DB.Delete(&model.ScPembelianBuku{}, "id = ?", id).Error
}

// GET BY USER ID
func (s *ScPembelianBukuService) GetAllScPembelianBukuByUserID(userID string) ([]model.ScPembelianBuku, error) {
	var list []model.ScPembelianBuku
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
