package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScLearningCenterService struct {
	DB *gorm.DB
}

func NewScLearningCenterService(db *gorm.DB) *ScLearningCenterService {
	return &ScLearningCenterService{DB: db}
}

// CREATE
func (s *ScLearningCenterService) CreateScLearningCenter(learningCenter *model.ScLearningCenter) error {
	if learningCenter.ID == "" {
		learningCenter.ID = uuid.New().String()
	}

	if learningCenter.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", learningCenter.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(learningCenter).Error
}

// READ ALL
func (s *ScLearningCenterService) GetAllScLearningCenter() ([]model.ScLearningCenter, error) {
	var learningCenters []model.ScLearningCenter
	result := s.DB.Find(&learningCenters)
	return learningCenters, result.Error
}

// READ BY ID
func (s *ScLearningCenterService) GetScLearningCenterByID(id string) (*model.ScLearningCenter, error) {
	var learningCenter model.ScLearningCenter
	result := s.DB.First(&learningCenter, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &learningCenter, result.Error
}

// UPDATE
func (s *ScLearningCenterService) UpdateScLearningCenter(learningCenter *model.ScLearningCenter) error {
	if learningCenter.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldLearningCenter model.ScLearningCenter
	if err := s.DB.First(&oldLearningCenter, "id = ?", learningCenter.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_learning_center not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if learningCenter.UserID != "" && learningCenter.UserID != oldLearningCenter.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", learningCenter.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = learningCenter.UserID
	}

	if learningCenter.Name != nil && (oldLearningCenter.Name == nil || *learningCenter.Name != *oldLearningCenter.Name) {
		updateData["name"] = learningCenter.Name
	}

	if learningCenter.GroupIDN != nil && (oldLearningCenter.GroupIDN == nil || *learningCenter.GroupIDN != *oldLearningCenter.GroupIDN) {
		updateData["group_idn"] = learningCenter.GroupIDN
	}

	if learningCenter.Principal != nil && (oldLearningCenter.Principal == nil || *learningCenter.Principal != *oldLearningCenter.Principal) {
		updateData["principal"] = learningCenter.Principal
	}

	if learningCenter.HomeTeacher != nil && (oldLearningCenter.HomeTeacher == nil || *learningCenter.HomeTeacher != *oldLearningCenter.HomeTeacher) {
		updateData["home_teacher"] = learningCenter.HomeTeacher
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScLearningCenter{}).Where("id = ?", learningCenter.ID).Updates(updateData).Error
}

// DELETE
func (s *ScLearningCenterService) DeleteScLearningCenter(id string) error {
	return s.DB.Delete(&model.ScLearningCenter{}, "id = ?", id).Error
}

func (s *ScLearningCenterService) GetAllScLearningCenterByUserID(userID string) ([]model.ScLearningCenter, error) {
	var list []model.ScLearningCenter
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
