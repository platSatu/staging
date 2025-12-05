package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KategoriPembayaranController struct {
	Service *service.KategoriPembayaranService
}

func NewKategoriPembayaranController(s *service.KategoriPembayaranService) *KategoriPembayaranController {
	return &KategoriPembayaranController{Service: s}
}

// CreateKategoriPembayaran
func (kc *KategoriPembayaranController) CreateKategoriPembayaran(c *gin.Context) {
	var req struct {
		NamaKategori string  `json:"nama_kategori" binding:"required"`
		Deskripsi    *string `json:"deskripsi"`
		Status       string  `json:"status"` // optional: active / inactive
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

	kategori := &model.KategoriPembayaran{
		UserID:       userID.(string),
		NamaKategori: req.NamaKategori,
		Deskripsi:    req.Deskripsi,
		Status:       req.Status,
	}

	if err := kc.Service.CreateKategoriPembayaran(kategori); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    kategori,
	})
}

// GetAllKategoriPembayaran (hanya milik user login)
func (kc *KategoriPembayaranController) GetAllKategoriPembayaran(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	kategoris, err := kc.Service.GetAllKategoriPembayaranByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    kategoris,
	})
}

// GetKategoriPembayaranByID
func (kc *KategoriPembayaranController) GetKategoriPembayaranByID(c *gin.Context) {
	id := c.Param("id")

	kategori, err := kc.Service.GetKategoriPembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if kategori == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "kategori pembayaran not found",
		})
		return
	}

	// hanya pemilik kategori boleh melihat
	userID, _ := c.Get("userID")
	if kategori.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    kategori,
	})
}

// UpdateKategoriPembayaran
func (kc *KategoriPembayaranController) UpdateKategoriPembayaran(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		NamaKategori string  `json:"nama_kategori"`
		Deskripsi    *string `json:"deskripsi"`
		Status       string  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// pastikan kategori itu milik user login
	existing, err := kc.Service.GetKategoriPembayaranByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "kategori pembayaran not found",
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

	updateData := &model.KategoriPembayaran{
		ID:           id,
		NamaKategori: req.NamaKategori,
		Deskripsi:    req.Deskripsi,
		Status:       req.Status,
	}

	if err := kc.Service.UpdateKategoriPembayaran(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := kc.Service.GetKategoriPembayaranByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteKategoriPembayaran
func (kc *KategoriPembayaranController) DeleteKategoriPembayaran(c *gin.Context) {
	id := c.Param("id")

	kategori, err := kc.Service.GetKategoriPembayaranByID(id)
	if err != nil || kategori == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "kategori pembayaran not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if kategori.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := kc.Service.DeleteKategoriPembayaran(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "kategori pembayaran deleted",
	})
}
