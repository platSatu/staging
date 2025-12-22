package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketFeeSettingController struct {
	Service *service.TicketFeeSettingService
}

func NewTicketFeeSettingController(s *service.TicketFeeSettingService) *TicketFeeSettingController {
	return &TicketFeeSettingController{Service: s}
}

func (tc *TicketFeeSettingController) GetAllTicketFeeSettingsPublic(c *gin.Context) {
	ticketFeeSettings, err := tc.Service.GetAllTicketFeeSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketFeeSettings,
	})
}

// CreateTicketFeeSetting
func (tc *TicketFeeSettingController) CreateTicketFeeSetting(c *gin.Context) {
	var ticketFeeSetting model.TicketFeeSetting
	if err := c.ShouldBindJSON(&ticketFeeSetting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTicketFeeSetting(&ticketFeeSetting); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    ticketFeeSetting,
	})
}

// GetAllTicketFeeSettings
func (tc *TicketFeeSettingController) GetAllTicketFeeSettings(c *gin.Context) {
	ticketFeeSettings, err := tc.Service.GetAllTicketFeeSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketFeeSettings,
	})
}

// GetTicketFeeSettingByID
func (tc *TicketFeeSettingController) GetTicketFeeSettingByID(c *gin.Context) {
	id := c.Param("id")
	ticketFeeSetting, err := tc.Service.GetTicketFeeSettingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if ticketFeeSetting == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Ticket fee setting not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketFeeSetting,
	})
}

// UpdateTicketFeeSetting
func (tc *TicketFeeSettingController) UpdateTicketFeeSetting(c *gin.Context) {
	id := c.Param("id")
	var updateData model.TicketFeeSetting
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTicketFeeSetting(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ticketFeeSetting, err := tc.Service.GetTicketFeeSettingByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated ticket fee setting",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    ticketFeeSetting,
	})
}

// DeleteTicketFeeSetting
func (tc *TicketFeeSettingController) DeleteTicketFeeSetting(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTicketFeeSetting(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ticket fee setting deleted",
	})
}
