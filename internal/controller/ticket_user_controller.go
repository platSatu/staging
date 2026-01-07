package controller

import (
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketUserController struct {
	ticketUserService *service.TicketUserService
}

func NewTicketUserController(ticketUserService *service.TicketUserService) *TicketUserController {
	return &TicketUserController{ticketUserService: ticketUserService}
}

// GetMyTicketQrcodes menampilkan ticket qrcodes milik user yang login
func (ctrl *TicketUserController) GetMyTicketQrcodes(c *gin.Context) {
	// Ambil userID dari context (di-set oleh AuthMiddleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	// Panggil service untuk mendapatkan data
	tickets, err := ctrl.ticketUserService.GetTicketQrcodesByUserID(userIDStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tickets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
