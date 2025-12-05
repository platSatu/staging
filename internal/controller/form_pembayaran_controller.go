package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FormPembayaranController struct {
	Service *service.FormPembayaranService
}

func NewFormPembayaranController(s *service.FormPembayaranService) *FormPembayaranController {
	return &FormPembayaranController{Service: s}
}

// CreateFormPembayaran
func (fc *FormPembayaranController) CreateFormPembayaran(c *gin.Context) {
	var form model.FormPembayaran
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := fc.Service.CreateFormPembayaran(&form); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    form,
	})
}

// GetAllFormPembayaran
func (fc *FormPembayaranController) GetAllFormPembayaran(c *gin.Context) {
	forms, err := fc.Service.GetAllFormPembayaran()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    forms,
	})
}

// GetFormPembayaranByID
func (fc *FormPembayaranController) GetFormPembayaranByID(c *gin.Context) {
	id := c.Param("id")
	form, err := fc.Service.GetFormPembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if form == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Form pembayaran not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    form,
	})
}

// UpdateFormPembayaran
func (fc *FormPembayaranController) UpdateFormPembayaran(c *gin.Context) {
	id := c.Param("id")
	var updateData model.FormPembayaran
	updateData.ID = id // Set ID dari param

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := fc.Service.UpdateFormPembayaran(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil data terbaru setelah update
	form, err := fc.Service.GetFormPembayaranByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated form pembayaran",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    form,
	})
}

// DeleteFormPembayaran
func (fc *FormPembayaranController) DeleteFormPembayaran(c *gin.Context) {
	id := c.Param("id")
	if err := fc.Service.DeleteFormPembayaran(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Form pembayaran deleted",
	})
}
