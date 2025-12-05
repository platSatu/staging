package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type CategoryPackagesController struct {
	Service *service.CategoryPackagesService
}

func NewCategoryPackagesController(s *service.CategoryPackagesService) *CategoryPackagesController {
	return &CategoryPackagesController{Service: s}
}

// CreateCategoryPackages
func (cc *CategoryPackagesController) CreateCategoryPackages(c *gin.Context) {
	var category model.CategoryPackages
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := cc.Service.CreateCategoryPackages(&category); err != nil {
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

// GetAllCategoryPackages
func (cc *CategoryPackagesController) GetAllCategoryPackages(c *gin.Context) {
	categories, err := cc.Service.GetAllCategoryPackages()
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

// GetCategoryPackagesByID
func (cc *CategoryPackagesController) GetCategoryPackagesByID(c *gin.Context) {
	id := c.Param("id")
	category, err := cc.Service.GetCategoryPackagesByID(id)
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

// GetCategoryPackagesByUser
func (cc *CategoryPackagesController) GetCategoryPackagesByUser(c *gin.Context) {
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

	categories, err := cc.Service.GetCategoryPackagesByUserID(userID)
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

// UpdateCategoryPackages
func (cc *CategoryPackagesController) UpdateCategoryPackages(c *gin.Context) {
	id := c.Param("id")
	var updateData model.CategoryPackages
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := cc.Service.UpdateCategoryPackages(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	category, err := cc.Service.GetCategoryPackagesByID(id)
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

// DeleteCategoryPackages
func (cc *CategoryPackagesController) DeleteCategoryPackages(c *gin.Context) {
	id := c.Param("id")
	if err := cc.Service.DeleteCategoryPackages(id); err != nil {
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
