package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentInvoicesController struct {
	Service *service.PaymentInvoicesService
}

func NewPaymentInvoicesController(s *service.PaymentInvoicesService) *PaymentInvoicesController {
	return &PaymentInvoicesController{Service: s}
}

// CreatePaymentInvoices
func (uc *PaymentInvoicesController) CreatePaymentInvoices(c *gin.Context) {
	var req struct {
		ParentID   string  `json:"parent_id" binding:"required"`
		CategoryID string  `json:"category_id" binding:"required"` // Diubah dari int ke string
		Amount     float64 `json:"amount" binding:"required,min=0"`
		DueDate    string  `json:"due_date" binding:"required"`
		Status     string  `json:"status" binding:"required,oneof=unpaid partial paid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid due_date format (use YYYY-MM-DD)"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	invoice := &model.PaymentInvoices{
		UserID:     userID.(string),
		ParentID:   req.ParentID,
		CategoryID: req.CategoryID, // Sekarang string
		Amount:     req.Amount,
		DueDate:    dueDate,
		Status:     req.Status,
	}

	if err := uc.Service.CreatePaymentInvoices(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": invoice})
}

// GetAllPaymentInvoices
func (uc *PaymentInvoicesController) GetAllPaymentInvoices(c *gin.Context) {
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

// GetPaymentInvoicesByID
func (uc *PaymentInvoicesController) GetPaymentInvoicesByID(c *gin.Context) {
	id := c.Param("id")
	invoice, err := uc.Service.GetPaymentInvoicesByID(id)
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

// UpdatePaymentInvoices
func (uc *PaymentInvoicesController) UpdatePaymentInvoices(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ParentID   string  `json:"parent_id"`
		CategoryID string  `json:"category_id"` // Diubah dari int ke string
		Amount     float64 `json:"amount"`
		DueDate    string  `json:"due_date"`
		Status     string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	var dueDate *time.Time
	if req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid due_date format (use YYYY-MM-DD)"})
			return
		}
		dueDate = &parsed
	}

	invoice := &model.PaymentInvoices{
		ID:         id,
		ParentID:   req.ParentID,
		CategoryID: req.CategoryID, // Sekarang string
		Amount:     req.Amount,
		Status:     req.Status,
	}
	if dueDate != nil {
		invoice.DueDate = *dueDate
	}

	if err := uc.Service.UpdatePaymentInvoices(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil data terbaru setelah update
	updatedInvoice, err := uc.Service.GetPaymentInvoicesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment invoice"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedInvoice})
}

// DeletePaymentInvoices
func (uc *PaymentInvoicesController) DeletePaymentInvoices(c *gin.Context) {
	id := c.Param("id")
	if err := uc.Service.DeletePaymentInvoices(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment invoice deleted"})
}
