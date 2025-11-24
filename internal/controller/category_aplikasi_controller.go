package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type CategoryAplikasiController struct {
	Service *service.CategoryAplikasiService
}

func NewCategoryAplikasiController(s *service.CategoryAplikasiService) *CategoryAplikasiController {
	return &CategoryAplikasiController{Service: s}
}

// CreateCategoryAplikasi
func (cc *CategoryAplikasiController) CreateCategoryAplikasi(c *gin.Context) {
	var category model.CategoryAplikasi
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := cc.Service.CreateCategoryAplikasi(&category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    category,
	})
}

// GetAllCategoryAplikasi
func (cc *CategoryAplikasiController) GetAllCategoryAplikasi(c *gin.Context) {
	categories, err := cc.Service.GetAllCategoryAplikasi()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}

// GetCategoryAplikasiByID
func (cc *CategoryAplikasiController) GetCategoryAplikasiByID(c *gin.Context) {
	id := c.Param("id")
	category, err := cc.Service.GetCategoryAplikasiByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if category == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Category not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
	})
}

// GetCategoryAplikasiByUser
func (cc *CategoryAplikasiController) GetCategoryAplikasiByUser(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := helper.GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	categories, err := cc.Service.GetCategoryAplikasiByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}

// UpdateCategoryAplikasi
func (cc *CategoryAplikasiController) UpdateCategoryAplikasi(c *gin.Context) {
	id := c.Param("id")
	var updateData model.CategoryAplikasi
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := cc.Service.UpdateCategoryAplikasi(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	category, err := cc.Service.GetCategoryAplikasiByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated category",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    category,
	})
}

// DeleteCategoryAplikasi
func (cc *CategoryAplikasiController) DeleteCategoryAplikasi(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteCategoryAplikasi(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Category deleted",
	})
}
