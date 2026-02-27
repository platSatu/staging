package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScTemplateHabitsController struct {
	Service *service.ScTemplateHabitsService
}

func NewScTemplateHabitsController(s *service.ScTemplateHabitsService) *ScTemplateHabitsController {
	return &ScTemplateHabitsController{Service: s}
}

// CreateScTemplateHabits
func (sc *ScTemplateHabitsController) CreateScTemplateHabits(c *gin.Context) {
	var req struct {
		Level        string `json:"level"`
		SubjectID    string `json:"subject_id"`
		SubSubjectID string `json:"sub_subject_id"`
		Grade        string `json:"grade"`
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

	template := &model.ScTemplateHabits{
		UserID:       userID.(string),
		Level:        req.Level,
		SubjectID:    req.SubjectID,
		SubSubjectID: req.SubSubjectID,
		Grade:        req.Grade,
	}

	if err := sc.Service.CreateScTemplateHabits(template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    template,
	})
}

// GetAllScTemplateHabits (hanya milik user login)
func (sc *ScTemplateHabitsController) GetAllScTemplateHabits(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	templates, err := sc.Service.GetAllScTemplateHabitsByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    templates,
	})
}

// GetScTemplateHabitsByID
func (sc *ScTemplateHabitsController) GetScTemplateHabitsByID(c *gin.Context) {
	id := c.Param("id")

	template, err := sc.Service.GetScTemplateHabitsByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if template == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_template_habits not found",
		})
		return
	}

	// hanya pemilik template boleh melihat
	userID, _ := c.Get("userID")
	if template.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    template,
	})
}

// UpdateScTemplateHabits
func (sc *ScTemplateHabitsController) UpdateScTemplateHabits(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Level        string `json:"level"`
		SubjectID    string `json:"subject_id"`
		SubSubjectID string `json:"sub_subject_id"`
		Grade        string `json:"grade"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan template itu milik user login
	existing, err := sc.Service.GetScTemplateHabitsByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_template_habits not found",
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

	updateData := &model.ScTemplateHabits{
		ID:           id,
		Level:        req.Level,
		SubjectID:    req.SubjectID,
		SubSubjectID: req.SubSubjectID,
		Grade:        req.Grade,
	}

	if err := sc.Service.UpdateScTemplateHabits(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScTemplateHabitsByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScTemplateHabits
func (sc *ScTemplateHabitsController) DeleteScTemplateHabits(c *gin.Context) {
	id := c.Param("id")

	template, err := sc.Service.GetScTemplateHabitsByID(id)
	if err != nil || template == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_template_habits not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if template.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScTemplateHabits(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_template_habits deleted",
	})
}