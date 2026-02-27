package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScSubjectListSemesterYearlyService struct {
	DB *gorm.DB
}

func NewScSubjectListSemesterYearlyService(db *gorm.DB) *ScSubjectListSemesterYearlyService {
	return &ScSubjectListSemesterYearlyService{DB: db}
}

// CREATE
func (s *ScSubjectListSemesterYearlyService) CreateScSubjectListSemesterYearly(subject *model.ScSubjectListSemesterYearly) error {
	if subject.ID == "" {
		subject.ID = uuid.New().String()
	}

	if subject.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", subject.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(subject).Error
}

// READ ALL
func (s *ScSubjectListSemesterYearlyService) GetAllScSubjectListSemesterYearly() ([]model.ScSubjectListSemesterYearly, error) {
	var subjects []model.ScSubjectListSemesterYearly
	result := s.DB.Find(&subjects)
	return subjects, result.Error
}

// READ BY ID
func (s *ScSubjectListSemesterYearlyService) GetScSubjectListSemesterYearlyByID(id string) (*model.ScSubjectListSemesterYearly, error) {
	var subject model.ScSubjectListSemesterYearly
	result := s.DB.First(&subject, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subject, result.Error
}

// UPDATE
func (s *ScSubjectListSemesterYearlyService) UpdateScSubjectListSemesterYearly(subject *model.ScSubjectListSemesterYearly) error {
	if subject.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldSubject model.ScSubjectListSemesterYearly
	if err := s.DB.First(&oldSubject, "id = ?", subject.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_subject_list_semester_yearly not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if subject.UserID != "" && subject.UserID != oldSubject.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", subject.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = subject.UserID
	}

	if subject.StartLevel != "" && subject.StartLevel != oldSubject.StartLevel {
		updateData["start_level"] = subject.StartLevel
	}

	if subject.EndLevel != "" && subject.EndLevel != oldSubject.EndLevel {
		updateData["end_level"] = subject.EndLevel
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScSubjectListSemesterYearly{}).Where("id = ?", subject.ID).Updates(updateData).Error
}

// DELETE
func (s *ScSubjectListSemesterYearlyService) DeleteScSubjectListSemesterYearly(id string) error {
	return s.DB.Delete(&model.ScSubjectListSemesterYearly{}, "id = ?", id).Error
}

func (s *ScSubjectListSemesterYearlyService) GetAllScSubjectListSemesterYearlyByUserID(userID string) ([]model.ScSubjectListSemesterYearly, error) {
	var list []model.ScSubjectListSemesterYearly
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}