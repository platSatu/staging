package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentPenaltySettingsController struct {
	Service *service.PaymentPenaltySettingsService
}

func NewPaymentPenaltySettingsController(s *service.PaymentPenaltySettingsService) *PaymentPenaltySettingsController {
	return &PaymentPenaltySettingsController{Service: s}
}

// CreatePaymentPenaltySettings
func (uc *PaymentPenaltySettingsController) CreatePaymentPenaltySettings(c *gin.Context) {
	var req struct {
		CategoryID       string   `json:"category_id" binding:"required"`
		PenaltyType      string   `json:"penalty_type" binding:"required,oneof=percent flat"`
		PercentValue     *float64 `json:"percent_value"`
		FlatValue        *float64 `json:"flat_value"`
		MaxPenaltyAmount float64  `json:"max_penalty_amount" binding:"required,min=0"`
		ApplyOn          string   `json:"apply_on" binding:"required,oneof=invoice installment both"`
		Active           string   `json:"active" binding:"omitempty,oneof=active inactive"` // Diubah ke string
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

	setting := &model.PaymentPenaltySettings{
		UserID:           userID.(string),
		CategoryID:       req.CategoryID,
		PenaltyType:      req.PenaltyType,
		PercentValue:     req.PercentValue,
		FlatValue:        req.FlatValue,
		MaxPenaltyAmount: req.MaxPenaltyAmount,
		ApplyOn:          req.ApplyOn,
		Active:           req.Active,
	}
	if setting.Active == "" {
		setting.Active = "active" // Default
	}

	if err := uc.Service.CreatePaymentPenaltySettings(setting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": setting})
}

// GetAllPaymentPenaltySettings
func (uc *PaymentPenaltySettingsController) GetAllPaymentPenaltySettings(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	settings, err := uc.Service.GetAllPaymentPenaltySettings(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": settings})
}

// GetPaymentPenaltySettingsByID
func (uc *PaymentPenaltySettingsController) GetPaymentPenaltySettingsByID(c *gin.Context) {
	id := c.Param("id")
	setting, err := uc.Service.GetPaymentPenaltySettingsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if setting == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment penalty setting not found"})
		return
	}

	// Pastikan setting milik user yang login
	userIDVal, exists := c.Get("userID")
	if !exists || setting.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": setting})
}

// UpdatePaymentPenaltySettings
func (uc *PaymentPenaltySettingsController) UpdatePaymentPenaltySettings(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		CategoryID       string   `json:"category_id"`
		PenaltyType      string   `json:"penalty_type" binding:"omitempty,oneof=percent flat"`
		PercentValue     *float64 `json:"percent_value"`
		FlatValue        *float64 `json:"flat_value"`
		MaxPenaltyAmount *float64 `json:"max_penalty_amount"` // Pointer untuk deteksi jika tidak dikirim
		ApplyOn          string   `json:"apply_on" binding:"omitempty,oneof=invoice installment both"`
		Active           string   `json:"active" binding:"omitempty,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}
	userID := userIDVal.(string)

	setting := &model.PaymentPenaltySettings{
		ID:           id,
		CategoryID:   req.CategoryID,
		PenaltyType:  req.PenaltyType,
		PercentValue: req.PercentValue,
		FlatValue:    req.FlatValue,
		ApplyOn:      req.ApplyOn,
		Active:       req.Active,
	}
	if req.MaxPenaltyAmount != nil {
		setting.MaxPenaltyAmount = *req.MaxPenaltyAmount
	}

	if err := uc.Service.UpdatePaymentPenaltySettings(setting, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil data terbaru setelah update
	updatedSetting, err := uc.Service.GetPaymentPenaltySettingsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment penalty setting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedSetting})
}

// DeletePaymentPenaltySettings
func (uc *PaymentPenaltySettingsController) DeletePaymentPenaltySettings(c *gin.Context) {
	id := c.Param("id")
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}
	userID := userIDVal.(string)

	if err := uc.Service.DeletePaymentPenaltySettings(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment penalty setting deleted"})
}
