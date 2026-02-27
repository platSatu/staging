package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScLearningCenterController struct {
	Service *service.ScLearningCenterService
}

func NewScLearningCenterController(s *service.ScLearningCenterService) *ScLearningCenterController {
	return &ScLearningCenterController{Service: s}
}

// CreateScLearningCenter
func (sc *ScLearningCenterController) CreateScLearningCenter(c *gin.Context) {
	var req struct {
		Name        *string `json:"name"`
		GroupIDN    *string `json:"group_idn"`
		Principal   *string `json:"principal"`
		HomeTeacher *string `json:"home_teacher"`
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

	learningCenter := &model.ScLearningCenter{
		UserID:      userID.(string),
		Name:        req.Name,
		GroupIDN:    req.GroupIDN,
		Principal:   req.Principal,
		HomeTeacher: req.HomeTeacher,
	}

	if err := sc.Service.CreateScLearningCenter(learningCenter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    learningCenter,
	})
}

// GetAllScLearningCenter (hanya milik user login)
func (sc *ScLearningCenterController) GetAllScLearningCenter(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	learningCenters, err := sc.Service.GetAllScLearningCenterByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    learningCenters,
	})
}

// GetScLearningCenterByID
func (sc *ScLearningCenterController) GetScLearningCenterByID(c *gin.Context) {
	id := c.Param("id")

	learningCenter, err := sc.Service.GetScLearningCenterByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if learningCenter == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_learning_center not found",
		})
		return
	}

	// hanya pemilik learning center boleh melihat
	userID, _ := c.Get("userID")
	if learningCenter.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    learningCenter,
	})
}

// UpdateScLearningCenter
func (sc *ScLearningCenterController) UpdateScLearningCenter(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        *string `json:"name"`
		GroupIDN    *string `json:"group_idn"`
		Principal   *string `json:"principal"`
		HomeTeacher *string `json:"home_teacher"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan learning center itu milik user login
	existing, err := sc.Service.GetScLearningCenterByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_learning_center not found",
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

	updateData := &model.ScLearningCenter{
		ID:          id,
		Name:        req.Name,
		GroupIDN:    req.GroupIDN,
		Principal:   req.Principal,
		HomeTeacher: req.HomeTeacher,
	}

	if err := sc.Service.UpdateScLearningCenter(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScLearningCenterByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScLearningCenter
func (sc *ScLearningCenterController) DeleteScLearningCenter(c *gin.Context) {
	id := c.Param("id")

	learningCenter, err := sc.Service.GetScLearningCenterByID(id)
	if err != nil || learningCenter == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_learning_center not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if learningCenter.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScLearningCenter(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_learning_center deleted",
	})
}
