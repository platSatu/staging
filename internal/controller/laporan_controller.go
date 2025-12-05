package controller

import (
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LaporanController struct {
	Service *service.LaporanService
}

func NewLaporanController(s *service.LaporanService) *LaporanController {
	return &LaporanController{Service: s}
}

func (lc *LaporanController) GetLaporanPerKategori(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	result, err := lc.Service.GetLaporanPerKategori(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

func (lc *LaporanController) GetDetailPerForm(c *gin.Context) {
	kewajibanID := c.Param("id") // <-- harus kewajiban_id
	parentIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}

	result, err := lc.Service.GetDetailPerForm(kewajibanID, parentIDVal.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}
