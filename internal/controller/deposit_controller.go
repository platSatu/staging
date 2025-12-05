package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DepositController struct {
	Service *service.DepositService
}

func NewDepositController(s *service.DepositService) *DepositController {
	return &DepositController{Service: s}
}

// CreateDeposit
func (dc *DepositController) CreateDeposit(c *gin.Context) {
	var deposit model.Deposit
	if err := c.ShouldBindJSON(&deposit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := dc.Service.CreateDeposit(&deposit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    deposit,
	})
}

// GetAllDeposits
func (dc *DepositController) GetAllDeposits(c *gin.Context) {
	deposits, err := dc.Service.GetAllDeposits()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deposits,
	})
}

// GetDepositByID
func (dc *DepositController) GetDepositByID(c *gin.Context) {
	id := c.Param("id")
	deposit, err := dc.Service.GetDepositByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if deposit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Deposit not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deposit,
	})
}

// GetDepositsByUser
func (dc *DepositController) GetDepositsByUser(c *gin.Context) {
	userID := c.Query("user_id") // Asumsi via query param, atau ganti dengan auth jika perlu
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "user_id is required",
		})
		return
	}

	deposits, err := dc.Service.GetDepositsByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deposits,
	})
}

// UpdateDeposit
func (dc *DepositController) UpdateDeposit(c *gin.Context) {
	id := c.Param("id")
	var updateData model.Deposit
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := dc.Service.UpdateDeposit(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	deposit, err := dc.Service.GetDepositByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated deposit",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    deposit,
	})
}

// DeleteDeposit
func (dc *DepositController) DeleteDeposit(c *gin.Context) {
	id := c.Param("id")
	if err := dc.Service.DeleteDeposit(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Deposit deleted",
	})
}
