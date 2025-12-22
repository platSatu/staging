// package controller

// import (
// 	"backend_go/internal/service"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// )

// type PurchaseController struct {
// 	Service *service.PurchaseService
// }

// func NewPurchaseController(s *service.PurchaseService) *PurchaseController {
// 	return &PurchaseController{Service: s}
// }

// // ProcessPurchase: Endpoint POST /ticket-purchases
// func (pc *PurchaseController) ProcessPurchase(c *gin.Context) {
// 	var req service.PurchaseRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

// 	// Ambil userID dari context (asumsi dari auth middleware)
// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"success": false,
// 			"error":   "User not authenticated",
// 		})
// 		return
// 	}

// 	if err := pc.Service.ProcessPurchase(&req, userID.(string)); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 		})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"success": true,
//			"message": "Pembelian berhasil diproses",
//		})
//	}
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
// Mendukung skenario public: user register, dapat token, klik link pembelian
func (pc *PurchaseController) ProcessPurchase(c *gin.Context) {
	var req struct {
		service.PurchaseRequest
		Token string `json:"token"` // token bisa dikirim via JSON body
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Request body tidak valid",
		})
		return
	}

	// Ambil token: prioritas dari body, fallback ke query parameter
	token := req.Token
	if token == "" {
		token = c.Query("token")
	}

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token pembelian wajib diisi",
		})
		return
	}

	// Validasi token via TicketRegisterService
	ticketRegisterService := service.NewTicketRegisterService(pc.Service.DB)
	ticketRegister, err := ticketRegisterService.ValidatePurchaseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil userID dari ticket register public
	userID := ticketRegister.UserID
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "User tidak valid untuk token ini",
		})
		return
	}

	// Proses pembelian via PurchaseService
	if err := pc.Service.ProcessPurchase(&req.PurchaseRequest, userID); err != nil {
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

// GetPurchaseByToken: Endpoint GET /purchase?token=...
func (pc *PurchaseController) GetPurchaseByToken(c *gin.Context) {
	token := c.Query("token") // Ambil token dari query parameter
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Token diperlukan",
		})
		return
	}

	purchases, err := pc.Service.GetPurchaseByToken(token)
	if err != nil {
		if err.Error() == "purchase tidak ditemukan" {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Purchase tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    purchases,
	})
}
