package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentUserController struct {
	Service *service.PaymentUserService
}

func NewPaymentUserController(s *service.PaymentUserService) *PaymentUserController {
	return &PaymentUserController{Service: s}
}

// CreatePaymentUser
func (uc *PaymentUserController) CreatePaymentUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	parentID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	user := &model.PaymentUser{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		ParentID: parentID.(string), // parent = user yang login
	}

	if err := uc.Service.CreatePaymentUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": user})
}

// GetAllPaymentUsers
func (uc *PaymentUserController) GetAllPaymentUsers(c *gin.Context) {
	// Ambil userID dari context (set di AuthMiddleware)
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	userID := userIDVal.(string)

	users, err := uc.Service.GetAllPaymentUsers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

// GetPaymentUserByID
func (uc *PaymentUserController) GetPaymentUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.Service.GetPaymentUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Payment user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": user})
}

// UpdatePaymentUser
func (uc *PaymentUserController) UpdatePaymentUser(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	user := &model.PaymentUser{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := uc.Service.UpdatePaymentUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil data terbaru setelah update
	updatedUser, err := uc.Service.GetPaymentUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated payment user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updatedUser})
}

// DeletePaymentUser
func (uc *PaymentUserController) DeletePaymentUser(c *gin.Context) {
	id := c.Param("id")
	if err := uc.Service.DeletePaymentUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Payment user deleted"})
}
