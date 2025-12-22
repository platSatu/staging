package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketKategoriController struct {
	Service *service.TicketKategoriService
}

func NewTicketKategoriController(s *service.TicketKategoriService) *TicketKategoriController {
	return &TicketKategoriController{Service: s}
}

func (tc *TicketKategoriController) GetAllTicketKategorisPublic(c *gin.Context) {
	ticketKategoris, err := tc.Service.GetAllTicketKategoris()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketKategoris,
	})
}

// CreateTicketKategori
func (tc *TicketKategoriController) CreateTicketKategori(c *gin.Context) {
	var ticketKategori model.TicketKategori
	if err := c.ShouldBindJSON(&ticketKategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketKategori(&ticketKategori); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketKategori,
	})
}

// GetAllTicketKategoris
func (tc *TicketKategoriController) GetAllTicketKategoris(c *gin.Context) {
	ticketKategoris, err := tc.Service.GetAllTicketKategoris()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketKategoris,
	})
}

// GetTicketKategoriByID
func (tc *TicketKategoriController) GetTicketKategoriByID(c *gin.Context) {
	id := c.Param("id")
	ticketKategori, err := tc.Service.GetTicketKategoriByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketKategori == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket kategori not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketKategori,
	})
}

// UpdateTicketKategori
func (tc *TicketKategoriController) UpdateTicketKategori(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketKategori
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketKategori(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketKategori, err := tc.Service.GetTicketKategoriByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket kategori",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketKategori,
	})
}

// DeleteTicketKategori
func (tc *TicketKategoriController) DeleteTicketKategori(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketKategori(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket kategori deleted",
	})
}
