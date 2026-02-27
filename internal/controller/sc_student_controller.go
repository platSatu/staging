package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ScStudentController struct {
	Service *service.ScStudentService
}

func NewScStudentController(s *service.ScStudentService) *ScStudentController {
	return &ScStudentController{Service: s}
}

// CreateScStudent
// func (sc *ScStudentController) CreateScStudent(c *gin.Context) {
// 	var req struct {
// 		Name          *string `json:"name"`
// 		Address       *string `json:"address"`
// 		Tin           *string `json:"tin"`
// 		Tags          *string `json:"tags"`
// 		StudentType   *string `json:"student_type"`
// 		LcID          *string `json:"lc_id"`
// 		LevelID       *string `json:"level_id"`
// 		Phone         *string `json:"phone"`
// 		Mobile        *string `json:"mobile"`
// 		Email         *string `json:"email"`
// 		Language      *string `json:"language"`
// 		StudentStatus *string `json:"student_status"`
// 		PartnerType   *string `json:"partner_type"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	// ambil user_id dari token
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"success": false,
// 			"error":   "unauthorized",
// 		})
// 		return
// 	}

// 	student := &model.ScStudent{
// 		UserID:        userID.(string),
// 		Name:          req.Name,
// 		Address:       req.Address,
// 		Tin:           req.Tin,
// 		Tags:          req.Tags,
// 		StudentType:   req.StudentType,
// 		LcID:          req.LcID,
// 		LevelID:       req.LevelID,
// 		Phone:         req.Phone,
// 		Mobile:        req.Mobile,
// 		Email:         req.Email,
// 		Language:      req.Language,
// 		StudentStatus: req.StudentStatus,
// 		PartnerType:   req.PartnerType,
// 	}

// 	if err := sc.Service.CreateScStudent(student); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

//		c.JSON(http.StatusCreated, gin.H{
//			"success": true,
//			"data":    student,
//		})
//	}
//
// CreateScStudent
// CreateScStudent
func (sc *ScStudentController) CreateScStudent(c *gin.Context) {
	var req struct {
		// Data untuk user (hanya full_name, email, password)
		FullName string `json:"full_name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`

		// Data untuk student
		Name          *string `json:"name"`
		Address       *string `json:"address"`
		Tin           *string `json:"tin"`
		Tags          *string `json:"tags"`
		StudentType   *string `json:"student_type"`
		LcID          *string `json:"lc_id"`
		LevelID       *string `json:"level_id"`
		Phone         *string `json:"phone"`
		Mobile        *string `json:"mobile"`
		Language      *string `json:"language"`
		StudentStatus *string `json:"student_status"`
		PartnerType   *string `json:"partner_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil user_id dari token (parent/user yang login)
	parentID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	// Generate username dari email
	username := strings.Split(req.Email, "@")[0]

	// Buat objek user baru (child) - password akan di-hash di service
	user := &model.User{
		Username: username,
		Email:    req.Email,
		Password: req.Password, // Kirim plain password, akan di-hash di service
		FullName: req.FullName,
		Status:   "active",
		Role:     "user",
		// ParentID akan diset di service
	}

	// Buat objek student
	// Student akan di-assign ke user yang login (parent), bukan ke user baru
	student := &model.ScStudent{
		Name:          req.Name,
		Address:       req.Address,
		Tin:           req.Tin,
		Tags:          req.Tags,
		StudentType:   req.StudentType,
		LcID:          req.LcID,
		LevelID:       req.LevelID,
		Phone:         req.Phone,
		Mobile:        req.Mobile,
		Email:         &req.Email,
		Language:      req.Language,
		StudentStatus: req.StudentStatus,
		PartnerType:   req.PartnerType,
		// UserID akan diset ke parentID (user yang login) di service
	}

	// Panggil service untuk create user dan student
	newUser, newStudent, err := sc.Service.CreateUserAndStudent(parentID.(string), user, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data": gin.H{
			"user":    newUser,
			"student": newStudent,
		},
	})
}

// GetAllScStudent (hanya milik user login)
func (sc *ScStudentController) GetAllScStudent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	students, err := sc.Service.GetAllScStudentByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    students,
	})
}

// GetScStudentByID
func (sc *ScStudentController) GetScStudentByID(c *gin.Context) {
	id := c.Param("id")

	student, err := sc.Service.GetScStudentByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if student == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_student not found",
		})
		return
	}

	// hanya pemilik student boleh melihat
	userID, _ := c.Get("userID")
	if student.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    student,
	})
}

// UpdateScStudent
func (sc *ScStudentController) UpdateScStudent(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name          *string `json:"name"`
		Address       *string `json:"address"`
		Tin           *string `json:"tin"`
		Tags          *string `json:"tags"`
		StudentType   *string `json:"student_type"`
		LcID          *string `json:"lc_id"`
		LevelID       *string `json:"level_id"`
		Phone         *string `json:"phone"`
		Mobile        *string `json:"mobile"`
		Email         *string `json:"email"`
		Language      *string `json:"language"`
		StudentStatus *string `json:"student_status"`
		PartnerType   *string `json:"partner_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan student itu milik user login
	existing, err := sc.Service.GetScStudentByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_student not found",
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

	updateData := &model.ScStudent{
		ID:            id,
		Name:          req.Name,
		Address:       req.Address,
		Tin:           req.Tin,
		Tags:          req.Tags,
		StudentType:   req.StudentType,
		LcID:          req.LcID,
		LevelID:       req.LevelID,
		Phone:         req.Phone,
		Mobile:        req.Mobile,
		Email:         req.Email,
		Language:      req.Language,
		StudentStatus: req.StudentStatus,
		PartnerType:   req.PartnerType,
	}

	if err := sc.Service.UpdateScStudent(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScStudentByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScStudent
func (sc *ScStudentController) DeleteScStudent(c *gin.Context) {
	id := c.Param("id")

	student, err := sc.Service.GetScStudentByID(id)
	if err != nil || student == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_student not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if student.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScStudent(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_student deleted",
	})
}
