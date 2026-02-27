package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScTemplateHabitsService struct {
	DB *gorm.DB
}

func NewScTemplateHabitsService(db *gorm.DB) *ScTemplateHabitsService {
	return &ScTemplateHabitsService{DB: db}
}

// CREATE
func (s *ScTemplateHabitsService) CreateScTemplateHabits(template *model.ScTemplateHabits) error {
	if template.ID == "" {
		template.ID = uuid.New().String()
	}

	if template.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", template.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(template).Error
}

// READ ALL
func (s *ScTemplateHabitsService) GetAllScTemplateHabits() ([]model.ScTemplateHabits, error) {
	var templates []model.ScTemplateHabits
	result := s.DB.Find(&templates)
	return templates, result.Error
}

// READ BY ID
func (s *ScTemplateHabitsService) GetScTemplateHabitsByID(id string) (*model.ScTemplateHabits, error) {
	var template model.ScTemplateHabits
	result := s.DB.First(&template, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &template, result.Error
}

// UPDATE
func (s *ScTemplateHabitsService) UpdateScTemplateHabits(template *model.ScTemplateHabits) error {
	if template.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTemplate model.ScTemplateHabits
	if err := s.DB.First(&oldTemplate, "id = ?", template.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_template_habits not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if template.UserID != "" && template.UserID != oldTemplate.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", template.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = template.UserID
	}

	if template.Level != "" && template.Level != oldTemplate.Level {
		updateData["level"] = template.Level
	}

	if template.SubjectID != "" && template.SubjectID != oldTemplate.SubjectID {
		updateData["subject_id"] = template.SubjectID
	}

	if template.SubSubjectID != "" && template.SubSubjectID != oldTemplate.SubSubjectID {
		updateData["sub_subject_id"] = template.SubSubjectID
	}

	if template.Grade != "" && template.Grade != oldTemplate.Grade {
		updateData["grade"] = template.Grade
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScTemplateHabits{}).Where("id = ?", template.ID).Updates(updateData).Error
}

// DELETE
func (s *ScTemplateHabitsService) DeleteScTemplateHabits(id string) error {
	return s.DB.Delete(&model.ScTemplateHabits{}, "id = ?", id).Error
}

func (s *ScTemplateHabitsService) GetAllScTemplateHabitsByUserID(userID string) ([]model.ScTemplateHabits, error) {
	var list []model.ScTemplateHabits
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}