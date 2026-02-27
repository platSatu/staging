package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScAcademicYearService struct {
	DB *gorm.DB
}

func NewScAcademicYearService(db *gorm.DB) *ScAcademicYearService {
	return &ScAcademicYearService{DB: db}
}

// CREATE
func (s *ScAcademicYearService) CreateScAcademicYear(academicYear *model.ScAcademicYear) error {
	if academicYear.ID == "" {
		academicYear.ID = uuid.New().String()
	}

	if academicYear.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi status enum
	if academicYear.Status != "open" && academicYear.Status != "active" && academicYear.Status != "closed" {
		return fmt.Errorf("status must be 'open', 'active', or 'closed'")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", academicYear.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(academicYear).Error
}

// READ ALL
func (s *ScAcademicYearService) GetAllScAcademicYear() ([]model.ScAcademicYear, error) {
	var academicYears []model.ScAcademicYear
	result := s.DB.Find(&academicYears)
	return academicYears, result.Error
}

// READ BY ID
func (s *ScAcademicYearService) GetScAcademicYearByID(id string) (*model.ScAcademicYear, error) {
	var academicYear model.ScAcademicYear
	result := s.DB.First(&academicYear, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &academicYear, result.Error
}

// UPDATE
func (s *ScAcademicYearService) UpdateScAcademicYear(academicYear *model.ScAcademicYear) error {
	if academicYear.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldAcademicYear model.ScAcademicYear
	if err := s.DB.First(&oldAcademicYear, "id = ?", academicYear.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_academic_year not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if academicYear.UserID != "" && academicYear.UserID != oldAcademicYear.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", academicYear.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = academicYear.UserID
	}

	if academicYear.Name != nil && (oldAcademicYear.Name == nil || *academicYear.Name != *oldAcademicYear.Name) {
		updateData["name"] = academicYear.Name
	}

	if academicYear.Status != "" && academicYear.Status != oldAcademicYear.Status {
		// Validasi status enum
		if academicYear.Status != "open" && academicYear.Status != "active" && academicYear.Status != "closed" {
			return fmt.Errorf("status must be 'open', 'active', or 'closed'")
		}
		updateData["status"] = academicYear.Status
	}

	if academicYear.BeginDate != nil && (oldAcademicYear.BeginDate == nil || !academicYear.BeginDate.Equal(*oldAcademicYear.BeginDate)) {
		updateData["begin_date"] = academicYear.BeginDate
	}

	if academicYear.EndDate != nil && (oldAcademicYear.EndDate == nil || !academicYear.EndDate.Equal(*oldAcademicYear.EndDate)) {
		updateData["end_date"] = academicYear.EndDate
	}

	if academicYear.AcademicYear != nil && (oldAcademicYear.AcademicYear == nil || *academicYear.AcademicYear != *oldAcademicYear.AcademicYear) {
		updateData["academic_year"] = academicYear.AcademicYear
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScAcademicYear{}).Where("id = ?", academicYear.ID).Updates(updateData).Error
}

// DELETE
func (s *ScAcademicYearService) DeleteScAcademicYear(id string) error {
	return s.DB.Delete(&model.ScAcademicYear{}, "id = ?", id).Error
}

func (s *ScAcademicYearService) GetAllScAcademicYearByUserID(userID string) ([]model.ScAcademicYear, error) {
	var list []model.ScAcademicYear
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
