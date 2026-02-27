package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScSubjectTypeGeneralSubjectService struct {
	DB *gorm.DB
}

func NewScSubjectTypeGeneralSubjectService(db *gorm.DB) *ScSubjectTypeGeneralSubjectService {
	return &ScSubjectTypeGeneralSubjectService{DB: db}
}

// CREATE
func (s *ScSubjectTypeGeneralSubjectService) CreateScSubjectTypeGeneralSubject(subject *model.ScSubjectTypeGeneralSubject) error {
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
func (s *ScSubjectTypeGeneralSubjectService) GetAllScSubjectTypeGeneralSubject() ([]model.ScSubjectTypeGeneralSubject, error) {
	var subjects []model.ScSubjectTypeGeneralSubject
	result := s.DB.Find(&subjects)
	return subjects, result.Error
}

// READ BY ID
func (s *ScSubjectTypeGeneralSubjectService) GetScSubjectTypeGeneralSubjectByID(id string) (*model.ScSubjectTypeGeneralSubject, error) {
	var subject model.ScSubjectTypeGeneralSubject
	result := s.DB.First(&subject, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subject, result.Error
}

// UPDATE
func (s *ScSubjectTypeGeneralSubjectService) UpdateScSubjectTypeGeneralSubject(subject *model.ScSubjectTypeGeneralSubject) error {
	if subject.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldSubject model.ScSubjectTypeGeneralSubject
	if err := s.DB.First(&oldSubject, "id = ?", subject.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_subject_type_general_subject not found")
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

	if subject.GpaWeight != nil && (oldSubject.GpaWeight == nil || *subject.GpaWeight != *oldSubject.GpaWeight) {
		updateData["gpa_weight"] = subject.GpaWeight
	}

	if subject.MinPassingScore != nil && (oldSubject.MinPassingScore == nil || *subject.MinPassingScore != *oldSubject.MinPassingScore) {
		updateData["min_passing_score"] = subject.MinPassingScore
	}

	if subject.StartPage != nil && (oldSubject.StartPage == nil || *subject.StartPage != *oldSubject.StartPage) {
		updateData["start_page"] = subject.StartPage
	}

	if subject.EndingPace != nil && (oldSubject.EndingPace == nil || *subject.EndingPace != *oldSubject.EndingPace) {
		updateData["ending_pace"] = subject.EndingPace
	}

	if subject.SubjectName != nil && (oldSubject.SubjectName == nil || *subject.SubjectName != *oldSubject.SubjectName) {
		updateData["subject_name"] = subject.SubjectName
	}

	if subject.ProductReference != nil && (oldSubject.ProductReference == nil || *subject.ProductReference != *oldSubject.ProductReference) {
		updateData["product_reference"] = subject.ProductReference
	}

	if subject.StartingPace != nil && (oldSubject.StartingPace == nil || *subject.StartingPace != *oldSubject.StartingPace) {
		updateData["starting_pace"] = subject.StartingPace
	}

	if subject.Unit != nil && (oldSubject.Unit == nil || *subject.Unit != *oldSubject.Unit) {
		updateData["unit"] = subject.Unit
	}

	if subject.IsAssignable != nil && (oldSubject.IsAssignable == nil || *subject.IsAssignable != *oldSubject.IsAssignable) {
		updateData["is_assignable"] = subject.IsAssignable
	}

	if subject.IsPace != nil && (oldSubject.IsPace == nil || *subject.IsPace != *oldSubject.IsPace) {
		updateData["is_pace"] = subject.IsPace
	}

	if subject.PacesNumber != nil && (oldSubject.PacesNumber == nil || *subject.PacesNumber != *oldSubject.PacesNumber) {
		updateData["paces_number"] = subject.PacesNumber
	}

	if subject.SubjectType != nil && (oldSubject.SubjectType == nil || *subject.SubjectType != *oldSubject.SubjectType) {
		updateData["subject_type"] = subject.SubjectType
	}

	if subject.EndPage != nil && (oldSubject.EndPage == nil || *subject.EndPage != *oldSubject.EndPage) {
		updateData["end_page"] = subject.EndPage
	}

	if subject.PerCredit != nil && (oldSubject.PerCredit == nil || *subject.PerCredit != *oldSubject.PerCredit) {
		updateData["per_credit"] = subject.PerCredit
	}

	if subject.TotalPages != nil && (oldSubject.TotalPages == nil || *subject.TotalPages != *oldSubject.TotalPages) {
		updateData["total_pages"] = subject.TotalPages
	}

	if subject.CourseName != nil && (oldSubject.CourseName == nil || *subject.CourseName != *oldSubject.CourseName) {
		updateData["course_name"] = subject.CourseName
	}

	if subject.PrevPage != nil && (oldSubject.PrevPage == nil || *subject.PrevPage != *oldSubject.PrevPage) {
		updateData["prev_page"] = subject.PrevPage
	}

	if subject.NextPage != nil && (oldSubject.NextPage == nil || *subject.NextPage != *oldSubject.NextPage) {
		updateData["next_page"] = subject.NextPage
	}

	if subject.Units != nil && (oldSubject.Units == nil || *subject.Units != *oldSubject.Units) {
		updateData["units"] = subject.Units
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScSubjectTypeGeneralSubject{}).Where("id = ?", subject.ID).Updates(updateData).Error
}

// DELETE
func (s *ScSubjectTypeGeneralSubjectService) DeleteScSubjectTypeGeneralSubject(id string) error {
	return s.DB.Delete(&model.ScSubjectTypeGeneralSubject{}, "id = ?", id).Error
}

// READ ALL BY USER ID
func (s *ScSubjectTypeGeneralSubjectService) GetAllScSubjectTypeGeneralSubjectByUserID(userID string) ([]model.ScSubjectTypeGeneralSubject, error) {
	var list []model.ScSubjectTypeGeneralSubject
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
