package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AturanDendaController struct {
	Service *service.AturanDendaService
}

func NewAturanDendaController(s *service.AturanDendaService) *AturanDendaController {
	return &AturanDendaController{Service: s}
}

// ====================== CREATE ======================
func (ac *AturanDendaController) CreateAturanDenda(c *gin.Context) {
	var req model.AturanDenda

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// Ambil user yang sedang login
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	req.UserID = userID.(string) // Set otomatis

	if err := ac.Service.CreateAturanDenda(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": req})
}

// ====================== GET ALL ======================
func (ac *AturanDendaController) GetAllAturanDenda(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	list, err := ac.Service.GetAllAturanDendaByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// ====================== GET BY ID ======================
func (ac *AturanDendaController) GetAturanDendaByID(c *gin.Context) {
	id := c.Param("id")

	aturan, err := ac.Service.GetAturanDendaByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	if aturan == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Aturan denda not found"})
		return
	}

	// Validasi kepemilikan data
	userID, exists := c.Get("userID")
	if !exists || aturan.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": aturan})
}

// ====================== UPDATE ======================
func (ac *AturanDendaController) UpdateAturanDenda(c *gin.Context) {
	id := c.Param("id")

	// Ambil data lama
	existing, err := ac.Service.GetAturanDendaByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Aturan denda not found"})
		return
	}

	// Validasi user pemilik
	userID, exists := c.Get("userID")
	if !exists || existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	var updateData model.AturanDenda
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	updateData.ID = id
	updateData.UserID = existing.UserID // jaga agar user tidak bisa mengubah user_id

	if err := ac.Service.UpdateAturanDenda(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	updated, _ := ac.Service.GetAturanDendaByID(id)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
}

// ====================== DELETE ======================
func (ac *AturanDendaController) DeleteAturanDenda(c *gin.Context) {
	id := c.Param("id")

	aturan, err := ac.Service.GetAturanDendaByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if aturan == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Aturan denda not found"})
		return
	}

	// Validasi kepemilikan
	userID, exists := c.Get("userID")
	if !exists || aturan.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "Forbidden"})
		return
	}

	if err := ac.Service.DeleteAturanDenda(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Aturan denda deleted"})
}
