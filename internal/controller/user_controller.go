package controller

import (
	"backend_go/helper"
	"backend_go/internal/model"
	"backend_go/internal/request"
	"backend_go/internal/service"
	"net/http"
	"strings"

	//"strings" // <- ini untuk TrimPrefix

	"github.com/gin-gonic/gin"
	//"backend_go/helper"
)

type UserController struct {
	Service *service.UserService
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{Service: s}
}

// CreateUser
// Tambahkan import jika perlu: "backend_go/internal/request"

func (uc *UserController) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest // Bind ke struct request terpisah
	if err := c.ShouldBindJSON(&req); err != nil {
		println("[ERROR] Binding JSON failed:", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Debug: Pastikan semua field diterima setelah binding
	println("[DEBUG] Binding successful")
	println("[DEBUG] Received FullName:", req.FullName)
	println("[DEBUG] Received Email:", req.Email)
	println("[DEBUG] Received Password length:", len(req.Password))

	// Map ke model.User
	user := &model.User{
		FullName: req.FullName,
		Email:    req.Email,
		Password: req.Password,
	}

	if err := uc.Service.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

// GetAllUsers
func (uc *UserController) GetAllUsers(c *gin.Context) {
	users, err := uc.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    users,
	})
}

// GetUserByID
func (uc *UserController) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := uc.Service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "User not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// UpdateUser
func (uc *UserController) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var updateData model.User
	updateData.ID = id // Set ID dari param SEBELUM bind, untuk menghindari kosong

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	if err := uc.Service.UpdateUser(&updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Ambil data terbaru setelah update
	user, err := uc.Service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch updated user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// DeleteUser
func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := uc.Service.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted",
	})
}

// GetProfileTerakhir
// func (uc *UserController) GetProfile(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")
// 	if authHeader == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
// 		return
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	userID, err := helper.GetUserIDFromToken(tokenString)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	user, err := uc.Service.GetUserByID(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if user == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"success": true,
//			"data": gin.H{
//				"id":        user.ID,
//				"full_name": user.FullName,
//				"email":     user.Email,
//				"username":  user.Username,
//				"role":      user.Role,
//			},
//		})
//	}
// func (uc *UserController) GetProfile(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")
// 	if authHeader == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
// 		return
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	userID, err := helper.GetUserIDFromToken(tokenString)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	user, err := uc.Service.GetUserByID(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if user == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// Ambil ParentID dan akses paket / aplikasi yang di daftarkan
// 	parentID := user.ParentID
// 	akses, _ := uc.Service.GetAksesByUser(user.ID) // misal ambil daftar aplikasi/menu

//		c.JSON(http.StatusOK, gin.H{
//			"success": true,
//			"data": gin.H{
//				"id":        user.ID,
//				"full_name": user.FullName,
//				"email":     user.Email,
//				"username":  user.Username,
//				"role":      user.Role,
//				"parent_id": parentID,
//				"akses":     akses, // array berisi aplikasi, status, dan menu yang boleh diakses
//			},
//		})
//	}
//
// 19-02-2026
// func (uc *UserController) GetProfile(c *gin.Context) {
// 	authHeader := c.GetHeader("Authorization")
// 	if authHeader == "" {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
// 		return
// 	}

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
// 	userID, err := helper.GetUserIDFromToken(tokenString)
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
// 		return
// 	}

// 	user, aplikasiList, hasParent, err := uc.Service.GetUserProfileWithAplikasi(userID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	if user == nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
// 		return
// 	}

// 	// Ambil list aplikasi_id yang aktif
// 	var aplikasiIDs []string
// 	for _, app := range aplikasiList {
// 		if app["aplikasi_id"] != nil {
// 			aplikasiIDs = append(aplikasiIDs, app["aplikasi_id"].(string))
// 		}
// 	}

// 	// Ambil ParentID
// 	parentID := user.ParentID

//		c.JSON(http.StatusOK, gin.H{
//			"success": true,
//			"data": gin.H{
//				"id":           user.ID,
//				"full_name":    user.FullName,
//				"email":        user.Email,
//				"username":     user.Username,
//				"role":         user.Role,
//				"status":       user.Status,
//				"parent_id":    parentID,
//				"has_parent":   hasParent,
//				"aplikasi":     aplikasiList,
//				"aplikasi_ids": aplikasiIDs,
//			},
//		})
//	}
func (uc *UserController) GetProfile(c *gin.Context) {
	// 1. Validasi Authorization Header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Authorization header missing",
		})
		return
	}

	// 2. Ekstrak userID dari token
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	userID, err := helper.GetUserIDFromToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid token",
		})
		return
	}

	// 3. Panggil service untuk ambil data profil
	user, responseData, err := uc.Service.GetUserProfileWithAplikasi(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// 4. Ekstrak data untuk response
	packageInfo := responseData["package"]
	aplikasiList := responseData["aplikasi"].([]map[string]interface{})
	menuAccess := responseData["menu_access"].(map[string]interface{})

	// 5. Ambil list aplikasi_ids
	var aplikasiIDs []string
	for _, app := range aplikasiList {
		if app["aplikasi_id"] != nil {
			aplikasiIDs = append(aplikasiIDs, app["aplikasi_id"].(string))
		}
	}

	// 6. Response sukses
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile retrieved successfully",
		"data": gin.H{
			"id":             user.ID,
			"full_name":      user.FullName,
			"email":          user.Email,
			"username":       user.Username,
			"role":           user.Role,
			"status":         user.Status,
			"parent_id":      responseData["parent_id"],
			"has_parent":     responseData["has_parent"],
			"aplikasi":       aplikasiList,
			"aplikasi_ids":   aplikasiIDs,
			"package_status": packageInfo,
			"menu_access":    menuAccess,
		},
	})
}

// func (uc *UserController) GetProfile(c *gin.Context) {
// 	userAny, exists := c.Get("user")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	user := userAny.(model.User)

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"data": gin.H{
// 			"id":        user.ID,
// 			"full_name": user.FullName,
// 			"email":     user.Email,
// 			"username":  user.Username,
// 			"role":      user.Role, // 🔥 INI YANG FRONTEND BUTUH
// 		},
// 	})
// }

// GET /users/children
func (uc *UserController) GetChildren(c *gin.Context) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Unauthorized"})
		return
	}
	userID := userIDVal.(string)

	var users []model.User
	if err := uc.Service.DB.Where("parent_id = ?", userID).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}
