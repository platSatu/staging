package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScSubjectTypeGeneralController struct {
	Service *service.ScSubjectTypeGeneralService
}

func NewScSubjectTypeGeneralController(s *service.ScSubjectTypeGeneralService) *ScSubjectTypeGeneralController {
	return &ScSubjectTypeGeneralController{Service: s}
}

// CreateScSubjectTypeGeneral
func (sc *ScSubjectTypeGeneralController) CreateScSubjectTypeGeneral(c *gin.Context) {
	var req struct {
		SubjectName *string `json:"subject_name"`
		Curriculum  *string `json:"curriculum"`
		Group       *string `json:"group"`
		Status      *string `json:"status"`
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

	subjectTypeGeneral := &model.ScSubjectTypeGeneral{
		UserID:      userID.(string),
		SubjectName: req.SubjectName,
		Curriculum:  req.Curriculum,
		Group:       req.Group,
		Status:      req.Status,
	}

	if err := sc.Service.CreateScSubjectTypeGeneral(subjectTypeGeneral); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    subjectTypeGeneral,
	})
}

// GetAllScSubjectTypeGeneral (hanya milik user login)
func (sc *ScSubjectTypeGeneralController) GetAllScSubjectTypeGeneral(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	subjectTypeGenerals, err := sc.Service.GetAllScSubjectTypeGeneralByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subjectTypeGenerals,
	})
}

// GetScSubjectTypeGeneralByID
func (sc *ScSubjectTypeGeneralController) GetScSubjectTypeGeneralByID(c *gin.Context) {
	id := c.Param("id")

	subjectTypeGeneral, err := sc.Service.GetScSubjectTypeGeneralByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if subjectTypeGeneral == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_type_general not found",
		})
		return
	}

	// hanya pemilik subject type general boleh melihat
	userID, _ := c.Get("userID")
	if subjectTypeGeneral.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    subjectTypeGeneral,
	})
}

// UpdateScSubjectTypeGeneral
func (sc *ScSubjectTypeGeneralController) UpdateScSubjectTypeGeneral(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		SubjectName *string `json:"subject_name"`
		Curriculum  *string `json:"curriculum"`
		Group       *string `json:"group"`
		Status      *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan subject type general itu milik user login
	existing, err := sc.Service.GetScSubjectTypeGeneralByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_type_general not found",
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

	updateData := &model.ScSubjectTypeGeneral{
		ID:          id,
		SubjectName: req.SubjectName,
		Curriculum:  req.Curriculum,
		Group:       req.Group,
		Status:      req.Status,
	}

	if err := sc.Service.UpdateScSubjectTypeGeneral(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScSubjectTypeGeneralByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScSubjectTypeGeneral
func (sc *ScSubjectTypeGeneralController) DeleteScSubjectTypeGeneral(c *gin.Context) {
	id := c.Param("id")

	subjectTypeGeneral, err := sc.Service.GetScSubjectTypeGeneralByID(id)
	if err != nil || subjectTypeGeneral == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_type_general not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if subjectTypeGeneral.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScSubjectTypeGeneral(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_subject_type_general deleted",
	})
}
