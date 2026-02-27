package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScSubjectTypeGeneralService struct {
	DB *gorm.DB
}

func NewScSubjectTypeGeneralService(db *gorm.DB) *ScSubjectTypeGeneralService {
	return &ScSubjectTypeGeneralService{DB: db}
}

// CREATE
func (s *ScSubjectTypeGeneralService) CreateScSubjectTypeGeneral(subjectTypeGeneral *model.ScSubjectTypeGeneral) error {
	if subjectTypeGeneral.ID == "" {
		subjectTypeGeneral.ID = uuid.New().String()
	}

	if subjectTypeGeneral.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", subjectTypeGeneral.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(subjectTypeGeneral).Error
}

// READ ALL
func (s *ScSubjectTypeGeneralService) GetAllScSubjectTypeGeneral() ([]model.ScSubjectTypeGeneral, error) {
	var subjectTypeGenerals []model.ScSubjectTypeGeneral
	result := s.DB.Find(&subjectTypeGenerals)
	return subjectTypeGenerals, result.Error
}

// READ BY ID
func (s *ScSubjectTypeGeneralService) GetScSubjectTypeGeneralByID(id string) (*model.ScSubjectTypeGeneral, error) {
	var subjectTypeGeneral model.ScSubjectTypeGeneral
	result := s.DB.First(&subjectTypeGeneral, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subjectTypeGeneral, result.Error
}

// UPDATE
func (s *ScSubjectTypeGeneralService) UpdateScSubjectTypeGeneral(subjectTypeGeneral *model.ScSubjectTypeGeneral) error {
	if subjectTypeGeneral.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldSubjectTypeGeneral model.ScSubjectTypeGeneral
	if err := s.DB.First(&oldSubjectTypeGeneral, "id = ?", subjectTypeGeneral.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_subject_type_general not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if subjectTypeGeneral.UserID != "" && subjectTypeGeneral.UserID != oldSubjectTypeGeneral.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", subjectTypeGeneral.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = subjectTypeGeneral.UserID
	}

	if subjectTypeGeneral.SubjectName != nil && (oldSubjectTypeGeneral.SubjectName == nil || *subjectTypeGeneral.SubjectName != *oldSubjectTypeGeneral.SubjectName) {
		updateData["subject_name"] = subjectTypeGeneral.SubjectName
	}

	if subjectTypeGeneral.Curriculum != nil && (oldSubjectTypeGeneral.Curriculum == nil || *subjectTypeGeneral.Curriculum != *oldSubjectTypeGeneral.Curriculum) {
		updateData["curriculum"] = subjectTypeGeneral.Curriculum
	}

	if subjectTypeGeneral.Group != nil && (oldSubjectTypeGeneral.Group == nil || *subjectTypeGeneral.Group != *oldSubjectTypeGeneral.Group) {
		updateData["group"] = subjectTypeGeneral.Group
	}

	if subjectTypeGeneral.Status != nil && (oldSubjectTypeGeneral.Status == nil || *subjectTypeGeneral.Status != *oldSubjectTypeGeneral.Status) {
		updateData["status"] = subjectTypeGeneral.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScSubjectTypeGeneral{}).Where("id = ?", subjectTypeGeneral.ID).Updates(updateData).Error
}

// DELETE
func (s *ScSubjectTypeGeneralService) DeleteScSubjectTypeGeneral(id string) error {
	return s.DB.Delete(&model.ScSubjectTypeGeneral{}, "id = ?", id).Error
}

func (s *ScSubjectTypeGeneralService) GetAllScSubjectTypeGeneralByUserID(userID string) ([]model.ScSubjectTypeGeneral, error) {
	var list []model.ScSubjectTypeGeneral
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
