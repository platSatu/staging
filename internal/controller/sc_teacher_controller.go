package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScTeacherController struct {
	Service *service.ScTeacherService
}

func NewScTeacherController(s *service.ScTeacherService) *ScTeacherController {
	return &ScTeacherController{Service: s}
}

// CreateScTeacher
func (sc *ScTeacherController) CreateScTeacher(c *gin.Context) {
	var req struct {
		Name   *string `json:"name"`
		Email  *string `json:"email"`
		Phone  *string `json:"phone"`
		Mobile *string `json:"mobile"`
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

	teacher := &model.ScTeacher{
		UserID: userID.(string),
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Mobile: req.Mobile,
	}

	if err := sc.Service.CreateScTeacher(teacher); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    teacher,
	})
}

// GetAllScTeacher (hanya milik user login)
func (sc *ScTeacherController) GetAllScTeacher(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	teachers, err := sc.Service.GetAllScTeacherByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    teachers,
	})
}

// GetScTeacherByID
func (sc *ScTeacherController) GetScTeacherByID(c *gin.Context) {
	id := c.Param("id")

	teacher, err := sc.Service.GetScTeacherByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if teacher == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_teacher not found",
		})
		return
	}

	// hanya pemilik teacher boleh melihat
	userID, _ := c.Get("userID")
	if teacher.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    teacher,
	})
}

// UpdateScTeacher
func (sc *ScTeacherController) UpdateScTeacher(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name   *string `json:"name"`
		Email  *string `json:"email"`
		Phone  *string `json:"phone"`
		Mobile *string `json:"mobile"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan teacher itu milik user login
	existing, err := sc.Service.GetScTeacherByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_teacher not found",
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

	updateData := &model.ScTeacher{
		ID:     id,
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		Mobile: req.Mobile,
	}

	if err := sc.Service.UpdateScTeacher(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScTeacherByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScTeacher
func (sc *ScTeacherController) DeleteScTeacher(c *gin.Context) {
	id := c.Param("id")

	teacher, err := sc.Service.GetScTeacherByID(id)
	if err != nil || teacher == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_teacher not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if teacher.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScTeacher(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_teacher deleted",
	})
}
