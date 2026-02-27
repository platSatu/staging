package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScSubjectTypeGeneralSubjectController struct {
	Service *service.ScSubjectTypeGeneralSubjectService
}

func NewScSubjectTypeGeneralSubjectController(s *service.ScSubjectTypeGeneralSubjectService) *ScSubjectTypeGeneralSubjectController {
	return &ScSubjectTypeGeneralSubjectController{Service: s}
}

// CreateScSubjectTypeGeneralSubject
func (sc *ScSubjectTypeGeneralSubjectController) CreateScSubjectTypeGeneralSubject(c *gin.Context) {
	var req struct {
		GpaWeight        *float64 `json:"gpa_weight"`
		MinPassingScore  *float64 `json:"min_passing_score"`
		StartPage        *int     `json:"start_page"`
		EndingPace       *int     `json:"ending_pace"`
		SubjectName      *string  `json:"subject_name"`
		ProductReference *string  `json:"product_reference"`
		StartingPace     *int     `json:"starting_pace"`
		Unit             *string  `json:"unit"`
		IsAssignable     *bool    `json:"is_assignable"`
		IsPace           *bool    `json:"is_pace"`
		PacesNumber      *int     `json:"paces_number"`
		SubjectType      *string  `json:"subject_type"`
		EndPage          *int     `json:"end_page"`
		PerCredit        *float64 `json:"per_credit"`
		TotalPages       *int     `json:"total_pages"`
		CourseName       *string  `json:"course_name"`
		PrevPage         *int     `json:"prev_page"`
		NextPage         *int     `json:"next_page"`
		Units            *int     `json:"units"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	subject := &model.ScSubjectTypeGeneralSubject{
		UserID:           userID.(string),
		GpaWeight:        req.GpaWeight,
		MinPassingScore:  req.MinPassingScore,
		StartPage:        req.StartPage,
		EndingPace:       req.EndingPace,
		SubjectName:      req.SubjectName,
		ProductReference: req.ProductReference,
		StartingPace:     req.StartingPace,
		Unit:             req.Unit,
		IsAssignable:     req.IsAssignable,
		IsPace:           req.IsPace,
		PacesNumber:      req.PacesNumber,
		SubjectType:      req.SubjectType,
		EndPage:          req.EndPage,
		PerCredit:        req.PerCredit,
		TotalPages:       req.TotalPages,
		CourseName:       req.CourseName,
		PrevPage:         req.PrevPage,
		NextPage:         req.NextPage,
		Units:            req.Units,
	}

	if err := sc.Service.CreateScSubjectTypeGeneralSubject(subject); err != nil {
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

// GetAllScSubjectTypeGeneralSubject (hanya milik user login)
func (sc *ScSubjectTypeGeneralSubjectController) GetAllScSubjectTypeGeneralSubject(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	subjects, err := sc.Service.GetAllScSubjectTypeGeneralSubjectByUserID(userID.(string))
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

// GetScSubjectTypeGeneralSubjectByID
func (sc *ScSubjectTypeGeneralSubjectController) GetScSubjectTypeGeneralSubjectByID(c *gin.Context) {
	id := c.Param("id")

	subject, err := sc.Service.GetScSubjectTypeGeneralSubjectByID(id)
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
			"error":   "sc_subject_type_general_subject not found",
		})
		return
	}

	// Hanya pemilik subject boleh melihat
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

// UpdateScSubjectTypeGeneralSubject
func (sc *ScSubjectTypeGeneralSubjectController) UpdateScSubjectTypeGeneralSubject(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		GpaWeight        *float64 `json:"gpa_weight"`
		MinPassingScore  *float64 `json:"min_passing_score"`
		StartPage        *int     `json:"start_page"`
		EndingPace       *int     `json:"ending_pace"`
		SubjectName      *string  `json:"subject_name"`
		ProductReference *string  `json:"product_reference"`
		StartingPace     *int     `json:"starting_pace"`
		Unit             *string  `json:"unit"`
		IsAssignable     *bool    `json:"is_assignable"`
		IsPace           *bool    `json:"is_pace"`
		PacesNumber      *int     `json:"paces_number"`
		SubjectType      *string  `json:"subject_type"`
		EndPage          *int     `json:"end_page"`
		PerCredit        *float64 `json:"per_credit"`
		TotalPages       *int     `json:"total_pages"`
		CourseName       *string  `json:"course_name"`
		PrevPage         *int     `json:"prev_page"`
		NextPage         *int     `json:"next_page"`
		Units            *int     `json:"units"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Pastikan subject itu milik user login
	existing, err := sc.Service.GetScSubjectTypeGeneralSubjectByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_type_general_subject not found",
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

	updateData := &model.ScSubjectTypeGeneralSubject{
		ID:               id,
		GpaWeight:        req.GpaWeight,
		MinPassingScore:  req.MinPassingScore,
		StartPage:        req.StartPage,
		EndingPace:       req.EndingPace,
		SubjectName:      req.SubjectName,
		ProductReference: req.ProductReference,
		StartingPace:     req.StartingPace,
		Unit:             req.Unit,
		IsAssignable:     req.IsAssignable,
		IsPace:           req.IsPace,
		PacesNumber:      req.PacesNumber,
		SubjectType:      req.SubjectType,
		EndPage:          req.EndPage,
		PerCredit:        req.PerCredit,
		TotalPages:       req.TotalPages,
		CourseName:       req.CourseName,
		PrevPage:         req.PrevPage,
		NextPage:         req.NextPage,
		Units:            req.Units,
	}

	if err := sc.Service.UpdateScSubjectTypeGeneralSubject(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScSubjectTypeGeneralSubjectByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScSubjectTypeGeneralSubject
func (sc *ScSubjectTypeGeneralSubjectController) DeleteScSubjectTypeGeneralSubject(c *gin.Context) {
	id := c.Param("id")

	subject, err := sc.Service.GetScSubjectTypeGeneralSubjectByID(id)
	if err != nil || subject == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_subject_type_general_subject not found",
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

	if err := sc.Service.DeleteScSubjectTypeGeneralSubject(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_subject_type_general_subject deleted",
	})
}
