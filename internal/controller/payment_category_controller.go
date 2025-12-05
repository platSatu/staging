package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PaymentCategoryController struct {
	Service *service.PaymentCategoryService
}

func NewPaymentCategoryController(s *service.PaymentCategoryService) *PaymentCategoryController {
	return &PaymentCategoryController{Service: s}
}

// helper untuk konversi date-only
func formatDate(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

// CreatePaymentCategory
func (uc *PaymentCategoryController) CreatePaymentCategory(c *gin.Context) {
	var req struct {
		Name             string     `json:"name" binding:"required"`
		Description      *string    `json:"description"`
		AllowPenalty     bool       `json:"allow_penalty"`
		AllowInstallment bool       `json:"allow_installment"`
		Status           string     `json:"status" binding:"omitempty,oneof=active inactive"`
		DueDateDefault   *time.Time `json:"due_date_default"`
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

	category := &model.PaymentCategory{
		UserID:           userID.(string),
		Name:             req.Name,
		Description:      req.Description,
		AllowPenalty:     req.AllowPenalty,
		AllowInstallment: req.AllowInstallment,
		Status:           req.Status,
		DueDateDefault:   req.DueDateDefault,
	}

	if err := uc.Service.CreatePaymentCategory(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	resp := map[string]interface{}{
		"id":                category.ID,
		"user_id":           category.UserID,
		"name":              category.Name,
		"description":       category.Description,
		"allow_penalty":     category.AllowPenalty,
		"allow_installment": category.AllowInstallment,
		"status":            category.Status,
		"due_date_default":  formatDate(category.DueDateDefault),
		"created_at":        category.CreatedAt,
		"updated_at":        category.UpdatedAt,
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": resp})
}

// GetAllPaymentCategories
func (uc *PaymentCategoryController) GetAllPaymentCategories(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	categories, err := uc.Service.GetAllPaymentCategories(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	resp := []map[string]interface{}{}
	for _, cat := range categories {
		resp = append(resp, map[string]interface{}{
			"id":                cat.ID,
			"user_id":           cat.UserID,
			"name":              cat.Name,
			"description":       cat.Description,
			"allow_penalty":     cat.AllowPenalty,
			"allow_installment": cat.AllowInstallment,
			"status":            cat.Status,
			"due_date_default":  formatDate(cat.DueDateDefault),
			"created_at":        cat.CreatedAt,
			"updated_at":        cat.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// GetPaymentCategoryByID
func (uc *PaymentCategoryController) GetPaymentCategoryByID(c *gin.Context) {
	id := c.Param("id")

	category, err := uc.Service.GetPaymentCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment category not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || category.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	resp := map[string]interface{}{
		"id":                category.ID,
		"user_id":           category.UserID,
		"name":              category.Name,
		"description":       category.Description,
		"allow_penalty":     category.AllowPenalty,
		"allow_installment": category.AllowInstallment,
		"status":            category.Status,
		"due_date_default":  formatDate(category.DueDateDefault),
		"created_at":        category.CreatedAt,
		"updated_at":        category.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// UpdatePaymentCategory
func (uc *PaymentCategoryController) UpdatePaymentCategory(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name             string     `json:"name"`
		Description      *string    `json:"description"`
		AllowPenalty     *bool      `json:"allow_penalty"`
		AllowInstallment *bool      `json:"allow_installment"`
		Status           string     `json:"status" binding:"omitempty,oneof=active inactive"`
		DueDateDefault   *time.Time `json:"due_date_default"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	category, err := uc.Service.GetPaymentCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment category not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || category.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	updateData := &model.PaymentCategory{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		Status:         req.Status,
		DueDateDefault: req.DueDateDefault,
	}

	if req.AllowPenalty != nil {
		updateData.AllowPenalty = *req.AllowPenalty
	} else {
		updateData.AllowPenalty = category.AllowPenalty
	}

	if req.AllowInstallment != nil {
		updateData.AllowInstallment = *req.AllowInstallment
	} else {
		updateData.AllowInstallment = category.AllowInstallment
	}

	if err := uc.Service.UpdatePaymentCategory(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	updatedCategory, _ := uc.Service.GetPaymentCategoryByID(id)

	resp := map[string]interface{}{
		"id":                updatedCategory.ID,
		"user_id":           updatedCategory.UserID,
		"name":              updatedCategory.Name,
		"description":       updatedCategory.Description,
		"allow_penalty":     updatedCategory.AllowPenalty,
		"allow_installment": updatedCategory.AllowInstallment,
		"status":            updatedCategory.Status,
		"due_date_default":  formatDate(updatedCategory.DueDateDefault),
		"created_at":        updatedCategory.CreatedAt,
		"updated_at":        updatedCategory.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": resp})
}

// DeletePaymentCategory
func (uc *PaymentCategoryController) DeletePaymentCategory(c *gin.Context) {
	id := c.Param("id")

	category, err := uc.Service.GetPaymentCategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment category not found"})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists || category.UserID != userIDVal.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := uc.Service.DeletePaymentCategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted successfully"})
}
