package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentFormController struct {
	Service *service.PaymentFormService
}

func NewPaymentFormController(s *service.PaymentFormService) *PaymentFormController {
	return &PaymentFormController{Service: s}
}

// CreatePaymentForm
func (uc *PaymentFormController) CreatePaymentForm(c *gin.Context) {
	var req struct {
		CategoryID        string  `json:"category_id" binding:"required"`
		Name              string  `json:"name" binding:"required"`
		Description       *string `json:"description"`
		BaseAmount        float64 `json:"base_amount" binding:"required,min=0"`
		EnableInstallment bool    `json:"enable_installment"`
		EnablePenalty     bool    `json:"enable_penalty"`
		Status            string  `json:"status" binding:"required,oneof=active inactive"`
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

	form := &model.PaymentForm{
		UserID:            userID.(string),
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       req.Description,
		BaseAmount:        req.BaseAmount,
		EnableInstallment: req.EnableInstallment,
		EnablePenalty:     req.EnablePenalty,
		Status:            req.Status,
	}

	if err := uc.Service.CreatePaymentForm(form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": form})
}

// GetAllPaymentForms
func (uc *PaymentFormController) GetAllPaymentForms(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	forms, err := uc.Service.GetAllPaymentForms(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": forms})
}

// GetPaymentFormByID
func (uc *PaymentFormController) GetPaymentFormByID(c *gin.Context) {
	id := c.Param("id")

	form, err := uc.Service.GetPaymentFormByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if form == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment form not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || form.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": form})
}

// UpdatePaymentForm
func (uc *PaymentFormController) UpdatePaymentForm(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		CategoryID        string  `json:"category_id"`
		Name              string  `json:"name"`
		Description       *string `json:"description"`
		BaseAmount        float64 `json:"base_amount"`
		EnableInstallment bool    `json:"enable_installment"`
		EnablePenalty     bool    `json:"enable_penalty"`
		Status            string  `json:"status" binding:"omitempty,oneof=active inactive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	form, err := uc.Service.GetPaymentFormByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if form == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment form not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || form.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	updateForm := &model.PaymentForm{
		ID:                id,
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       req.Description,
		BaseAmount:        req.BaseAmount,
		EnableInstallment: req.EnableInstallment,
		EnablePenalty:     req.EnablePenalty,
		Status:            req.Status,
	}

	if err := uc.Service.UpdatePaymentForm(updateForm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	updatedForm, err := uc.Service.GetPaymentFormByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment form"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedForm})
}

// DeletePaymentForm
func (uc *PaymentFormController) DeletePaymentForm(c *gin.Context) {
	id := c.Param("id")

	form, err := uc.Service.GetPaymentFormByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if form == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment form not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || form.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := uc.Service.DeletePaymentForm(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment form deleted"})
}
