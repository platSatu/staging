package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentPenaltiesService struct {
	DB *gorm.DB
}

func NewPaymentPenaltiesService(db *gorm.DB) *PaymentPenaltiesService {
	return &PaymentPenaltiesService{DB: db}
}

// CREATE
func (s *PaymentPenaltiesService) CreatePaymentPenalties(penalty *model.PaymentPenalty) error {
	if penalty.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Generate UUID for ID before creating
	penalty.ID = uuid.New().String()

	err := s.DB.Create(penalty).Error
	if err != nil {
		return err
	}
	return nil
}

// READ ALL (filtered by user_id)
func (s *PaymentPenaltiesService) GetAllPaymentPenalties(userID string) ([]model.PaymentPenalty, error) {
	var penalties []model.PaymentPenalty
	err := s.DB.Where("user_id = ?", userID).Find(&penalties).Error
	return penalties, err
}

// READ BY ID
func (s *PaymentPenaltiesService) GetPaymentPenaltiesByID(id string) (*model.PaymentPenalty, error) {
	var penalty model.PaymentPenalty
	result := s.DB.First(&penalty, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &penalty, result.Error
}

// UPDATE
func (s *PaymentPenaltiesService) UpdatePaymentPenalties(penalty *model.PaymentPenalty) error {
	if penalty.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldPenalty model.PaymentPenalty
	if err := s.DB.First(&oldPenalty, "id = ?", penalty.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment penalty not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if penalty.CategoryID != "" {
		updateData["category_id"] = penalty.CategoryID
	}
	if penalty.Name != "" {
		updateData["name"] = penalty.Name
	}
	if penalty.Description != nil {
		updateData["description"] = penalty.Description
	}
	if penalty.PenaltyType != "" {
		updateData["penalty_type"] = penalty.PenaltyType
	}
	if penalty.FlatValue != nil {
		updateData["flat_value"] = penalty.FlatValue
	}
	if penalty.PercentValue != nil {
		updateData["percent_value"] = penalty.PercentValue
	}
	if penalty.MaxPenalty != 0 {
		updateData["max_penalty"] = penalty.MaxPenalty
	}
	if penalty.ApplyOn != "" {
		updateData["apply_on"] = penalty.ApplyOn
	}
	if penalty.Active != "" {
		updateData["active"] = penalty.Active
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	return s.DB.Model(&model.PaymentPenalty{}).Where("id = ?", penalty.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentPenaltiesService) DeletePaymentPenalties(id string) error {
	if id == "" {
		return fmt.Errorf("ID is required for delete")
	}
	return s.DB.Delete(&model.PaymentPenalty{}, "id = ?", id).Error
}
