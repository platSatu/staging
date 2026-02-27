package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScAcademicProjectionController struct {
	Service *service.ScAcademicProjectionService
}

func NewScAcademicProjectionController(s *service.ScAcademicProjectionService) *ScAcademicProjectionController {
	return &ScAcademicProjectionController{Service: s}
}

// CreateScAcademicProjection
func (sc *ScAcademicProjectionController) CreateScAcademicProjection(c *gin.Context) {
	var req struct {
		StudentID    *string `json:"student_id"`
		AcademicYear *string `json:"academic_year"`
		Semester     *string `json:"semester"`
		QuarterID    *string `json:"quarter_id"`
		IsSplitIso   *bool   `json:"is_split_iso"`
		Level        *string `json:"level"`
		LcID         *string `json:"lc_id"`
		TotalSchool  *string `json:"total_school"` // Ubah ke *string
		TotalPages   *string `json:"total_pages"`  // Ubah ke *string
		Status       *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Parse total_school dan total_pages ke int
	var totalSchool, totalPages *int
	if req.TotalSchool != nil && *req.TotalSchool != "" {
		if parsed, err := strconv.Atoi(*req.TotalSchool); err == nil {
			totalSchool = &parsed
		}
	}
	if req.TotalPages != nil && *req.TotalPages != "" {
		if parsed, err := strconv.Atoi(*req.TotalPages); err == nil {
			totalPages = &parsed
		}
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

	projection := &model.ScAcademicProjection{
		UserID:       userID.(string),
		StudentID:    req.StudentID,
		AcademicYear: req.AcademicYear,
		Semester:     req.Semester,
		QuarterID:    req.QuarterID,
		IsSplitIso:   req.IsSplitIso,
		Level:        req.Level,
		LcID:         req.LcID,
		TotalSchool:  totalSchool,
		TotalPages:   totalPages,
		Status:       req.Status,
	}

	if err := sc.Service.CreateScAcademicProjection(projection); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    projection,
	})
}

// GetAllScAcademicProjection (hanya milik user login)
func (sc *ScAcademicProjectionController) GetAllScAcademicProjection(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	projections, err := sc.Service.GetAllScAcademicProjectionByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projections,
	})
}

// GetScAcademicProjectionByID
func (sc *ScAcademicProjectionController) GetScAcademicProjectionByID(c *gin.Context) {
	id := c.Param("id")

	projection, err := sc.Service.GetScAcademicProjectionByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if projection == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_projection not found",
		})
		return
	}

	// hanya pemilik projection boleh melihat
	userID, _ := c.Get("userID")
	if projection.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projection,
	})
}

// UpdateScAcademicProjection
func (sc *ScAcademicProjectionController) UpdateScAcademicProjection(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		StudentID    *string `json:"student_id"`
		AcademicYear *string `json:"academic_year"`
		Semester     *string `json:"semester"`
		QuarterID    *string `json:"quarter_id"`
		IsSplitIso   *bool   `json:"is_split_iso"`
		Level        *string `json:"level"`
		LcID         *string `json:"lc_id"`
		TotalSchool  *string `json:"total_school"` // Ubah ke *string
		TotalPages   *string `json:"total_pages"`  // Ubah ke *string
		Status       *string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Parse total_school dan total_pages ke int
	var totalSchool, totalPages *int
	if req.TotalSchool != nil && *req.TotalSchool != "" {
		if parsed, err := strconv.Atoi(*req.TotalSchool); err == nil {
			totalSchool = &parsed
		}
	}
	if req.TotalPages != nil && *req.TotalPages != "" {
		if parsed, err := strconv.Atoi(*req.TotalPages); err == nil {
			totalPages = &parsed
		}
	}

	// pastikan projection itu milik user login
	existing, err := sc.Service.GetScAcademicProjectionByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_projection not found",
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

	updateData := &model.ScAcademicProjection{
		ID:           id,
		StudentID:    req.StudentID,
		AcademicYear: req.AcademicYear,
		Semester:     req.Semester,
		QuarterID:    req.QuarterID,
		IsSplitIso:   req.IsSplitIso,
		Level:        req.Level,
		LcID:         req.LcID,
		TotalSchool:  totalSchool,
		TotalPages:   totalPages,
		Status:       req.Status,
	}

	if err := sc.Service.UpdateScAcademicProjection(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScAcademicProjectionByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScAcademicProjection
func (sc *ScAcademicProjectionController) DeleteScAcademicProjection(c *gin.Context) {
	id := c.Param("id")

	projection, err := sc.Service.GetScAcademicProjectionByID(id)
	if err != nil || projection == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_projection not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if projection.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScAcademicProjection(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_academic_projection deleted",
	})
}

// GetScAcademicProjectionByUser returns academic projection based on user login (student_id)
func (sc *ScAcademicProjectionController) GetScAcademicProjectionByUser(c *gin.Context) {
	// Ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Tampilkan semua projection dimana student_id = user yang login
	projections, err := sc.Service.GetScAcademicProjectionByStudentID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    projections,
	})
}
