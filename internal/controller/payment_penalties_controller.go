package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentPenaltiesController struct {
	Service *service.PaymentPenaltiesService
}

func NewPaymentPenaltiesController(s *service.PaymentPenaltiesService) *PaymentPenaltiesController {
	return &PaymentPenaltiesController{Service: s}
}

// CreatePaymentPenalties
func (uc *PaymentPenaltiesController) CreatePaymentPenalties(c *gin.Context) {
	var req struct {
		CategoryID   string   `json:"category_id" binding:"required"`
		Name         string   `json:"name" binding:"required"`
		Description  *string  `json:"description"`
		PenaltyType  string   `json:"penalty_type" binding:"required,oneof=flat percent"`
		FlatValue    *float64 `json:"flat_value"`
		PercentValue *float64 `json:"percent_value"`
		MaxPenalty   float64  `json:"max_penalty" binding:"required,min=0"`
		ApplyOn      string   `json:"apply_on" binding:"required,oneof=invoice installment both"`
		Active       string   `json:"active" binding:"required,oneof=active inactive"`
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

	penalty := &model.PaymentPenalty{
		UserID:       userID.(string),
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Description:  req.Description,
		PenaltyType:  req.PenaltyType,
		FlatValue:    req.FlatValue,
		PercentValue: req.PercentValue,
		MaxPenalty:   req.MaxPenalty,
		ApplyOn:      req.ApplyOn,
		Active:       req.Active,
	}

	if err := uc.Service.CreatePaymentPenalties(penalty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": penalty})
}

// GetAllPaymentPenalties
func (uc *PaymentPenaltiesController) GetAllPaymentPenalties(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	penalties, err := uc.Service.GetAllPaymentPenalties(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": penalties})
}

// GetPaymentPenaltiesByID
func (uc *PaymentPenaltiesController) GetPaymentPenaltiesByID(c *gin.Context) {
	id := c.Param("id")

	penalty, err := uc.Service.GetPaymentPenaltiesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if penalty == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment penalty not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || penalty.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": penalty})
}

// UpdatePaymentPenalties
func (uc *PaymentPenaltiesController) UpdatePaymentPenalties(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CategoryID   string   `json:"category_id"`
		Name         string   `json:"name"`
		Description  *string  `json:"description"`
		PenaltyType  string   `json:"penalty_type" binding:"omitempty,oneof=flat percent"`
		FlatValue    *float64 `json:"flat_value"`
		PercentValue *float64 `json:"percent_value"`
		MaxPenalty   float64  `json:"max_penalty" binding:"omitempty,min=0"`
		ApplyOn      string   `json:"apply_on" binding:"omitempty,oneof=invoice installment both"`
		Active       string   `json:"active" binding:"omitempty,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	penalty, err := uc.Service.GetPaymentPenaltiesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if penalty == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment penalty not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || penalty.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	updatePenalty := &model.PaymentPenalty{
		ID:           id,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Description:  req.Description,
		PenaltyType:  req.PenaltyType,
		FlatValue:    req.FlatValue,
		PercentValue: req.PercentValue,
		MaxPenalty:   req.MaxPenalty,
		ApplyOn:      req.ApplyOn,
		Active:       req.Active,
	}

	if err := uc.Service.UpdatePaymentPenalties(updatePenalty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	updatedPenalty, err := uc.Service.GetPaymentPenaltiesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment penalty"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedPenalty})
}

// DeletePaymentPenalties
func (uc *PaymentPenaltiesController) DeletePaymentPenalties(c *gin.Context) {
	id := c.Param("id")

	penalty, err := uc.Service.GetPaymentPenaltiesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if penalty == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment penalty not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || penalty.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := uc.Service.DeletePaymentPenalties(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment penalty deleted"})
}
