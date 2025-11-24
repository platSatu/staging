package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type JenisTiketController struct {
	Service *service.JenisTiketService
}

func NewJenisTiketController(s *service.JenisTiketService) *JenisTiketController {
	return &JenisTiketController{Service: s}
}

// CreateJenisTiket
func (jtc *JenisTiketController) CreateJenisTiket(c *gin.Context) {
	var jenisTiket model.JenisTiket
	if err := c.ShouldBindJSON(&jenisTiket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := jtc.Service.CreateJenisTiket(&jenisTiket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    jenisTiket,
	})
}

// GetAllJenisTiket
func (jtc *JenisTiketController) GetAllJenisTiket(c *gin.Context) {
	jenisTikets, err := jtc.Service.GetAllJenisTiket()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jenisTikets,
	})
}

// GetJenisTiketByID
func (jtc *JenisTiketController) GetJenisTiketByID(c *gin.Context) {
	id := c.Param("id")
	jenisTiket, err := jtc.Service.GetJenisTiketByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if jenisTiket == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Jenis tiket not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jenisTiket,
	})
}

// GetJenisTiketByUser
func (jtc *JenisTiketController) GetJenisTiketByUser(c *gin.Context) {
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

	jenisTikets, err := jtc.Service.GetJenisTiketByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jenisTikets,
	})
}

// GetJenisTiketByEvent
func (jtc *JenisTiketController) GetJenisTiketByEvent(c *gin.Context) {
	eventID := c.Param("event_id")
	jenisTikets, err := jtc.Service.GetJenisTiketByEventID(eventID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jenisTikets,
	})
}

// UpdateJenisTiket
func (jtc *JenisTiketController) UpdateJenisTiket(c *gin.Context) {
	id := c.Param("id")
	var updateData model.JenisTiket
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := jtc.Service.UpdateJenisTiket(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	jenisTiket, err := jtc.Service.GetJenisTiketByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated jenis tiket",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    jenisTiket,
	})
}

// DeleteJenisTiket
func (jtc *JenisTiketController) DeleteJenisTiket(c *gin.Context) {
	id := c.Param("id")
	if err := jtc.Service.DeleteJenisTiket(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Jenis tiket deleted",
	})
}
