package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScHabitsController struct {
	Service *service.ScHabitsService
}

func NewScHabitsController(s *service.ScHabitsService) *ScHabitsController {
	return &ScHabitsController{Service: s}
}

// CreateScHabits
func (sc *ScHabitsController) CreateScHabits(c *gin.Context) {
	var req struct {
		Subject     string  `json:"subject" binding:"required"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	habit := &model.ScHabits{
		UserID:      userID.(string),
		Subject:     req.Subject,
		Description: req.Description,
	}

	if err := sc.Service.CreateScHabits(habit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    habit,
	})
}

// GetAllScHabits (hanya milik user login)
func (sc *ScHabitsController) GetAllScHabits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	habits, err := sc.Service.GetAllScHabitsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    habits,
	})
}

// GetScHabitsByID
func (sc *ScHabitsController) GetScHabitsByID(c *gin.Context) {
	id := c.Param("id")

	habit, err := sc.Service.GetScHabitsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if habit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_habits not found",
		})
		return
	}

	// hanya pemilik habit boleh melihat
	userID, _ := c.Get("userID")
	if habit.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    habit,
	})
}

// UpdateScHabits
func (sc *ScHabitsController) UpdateScHabits(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Subject     string  `json:"subject"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan habit itu milik user login
	existing, err := sc.Service.GetScHabitsByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_habits not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	updateData := &model.ScHabits{
		ID:          id,
		Subject:     req.Subject,
		Description: req.Description,
	}

	if err := sc.Service.UpdateScHabits(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScHabitsByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScHabits
func (sc *ScHabitsController) DeleteScHabits(c *gin.Context) {
	id := c.Param("id")

	habit, err := sc.Service.GetScHabitsByID(id)
	if err != nil || habit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_habits not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if habit.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScHabits(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_habits deleted",
	})
}