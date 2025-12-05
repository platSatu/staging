package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type VoucherController struct {
	Service *service.VoucherService
}

func NewVoucherController(s *service.VoucherService) *VoucherController {
	return &VoucherController{Service: s}
}

// CreateVoucher
func (vc *VoucherController) CreateVoucher(c *gin.Context) {
	var voucher model.Voucher
	if err := c.ShouldBindJSON(&voucher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := vc.Service.CreateVoucher(&voucher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    voucher,
	})
}

// GetAllVouchers
func (vc *VoucherController) GetAllVouchers(c *gin.Context) {
	vouchers, err := vc.Service.GetAllVouchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vouchers,
	})
}

// GetVoucherByID
func (vc *VoucherController) GetVoucherByID(c *gin.Context) {
	id := c.Param("id")
	voucher, err := vc.Service.GetVoucherByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if voucher == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Voucher not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    voucher,
	})
}

// GetVouchersByUser
func (vc *VoucherController) GetVouchersByUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := helper.GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	vouchers, err := vc.Service.GetVouchersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vouchers,
	})
}

// UpdateVoucher
func (vc *VoucherController) UpdateVoucher(c *gin.Context) {
	id := c.Param("id")
	var updateData model.Voucher
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := vc.Service.UpdateVoucher(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	voucher, err := vc.Service.GetVoucherByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated voucher",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    voucher,
	})
}

// DeleteVoucher
func (vc *VoucherController) DeleteVoucher(c *gin.Context) {
	id := c.Param("id")
	if err := vc.Service.DeleteVoucher(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Voucher deleted",
	})
}
