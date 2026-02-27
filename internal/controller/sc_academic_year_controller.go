package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ScAcademicYearController struct {
	Service *service.ScAcademicYearService
}

func NewScAcademicYearController(s *service.ScAcademicYearService) *ScAcademicYearController {
	return &ScAcademicYearController{Service: s}
}

// CreateScAcademicYear
func (sc *ScAcademicYearController) CreateScAcademicYear(c *gin.Context) {
	var req struct {
		Name         *string `json:"name"`
		Status       string  `json:"status"`     // Required enum
		BeginDate    *string `json:"begin_date"` // Format: YYYY-MM-DD
		EndDate      *string `json:"end_date"`   // Format: YYYY-MM-DD
		AcademicYear *string `json:"academic_year"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Parse dates
	var beginDate, endDate *time.Time
	if req.BeginDate != nil && *req.BeginDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.BeginDate); err == nil {
			beginDate = &parsed
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &parsed
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

	academicYear := &model.ScAcademicYear{
		UserID:       userID.(string),
		Name:         req.Name,
		Status:       req.Status,
		BeginDate:    beginDate,
		EndDate:      endDate,
		AcademicYear: req.AcademicYear,
	}

	if err := sc.Service.CreateScAcademicYear(academicYear); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    academicYear,
	})
}

// GetAllScAcademicYear (hanya milik user login)
func (sc *ScAcademicYearController) GetAllScAcademicYear(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	academicYears, err := sc.Service.GetAllScAcademicYearByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    academicYears,
	})
}

// GetScAcademicYearByID
func (sc *ScAcademicYearController) GetScAcademicYearByID(c *gin.Context) {
	id := c.Param("id")

	academicYear, err := sc.Service.GetScAcademicYearByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if academicYear == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_year not found",
		})
		return
	}

	// hanya pemilik academic year boleh melihat
	userID, _ := c.Get("userID")
	if academicYear.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    academicYear,
	})
}

// UpdateScAcademicYear
func (sc *ScAcademicYearController) UpdateScAcademicYear(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name         *string `json:"name"`
		Status       *string `json:"status"`     // Optional enum
		BeginDate    *string `json:"begin_date"` // Format: YYYY-MM-DD
		EndDate      *string `json:"end_date"`   // Format: YYYY-MM-DD
		AcademicYear *string `json:"academic_year"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Parse dates
	var beginDate, endDate *time.Time
	if req.BeginDate != nil && *req.BeginDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.BeginDate); err == nil {
			beginDate = &parsed
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &parsed
		}
	}

	// pastikan academic year itu milik user login
	existing, err := sc.Service.GetScAcademicYearByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_year not found",
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

	// Handle status dengan if-else (perbaikan dari ternary operator)
	var status string
	if req.Status != nil {
		status = *req.Status
	} else {
		status = existing.Status
	}

	updateData := &model.ScAcademicYear{
		ID:           id,
		Name:         req.Name,
		Status:       status,
		BeginDate:    beginDate,
		EndDate:      endDate,
		AcademicYear: req.AcademicYear,
	}

	if err := sc.Service.UpdateScAcademicYear(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScAcademicYearByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScAcademicYear
func (sc *ScAcademicYearController) DeleteScAcademicYear(c *gin.Context) {
	id := c.Param("id")

	academicYear, err := sc.Service.GetScAcademicYearByID(id)
	if err != nil || academicYear == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_year not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if academicYear.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScAcademicYear(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_academic_year deleted",
	})
}
