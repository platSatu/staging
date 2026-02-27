package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScGradeController struct {
	Service *service.ScGradeService
}

func NewScGradeController(s *service.ScGradeService) *ScGradeController {
	return &ScGradeController{Service: s}
}

// CreateScGrade
func (sc *ScGradeController) CreateScGrade(c *gin.Context) {
	var req struct {
		GradeName string `json:"grade_name"`
		MinScore  string `json:"min_score"`
		MaxScore  string `json:"max_score"`
		Status    string `json:"status"`
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

	grade := &model.ScGrade{
		UserID:    userID.(string),
		GradeName: req.GradeName,
		MinScore:  req.MinScore,
		MaxScore:  req.MaxScore,
		Status:    req.Status,
	}

	if err := sc.Service.CreateScGrade(grade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    grade,
	})
}

// GetAllScGrade (hanya milik user login)
func (sc *ScGradeController) GetAllScGrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	grades, err := sc.Service.GetAllScGradeByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    grades,
	})
}

// GetScGradeByID
func (sc *ScGradeController) GetScGradeByID(c *gin.Context) {
	id := c.Param("id")

	grade, err := sc.Service.GetScGradeByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if grade == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_grade not found",
		})
		return
	}

	// hanya pemilik grade boleh melihat
	userID, _ := c.Get("userID")
	if grade.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    grade,
	})
}

// UpdateScGrade
func (sc *ScGradeController) UpdateScGrade(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		GradeName string `json:"grade_name"`
		MinScore  string `json:"min_score"`
		MaxScore  string `json:"max_score"`
		Status    string `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan grade itu milik user login
	existing, err := sc.Service.GetScGradeByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_grade not found",
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

	updateData := &model.ScGrade{
		ID:        id,
		GradeName: req.GradeName,
		MinScore:  req.MinScore,
		MaxScore:  req.MaxScore,
		Status:    req.Status,
	}

	if err := sc.Service.UpdateScGrade(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScGradeByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScGrade
func (sc *ScGradeController) DeleteScGrade(c *gin.Context) {
	id := c.Param("id")

	grade, err := sc.Service.GetScGradeByID(id)
	if err != nil || grade == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_grade not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if grade.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScGrade(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_grade deleted",
	})
}