package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ScPembelianBukuController struct {
	Service *service.ScPembelianBukuService
}

func NewScPembelianBukuController(s *service.ScPembelianBukuService) *ScPembelianBukuController {
	return &ScPembelianBukuController{Service: s}
}

// CreateScPembelianBuku
func (sc *ScPembelianBukuController) CreateScPembelianBuku(c *gin.Context) {
	var req struct {
		ParentID    *string `json:"parent_id"`
		Subject     *string `json:"subject"`
		PacesNumber *string `json:"paces_number"`
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

	pembelianBuku := &model.ScPembelianBuku{
		UserID:      userID.(string),
		ParentID:    req.ParentID,
		Subject:     req.Subject,
		PacesNumber: req.PacesNumber,
	}

	if err := sc.Service.CreateScPembelianBuku(pembelianBuku); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    pembelianBuku,
	})
}

// GetAllScPembelianBuku (hanya milik user login)
func (sc *ScPembelianBukuController) GetAllScPembelianBuku(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	pembelianBukus, err := sc.Service.GetAllScPembelianBukuByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pembelianBukus,
	})
}

// GetScPembelianBukuByID
func (sc *ScPembelianBukuController) GetScPembelianBukuByID(c *gin.Context) {
	id := c.Param("id")

	pembelianBuku, err := sc.Service.GetScPembelianBukuByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if pembelianBuku == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_pembelian_buku not found",
		})
		return
	}

	// Hanya pemilik pembelian buku boleh melihat
	userID, _ := c.Get("userID")
	if pembelianBuku.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pembelianBuku,
	})
}

// UpdateScPembelianBuku
func (sc *ScPembelianBukuController) UpdateScPembelianBuku(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		ParentID    *string `json:"parent_id"`
		Subject     *string `json:"subject"`
		PacesNumber *string `json:"paces_number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Pastikan pembelian buku itu milik user login
	existing, err := sc.Service.GetScPembelianBukuByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_pembelian_buku not found",
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

	updateData := &model.ScPembelianBuku{
		ID:          id,
		ParentID:    req.ParentID,
		Subject:     req.Subject,
		PacesNumber: req.PacesNumber,
	}

	if err := sc.Service.UpdateScPembelianBuku(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScPembelianBukuByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScPembelianBuku
func (sc *ScPembelianBukuController) DeleteScPembelianBuku(c *gin.Context) {
	id := c.Param("id")

	pembelianBuku, err := sc.Service.GetScPembelianBukuByID(id)
	if err != nil || pembelianBuku == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_pembelian_buku not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if pembelianBuku.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScPembelianBuku(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_pembelian_buku deleted",
	})
}
