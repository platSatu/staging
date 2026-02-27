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
// CreateDeposit
func (dc *DepositController) CreateDeposit(c *gin.Context) {
	var req struct {
		Kredit          float64 `json:"kredit"`
		Debit           float64 `json:"debit"`
		TransaksiStatus string  `json:"transaksi_status"`
		TransaksiMethod string  `json:"transaksi_method"`
		Keterangan      string  `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	deposit := &model.Deposit{
		UserID:          userID.(string),
		Kredit:          req.Kredit,
		Debit:           req.Debit,
		TransaksiStatus: req.TransaksiStatus,
		TransaksiMethod: req.TransaksiMethod,
		Keterangan:      req.Keterangan,
	}

	if err := dc.Service.CreateDeposit(deposit); err != nil {
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

// GetAllDeposits (hanya milik user login)
func (dc *DepositController) GetAllDeposits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	deposits, err := dc.Service.GetDepositsByUserID(userID.(string))
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

	// hanya pemilik deposit boleh melihat
	userID, _ := c.Get("userID")
	if deposit.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
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

	var req struct {
		Kredit          float64 `json:"kredit"`
		Debit           float64 `json:"debit"`
		TransaksiStatus string  `json:"transaksi_status"`
		TransaksiMethod string  `json:"transaksi_method"`
		Keterangan      string  `json:"keterangan"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan deposit itu milik user login
	existing, err := dc.Service.GetDepositByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Deposit not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	updateData := &model.Deposit{
		ID:              id,
		Kredit:          req.Kredit,
		Debit:           req.Debit,
		TransaksiStatus: req.TransaksiStatus,
		TransaksiMethod: req.TransaksiMethod,
		Keterangan:      req.Keterangan,
	}

	if err := dc.Service.UpdateDeposit(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := dc.Service.GetDepositByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteDeposit
func (dc *DepositController) DeleteDeposit(c *gin.Context) {
	id := c.Param("id")

	deposit, err := dc.Service.GetDepositByID(id)
	if err != nil || deposit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Deposit not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if deposit.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

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
