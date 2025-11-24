package controller

import (
	"backend_go/internal/request"
	"backend_go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService *service.AuthService
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req request.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	user, err := c.AuthService.Register(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req request.LoginRequest

	// Bind request
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "format request tidak valid",
		})
		return
	}

	// Panggil service login
	res, err := c.AuthService.Login(req)
	if err != nil {
		// Login gagal → jangan buat cookie sama sekali
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil refresh token dari service response
	refreshToken, ok := res["refresh_token"].(string)
	if !ok || refreshToken == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "refresh token tidak tersedia"})
		return
	}

	// Set cookie refresh token 7 hari
	cookieMaxAge := 7 * 24 * 60 * 60
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		cookieMaxAge,
		"/",
		"localhost", // ganti dengan domain production jika perlu
		false,
		true, // httpOnly
	)

	// Kirim access token ke frontend (bisa dipakai sementara untuk header Authorization)
	accessToken, ok := res["access_token"].(string)
	if !ok || accessToken == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "access token tidak tersedia"})
		return
	}

	// Response sukses
	ctx.JSON(http.StatusOK, gin.H{
		"success":      true,
		"access_token": accessToken,
	})
}

// Refresh token endpoint, baca cookie
func (c *AuthController) Refresh(ctx *gin.Context) {
	var req RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil || req.RefreshToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token tidak ditemukan"})
		return
	}

	res, err := c.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success":      true,
		"access_token": res["access_token"],
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	var req request.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	err := c.AuthService.Logout(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout berhasil"})
}
