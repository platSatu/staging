package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScTeacherService struct {
	DB *gorm.DB
}

func NewScTeacherService(db *gorm.DB) *ScTeacherService {
	return &ScTeacherService{DB: db}
}

// CREATE
func (s *ScTeacherService) CreateScTeacher(teacher *model.ScTeacher) error {
	if teacher.ID == "" {
		teacher.ID = uuid.New().String()
	}

	if teacher.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", teacher.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(teacher).Error
}

// READ ALL
func (s *ScTeacherService) GetAllScTeacher() ([]model.ScTeacher, error) {
	var teachers []model.ScTeacher
	result := s.DB.Find(&teachers)
	return teachers, result.Error
}

// READ BY ID
func (s *ScTeacherService) GetScTeacherByID(id string) (*model.ScTeacher, error) {
	var teacher model.ScTeacher
	result := s.DB.First(&teacher, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &teacher, result.Error
}

// UPDATE
func (s *ScTeacherService) UpdateScTeacher(teacher *model.ScTeacher) error {
	if teacher.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTeacher model.ScTeacher
	if err := s.DB.First(&oldTeacher, "id = ?", teacher.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_teacher not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if teacher.UserID != "" && teacher.UserID != oldTeacher.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", teacher.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = teacher.UserID
	}

	if teacher.Name != nil && (oldTeacher.Name == nil || *teacher.Name != *oldTeacher.Name) {
		updateData["name"] = teacher.Name
	}

	if teacher.Email != nil && (oldTeacher.Email == nil || *teacher.Email != *oldTeacher.Email) {
		updateData["email"] = teacher.Email
	}

	if teacher.Phone != nil && (oldTeacher.Phone == nil || *teacher.Phone != *oldTeacher.Phone) {
		updateData["phone"] = teacher.Phone
	}

	if teacher.Mobile != nil && (oldTeacher.Mobile == nil || *teacher.Mobile != *oldTeacher.Mobile) {
		updateData["mobile"] = teacher.Mobile
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScTeacher{}).Where("id = ?", teacher.ID).Updates(updateData).Error
}

// DELETE
func (s *ScTeacherService) DeleteScTeacher(id string) error {
	return s.DB.Delete(&model.ScTeacher{}, "id = ?", id).Error
}

func (s *ScTeacherService) GetAllScTeacherByUserID(userID string) ([]model.ScTeacher, error) {
	var list []model.ScTeacher
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
