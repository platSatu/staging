package service

import (
	"backend_go/internal/model"
	"backend_go/internal/request"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB          *gorm.DB
	UserService *UserService
}

func NewAuthService(db *gorm.DB, userService *UserService) *AuthService {
	return &AuthService{
		DB:          db,
		UserService: userService,
	}
}

var accessTokenTTL = time.Minute * 15
var refreshTokenTTL = time.Hour * 24 * 7 // 7 hari

// ================= Register =================
func (s *AuthService) Register(req request.RegisterRequest) (*model.User, error) {
	req.Password = strings.TrimSpace(req.Password)
	if req.Password == "" {
		return nil, errors.New("password tidak boleh kosong")
	}

	// Cek apakah email sudah terdaftar
	var existingUser model.User
	if err := s.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// Kirim password plain ke UserService, CreateUser akan handle hashing
	user := &model.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password, // Kirim sebagai plain, akan di-hash di CreateUser
	}

	if err := s.UserService.CreateUser(user); err != nil {
		return nil, err
	}

	println("[DEBUG] User registered:", user.Email)
	return user, nil
}

// ================= JWT =================
func (s *AuthService) generateAccessToken(user model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "b6c0f8a23e9f4d7e8b3a2d1f0c9e4a7d9b0c1e2f3d4a5b6c7d8e9f0a1b2c3d4e" // Fallback dari .env Anda
		println("[WARNING] JWT_SECRET not set, using default from .env")
	}
	if len(secret) < 32 {
		return "", errors.New("JWT_SECRET harus minimal 256-bit / 32 karakter")
	}

	println("[DEBUG] Auth Service: Using JWT secret:", secret[:10]+"...") // Log awal secret untuk debug

	// Claims minimalis: hanya sub (userID) dan exp, tanpa email atau info sensitive
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(accessTokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ================= Refresh Token =================
func (s *AuthService) generateRefreshToken(user model.User) (string, time.Time, error) {
	token := uuid.New().String() + "." + uuid.New().String()
	expires := time.Now().Add(refreshTokenTTL)

	refresh := model.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expires,
		Revoked:   false,
	}

	if err := s.DB.Create(&refresh).Error; err != nil {
		return "", time.Time{}, err
	}

	return token, expires, nil
}

// ================= Login =================
func (s *AuthService) Login(req request.LoginRequest) (map[string]interface{}, error) {
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" {
		return nil, errors.New("email tidak boleh kosong")
	}
	if req.Password == "" {
		return nil, errors.New("password tidak boleh kosong")
	}

	println("[DEBUG] Login attempt for email:", req.Email)

	var user model.User
	if err := s.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("email tidak terdaftar")
	}

	// Jika password hash kosong → error
	if len(user.Password) == 0 {
		return nil, errors.New("password pada server tidak valid")
	}

	// Validasi password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("password salah")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, errors.New("gagal membuat access token")
	}

	refreshToken, expiresAt, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, errors.New("gagal membuat refresh token")
	}

	// Sukses
	return map[string]interface{}{
		"success":       true,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_at":    expiresAt,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (map[string]interface{}, error) {
	var rt model.RefreshToken
	if err := s.DB.Where("token = ? AND revoked = ?", refreshToken, false).First(&rt).Error; err != nil {
		return nil, errors.New("refresh token tidak valid")
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	var user model.User
	if err := s.DB.First(&user, "id = ?", rt.UserID).Error; err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"access_token": accessToken,
	}, nil
}

// ================= Logout =================
func (s *AuthService) Logout(refreshToken string) error {
	return s.DB.Model(&model.RefreshToken{}).
		Where("token = ? AND revoked = ?", refreshToken, false).
		Update("revoked", true).Error
}
