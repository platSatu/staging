package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketMethodPembayaranController struct {
	Service *service.TicketMethodPembayaranService
}

func NewTicketMethodPembayaranController(s *service.TicketMethodPembayaranService) *TicketMethodPembayaranController {
	return &TicketMethodPembayaranController{Service: s}
}

// CreateTicketMethodPembayaran
func (tc *TicketMethodPembayaranController) CreateTicketMethodPembayaran(c *gin.Context) {
	var ticketMethodPembayaran model.TicketMethodPembayaran
	if err := c.ShouldBindJSON(&ticketMethodPembayaran); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketMethodPembayaran(&ticketMethodPembayaran); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketMethodPembayaran,
	})
}

// GetAllTicketMethodPembayarans
func (tc *TicketMethodPembayaranController) GetAllTicketMethodPembayarans(c *gin.Context) {
	ticketMethodPembayarans, err := tc.Service.GetAllTicketMethodPembayarans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketMethodPembayarans,
	})
}

// GetTicketMethodPembayaranByID
func (tc *TicketMethodPembayaranController) GetTicketMethodPembayaranByID(c *gin.Context) {
	id := c.Param("id")
	ticketMethodPembayaran, err := tc.Service.GetTicketMethodPembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketMethodPembayaran == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket method pembayaran not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketMethodPembayaran,
	})
}

// UpdateTicketMethodPembayaran
func (tc *TicketMethodPembayaranController) UpdateTicketMethodPembayaran(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketMethodPembayaran
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketMethodPembayaran(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketMethodPembayaran, err := tc.Service.GetTicketMethodPembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket method pembayaran",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketMethodPembayaran,
	})
}

// DeleteTicketMethodPembayaran
func (tc *TicketMethodPembayaranController) DeleteTicketMethodPembayaran(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketMethodPembayaran(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket method pembayaran deleted",
	})
}