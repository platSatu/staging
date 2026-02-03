package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"fmt"
	"net/http"
	"strings"

	"backend_go/helper"

	"github.com/gin-gonic/gin"
)

type ProfileController struct {
	Service *service.ProfileService
}

func NewProfileController(s *service.ProfileService) *ProfileController {
	return &ProfileController{Service: s}
}

// Helper untuk ekstrak userID dari token (tanpa role)
func (pc *ProfileController) getUserFromToken(c *gin.Context) (userID string, err error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err = helper.GetUserIDFromToken(tokenString)
	if err != nil {
		return "", err
	}
	return userID, nil
}

// CreateProfile
func (pc *ProfileController) CreateProfile(c *gin.Context) {
	userID, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID from token"})
		return
	}

	var profile model.Profile
	if err := c.ShouldBindJSON(&profile); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := pc.Service.CreateProfile(&profile, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": profile})
}

// GetAllProfiles - Sekarang alias untuk GetProfilesByUser (user lihat profilnya sendiri)
func (pc *ProfileController) GetAllProfiles(c *gin.Context) {
	userID, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	profiles, err := pc.Service.GetAllProfiles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profiles})
}

// GetProfileByID
func (pc *ProfileController) GetProfileByID(c *gin.Context) {
	userID, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile ID is required"})
		return
	}

	profile, err := pc.Service.GetProfileByID(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if profile == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Profile not found or access denied"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// GetProfilesByUser
// Di GetProfilesByUser
func (pc *ProfileController) GetProfilesByUser(c *gin.Context) {
	userID, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID from token"})
		return
	}

	// Cek apakah user sudah punya profile
	hasProfile, err := pc.Service.CheckUserHasProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if !hasProfile {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Data not found: User belum membuat profile"})
		return
	}

	profiles, err := pc.Service.GetProfilesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": profiles})
}

// UpdateProfile
func (pc *ProfileController) UpdateProfile(c *gin.Context) {
	userID, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile ID is required"})
		return
	}

	var updateData model.Profile
	updateData.ID = id

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	if err := pc.Service.UpdateProfile(&updateData, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	profile, err := pc.Service.GetProfileByID(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to fetch updated profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": profile})
}

// DeleteProfile
func (pc *ProfileController) DeleteProfile(c *gin.Context) {
	userID, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "profile ID is required"})
		return
	}

	if err := pc.Service.DeleteProfile(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Profile deleted"})
}

// Tambahkan di akhir file, sebelum export
// CheckUserHasProfile - Mengecek apakah user sudah punya profile
func (pc *ProfileController) CheckUserHasProfile(c *gin.Context) {
	userIDFromQuery := c.Query("user_id")
	if userIDFromQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "user_id query parameter is required"})
		return
	}

	// Ekstrak userID dari token untuk verifikasi (pastikan hanya user sendiri yang bisa cek)
	userIDFromToken, err := pc.getUserFromToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if userIDFromToken != userIDFromQuery {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Access denied: You can only check your own profile"})
		return
	}

	hasProfile, err := pc.Service.CheckUserHasProfile(userIDFromQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "hasProfile": hasProfile})
}
