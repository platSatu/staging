package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CicilanUserController struct {
	Service *service.CicilanUserService
}

func NewCicilanUserController(s *service.CicilanUserService) *CicilanUserController {
	return &CicilanUserController{Service: s}
}

// CREATE — selalu pakai userID dari token
func (cc *CicilanUserController) CreateCicilanUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	var req model.CicilanUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	// override userID supaya tidak bisa curang membuat cicilan untuk user lain
	req.UserID = userID.(string)

	if err := cc.Service.CreateCicilanUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": req})
}

// GET ALL — hanya milik user login
func (cc *CicilanUserController) GetAllCicilanUser(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "unauthorized"})
		return
	}

	cicilans, err := cc.Service.GetAllCicilanUserByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": cicilans})
}

// GET BY ID — hanya pemilik boleh melihat
func (cc *CicilanUserController) GetCicilanUserByID(c *gin.Context) {
	id := c.Param("id")

	cicilan, err := cc.Service.GetCicilanUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}
	if cicilan == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "cicilan not found"})
		return
	}

	userID, _ := c.Get("userID")
	if cicilan.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "forbidden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": cicilan})
}

// UPDATE — hanya pemilik boleh edit
func (cc *CicilanUserController) UpdateCicilanUser(c *gin.Context) {
	id := c.Param("id")

	existing, err := cc.Service.GetCicilanUserByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "cicilan not found"})
		return
	}

	userID, _ := c.Get("userID")
	if existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "forbidden"})
		return
	}

	var req model.CicilanUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	req.ID = id                  // pastikan ID benar
	req.UserID = existing.UserID // cegah manipulasi

	if err := cc.Service.UpdateCicilanUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	updated, _ := cc.Service.GetCicilanUserByID(id)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": updated})
}

// DELETE — hanya pemilik boleh hapus
func (cc *CicilanUserController) DeleteCicilanUser(c *gin.Context) {
	id := c.Param("id")

	existing, err := cc.Service.GetCicilanUserByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "cicilan not found"})
		return
	}

	userID, _ := c.Get("userID")
	if existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"success": false, "error": "forbidden"})
		return
	}

	if err := cc.Service.DeleteCicilanUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "cicilan deleted"})
}

func (cc *CicilanUserController) GetParentList(c *gin.Context) {
	userIDRaw, _ := c.Get("userID")
	userID := userIDRaw.(string)

	result, err := cc.Service.GetParentSummary(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

func (cc *CicilanUserController) GetCicilanByParentID(c *gin.Context) {
	parentID := c.Param("parent_id")
	userIDRaw, _ := c.Get("userID")
	userID := userIDRaw.(string)

	result, err := cc.Service.GetCicilanByParentID(userID, parentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}
