package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentInstallmentsService struct {
	DB *gorm.DB
}

func NewPaymentInstallmentsService(db *gorm.DB) *PaymentInstallmentsService {
	return &PaymentInstallmentsService{DB: db}
}

// CREATE
func (s *PaymentInstallmentsService) CreatePaymentInstallments(installment *model.PaymentInstallments) error {
	if installment.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if installment.ID == "" {
		installment.ID = uuid.New().String()
	}

	err := s.DB.Create(installment).Error
	if err != nil {
		return err
	}
	return nil
}

// READ ALL (filtered by user_id)
func (s *PaymentInstallmentsService) GetAllPaymentInstallments(userID string) ([]model.PaymentInstallments, error) {
	var installments []model.PaymentInstallments
	err := s.DB.Where("user_id = ?", userID).Find(&installments).Error
	return installments, err
}

// READ BY ID
func (s *PaymentInstallmentsService) GetPaymentInstallmentsByID(id string) (*model.PaymentInstallments, error) {
	var installment model.PaymentInstallments
	result := s.DB.First(&installment, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &installment, result.Error
}

// UPDATE
func (s *PaymentInstallmentsService) UpdatePaymentInstallments(installment *model.PaymentInstallments) error {
	if installment.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldInstallment model.PaymentInstallments
	if err := s.DB.First(&oldInstallment, "id = ?", installment.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment installment not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if installment.ParentID != "" {
		updateData["parent_id"] = installment.ParentID
	}
	if installment.FormID != "" {
		updateData["form_id"] = installment.FormID
	}
	if installment.InstallmentNumber != 0 {
		updateData["installment_number"] = installment.InstallmentNumber
	}
	if installment.Amount != 0 {
		updateData["amount"] = installment.Amount
	}
	if !installment.DueDate.IsZero() {
		updateData["due_date"] = installment.DueDate
	}
	if installment.Status != "" {
		updateData["status"] = installment.Status
	}
	if installment.PenaltyAmount != 0 {
		updateData["penalty_amount"] = installment.PenaltyAmount
	}
	if installment.PaidAt != nil {
		updateData["paid_at"] = installment.PaidAt
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	return s.DB.Model(&model.PaymentInstallments{}).Where("id = ?", installment.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentInstallmentsService) DeletePaymentInstallments(id string) error {
	if id == "" {
		return fmt.Errorf("ID is required for delete")
	}
	return s.DB.Delete(&model.PaymentInstallments{}, "id = ?", id).Error
}
