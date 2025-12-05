package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentInvoicesService struct {
	DB *gorm.DB
}

func NewPaymentInvoicesService(db *gorm.DB) *PaymentInvoicesService {
	return &PaymentInvoicesService{DB: db}
}

// CREATE
func (s *PaymentInvoicesService) CreatePaymentInvoices(invoice *model.PaymentInvoices) error {
	if invoice.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if invoice.ID == "" {
		invoice.ID = uuid.New().String()
	}

	err := s.DB.Create(invoice).Error
	if err != nil {
		return err
	}
	return nil
}

// READ ALL (filtered by user_id)
func (s *PaymentInvoicesService) GetAllPaymentInvoices(userID string) ([]model.PaymentInvoices, error) {
	var invoices []model.PaymentInvoices
	err := s.DB.Where("user_id = ?", userID).Find(&invoices).Error
	return invoices, err
}

// READ BY ID
func (s *PaymentInvoicesService) GetPaymentInvoicesByID(id string) (*model.PaymentInvoices, error) {
	var invoice model.PaymentInvoices
	result := s.DB.First(&invoice, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &invoice, result.Error
}

// UPDATE
func (s *PaymentInvoicesService) UpdatePaymentInvoices(invoice *model.PaymentInvoices) error {
	if invoice.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldInvoice model.PaymentInvoices
	if err := s.DB.First(&oldInvoice, "id = ?", invoice.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment invoice not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if invoice.ParentID != "" {
		updateData["parent_id"] = invoice.ParentID
	}
	if invoice.CategoryID != "" { // Diubah dari != 0 ke != "" karena sekarang string
		updateData["category_id"] = invoice.CategoryID
	}
	if invoice.Amount != 0 {
		updateData["amount"] = invoice.Amount
	}
	if !invoice.DueDate.IsZero() {
		updateData["due_date"] = invoice.DueDate
	}
	if invoice.Status != "" {
		updateData["status"] = invoice.Status
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.PaymentInvoices{}).Where("id = ?", invoice.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentInvoicesService) DeletePaymentInvoices(id string) error {
	return s.DB.Delete(&model.PaymentInvoices{}, "id = ?", id).Error
}
