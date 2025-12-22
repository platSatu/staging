package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketTemplateController struct {
	Service *service.TicketTemplateService
}

func NewTicketTemplateController(s *service.TicketTemplateService) *TicketTemplateController {
	return &TicketTemplateController{Service: s}
}

// CreateTicketTemplate
func (tc *TicketTemplateController) CreateTicketTemplate(c *gin.Context) {
	var ticketTemplate model.TicketTemplate
	if err := c.ShouldBindJSON(&ticketTemplate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketTemplate(&ticketTemplate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketTemplate,
	})
}

// GetAllTicketTemplates
func (tc *TicketTemplateController) GetAllTicketTemplates(c *gin.Context) {
	ticketTemplates, err := tc.Service.GetAllTicketTemplates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketTemplates,
	})
}

// GetTicketTemplateByID
func (tc *TicketTemplateController) GetTicketTemplateByID(c *gin.Context) {
	id := c.Param("id")
	ticketTemplate, err := tc.Service.GetTicketTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketTemplate == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket template not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketTemplate,
	})
}

// UpdateTicketTemplate
func (tc *TicketTemplateController) UpdateTicketTemplate(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketTemplate
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketTemplate(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketTemplate, err := tc.Service.GetTicketTemplateByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket template",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketTemplate,
	})
}

// DeleteTicketTemplate
func (tc *TicketTemplateController) DeleteTicketTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketTemplate(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket template deleted",
	})
}