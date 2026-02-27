// internal/controller/type_user_aplikasi_controller.go

package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TypeUserAplikasiController struct {
	Service *service.TypeUserAplikasiService
}

func NewTypeUserAplikasiController(s *service.TypeUserAplikasiService) *TypeUserAplikasiController {
	return &TypeUserAplikasiController{Service: s}
}

// CREATE - Membuat TypeUserAplikasi baru
func (c *TypeUserAplikasiController) CreateTypeUserAplikasi(ctx *gin.Context) {
	var req struct {
		UserID     string `json:"user_id" binding:"required"`
		ParentID   string `json:"parent_id" binding:"required"`
		AplikasiID string `json:"aplikasi_id" binding:"required"`
		Status     string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	tua := &model.TypeUserAplikasi{
		UserID:     req.UserID,
		ParentID:   req.ParentID,
		AplikasiID: req.AplikasiID,
		Status:     req.Status,
	}

	if err := c.Service.CreateTypeUserAplikasi(tua); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    tua,
	})
}

// READ ALL - Ambil semua TypeUserAplikasi
func (c *TypeUserAplikasiController) GetAllTypeUserAplikasi(ctx *gin.Context) {
	data, err := c.Service.GetAllTypeUserAplikasi()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// READ BY ID - Ambil TypeUserAplikasi by ID
func (c *TypeUserAplikasiController) GetTypeUserAplikasiByID(ctx *gin.Context) {
	id := ctx.Param("id")

	tua, err := c.Service.GetTypeUserAplikasiByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if tua == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "type user aplikasi not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    tua,
	})
}

// READ BY USER ID - Ambil TypeUserAplikasi by UserID
func (c *TypeUserAplikasiController) GetTypeUserAplikasiByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	data, err := c.Service.GetTypeUserAplikasiByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// UPDATE - Update TypeUserAplikasi
func (c *TypeUserAplikasiController) UpdateTypeUserAplikasi(ctx *gin.Context) {
	id := ctx.Param("id")

	var req struct {
		UserID     string `json:"user_id"`
		ParentID   string `json:"parent_id"`
		AplikasiID string `json:"aplikasi_id"`
		Status     string `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	tua := &model.TypeUserAplikasi{
		ID:         id,
		UserID:     req.UserID,
		ParentID:   req.ParentID,
		AplikasiID: req.AplikasiID,
		Status:     req.Status,
	}

	if err := c.Service.UpdateTypeUserAplikasi(tua); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := c.Service.GetTypeUserAplikasiByID(id)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DELETE - Hapus TypeUserAplikasi
func (c *TypeUserAplikasiController) DeleteTypeUserAplikasi(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := c.Service.DeleteTypeUserAplikasi(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "type user aplikasi deleted",
	})
}
