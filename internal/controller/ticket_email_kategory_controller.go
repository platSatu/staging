package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketEmailKategoryController struct {
	Service *service.TicketEmailKategoryService
}

func NewTicketEmailKategoryController(s *service.TicketEmailKategoryService) *TicketEmailKategoryController {
	return &TicketEmailKategoryController{Service: s}
}

// CreateTicketEmailKategory
func (tc *TicketEmailKategoryController) CreateTicketEmailKategory(c *gin.Context) {
	var ticketEmailKategory model.TicketEmailKategory
	if err := c.ShouldBindJSON(&ticketEmailKategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketEmailKategory(&ticketEmailKategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketEmailKategory,
	})
}

// GetAllTicketEmailKategories
func (tc *TicketEmailKategoryController) GetAllTicketEmailKategories(c *gin.Context) {
	ticketEmailKategories, err := tc.Service.GetAllTicketEmailKategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEmailKategories,
	})
}

// GetTicketEmailKategoryByID
func (tc *TicketEmailKategoryController) GetTicketEmailKategoryByID(c *gin.Context) {
	id := c.Param("id")
	ticketEmailKategory, err := tc.Service.GetTicketEmailKategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketEmailKategory == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket email kategory not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEmailKategory,
	})
}

// UpdateTicketEmailKategory
func (tc *TicketEmailKategoryController) UpdateTicketEmailKategory(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketEmailKategory
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketEmailKategory(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketEmailKategory, err := tc.Service.GetTicketEmailKategoryByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket email kategory",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketEmailKategory,
	})
}

// DeleteTicketEmailKategory
func (tc *TicketEmailKategoryController) DeleteTicketEmailKategory(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketEmailKategory(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket email kategory deleted",
	})
}