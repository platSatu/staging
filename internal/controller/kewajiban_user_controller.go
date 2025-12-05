package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type KewajibanUserController struct {
	Service *service.KewajibanUserService
}

func NewKewajibanUserController(s *service.KewajibanUserService) *KewajibanUserController {
	return &KewajibanUserController{Service: s}
}

// ===============================
// CreateKewajibanUser
// ===============================
func (kc *KewajibanUserController) CreateKewajibanUser(c *gin.Context) {
	var kewajiban model.KewajibanUser

	if err := c.ShouldBindJSON(&kewajiban); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil user ID dari JWT
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	// pastikan userID tidak bisa dipalsukan
	kewajiban.UserID = userID.(string)

	if err := kc.Service.CreateKewajibanUser(&kewajiban); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": kewajiban})
}

// ===============================
// GetAllKewajibanUser (HANYA MILIK USER LOGIN)
// ===============================
func (kc *KewajibanUserController) GetAllKewajibanUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	kewajibans, err := kc.Service.GetAllKewajibanUserByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": kewajibans})
}

// ===============================
// GetKewajibanUserByID (HANYA BOLEH AMBIL MILIK SENDIRI)
// ===============================
func (kc *KewajibanUserController) GetKewajibanUserByID(c *gin.Context) {
	id := c.Param("id")

	kewajiban, err := kc.Service.GetKewajibanUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if kewajiban == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Kewajiban user not found"})
		return
	}

	userID, _ := c.Get("userID")
	if kewajiban.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": kewajiban})
}

// ===============================
// UpdateKewajibanUser
// ===============================
func (kc *KewajibanUserController) UpdateKewajibanUser(c *gin.Context) {
	id := c.Param("id")

	// cek apakah data ada
	existing, err := kc.Service.GetKewajibanUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Kewajiban user not found"})
		return
	}

	// cek user pemilik data
	userID, _ := c.Get("userID")
	if existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	// proses update
	var updateData model.KewajibanUser
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	updateData.ID = id
	updateData.UserID = existing.UserID // lock user id

	if err := kc.Service.UpdateKewajibanUser(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	newData, _ := kc.Service.GetKewajibanUserByID(id)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": newData})
}

// ===============================
// DeleteKewajibanUser
// ===============================
func (kc *KewajibanUserController) DeleteKewajibanUser(c *gin.Context) {
	id := c.Param("id")

	kewajiban, err := kc.Service.GetKewajibanUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if kewajiban == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Kewajiban user not found"})
		return
	}

	userID, _ := c.Get("userID")
	if kewajiban.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := kc.Service.DeleteKewajibanUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Deleted successfully"})
}
