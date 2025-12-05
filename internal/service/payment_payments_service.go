package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentPaymentsService struct {
	DB *gorm.DB
}

func NewPaymentPaymentsService(db *gorm.DB) *PaymentPaymentsService {
	return &PaymentPaymentsService{DB: db}
}

// CREATE
func (s *PaymentPaymentsService) CreatePaymentPayments(payment *model.PaymentPayments) error {
	if payment.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}

	err := s.DB.Create(payment).Error
	if err != nil {
		return err
	}
	return nil
}

// READ ALL (filtered by user_id)
func (s *PaymentPaymentsService) GetAllPaymentPayments(userID string) ([]model.PaymentPayments, error) {
	var payments []model.PaymentPayments
	err := s.DB.Where("user_id = ?", userID).Find(&payments).Error
	return payments, err
}

// READ BY ID
func (s *PaymentPaymentsService) GetPaymentPaymentsByID(id string) (*model.PaymentPayments, error) {
	var payment model.PaymentPayments
	result := s.DB.First(&payment, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &payment, result.Error
}

// UPDATE
func (s *PaymentPaymentsService) UpdatePaymentPayments(payment *model.PaymentPayments) error {
	if payment.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldPayment model.PaymentPayments
	if err := s.DB.First(&oldPayment, "id = ?", payment.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment payment not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if payment.ParentID != "" {
		updateData["parent_id"] = payment.ParentID
	}
	if payment.InvoiceID != "" {
		updateData["invoice_id"] = payment.InvoiceID
	}
	if payment.InstallmentID != "" {
		updateData["installment_id"] = payment.InstallmentID
	}
	if payment.Amount != 0 {
		updateData["amount"] = payment.Amount
	}
	if payment.PaymentMethod != "" {
		updateData["payment_method"] = payment.PaymentMethod
	}
	if !payment.PaymentDate.IsZero() {
		updateData["payment_date"] = payment.PaymentDate
	}
	if payment.Status != "" {
		updateData["status"] = payment.Status
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.PaymentPayments{}).Where("id = ?", payment.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentPaymentsService) DeletePaymentPayments(id string) error {
	return s.DB.Delete(&model.PaymentPayments{}, "id = ?", id).Error
}
