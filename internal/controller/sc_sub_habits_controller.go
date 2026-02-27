package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScSubHabitsController struct {
	Service *service.ScSubHabitsService
}

func NewScSubHabitsController(s *service.ScSubHabitsService) *ScSubHabitsController {
	return &ScSubHabitsController{Service: s}
}

// CreateScSubHabits
func (sc *ScSubHabitsController) CreateScSubHabits(c *gin.Context) {
	var req struct {
		HabitsID    string  `json:"habits_id" binding:"required"`
		Subject     *string `json:"subject"`
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

	subHabit := &model.ScSubHabits{
		UserID:      userID.(string),
		HabitsID:    req.HabitsID,
		Subject:     req.Subject,
		Description: req.Description,
	}

	if err := sc.Service.CreateScSubHabits(subHabit); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    subHabit,
	})
}

// GetAllScSubHabits (hanya milik user login)
func (sc *ScSubHabitsController) GetAllScSubHabits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	subHabits, err := sc.Service.GetAllScSubHabitsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subHabits,
	})
}

// GetScSubHabitsByID
func (sc *ScSubHabitsController) GetScSubHabitsByID(c *gin.Context) {
	id := c.Param("id")

	subHabit, err := sc.Service.GetScSubHabitsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if subHabit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_sub_habits not found",
		})
		return
	}

	// hanya pemilik subHabit boleh melihat
	userID, _ := c.Get("userID")
	if subHabit.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subHabit,
	})
}

// UpdateScSubHabits
func (sc *ScSubHabitsController) UpdateScSubHabits(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		HabitsID    string  `json:"habits_id"`
		Subject     *string `json:"subject"`
		Description *string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan subHabit itu milik user login
	existing, err := sc.Service.GetScSubHabitsByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_sub_habits not found",
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

	updateData := &model.ScSubHabits{
		ID:          id,
		HabitsID:    req.HabitsID,
		Subject:     req.Subject,
		Description: req.Description,
	}

	if err := sc.Service.UpdateScSubHabits(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScSubHabitsByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScSubHabits
func (sc *ScSubHabitsController) DeleteScSubHabits(c *gin.Context) {
	id := c.Param("id")

	subHabit, err := sc.Service.GetScSubHabitsByID(id)
	if err != nil || subHabit == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_sub_habits not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if subHabit.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScSubHabits(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_sub_habits deleted",
	})
}