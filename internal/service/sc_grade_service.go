package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScGradeService struct {
	DB *gorm.DB
}

func NewScGradeService(db *gorm.DB) *ScGradeService {
	return &ScGradeService{DB: db}
}

// CREATE
func (s *ScGradeService) CreateScGrade(grade *model.ScGrade) error {
	if grade.ID == "" {
		grade.ID = uuid.New().String()
	}

	if grade.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", grade.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(grade).Error
}

// READ ALL
func (s *ScGradeService) GetAllScGrade() ([]model.ScGrade, error) {
	var grades []model.ScGrade
	result := s.DB.Find(&grades)
	return grades, result.Error
}

// READ BY ID
func (s *ScGradeService) GetScGradeByID(id string) (*model.ScGrade, error) {
	var grade model.ScGrade
	result := s.DB.First(&grade, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &grade, result.Error
}

// UPDATE
func (s *ScGradeService) UpdateScGrade(grade *model.ScGrade) error {
	if grade.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldGrade model.ScGrade
	if err := s.DB.First(&oldGrade, "id = ?", grade.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_grade not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if grade.UserID != "" && grade.UserID != oldGrade.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", grade.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = grade.UserID
	}

	if grade.GradeName != "" && grade.GradeName != oldGrade.GradeName {
		updateData["grade_name"] = grade.GradeName
	}

	if grade.MinScore != "" && grade.MinScore != oldGrade.MinScore {
		updateData["min_score"] = grade.MinScore
	}

	if grade.MaxScore != "" && grade.MaxScore != oldGrade.MaxScore {
		updateData["max_score"] = grade.MaxScore
	}

	if grade.Status != "" && grade.Status != oldGrade.Status {
		updateData["status"] = grade.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScGrade{}).Where("id = ?", grade.ID).Updates(updateData).Error
}

// DELETE
func (s *ScGradeService) DeleteScGrade(id string) error {
	return s.DB.Delete(&model.ScGrade{}, "id = ?", id).Error
}

func (s *ScGradeService) GetAllScGradeByUserID(userID string) ([]model.ScGrade, error) {
	var list []model.ScGrade
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}