package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type PackagesController struct {
	Service *service.PackagesService
}

func NewPackagesController(s *service.PackagesService) *PackagesController {
	return &PackagesController{Service: s}
}

// CreatePackages
func (pc *PackagesController) CreatePackages(c *gin.Context) {
	var pkg model.Packages
	if err := c.ShouldBindJSON(&pkg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := pc.Service.CreatePackages(&pkg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    pkg,
	})
}

// GetAllPackages
func (pc *PackagesController) GetAllPackages(c *gin.Context) {
	packages, err := pc.Service.GetAllPackages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    packages,
	})
}

// GetPackagesByID
func (pc *PackagesController) GetPackagesByID(c *gin.Context) {
	id := c.Param("id")
	pkg, err := pc.Service.GetPackagesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if pkg == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Packages not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pkg,
	})
}

// GetPackagesByUser
func (pc *PackagesController) GetPackagesByUser(c *gin.Context) {
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

	packages, err := pc.Service.GetPackagesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    packages,
	})
}

// UpdatePackages
func (pc *PackagesController) UpdatePackages(c *gin.Context) {
	id := c.Param("id")
	var updateData model.Packages
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := pc.Service.UpdatePackages(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	pkg, err := pc.Service.GetPackagesByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated packages",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    pkg,
	})
}

// DeletePackages
func (pc *PackagesController) DeletePackages(c *gin.Context) {
	id := c.Param("id")
	if err := pc.Service.DeletePackages(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Packages deleted",
	})
}
