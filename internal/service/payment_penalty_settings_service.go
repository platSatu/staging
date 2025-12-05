package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentPenaltySettingsService struct {
	DB *gorm.DB
}

func NewPaymentPenaltySettingsService(db *gorm.DB) *PaymentPenaltySettingsService {
	return &PaymentPenaltySettingsService{DB: db}
}

// CREATE
func (s *PaymentPenaltySettingsService) CreatePaymentPenaltySettings(setting *model.PaymentPenaltySettings) error {
	if setting.UserID == "" {
		return fmt.Errorf("user_id is required")
	}
	if setting.CategoryID == "" {
		return fmt.Errorf("category_id is required")
	}
	if setting.PenaltyType == "" {
		return fmt.Errorf("penalty_type is required")
	}
	if setting.PenaltyType != "percent" && setting.PenaltyType != "flat" {
		return fmt.Errorf("penalty_type must be 'percent' or 'flat'")
	}
	if setting.PenaltyType == "percent" && setting.PercentValue == nil {
		return fmt.Errorf("percent_value is required when penalty_type is 'percent'")
	}
	if setting.PenaltyType == "flat" && setting.FlatValue == nil {
		return fmt.Errorf("flat_value is required when penalty_type is 'flat'")
	}
	if setting.MaxPenaltyAmount <= 0 {
		return fmt.Errorf("max_penalty_amount must be greater than 0")
	}
	if setting.ApplyOn == "" {
		return fmt.Errorf("apply_on is required")
	}
	if setting.ApplyOn != "invoice" && setting.ApplyOn != "installment" && setting.ApplyOn != "both" {
		return fmt.Errorf("apply_on must be 'invoice', 'installment', or 'both'")
	}
	if setting.Active == "" {
		setting.Active = "active" // Default jika tidak di-set
	}
	if setting.Active != "active" && setting.Active != "inactive" {
		return fmt.Errorf("active must be 'active' or 'inactive'")
	}

	if setting.ID == "" {
		setting.ID = uuid.New().String()
	}

	err := s.DB.Create(setting).Error
	if err != nil {
		return err
	}
	return nil
}

// READ ALL (filtered by user_id)
func (s *PaymentPenaltySettingsService) GetAllPaymentPenaltySettings(userID string) ([]model.PaymentPenaltySettings, error) {
	var settings []model.PaymentPenaltySettings
	err := s.DB.Where("user_id = ?", userID).Find(&settings).Error
	return settings, err
}

// READ BY ID
func (s *PaymentPenaltySettingsService) GetPaymentPenaltySettingsByID(id string) (*model.PaymentPenaltySettings, error) {
	var setting model.PaymentPenaltySettings
	result := s.DB.First(&setting, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &setting, result.Error
}

// UPDATE
func (s *PaymentPenaltySettingsService) UpdatePaymentPenaltySettings(setting *model.PaymentPenaltySettings, userID string) error {
	if setting.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldSetting model.PaymentPenaltySettings
	if err := s.DB.First(&oldSetting, "id = ?", setting.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment penalty setting not found")
		}
		return err
	}

	// Pastikan hanya pemilik yang bisa update
	if oldSetting.UserID != userID {
		return fmt.Errorf("unauthorized: you can only update your own settings")
	}

	updateData := map[string]interface{}{}

	if setting.CategoryID != "" {
		updateData["category_id"] = setting.CategoryID
	}
	if setting.PenaltyType != "" {
		if setting.PenaltyType != "percent" && setting.PenaltyType != "flat" {
			return fmt.Errorf("penalty_type must be 'percent' or 'flat'")
		}
		updateData["penalty_type"] = setting.PenaltyType
		// Validasi ulang berdasarkan penalty_type baru
		if setting.PenaltyType == "percent" && setting.PercentValue == nil {
			return fmt.Errorf("percent_value is required when penalty_type is 'percent'")
		}
		if setting.PenaltyType == "flat" && setting.FlatValue == nil {
			return fmt.Errorf("flat_value is required when penalty_type is 'flat'")
		}
	}
	if setting.PercentValue != nil {
		updateData["percent_value"] = setting.PercentValue
	}
	if setting.FlatValue != nil {
		updateData["flat_value"] = setting.FlatValue
	}
	if setting.MaxPenaltyAmount > 0 {
		updateData["max_penalty_amount"] = setting.MaxPenaltyAmount
	}
	if setting.ApplyOn != "" {
		if setting.ApplyOn != "invoice" && setting.ApplyOn != "installment" && setting.ApplyOn != "both" {
			return fmt.Errorf("apply_on must be 'invoice', 'installment', or 'both'")
		}
		updateData["apply_on"] = setting.ApplyOn
	}
	if setting.Active != "" {
		if setting.Active != "active" && setting.Active != "inactive" {
			return fmt.Errorf("active must be 'active' or 'inactive'")
		}
		updateData["active"] = setting.Active
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.PaymentPenaltySettings{}).Where("id = ?", setting.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentPenaltySettingsService) DeletePaymentPenaltySettings(id string, userID string) error {
	var setting model.PaymentPenaltySettings
	if err := s.DB.First(&setting, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment penalty setting not found")
		}
		return err
	}
	// Pastikan hanya pemilik yang bisa delete
	if setting.UserID != userID {
		return fmt.Errorf("unauthorized: you can only delete your own settings")
	}
	return s.DB.Delete(&model.PaymentPenaltySettings{}, "id = ?", id).Error
}
