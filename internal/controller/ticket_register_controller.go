package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TicketRegisterController struct {
	Service *service.TicketRegisterService
}

func NewTicketRegisterController(s *service.TicketRegisterService) *TicketRegisterController {
	return &TicketRegisterController{Service: s}
}

func (tc *TicketRegisterController) CreatePublicTicketRegister(c *gin.Context) {
	var ticketRegister model.TicketRegister

	// Bind JSON
	if err := c.ShouldBindJSON(&ticketRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"fields": map[string]string{
				"general": "Request tidak valid",
			},
		})
		return
	}

	// Validasi & simpan via service
	errorsMap := tc.Service.CreatePublicTicketRegister(&ticketRegister)
	if errorsMap != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"fields":  errorsMap,
		})
		return
	}

	// Build relative purchase link
	purchaseURL := "/purchase?token=" + ticketRegister.PurchaseToken

	// Response sukses
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registrasi berhasil",
		"data": gin.H{
			"purchase_url": purchaseURL,
			"expired_at":   ticketRegister.TokenExpiredAt,
		},
	})
}

// ValidatePurchaseToken memeriksa token pembelian
func (tc *TicketRegisterController) ValidatePurchaseToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Token wajib diisi",
		})
		return
	}

	// Panggil service untuk validasi
	tr, err := tc.Service.ValidatePurchaseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// Pastikan token belum expired
	if tr.TokenExpiredAt != nil && time.Now().After(*tr.TokenExpiredAt) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Token sudah kadaluarsa",
		})
		return
	}

	// Token valid
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"id":        tr.ID,
			"email":     tr.Email,
			"full_name": tr.FullName,
		},
	})
}

// CreateTicketRegister
func (tc *TicketRegisterController) CreateTicketRegister(c *gin.Context) {
	var ticketRegister model.TicketRegister
	if err := c.ShouldBindJSON(&ticketRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketRegister(&ticketRegister); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketRegister,
	})
}

// GetAllTicketRegisters
func (tc *TicketRegisterController) GetAllTicketRegisters(c *gin.Context) {
	ticketRegisters, err := tc.Service.GetAllTicketRegisters()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketRegisters,
	})
}

// GetTicketRegisterByID
func (tc *TicketRegisterController) GetTicketRegisterByID(c *gin.Context) {
	id := c.Param("id")
	ticketRegister, err := tc.Service.GetTicketRegisterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketRegister == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket register not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketRegister,
	})
}

// UpdateTicketRegister
func (tc *TicketRegisterController) UpdateTicketRegister(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketRegister
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketRegister(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketRegister, err := tc.Service.GetTicketRegisterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket register",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketRegister,
	})
}

// DeleteTicketRegister
func (tc *TicketRegisterController) DeleteTicketRegister(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketRegister(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket register deleted",
	})
}
