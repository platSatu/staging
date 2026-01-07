package controller

import (
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PurchaseController struct {
	Service *service.PurchaseService
}

func NewPurchaseController(s *service.PurchaseService) *PurchaseController {
	return &PurchaseController{Service: s}
}

// ProcessPurchase: Endpoint POST /ticket-purchases
// Sekarang menggunakan autentikasi login: user harus login terlebih dahulu, tidak perlu register/token lagi
func (pc *PurchaseController) ProcessPurchase(c *gin.Context) {
	var req service.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Request body tidak valid",
		})
		return
	}

	// Ambil userID dari context (asumsi dari auth middleware setelah login)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "Silakan login terlebih dahulu.",
		})
		return
	}

	// Proses pembelian via PurchaseService dengan userID dari login
	if err := pc.Service.ProcessPurchase(&req, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pembelian berhasil diproses",
	})
}

// GetPurchaseByToken: Endpoint ini mungkin tidak diperlukan lagi karena logika berubah ke login
// Jika masih ingin digunakan untuk skenario lain, bisa dipertahankan atau dihapus
// Untuk konsistensi, saya hapus karena logika baru tidak menggunakan token
// func (pc *PurchaseController) GetPurchaseByToken(c *gin.Context) {
// 	token := c.Query("token") // Ambil token dari query parameter
// 	if token == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   "Token diperlukan",
// 		})
// 		return
// 	}

// 	purchases, err := pc.Service.GetPurchaseByToken(token)
// 	if err != nil {
// 		if err.Error() == "purchase tidak ditemukan" {
// 			c.JSON(http.StatusNotFound, gin.H{
// 				"success": false,
// 				"error":   "Purchase tidak ditemukan",
// 			})
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data":    purchases,
// 	})
// }
