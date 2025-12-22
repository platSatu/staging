package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketJenisQuantityController struct {
	Service *service.TicketJenisQuantityService
}

func NewTicketJenisQuantityController(s *service.TicketJenisQuantityService) *TicketJenisQuantityController {
	return &TicketJenisQuantityController{Service: s}
}

func (tc *TicketJenisQuantityController) GetAllTicketJenisQuantitiesPublic(c *gin.Context) {
	ticketJenisQuantities, err := tc.Service.GetAllTicketJenisQuantities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketJenisQuantities,
	})
}

// CreateTicketJenisQuantity
func (tc *TicketJenisQuantityController) CreateTicketJenisQuantity(c *gin.Context) {
	var ticketJenisQuantity model.TicketJenisQuantity
	if err := c.ShouldBindJSON(&ticketJenisQuantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketJenisQuantity(&ticketJenisQuantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketJenisQuantity,
	})
}

// GetAllTicketJenisQuantities
func (tc *TicketJenisQuantityController) GetAllTicketJenisQuantities(c *gin.Context) {
	ticketJenisQuantities, err := tc.Service.GetAllTicketJenisQuantities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketJenisQuantities,
	})
}

// GetTicketJenisQuantityByID
func (tc *TicketJenisQuantityController) GetTicketJenisQuantityByID(c *gin.Context) {
	id := c.Param("id")
	ticketJenisQuantity, err := tc.Service.GetTicketJenisQuantityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketJenisQuantity == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket jenis quantity not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketJenisQuantity,
	})
}

// UpdateTicketJenisQuantity
func (tc *TicketJenisQuantityController) UpdateTicketJenisQuantity(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketJenisQuantity
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketJenisQuantity(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketJenisQuantity, err := tc.Service.GetTicketJenisQuantityByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket jenis quantity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketJenisQuantity,
	})
}

// DeleteTicketJenisQuantity
func (tc *TicketJenisQuantityController) DeleteTicketJenisQuantity(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketJenisQuantity(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket jenis quantity deleted",
	})
}
