package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"strings"

	"github.com/gin-gonic/gin"

	"backend_go/helper"
)

type ProfileController struct {
	Service *service.ProfileService
}

func NewProfileController(s *service.ProfileService) *ProfileController {
	return &ProfileController{Service: s}
}

// CreateProfile
func (pc *ProfileController) CreateProfile(c *gin.Context) {
	var profile model.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := pc.Service.CreateProfile(&profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    profile,
	})
}

// GetAllProfiles
func (pc *ProfileController) GetAllProfiles(c *gin.Context) {
	profiles, err := pc.Service.GetAllProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profiles,
	})
}

// GetProfileByID
func (pc *ProfileController) GetProfileByID(c *gin.Context) {
	id := c.Param("id")
	profile, err := pc.Service.GetProfileByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Profile not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// GetProfilesByUser
func (pc *ProfileController) GetProfilesByUser(c *gin.Context) {
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

	profiles, err := pc.Service.GetProfilesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profiles,
	})
}

// UpdateProfile
func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	id := c.Param("id")
	var updateData model.Profile
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := pc.Service.UpdateProfile(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	profile, err := pc.Service.GetProfileByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    profile,
	})
}

// DeleteProfile
func (pc *ProfileController) DeleteProfile(c *gin.Context) {
	id := c.Param("id")
	if err := pc.Service.DeleteProfile(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile deleted",
	})
}
