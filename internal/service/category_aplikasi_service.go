package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryAplikasiService struct {
	DB *gorm.DB
}

func NewCategoryAplikasiService(db *gorm.DB) *CategoryAplikasiService {
	return &CategoryAplikasiService{DB: db}
}

// CREATE
func (s *CategoryAplikasiService) CreateCategoryAplikasi(category *model.CategoryAplikasi) error {
	if category.ID == "" {
		category.ID = uuid.New().String()
	}

	if category.Status == "" {
		category.Status = "active"
	}

	// Validasi user_id ada di tabel users (opsional, jika ingin foreign key constraint)
	var user model.User
	if err := s.DB.First(&user, "id = ?", category.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	return s.DB.Create(category).Error
}

// READ ALL
func (s *CategoryAplikasiService) GetAllCategoryAplikasi() ([]model.CategoryAplikasi, error) {
	var categories []model.CategoryAplikasi
	result := s.DB.Find(&categories)
	return categories, result.Error
}

// READ BY ID
func (s *CategoryAplikasiService) GetCategoryAplikasiByID(id string) (*model.CategoryAplikasi, error) {
	var category model.CategoryAplikasi
	result := s.DB.First(&category, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, result.Error
}

// READ BY USER ID
func (s *CategoryAplikasiService) GetCategoryAplikasiByUserID(userID string) ([]model.CategoryAplikasi, error) {
	var categories []model.CategoryAplikasi
	result := s.DB.Where("user_id = ?", userID).Find(&categories)
	return categories, result.Error
}

// UPDATE
func (s *CategoryAplikasiService) UpdateCategoryAplikasi(category *model.CategoryAplikasi) error {
	if category.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldCategory model.CategoryAplikasi
	if err := s.DB.First(&oldCategory, "id = ?", category.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("category not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if category.Name != "" && category.Name != oldCategory.Name {
		updateData["name"] = category.Name
	}

	if category.Status != "" && category.Status != oldCategory.Status {
		updateData["status"] = category.Status
	}

	if category.UserID != "" && category.UserID != oldCategory.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", category.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = category.UserID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.CategoryAplikasi{}).Where("id = ?", category.ID).Updates(updateData).Error
}

// DELETE
func (s *CategoryAplikasiService) DeleteCategoryAplikasi(id string) error {
	return s.DB.Delete(&model.CategoryAplikasi{}, "id = ?", id).Error
}
