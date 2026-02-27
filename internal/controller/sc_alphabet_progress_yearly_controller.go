package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScAlphabetProgressYearlyController struct {
	Service *service.ScAlphabetProgressYearlyService
}

func NewScAlphabetProgressYearlyController(s *service.ScAlphabetProgressYearlyService) *ScAlphabetProgressYearlyController {
	return &ScAlphabetProgressYearlyController{Service: s}
}

// CreateScAlphabetProgressYearly
func (sc *ScAlphabetProgressYearlyController) CreateScAlphabetProgressYearly(c *gin.Context) {
	var req struct {
		LVertical   string `json:"l_vertical"`
		LHorizontal string `json:"l_horizontal"`
		Score       string `json:"score"`
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

	progress := &model.ScAlphabetProgressYearly{
		UserID:      userID.(string),
		LVertical:   req.LVertical,
		LHorizontal: req.LHorizontal,
		Score:       req.Score,
	}

	if err := sc.Service.CreateScAlphabetProgressYearly(progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    progress,
	})
}

// GetAllScAlphabetProgressYearly (hanya milik user login)
func (sc *ScAlphabetProgressYearlyController) GetAllScAlphabetProgressYearly(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	progresses, err := sc.Service.GetAllScAlphabetProgressYearlyByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    progresses,
	})
}

// GetScAlphabetProgressYearlyByID
func (sc *ScAlphabetProgressYearlyController) GetScAlphabetProgressYearlyByID(c *gin.Context) {
	id := c.Param("id")

	progress, err := sc.Service.GetScAlphabetProgressYearlyByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if progress == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_alphabet_progress_yearly not found",
		})
		return
	}

	// hanya pemilik progress boleh melihat
	userID, _ := c.Get("userID")
	if progress.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    progress,
	})
}

// UpdateScAlphabetProgressYearly
func (sc *ScAlphabetProgressYearlyController) UpdateScAlphabetProgressYearly(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		LVertical   string `json:"l_vertical"`
		LHorizontal string `json:"l_horizontal"`
		Score       string `json:"score"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan progress itu milik user login
	existing, err := sc.Service.GetScAlphabetProgressYearlyByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_alphabet_progress_yearly not found",
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

	updateData := &model.ScAlphabetProgressYearly{
		ID:          id,
		LVertical:   req.LVertical,
		LHorizontal: req.LHorizontal,
		Score:       req.Score,
	}

	if err := sc.Service.UpdateScAlphabetProgressYearly(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScAlphabetProgressYearlyByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScAlphabetProgressYearly
func (sc *ScAlphabetProgressYearlyController) DeleteScAlphabetProgressYearly(c *gin.Context) {
	id := c.Param("id")

	progress, err := sc.Service.GetScAlphabetProgressYearlyByID(id)
	if err != nil || progress == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_alphabet_progress_yearly not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if progress.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScAlphabetProgressYearly(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_alphabet_progress_yearly deleted",
	})
}