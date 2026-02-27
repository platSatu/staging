package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScSubjectListSemesterService struct {
	DB *gorm.DB
}

func NewScSubjectListSemesterService(db *gorm.DB) *ScSubjectListSemesterService {
	return &ScSubjectListSemesterService{DB: db}
}

// CREATE
func (s *ScSubjectListSemesterService) CreateScSubjectListSemester(subject *model.ScSubjectListSemester) error {
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
func (s *ScSubjectListSemesterService) GetAllScSubjectListSemester() ([]model.ScSubjectListSemester, error) {
	var subjects []model.ScSubjectListSemester
	result := s.DB.Find(&subjects)
	return subjects, result.Error
}

// READ BY ID
func (s *ScSubjectListSemesterService) GetScSubjectListSemesterByID(id string) (*model.ScSubjectListSemester, error) {
	var subject model.ScSubjectListSemester
	result := s.DB.First(&subject, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subject, result.Error
}

// UPDATE
func (s *ScSubjectListSemesterService) UpdateScSubjectListSemester(subject *model.ScSubjectListSemester) error {
	if subject.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldSubject model.ScSubjectListSemester
	if err := s.DB.First(&oldSubject, "id = ?", subject.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_subject_list_semester not found")
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

	if subject.No != "" && subject.No != oldSubject.No {
		updateData["no"] = subject.No
	}

	if subject.Subject != "" && subject.Subject != oldSubject.Subject {
		updateData["subject"] = subject.Subject
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScSubjectListSemester{}).Where("id = ?", subject.ID).Updates(updateData).Error
}

// DELETE
func (s *ScSubjectListSemesterService) DeleteScSubjectListSemester(id string) error {
	return s.DB.Delete(&model.ScSubjectListSemester{}, "id = ?", id).Error
}

func (s *ScSubjectListSemesterService) GetAllScSubjectListSemesterByUserID(userID string) ([]model.ScSubjectListSemester, error) {
	var list []model.ScSubjectListSemester
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}