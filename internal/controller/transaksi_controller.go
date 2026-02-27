package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransaksiController struct {
	Service *service.TransaksiService
}

func NewTransaksiController(s *service.TransaksiService) *TransaksiController {
	return &TransaksiController{Service: s}
}

// CreateTransaksi
func (tc *TransaksiController) CreateTransaksi(c *gin.Context) {
	var transaksi model.Transaksi
	if err := c.ShouldBindJSON(&transaksi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTransaksi(&transaksi); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    transaksi,
	})
}

// GetAllTransaksi
func (tc *TransaksiController) GetAllTransaksi(c *gin.Context) {
	transaksis, err := tc.Service.GetAllTransaksi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaksis,
	})
}

// GetTransaksiByID
func (tc *TransaksiController) GetTransaksiByID(c *gin.Context) {
	id := c.Param("id")
	transaksi, err := tc.Service.GetTransaksiByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if transaksi == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Transaksi not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaksi,
	})
}

// UpdateTransaksi
func (tc *TransaksiController) UpdateTransaksi(c *gin.Context) {
	id := c.Param("id")
	var updateData model.Transaksi
	updateData.ID = id // Set ID dari param

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTransaksi(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil data terbaru setelah update
	transaksi, err := tc.Service.GetTransaksiByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated transaksi",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaksi,
	})
}

// DeleteTransaksi
func (tc *TransaksiController) DeleteTransaksi(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTransaksi(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Transaksi deleted",
	})
}

// GetTransaksiByUser returns transaksi based on user login
// GetTransaksiByUser returns transaksi based on parent_id
func (tc *TransaksiController) GetTransaksiByUser(c *gin.Context) {
	// Ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Langsung gunakan userID sebagai parent_id
	// Tampilkan semua transaksi dimana parent_id = user yang login
	transaksis, err := tc.Service.GetTransaksiByParentID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    transaksis,
	})
}
