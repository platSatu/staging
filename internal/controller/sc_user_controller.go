// internal/controller/sc_user_controller.go

package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SCUserController struct {
	Service *service.SCUserService
}

func NewSCUserController(s *service.SCUserService) *SCUserController {
	return &SCUserController{Service: s}
}

// CreateSCUser - Membuat SC user baru
func (c *SCUserController) CreateSCUser(ctx *gin.Context) {
	// Ambil userID dari token (middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	var req struct {
		Email       string   `json:"email" binding:"required"`
		Password    string   `json:"password" binding:"required"`
		FullName    string   `json:"full_name" binding:"required"`
		Username    string   `json:"username"`
		Role        string   `json:"role"`
		KodeReferal *string  `json:"kode_referal"`
		Saldo       *float64 `json:"saldo"`
		AplikasiID  string   `json:"aplikasi_id"` // Untuk role student
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	scUser := &model.User{
		Email:       req.Email,
		Password:    req.Password,
		FullName:    req.FullName,
		Username:    req.Username,
		Role:        req.Role,
		KodeReferal: req.KodeReferal,
		Saldo:       req.Saldo,
	}

	// Panggil service dengan parentID dan aplikasiID
	if err := c.Service.CreateSCUser(scUser, userID.(string), req.AplikasiID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    scUser,
	})
}

// GetAllSCUsers - Ambil semua SC user
func (c *SCUserController) GetAllSCUsers(ctx *gin.Context) {
	users, err := c.Service.GetAllSCUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// GetSCUserByID - Ambil SC user by ID
func (c *SCUserController) GetSCUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.Service.GetSCUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc user not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// GetSCUsersByParentID - Ambil semua SC user dari parent tertentu
func (c *SCUserController) GetSCUsersByParentID(ctx *gin.Context) {
	// Ambil userID dari token (middleware)
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	users, err := c.Service.GetSCUsersByParentID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// UpdateSCUser - Update SC user
func (c *SCUserController) UpdateSCUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		FullName    string   `json:"full_name"`
		Username    string   `json:"username"`
		Password    string   `json:"password"`
		Status      string   `json:"status"`
		Role        string   `json:"role"`
		KodeReferal *string  `json:"kode_referal"`
		Saldo       *float64 `json:"saldo"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	scUser := &model.User{
		ID:          id,
		FullName:    req.FullName,
		Username:    req.Username,
		Password:    req.Password,
		Status:      req.Status,
		Role:        req.Role,
		KodeReferal: req.KodeReferal,
		Saldo:       req.Saldo,
	}

	if err := c.Service.UpdateSCUser(scUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := c.Service.GetSCUserByID(id)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteSCUser - Hapus SC user
func (c *SCUserController) DeleteSCUser(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.Service.DeleteSCUser(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc user deleted",
	})
}

// internal/controller/sc_user_controller.go

// GetMySCUsers - Ambil semua SC user yang dibuat oleh user yang login
func (c *SCUserController) GetMySCUsers(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	users, err := c.Service.GetSCUsersByParentID(userID.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}
