package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketResellerSettingController struct {
	Service *service.TicketResellerSettingService
}

func NewTicketResellerSettingController(s *service.TicketResellerSettingService) *TicketResellerSettingController {
	return &TicketResellerSettingController{Service: s}
}

// CreateTicketResellerSetting
func (tc *TicketResellerSettingController) CreateTicketResellerSetting(c *gin.Context) {
	var ticketResellerSetting model.TicketResellerSetting
	if err := c.ShouldBindJSON(&ticketResellerSetting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketResellerSetting(&ticketResellerSetting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketResellerSetting,
	})
}

// GetAllTicketResellerSettings
func (tc *TicketResellerSettingController) GetAllTicketResellerSettings(c *gin.Context) {
	ticketResellerSettings, err := tc.Service.GetAllTicketResellerSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketResellerSettings,
	})
}

// GetTicketResellerSettingByID
func (tc *TicketResellerSettingController) GetTicketResellerSettingByID(c *gin.Context) {
	id := c.Param("id")
	ticketResellerSetting, err := tc.Service.GetTicketResellerSettingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketResellerSetting == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket reseller setting not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketResellerSetting,
	})
}

// UpdateTicketResellerSetting
func (tc *TicketResellerSettingController) UpdateTicketResellerSetting(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketResellerSetting
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketResellerSetting(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketResellerSetting, err := tc.Service.GetTicketResellerSettingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket reseller setting",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketResellerSetting,
	})
}

// DeleteTicketResellerSetting
func (tc *TicketResellerSettingController) DeleteTicketResellerSetting(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketResellerSetting(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket reseller setting deleted",
	})
}
