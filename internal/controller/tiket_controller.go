package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type TiketController struct {
	Service *service.TiketService
}

func NewTiketController(s *service.TiketService) *TiketController {
	return &TiketController{Service: s}
}

// CreateTiket
func (tc *TiketController) CreateTiket(c *gin.Context) {
	var tiket model.Tiket
	if err := c.ShouldBindJSON(&tiket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.CreateTiket(&tiket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    tiket,
	})
}

// GetAllTiket
func (tc *TiketController) GetAllTiket(c *gin.Context) {
	tikets, err := tc.Service.GetAllTiket()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tikets,
	})
}

// GetTiketByID
func (tc *TiketController) GetTiketByID(c *gin.Context) {
	id := c.Param("id")
	tiket, err := tc.Service.GetTiketByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if tiket == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Tiket not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tiket,
	})
}

// GetTiketByUser
func (tc *TiketController) GetTiketByUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := helper.GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	tikets, err := tc.Service.GetTiketByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tikets,
	})
}

// GetTiketByEvent
func (tc *TiketController) GetTiketByEvent(c *gin.Context) {
	eventID := c.Param("event_id")
	tikets, err := tc.Service.GetTiketByEventID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tikets,
	})
}

// GetTiketByKodeBooking
func (tc *TiketController) GetTiketByKodeBooking(c *gin.Context) {
	kodeBooking := c.Param("kode_booking")
	tiket, err := tc.Service.GetTiketByKodeBooking(kodeBooking)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if tiket == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Tiket not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tiket,
	})
}

// UpdateTiket
func (tc *TiketController) UpdateTiket(c *gin.Context) {
	id := c.Param("id")
	var updateData model.Tiket
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := tc.Service.UpdateTiket(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	tiket, err := tc.Service.GetTiketByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated tiket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tiket,
	})
}

// DeleteTiket
func (tc *TiketController) DeleteTiket(c *gin.Context) {
	id := c.Param("id")
	if err := tc.Service.DeleteTiket(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tiket deleted",
	})
}
