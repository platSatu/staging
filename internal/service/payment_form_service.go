package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentFormService struct {
	DB *gorm.DB
}

func NewPaymentFormService(db *gorm.DB) *PaymentFormService {
	return &PaymentFormService{DB: db}
}

// CREATE
func (s *PaymentFormService) CreatePaymentForm(form *model.PaymentForm) error {
	if form.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if form.ID == "" {
		form.ID = uuid.New().String()
	}

	return s.DB.Create(form).Error
}

// READ ALL (filtered by user_id)
func (s *PaymentFormService) GetAllPaymentForms(userID string) ([]model.PaymentForm, error) {
	var forms []model.PaymentForm
	err := s.DB.Where("user_id = ?", userID).Find(&forms).Error
	return forms, err
}

// READ BY ID
func (s *PaymentFormService) GetPaymentFormByID(id string) (*model.PaymentForm, error) {
	var form model.PaymentForm
	result := s.DB.First(&form, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &form, result.Error
}

// UPDATE
func (s *PaymentFormService) UpdatePaymentForm(form *model.PaymentForm) error {
	if form.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldForm model.PaymentForm
	if err := s.DB.First(&oldForm, "id = ?", form.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment form not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if form.CategoryID != "" {
		updateData["category_id"] = form.CategoryID
	}
	if form.Name != "" {
		updateData["name"] = form.Name
	}
	if form.Description != nil {
		updateData["description"] = form.Description
	}
	if form.BaseAmount != 0 {
		updateData["base_amount"] = form.BaseAmount
	}
	// For booleans, always update if provided (since false is valid)
	updateData["enable_installment"] = form.EnableInstallment
	updateData["enable_penalty"] = form.EnablePenalty
	if form.Status != "" {
		updateData["status"] = form.Status
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	return s.DB.Model(&model.PaymentForm{}).Where("id = ?", form.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentFormService) DeletePaymentForm(id string) error {
	if id == "" {
		return fmt.Errorf("ID is required for delete")
	}
	return s.DB.Delete(&model.PaymentForm{}, "id = ?", id).Error
}
