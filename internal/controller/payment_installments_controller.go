package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentInstallmentsController struct {
	Service *service.PaymentInstallmentsService
}

func NewPaymentInstallmentsController(s *service.PaymentInstallmentsService) *PaymentInstallmentsController {
	return &PaymentInstallmentsController{Service: s}
}

// CreatePaymentInstallments
func (uc *PaymentInstallmentsController) CreatePaymentInstallments(c *gin.Context) {
	var req struct {
		ParentID          string  `json:"parent_id" binding:"required"`
		FormID            string  `json:"form_id" binding:"required"`
		InstallmentNumber int     `json:"installment_number" binding:"required,min=1"`
		Amount            float64 `json:"amount" binding:"required,min=0"`
		DueDate           string  `json:"due_date" binding:"required"`
		Status            string  `json:"status" binding:"required,oneof=unpaid paid late"`
		PenaltyAmount     float64 `json:"penalty_amount"`
		PaidAt            *string `json:"paid_at"`
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

	var paidAt *time.Time
	if req.PaidAt != nil && *req.PaidAt != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", *req.PaidAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid paid_at format (use YYYY-MM-DDTHH:MM:SS)"})
			return
		}
		paidAt = &parsed
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	installment := &model.PaymentInstallments{
		UserID:            userID.(string),
		ParentID:          req.ParentID,
		FormID:            req.FormID,
		InstallmentNumber: req.InstallmentNumber,
		Amount:            req.Amount,
		DueDate:           dueDate,
		Status:            req.Status,
		PenaltyAmount:     req.PenaltyAmount,
		PaidAt:            paidAt,
	}

	if err := uc.Service.CreatePaymentInstallments(installment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": installment})
}

// GetAllPaymentInstallments
func (uc *PaymentInstallmentsController) GetAllPaymentInstallments(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	installments, err := uc.Service.GetAllPaymentInstallments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": installments})
}

// GetPaymentInstallmentsByID
func (uc *PaymentInstallmentsController) GetPaymentInstallmentsByID(c *gin.Context) {
	id := c.Param("id")
	installment, err := uc.Service.GetPaymentInstallmentsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if installment == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment installment not found"})
		return
	}

	// Pastikan installment milik user yang login
	userIDVal, exists := c.Get("userID")
	if !exists || installment.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": installment})
}

// UpdatePaymentInstallments
func (uc *PaymentInstallmentsController) UpdatePaymentInstallments(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		ParentID          string  `json:"parent_id"`
		FormID            string  `json:"form_id"`
		InstallmentNumber int     `json:"installment_number"`
		Amount            float64 `json:"amount"`
		DueDate           string  `json:"due_date"`
		Status            string  `json:"status"`
		PenaltyAmount     float64 `json:"penalty_amount"`
		PaidAt            *string `json:"paid_at"`
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

	var paidAt *time.Time
	if req.PaidAt != nil && *req.PaidAt != "" {
		parsed, err := time.Parse("2006-01-02T15:04:05", *req.PaidAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Invalid paid_at format (use YYYY-MM-DDTHH:MM:SS)"})
			return
		}
		paidAt = &parsed
	}

	installment := &model.PaymentInstallments{
		ID:                id,
		ParentID:          req.ParentID,
		FormID:            req.FormID,
		InstallmentNumber: req.InstallmentNumber,
		Amount:            req.Amount,
		Status:            req.Status,
		PenaltyAmount:     req.PenaltyAmount,
	}
	if dueDate != nil {
		installment.DueDate = *dueDate
	}
	installment.PaidAt = paidAt

	if err := uc.Service.UpdatePaymentInstallments(installment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil data terbaru setelah update
	updatedInstallment, err := uc.Service.GetPaymentInstallmentsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment installment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedInstallment})
}

// DeletePaymentInstallments
func (uc *PaymentInstallmentsController) DeletePaymentInstallments(c *gin.Context) {
	id := c.Param("id")

	installment, err := uc.Service.GetPaymentInstallmentsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if installment == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment installment not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || installment.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := uc.Service.DeletePaymentInstallments(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment installment deleted"})
}
