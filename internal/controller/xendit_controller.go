package controller

import (
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type XenditController struct {
	XenditService *service.XenditService
}

func NewXenditController(xs *service.XenditService) *XenditController {
	return &XenditController{XenditService: xs}
}

// CreatePayment - Buat payment untuk transaksi
func (xc *XenditController) CreatePayment(c *gin.Context) {
	// Ambil userID dari token (hanya untuk cek auth, tidak digunakan)
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Parse request
	var req struct {
		TransaksiID string  `json:"transaksi_id" binding:"required"`
		Amount      float64 `json:"amount" binding:"required"`
		Email       string  `json:"email" binding:"required"`
		Description string  `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Validasi amount
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "amount must be greater than 0",
		})
		return
	}

	// Validasi email
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "email is required",
		})
		return
	}

	// Panggil service untuk buat invoice
	result, err := xc.XenditService.CreateInvoice(
		req.TransaksiID,
		req.Amount,
		req.Email,
		req.Description,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Invoice created successfully",
		"data":    result,
	})
}

// Callback - Handle callback dari Xendit
func (xc *XenditController) Callback(c *gin.Context) {
	xc.XenditService.HandleCallback(c)
}
