package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TypeUserAplikasiService struct {
	DB *gorm.DB
}

func NewTypeUserAplikasiService(db *gorm.DB) *TypeUserAplikasiService {
	return &TypeUserAplikasiService{DB: db}
}

// CREATE
func (s *TypeUserAplikasiService) CreateTypeUserAplikasi(t *model.TypeUserAplikasi) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	if t.Status == "" {
		t.Status = "active"
	}

	// Validasi user_id ada di tabel users (opsional, jika ingin foreign key constraint)
	var user model.User
	if err := s.DB.First(&user, "id = ?", t.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// Validasi parent_id ada di tabel type_user_aplikasi (opsional, jika ingin foreign key constraint)
	var parent model.TypeUserAplikasi
	if err := s.DB.First(&parent, "id = ?", t.ParentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("parent not found")
		}
		return err
	}

	// Validasi aplikasi_id ada di tabel category_aplikasi (opsional, jika ingin foreign key constraint)
	var category model.CategoryAplikasi
	if err := s.DB.First(&category, "id = ?", t.AplikasiID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("category not found")
		}
		return err
	}

	return s.DB.Create(t).Error
}

// READ ALL
func (s *TypeUserAplikasiService) GetAllTypeUserAplikasi() ([]model.TypeUserAplikasi, error) {
	var types []model.TypeUserAplikasi
	result := s.DB.Find(&types)
	return types, result.Error
}

// READ BY ID
func (s *TypeUserAplikasiService) GetTypeUserAplikasiByID(id string) (*model.TypeUserAplikasi, error) {
	var t model.TypeUserAplikasi
	result := s.DB.First(&t, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &t, result.Error
}

// READ BY USER ID
func (s *TypeUserAplikasiService) GetTypeUserAplikasiByUserID(userID string) ([]model.TypeUserAplikasi, error) {
	var types []model.TypeUserAplikasi
	result := s.DB.Where("user_id = ?", userID).Find(&types)
	return types, result.Error
}

// UPDATE
func (s *TypeUserAplikasiService) UpdateTypeUserAplikasi(t *model.TypeUserAplikasi) error {
	if t.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldType model.TypeUserAplikasi
	if err := s.DB.First(&oldType, "id = ?", t.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("type user aplikasi not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if t.UserID != "" && t.UserID != oldType.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", t.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = t.UserID
	}

	if t.ParentID != "" && t.ParentID != oldType.ParentID {
		// Validasi parent_id baru ada
		var parent model.TypeUserAplikasi
		if err := s.DB.First(&parent, "id = ?", t.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent not found")
			}
			return err
		}
		updateData["parent_id"] = t.ParentID
	}

	if t.AplikasiID != "" && t.AplikasiID != oldType.AplikasiID {
		// Validasi aplikasi_id baru ada
		var category model.CategoryAplikasi
		if err := s.DB.First(&category, "id = ?", t.AplikasiID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("category not found")
			}
			return err
		}
		updateData["aplikasi_id"] = t.AplikasiID
	}

	if t.Status != "" && t.Status != oldType.Status {
		updateData["status"] = t.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TypeUserAplikasi{}).Where("id = ?", t.ID).Updates(updateData).Error
}

// DELETE
func (s *TypeUserAplikasiService) DeleteTypeUserAplikasi(id string) error {
	return s.DB.Delete(&model.TypeUserAplikasi{}, "id = ?", id).Error
}
