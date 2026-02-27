package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScSubjectListSemesterController struct {
	Service *service.ScSubjectListSemesterService
}

func NewScSubjectListSemesterController(s *service.ScSubjectListSemesterService) *ScSubjectListSemesterController {
	return &ScSubjectListSemesterController{Service: s}
}

// CreateScSubjectListSemester
func (sc *ScSubjectListSemesterController) CreateScSubjectListSemester(c *gin.Context) {
	var req struct {
		No      string `json:"no"`
		Subject string `json:"subject"`
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

	subject := &model.ScSubjectListSemester{
		UserID:  userID.(string),
		No:      req.No,
		Subject: req.Subject,
	}

	if err := sc.Service.CreateScSubjectListSemester(subject); err != nil {
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

// GetAllScSubjectListSemester (hanya milik user login)
func (sc *ScSubjectListSemesterController) GetAllScSubjectListSemester(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	subjects, err := sc.Service.GetAllScSubjectListSemesterByUserID(userID.(string))
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

// GetScSubjectListSemesterByID
func (sc *ScSubjectListSemesterController) GetScSubjectListSemesterByID(c *gin.Context) {
	id := c.Param("id")

	subject, err := sc.Service.GetScSubjectListSemesterByID(id)
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
			"error":   "sc_subject_list_semester not found",
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

// UpdateScSubjectListSemester
func (sc *ScSubjectListSemesterController) UpdateScSubjectListSemester(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		No      string `json:"no"`
		Subject string `json:"subject"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan subject itu milik user login
	existing, err := sc.Service.GetScSubjectListSemesterByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_list_semester not found",
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

	updateData := &model.ScSubjectListSemester{
		ID:      id,
		No:      req.No,
		Subject: req.Subject,
	}

	if err := sc.Service.UpdateScSubjectListSemester(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScSubjectListSemesterByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScSubjectListSemester
func (sc *ScSubjectListSemesterController) DeleteScSubjectListSemester(c *gin.Context) {
	id := c.Param("id")

	subject, err := sc.Service.GetScSubjectListSemesterByID(id)
	if err != nil || subject == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_list_semester not found",
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

	if err := sc.Service.DeleteScSubjectListSemester(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_subject_list_semester deleted",
	})
}