package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentInvoiceService struct {
	DB *gorm.DB
}

func NewPaymentInvoiceService(db *gorm.DB) *PaymentInvoiceService {
	return &PaymentInvoiceService{DB: db}
}

// CREATE
func (s *PaymentInvoiceService) CreatePaymentInvoice(invoice *model.PaymentInvoice) error {
	// Pastikan UserID ada
	if invoice.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Generate ID jika kosong
	if invoice.ID == "" {
		invoice.ID = uuid.New().String()
	}

	// Set default PaymentStatus jika kosong
	if invoice.PaymentStatus == "" {
		invoice.PaymentStatus = "unpaid"
	}

	// Set default AmountPaid & OutstandingAmount jika belum diisi
	if invoice.AmountPaid == 0 {
		invoice.AmountPaid = 0.00
	}
	if invoice.OutstandingAmount == 0 {
		invoice.OutstandingAmount = invoice.Amount
	}

	// Set default EnableInstallment & EnablePenalty jika belum diisi
	// (Gorm default boolean sudah false, tapi ini memastikan)
	// invoice.EnableInstallment = invoice.EnableInstallment
	// invoice.EnablePenalty = invoice.EnablePenalty

	// Simpan ke DB
	err := s.DB.Create(invoice).Error
	if err != nil {
		return err
	}

	return nil
}

// READ ALL (filtered by user_id)
func (s *PaymentInvoiceService) GetAllPaymentInvoices(userID string) ([]model.PaymentInvoice, error) {
	var invoices []model.PaymentInvoice
	err := s.DB.Where("user_id = ?", userID).Find(&invoices).Error
	return invoices, err
}

// READ BY ID
func (s *PaymentInvoiceService) GetPaymentInvoiceByID(id string) (*model.PaymentInvoice, error) {
	var invoice model.PaymentInvoice
	result := s.DB.First(&invoice, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &invoice, result.Error
}

// UPDATE
func (s *PaymentInvoiceService) UpdatePaymentInvoice(invoice *model.PaymentInvoice) error {
	if invoice.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldInvoice model.PaymentInvoice
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
	if invoice.CategoryID != "" {
		updateData["category_id"] = invoice.CategoryID
	}
	if invoice.FormID != "" {
		updateData["form_id"] = invoice.FormID
	}
	if invoice.PenaltyID != nil {
		updateData["penalty_id"] = invoice.PenaltyID
	}
	if invoice.InstallmentID != nil {
		updateData["installment_id"] = invoice.InstallmentID
	}
	if invoice.Name != "" {
		updateData["name"] = invoice.Name
	}
	if invoice.Description != nil {
		updateData["description"] = invoice.Description
	}
	if invoice.Amount != 0 {
		updateData["amount"] = invoice.Amount
	}
	if invoice.AmountPaid != 0 {
		updateData["amount_paid"] = invoice.AmountPaid
	}
	if invoice.OutstandingAmount != 0 {
		updateData["outstanding_amount"] = invoice.OutstandingAmount
	}
	if !invoice.DueDate.IsZero() {
		updateData["due_date"] = invoice.DueDate
	}
	updateData["enable_installment"] = invoice.EnableInstallment
	updateData["enable_penalty"] = invoice.EnablePenalty
	if invoice.PaymentStatus != "" {
		updateData["payment_status"] = invoice.PaymentStatus
	}
	if invoice.PaymentMethod != nil {
		updateData["payment_method"] = invoice.PaymentMethod
	}
	if invoice.PaymentDate != nil {
		updateData["payment_date"] = invoice.PaymentDate
	}
	if invoice.OrderID != nil {
		updateData["order_id"] = invoice.OrderID
	}
	if invoice.Notes != nil {
		updateData["notes"] = invoice.Notes
	}

	if len(updateData) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	return s.DB.Model(&model.PaymentInvoice{}).Where("id = ?", invoice.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentInvoiceService) DeletePaymentInvoice(id string) error {
	if id == "" {
		return fmt.Errorf("ID is required for delete")
	}
	return s.DB.Delete(&model.PaymentInvoice{}, "id = ?", id).Error
}
