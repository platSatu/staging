package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketQrcodeController struct {
	Service *service.TicketQrcodeService
}

func NewTicketQrcodeController(s *service.TicketQrcodeService) *TicketQrcodeController {
	return &TicketQrcodeController{Service: s}
}

// CreateTicketQrcode
func (tc *TicketQrcodeController) CreateTicketQrcode(c *gin.Context) {
	var ticketQrcode model.TicketQrcode
	if err := c.ShouldBindJSON(&ticketQrcode); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketQrcode(&ticketQrcode); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketQrcode,
	})
}

// GetAllTicketQrcodes
func (tc *TicketQrcodeController) GetAllTicketQrcodes(c *gin.Context) {
	ticketQrcodes, err := tc.Service.GetAllTicketQrcodes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketQrcodes,
	})
}

// GetTicketQrcodeByID
func (tc *TicketQrcodeController) GetTicketQrcodeByID(c *gin.Context) {
	id := c.Param("id")
	ticketQrcode, err := tc.Service.GetTicketQrcodeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketQrcode == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket qrcode not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketQrcode,
	})
}

// UpdateTicketQrcode
func (tc *TicketQrcodeController) UpdateTicketQrcode(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketQrcode
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketQrcode(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketQrcode, err := tc.Service.GetTicketQrcodeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket qrcode",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketQrcode,
	})
}

// DeleteTicketQrcode
func (tc *TicketQrcodeController) DeleteTicketQrcode(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketQrcode(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket qrcode deleted",
	})
}