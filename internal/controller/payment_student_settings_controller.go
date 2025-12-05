package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentInvoiceController struct {
	Service *service.PaymentInvoiceService
}

func NewPaymentInvoiceController(s *service.PaymentInvoiceService) *PaymentInvoiceController {
	return &PaymentInvoiceController{Service: s}
}

// CreatePaymentInvoice
func (uc *PaymentInvoiceController) CreatePaymentInvoice(c *gin.Context) {
	var req struct {
		ParentID          string  `json:"parent_id" binding:"required"`
		CategoryID        string  `json:"category_id" binding:"required"`
		FormID            string  `json:"form_id" binding:"required"`
		PenaltyID         *string `json:"penalty_id"`
		InstallmentID     *string `json:"installment_id"`
		Name              string  `json:"name" binding:"required"`
		Description       *string `json:"description"`
		Amount            float64 `json:"amount" binding:"required"`
		AmountPaid        float64 `json:"amount_paid"`
		OutstandingAmount float64 `json:"outstanding_amount"`
		DueDate           string  `json:"due_date" binding:"required"`
		EnableInstallment bool    `json:"enable_installment"`
		EnablePenalty     bool    `json:"enable_penalty"`
		PaymentStatus     string  `json:"payment_status" binding:"omitempty,oneof=unpaid partial paid"`
		PaymentMethod     *string `json:"payment_method"`
		PaymentDate       *string `json:"payment_date"`
		OrderID           *string `json:"order_id"`
		Notes             *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	invoice := &model.PaymentInvoice{
		UserID:            userID.(string),
		ParentID:          req.ParentID,
		CategoryID:        req.CategoryID,
		FormID:            req.FormID,
		PenaltyID:         req.PenaltyID,
		InstallmentID:     req.InstallmentID,
		Name:              req.Name,
		Description:       req.Description,
		Amount:            req.Amount,
		AmountPaid:        req.AmountPaid,
		OutstandingAmount: req.OutstandingAmount,
		DueDate:           parseDate(req.DueDate),
		EnableInstallment: req.EnableInstallment,
		EnablePenalty:     req.EnablePenalty,
		PaymentStatus:     req.PaymentStatus,
		PaymentMethod:     req.PaymentMethod,
		PaymentDate:       parseDateTime(req.PaymentDate),
		OrderID:           req.OrderID,
		Notes:             req.Notes,
	}

	if err := uc.Service.CreatePaymentInvoice(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": invoice})
}

// GetAllPaymentInvoices
func (uc *PaymentInvoiceController) GetAllPaymentInvoices(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	invoices, err := uc.Service.GetAllPaymentInvoices(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": invoices})
}

// GetPaymentInvoiceByID
func (uc *PaymentInvoiceController) GetPaymentInvoiceByID(c *gin.Context) {
	id := c.Param("id")
	invoice, err := uc.Service.GetPaymentInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if invoice == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment invoice not found"})
		return
	}

	// Pastikan invoice milik user yang login
	userIDVal, exists := c.Get("userID")
	if !exists || invoice.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": invoice})
}

// UpdatePaymentInvoice
func (uc *PaymentInvoiceController) UpdatePaymentInvoice(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ParentID          string  `json:"parent_id"`
		CategoryID        string  `json:"category_id"`
		FormID            string  `json:"form_id"`
		PenaltyID         *string `json:"penalty_id"`
		InstallmentID     *string `json:"installment_id"`
		Name              string  `json:"name"`
		Description       *string `json:"description"`
		Amount            float64 `json:"amount"`
		AmountPaid        float64 `json:"amount_paid"`
		OutstandingAmount float64 `json:"outstanding_amount"`
		DueDate           string  `json:"due_date"`
		EnableInstallment bool    `json:"enable_installment"`
		EnablePenalty     bool    `json:"enable_penalty"`
		PaymentStatus     string  `json:"payment_status" binding:"omitempty,oneof=unpaid partial paid"`
		PaymentMethod     *string `json:"payment_method"`
		PaymentDate       *string `json:"payment_date"`
		OrderID           *string `json:"order_id"`
		Notes             *string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	invoice, err := uc.Service.GetPaymentInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if invoice == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment invoice not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || invoice.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	updateInvoice := &model.PaymentInvoice{
		ID:                id,
		ParentID:          req.ParentID,
		CategoryID:        req.CategoryID,
		FormID:            req.FormID,
		PenaltyID:         req.PenaltyID,
		InstallmentID:     req.InstallmentID,
		Name:              req.Name,
		Description:       req.Description,
		Amount:            req.Amount,
		AmountPaid:        req.AmountPaid,
		OutstandingAmount: req.OutstandingAmount,
		DueDate:           parseDate(req.DueDate),
		EnableInstallment: req.EnableInstallment,
		EnablePenalty:     req.EnablePenalty,
		PaymentStatus:     req.PaymentStatus,
		PaymentMethod:     req.PaymentMethod,
		PaymentDate:       parseDateTime(req.PaymentDate),
		OrderID:           req.OrderID,
		Notes:             req.Notes,
	}

	if err := uc.Service.UpdatePaymentInvoice(updateInvoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	updatedInvoice, err := uc.Service.GetPaymentInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedInvoice})
}

// DeletePaymentInvoice
func (uc *PaymentInvoiceController) DeletePaymentInvoice(c *gin.Context) {
	id := c.Param("id")

	invoice, err := uc.Service.GetPaymentInvoiceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if invoice == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment invoice not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || invoice.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := uc.Service.DeletePaymentInvoice(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment invoice deleted"})
}

// Helper functions
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Time{} // return zero time if parsing fails
	}
	return t
}

func parseDateTime(dateTimeStr *string) *time.Time {
	if dateTimeStr == nil || *dateTimeStr == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02T15:04:05", *dateTimeStr)
	if err != nil {
		return nil // return nil if parsing fails
	}
	return &t
}
