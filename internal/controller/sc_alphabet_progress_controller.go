package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScAlphabetProgressController struct {
	Service *service.ScAlphabetProgressService
}

func NewScAlphabetProgressController(s *service.ScAlphabetProgressService) *ScAlphabetProgressController {
	return &ScAlphabetProgressController{Service: s}
}

// CreateScAlphabetProgress
func (sc *ScAlphabetProgressController) CreateScAlphabetProgress(c *gin.Context) {
	var req struct {
		Level       string `json:"level"`
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

	progress := &model.ScAlphabetProgress{
		UserID:      userID.(string),
		Level:       req.Level,
		LVertical:   req.LVertical,
		LHorizontal: req.LHorizontal,
		Score:       req.Score,
	}

	if err := sc.Service.CreateScAlphabetProgress(progress); err != nil {
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

// GetAllScAlphabetProgress (hanya milik user login)
func (sc *ScAlphabetProgressController) GetAllScAlphabetProgress(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	progresses, err := sc.Service.GetAllScAlphabetProgressByUserID(userID.(string))
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

// GetScAlphabetProgressByID
func (sc *ScAlphabetProgressController) GetScAlphabetProgressByID(c *gin.Context) {
	id := c.Param("id")

	progress, err := sc.Service.GetScAlphabetProgressByID(id)
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
			"error":   "sc_alphabet_progress not found",
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

// UpdateScAlphabetProgress
func (sc *ScAlphabetProgressController) UpdateScAlphabetProgress(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Level       string `json:"level"`
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
	existing, err := sc.Service.GetScAlphabetProgressByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_alphabet_progress not found",
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

	updateData := &model.ScAlphabetProgress{
		ID:          id,
		Level:       req.Level,
		LVertical:   req.LVertical,
		LHorizontal: req.LHorizontal,
		Score:       req.Score,
	}

	if err := sc.Service.UpdateScAlphabetProgress(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScAlphabetProgressByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScAlphabetProgress
func (sc *ScAlphabetProgressController) DeleteScAlphabetProgress(c *gin.Context) {
	id := c.Param("id")

	progress, err := sc.Service.GetScAlphabetProgressByID(id)
	if err != nil || progress == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_alphabet_progress not found",
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

	if err := sc.Service.DeleteScAlphabetProgress(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_alphabet_progress deleted",
	})
}