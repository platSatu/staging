package controller

import (
	"backend_go/internal/service" // Ganti dengan path import service Anda
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

type RegistrasiUlangController struct {
	Service *service.RegistrasiUlangService
}

func NewRegistrasiUlangController(svc *service.RegistrasiUlangService) *RegistrasiUlangController {
	return &RegistrasiUlangController{Service: svc}
}

func (c *RegistrasiUlangController) RegistrasiUlang(ctx *gin.Context) {
	var req struct {
		Qrcode string `json:"qrcode" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "qrcode wajib diisi"})
		return
	}

	// Validasi: qrcode hanya boleh berupa angka
	if !regexp.MustCompile(`^[0-9]+$`).MatchString(req.Qrcode) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "qrcode hanya boleh berupa angka"})
		return
	}

	// Dapatkan detail scanner dari context (sesuaikan dengan implementasi auth Anda)
	scannedByUser := ctx.GetString("userID")       // Misalnya dari JWT middleware (ganti jika key berbeda)
	scannedByDevice := ctx.GetHeader("User-Agent") // Device dari header User-Agent
	ipAddress := ctx.ClientIP()                    // IP address dari client
	browser := ctx.GetHeader("User-Agent")         // Browser dari header User-Agent (atau parse lebih detail jika perlu)

	// Panggil service dengan parameter tambahan
	kategori, err := c.Service.RegistrasiUlang(req.Qrcode, scannedByUser, scannedByDevice, ipAddress, browser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Registrasi ulang berhasil", "kategori": kategori})
}

// Method baru: GetAllTicketKategori - Mengambil dan mengembalikan seluruh data ticket kategori
func (c *RegistrasiUlangController) GetAllTicketKategori(ctx *gin.Context) {
	categories, err := c.Service.GetAllTicketKategori()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data ticket kategori"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": categories})
}
