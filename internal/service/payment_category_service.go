package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentCategoryService struct {
	DB *gorm.DB
}

func NewPaymentCategoryService(db *gorm.DB) *PaymentCategoryService {
	return &PaymentCategoryService{DB: db}
}

func (s *PaymentCategoryService) CreatePaymentCategory(category *model.PaymentCategory) error {
	category.ID = uuid.New().String()
	return s.DB.Create(category).Error
}

func (s *PaymentCategoryService) GetAllPaymentCategories(userID string) ([]model.PaymentCategory, error) {
	var list []model.PaymentCategory
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (s *PaymentCategoryService) GetPaymentCategoryByID(id string) (*model.PaymentCategory, error) {
	var category model.PaymentCategory
	result := s.DB.First(&category, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, result.Error
}

func (s *PaymentCategoryService) UpdatePaymentCategory(category *model.PaymentCategory) error {
	if category.ID == "" {
		return fmt.Errorf("ID is required")
	}

	updateData := map[string]interface{}{}

	if category.Name != "" {
		updateData["name"] = category.Name
	}
	if category.Description != nil {
		updateData["description"] = category.Description
	}
	updateData["allow_penalty"] = category.AllowPenalty
	updateData["allow_installment"] = category.AllowInstallment

	if category.Status != "" {
		updateData["status"] = category.Status
	}
	if category.DueDateDefault != nil {
		updateData["due_date_default"] = category.DueDateDefault
	}

	return s.DB.Model(&model.PaymentCategory{}).Where("id = ?", category.ID).Updates(updateData).Error
}

func (s *PaymentCategoryService) DeletePaymentCategory(id string) error {
	return s.DB.Delete(&model.PaymentCategory{}, "id = ?", id).Error
}
