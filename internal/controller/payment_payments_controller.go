package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentPaymentsController struct {
	Service *service.PaymentPaymentsService
}

func NewPaymentPaymentsController(s *service.PaymentPaymentsService) *PaymentPaymentsController {
	return &PaymentPaymentsController{Service: s}
}

// CreatePaymentPayments
func (uc *PaymentPaymentsController) CreatePaymentPayments(c *gin.Context) {
	var req struct {
		ParentID      string  `json:"parent_id" binding:"required"`
		InvoiceID     string  `json:"invoice_id" binding:"required"`
		InstallmentID string  `json:"installment_id" binding:"required"`
		Amount        float64 `json:"amount" binding:"required,min=0"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		PaymentDate   string  `json:"payment_date" binding:"required"`
		Status        string  `json:"status" binding:"required,oneof=pending success failed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	paymentDate, err := time.Parse("2006-01-02T15:04:05", req.PaymentDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid payment_date format (use YYYY-MM-DDTHH:MM:SS)"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	payment := &model.PaymentPayments{
		UserID:        userID.(string),
		ParentID:      req.ParentID,
		InvoiceID:     req.InvoiceID,
		InstallmentID: req.InstallmentID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentDate:   paymentDate,
		Status:        req.Status,
	}

	if err := uc.Service.CreatePaymentPayments(payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": payment})
}

// GetAllPaymentPayments
func (uc *PaymentPaymentsController) GetAllPaymentPayments(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	payments, err := uc.Service.GetAllPaymentPayments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payments})
}

// GetPaymentPaymentsByID
func (uc *PaymentPaymentsController) GetPaymentPaymentsByID(c *gin.Context) {
	id := c.Param("id")
	payment, err := uc.Service.GetPaymentPaymentsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if payment == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment payment not found"})
		return
	}

	// Pastikan payment milik user yang login
	userIDVal, exists := c.Get("userID")
	if !exists || payment.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": payment})
}

// UpdatePaymentPayments
func (uc *PaymentPaymentsController) UpdatePaymentPayments(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ParentID      string  `json:"parent_id"`
		InvoiceID     string  `json:"invoice_id"`
		InstallmentID string  `json:"installment_id"`
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		PaymentDate   string  `json:"payment_date"`
		Status        string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	var paymentDate *time.Time
	if req.PaymentDate != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", req.PaymentDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid payment_date format (use YYYY-MM-DDTHH:MM:SS)"})
			return
		}
		paymentDate = &parsed
	}

	payment := &model.PaymentPayments{
		ID:            id,
		ParentID:      req.ParentID,
		InvoiceID:     req.InvoiceID,
		InstallmentID: req.InstallmentID,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		Status:        req.Status,
	}
	if paymentDate != nil {
		payment.PaymentDate = *paymentDate
	}

	if err := uc.Service.UpdatePaymentPayments(payment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil data terbaru setelah update
	updatedPayment, err := uc.Service.GetPaymentPaymentsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedPayment})
}

// DeletePaymentPayments
func (uc *PaymentPaymentsController) DeletePaymentPayments(c *gin.Context) {
	id := c.Param("id")
	if err := uc.Service.DeletePaymentPayments(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment payment deleted"})
}
