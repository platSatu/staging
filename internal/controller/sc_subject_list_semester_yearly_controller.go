package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScSubjectListSemesterYearlyController struct {
	Service *service.ScSubjectListSemesterYearlyService
}

func NewScSubjectListSemesterYearlyController(s *service.ScSubjectListSemesterYearlyService) *ScSubjectListSemesterYearlyController {
	return &ScSubjectListSemesterYearlyController{Service: s}
}

// CreateScSubjectListSemesterYearly
func (sc *ScSubjectListSemesterYearlyController) CreateScSubjectListSemesterYearly(c *gin.Context) {
	var req struct {
		StartLevel string `json:"start_level"`
		EndLevel   string `json:"end_level"`
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

	subject := &model.ScSubjectListSemesterYearly{
		UserID:     userID.(string),
		StartLevel: req.StartLevel,
		EndLevel:   req.EndLevel,
	}

	if err := sc.Service.CreateScSubjectListSemesterYearly(subject); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    subject,
	})
}

// GetAllScSubjectListSemesterYearly (hanya milik user login)
func (sc *ScSubjectListSemesterYearlyController) GetAllScSubjectListSemesterYearly(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	subjects, err := sc.Service.GetAllScSubjectListSemesterYearlyByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subjects,
	})
}

// GetScSubjectListSemesterYearlyByID
func (sc *ScSubjectListSemesterYearlyController) GetScSubjectListSemesterYearlyByID(c *gin.Context) {
	id := c.Param("id")

	subject, err := sc.Service.GetScSubjectListSemesterYearlyByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if subject == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_list_semester_yearly not found",
		})
		return
	}

	// hanya pemilik subject boleh melihat
	userID, _ := c.Get("userID")
	if subject.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subject,
	})
}

// UpdateScSubjectListSemesterYearly
func (sc *ScSubjectListSemesterYearlyController) UpdateScSubjectListSemesterYearly(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		StartLevel string `json:"start_level"`
		EndLevel   string `json:"end_level"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan subject itu milik user login
	existing, err := sc.Service.GetScSubjectListSemesterYearlyByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_list_semester_yearly not found",
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

	updateData := &model.ScSubjectListSemesterYearly{
		ID:         id,
		StartLevel: req.StartLevel,
		EndLevel:   req.EndLevel,
	}

	if err := sc.Service.UpdateScSubjectListSemesterYearly(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScSubjectListSemesterYearlyByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScSubjectListSemesterYearly
func (sc *ScSubjectListSemesterYearlyController) DeleteScSubjectListSemesterYearly(c *gin.Context) {
	id := c.Param("id")

	subject, err := sc.Service.GetScSubjectListSemesterYearlyByID(id)
	if err != nil || subject == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_list_semester_yearly not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if subject.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScSubjectListSemesterYearly(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_subject_list_semester_yearly deleted",
	})
}