package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketVoucherController struct {
	Service *service.TicketVoucherService
}

func NewTicketVoucherController(s *service.TicketVoucherService) *TicketVoucherController {
	return &TicketVoucherController{Service: s}
}

// CreateTicketVoucher
func (tc *TicketVoucherController) CreateTicketVoucher(c *gin.Context) {
	var ticketVoucher model.TicketVoucher
	if err := c.ShouldBindJSON(&ticketVoucher); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketVoucher(&ticketVoucher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketVoucher,
	})
}

// GetAllTicketVouchers
func (tc *TicketVoucherController) GetAllTicketVouchers(c *gin.Context) {
	ticketVouchers, err := tc.Service.GetAllTicketVouchers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketVouchers,
	})
}

// GetTicketVoucherByID
func (tc *TicketVoucherController) GetTicketVoucherByID(c *gin.Context) {
	id := c.Param("id")
	ticketVoucher, err := tc.Service.GetTicketVoucherByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketVoucher == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket voucher not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketVoucher,
	})
}

// UpdateTicketVoucher
func (tc *TicketVoucherController) UpdateTicketVoucher(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketVoucher
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketVoucher(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketVoucher, err := tc.Service.GetTicketVoucherByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket voucher",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketVoucher,
	})
}

// DeleteTicketVoucher
func (tc *TicketVoucherController) DeleteTicketVoucher(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketVoucher(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket voucher deleted",
	})
}

func (tc *TicketVoucherController) CheckVoucherAvailability(c *gin.Context) {
	kodeVoucher := c.Param("kodeVoucher")
	if kodeVoucher == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Kode voucher diperlukan",
		})
		return
	}

	voucher, err := tc.Service.CheckVoucherAvailability(kodeVoucher)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(), // Pesan error spesifik dari service
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Voucher tersedia",
		"data":    voucher,
	})
}
