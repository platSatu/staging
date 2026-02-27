package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	rand.Seed(time.Now().UnixNano())
	return &UserService{DB: db}
}

// CREATE
func (s *UserService) CreateUser(user *model.User) error {
	log.Printf("[DEBUG] Starting CreateUser for email: %s", user.Email)
	log.Printf("[DEBUG] Initial user.Password length: %d", len(user.Password)) // Verifikasi panjang password awal

	// Perbaikan: Pindahkan cek password ke awal untuk menghindari proses yang tidak perlu jika password kosong
	if user.Password == "" {
		log.Printf("[ERROR] Password is required but not provided (length: %d)", len(user.Password))
		return fmt.Errorf("password is required")
	}

	if user.ID == "" {
		user.ID = uuid.New().String()
		log.Printf("[DEBUG] Generated new ID: %s", user.ID)
	} else {
		log.Printf("[DEBUG] Using provided ID: %s", user.ID)
	}

	if user.Status == "" {
		user.Status = "active"
		log.Printf("[DEBUG] Set default status to: %s", user.Status)
	} else {
		log.Printf("[DEBUG] Using provided status: %s", user.Status)
	}

	if user.Username == "" && user.FullName != "" {
		baseUsername := strings.ReplaceAll(strings.ToLower(user.FullName), " ", "")
		username := baseUsername
		counter := 1
		log.Printf("[DEBUG] Generating username from FullName: %s, baseUsername: %s", user.FullName, baseUsername)

		for {
			var count int64
			err := s.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
			if err != nil {
				log.Printf("[ERROR] Failed to check username uniqueness: %v", err)
				return fmt.Errorf("failed to check username uniqueness: %v", err)
			}
			if count == 0 {
				log.Printf("[DEBUG] Username available: %s", username)
				break
			}
			username = fmt.Sprintf("%s%d", baseUsername, counter)
			counter++
			log.Printf("[DEBUG] Username conflict, trying: %s", username)
		}
		user.Username = username
		log.Printf("[DEBUG] Final username set: %s", user.Username)
	} else {
		log.Printf("[DEBUG] Username already provided or FullName empty: %s", user.Username)
	}

	if user.KodeReferal == nil || *user.KodeReferal == "" {
		log.Printf("[DEBUG] Generating KodeReferal for username: %s", user.Username)
		code, err := s.generateKodeReferal(user.Username)
		if err != nil {
			log.Printf("[ERROR] Failed to generate KodeReferal: %v", err)
			return err
		}
		user.KodeReferal = &code
		log.Printf("[DEBUG] Generated KodeReferal: %s", code)
	} else {
		log.Printf("[DEBUG] KodeReferal already provided: %s", *user.KodeReferal)
	}

	// Hash password setelah semua validasi awal
	log.Printf("[DEBUG] Hashing password, original length: %d", len(user.Password))
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Failed to hash password: %v", err)
		return fmt.Errorf("failed to hash password: %v", err)
	}
	log.Printf("[DEBUG] Password hashed, hashed length: %d", len(hashed))
	// Perbaikan: Tambahkan validasi panjang hashed password untuk memastikan kompatibilitas dengan field DB (asumsikan varchar(255))
	if len(hashed) > 255 {
		log.Printf("[ERROR] Hashed password length %d exceeds database field limit (255)", len(hashed))
		return fmt.Errorf("hashed password length exceeds database field limit")
	}
	user.Password = string(hashed)
	log.Printf("[DEBUG] Password hashed successfully, final length: %d", len(user.Password))

	// Cek duplicate email sebelum create
	log.Printf("[DEBUG] Checking for duplicate email: %s", user.Email)
	var count int64
	err = s.DB.Model(&model.User{}).Where("email = ?", user.Email).Count(&count).Error
	if err != nil {
		log.Printf("[ERROR] Failed to check email uniqueness: %v", err)
		return fmt.Errorf("failed to check email uniqueness: %v", err)
	}
	if count > 0 {
		log.Printf("[ERROR] Duplicate email found: %s", user.Email)
		return fmt.Errorf("user dengan email sudah ada")
	}
	log.Printf("[DEBUG] Email is unique: %s", user.Email)

	log.Printf("[DEBUG] Creating user with ID: %s, Email: %s, Username: %s, Password length: %d", user.ID, user.Email, user.Username, len(user.Password))
	err = s.DB.Create(user).Error
	if err != nil {
		log.Printf("[ERROR] Failed to create user: %v", err)
		return err
	}
	log.Printf("[DEBUG] User created successfully with ID: %s", user.ID)
	return nil
}

// READ ALL
func (s *UserService) GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := s.DB.Order("created_at DESC").Find(&users)
	return users, result.Error
}

// READ BY ID
func (s *UserService) GetUserByID(id string) (*model.User, error) {
	var user model.User
	result := s.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// UPDATE
func (s *UserService) UpdateUser(user *model.User) error {
	// Pastikan ID di-set dan tidak kosong
	if user.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldUser model.User
	if err := s.DB.First(&oldUser, "id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if user.FullName != "" && user.FullName != oldUser.FullName {
		updateData["full_name"] = user.FullName
		baseUsername := strings.ReplaceAll(strings.ToLower(user.FullName), " ", "")
		username := baseUsername
		counter := 1
		for {
			var count int64
			s.DB.Model(&model.User{}).Where("username = ? AND id != ?", username, user.ID).Count(&count)
			if count == 0 {
				break
			}
			username = fmt.Sprintf("%s%d", baseUsername, counter)
			counter++
		}
		updateData["username"] = username
	}

	// Cek duplicate email HANYA jika email dikirim, tidak kosong, dan berbeda dari yang lama
	if user.Email != "" && user.Email != oldUser.Email {
		var count int64
		s.DB.Model(&model.User{}).Where("email = ? AND id != ?", user.Email, user.ID).Count(&count)
		if count > 0 {
			return fmt.Errorf("user dengan email sudah ada")
		}
		updateData["email"] = user.Email
	}

	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		// Perbaikan: Tambahkan validasi panjang hashed password untuk memastikan kompatibilitas dengan field DB (asumsikan varchar(255))
		if len(hashed) > 255 {
			return fmt.Errorf("hashed password length exceeds database field limit")
		}
		updateData["password"] = string(hashed) // Kolom DB "password"
	}

	// Update field lain yang dikirim (boleh duplikat)
	if user.Status != "" && user.Status != oldUser.Status {
		updateData["status"] = user.Status
	}
	if user.Role != "" && user.Role != oldUser.Role {
		updateData["role"] = user.Role
	}
	if user.LastLogin != nil {
		updateData["last_login"] = user.LastLogin
	}
	if user.LastLogout != nil {
		updateData["last_logout"] = user.LastLogout
	}
	if user.ParentID != nil { // Diubah dari PartnerID
		updateData["parent_id"] = user.ParentID // Diubah dari "partner_id"
	}
	if user.KodeReferal != nil && *user.KodeReferal != "" {
		updateData["kode_referal"] = user.KodeReferal
	}
	if user.EmailVerifiedAt != nil {
		updateData["email_verified_at"] = user.EmailVerifiedAt
	}
	if user.RememberToken != nil {
		updateData["remember_token"] = user.RememberToken
	}
	if user.Saldo != nil {
		updateData["saldo"] = user.Saldo
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.User{}).Where("id = ?", user.ID).Updates(updateData).Error
}

// DELETE
func (s *UserService) DeleteUser(id string) error {
	return s.DB.Delete(&model.User{}, "id = ?", id).Error
}

// Generate kode referal unik
func (s *UserService) generateKodeReferal(username string) (string, error) {
	prefix := strings.ToUpper(username)
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}

	for i := 0; i < 100; i++ {
		suffix := fmt.Sprintf("%06d", rand.Intn(1000000))
		code := prefix + suffix

		var count int64
		s.DB.Model(&model.User{}).Where("kode_referal = ?", code).Count(&count)
		if count == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("gagal generate kode referal unik setelah 100 percobaan")
}

// GetUserProfileWithAplikasi - Ambil profile user dengan aplikasi dan menu access
// GetUserProfileWithAplikasi
func (s *UserService) GetUserProfileWithAplikasi(userID string) (*model.User, map[string]interface{}, error) {
	var user model.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("user not found")
		}
		return nil, nil, err
	}

	hasParent := user.ParentID != nil && *user.ParentID != ""

	response := map[string]interface{}{
		"user":         user,
		"parent_id":    user.ParentID,
		"has_parent":   hasParent,
		"package":      nil,
		"aplikasi":     []map[string]interface{}{},
		"aplikasi_ids": []string{},
		"menu_access":  map[string]interface{}{},
	}

	// ✅ Hanya proses aplikasi, paket, dan menu access untuk role "user"
	// Admin tidak memerlukan ini
	if user.Role == "user" {
		// Ambil daftar aplikasi user
		aplikasiList, aplikasiIDs, err := s.GetUserApplications(userID)
		if err != nil {
			return nil, nil, err
		}
		response["aplikasi"] = aplikasiList
		response["aplikasi_ids"] = aplikasiIDs

		// Ambil info paket dan menu access
		packageInfo := s.CheckUserActivePackage(userID)
		menuAccess := s.GetUserMenuAccess(hasParent, packageInfo)

		response["package"] = packageInfo
		response["menu_access"] = menuAccess
	} else if user.Role == "admin" {
		// Admin mendapat akses penuh
		response["menu_access"] = s.GetAdminMenuAccess()
	} else {
		return nil, nil, fmt.Errorf("unknown role: %s", user.Role)
	}

	return &user, response, nil
}

// GetUserApplications - Ambil daftar aplikasi user DAN return IDs
func (s *UserService) GetUserApplications(userID string) ([]map[string]interface{}, []string, error) {
	var aplikasiList []map[string]interface{}
	var aplikasiIDs []string

	rows, err := s.DB.Table("type_user_aplikasi tua").
		Select("tua.id, tua.user_id, tua.parent_id, tua.aplikasi_id, tua.status, ca.name as aplikasi_name").
		Joins("LEFT JOIN category_aplikasi ca ON tua.aplikasi_id = ca.id").
		Where("tua.user_id = ?", userID).
		Where("tua.status = ?", "active").
		Rows()

	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tuaID, userID, parentID, aplikasiID, status, aplikasiName string

		err := rows.Scan(&tuaID, &userID, &parentID, &aplikasiID, &status, &aplikasiName)
		if err != nil {
			continue
		}

		aplikasi := map[string]interface{}{
			"id":            tuaID,
			"user_id":       userID,
			"parent_id":     parentID,
			"aplikasi_id":   aplikasiID,
			"aplikasi_name": aplikasiName,
			"status":        status,
		}

		aplikasiList = append(aplikasiList, aplikasi)
		aplikasiIDs = append(aplikasiIDs, aplikasiID)
	}

	return aplikasiList, aplikasiIDs, nil
}

// GetAdminMenuAccess - Menu akses untuk admin (akses penuh)
func (s *UserService) GetAdminMenuAccess() map[string]interface{} {
	return map[string]interface{}{
		"can_access_all":     true,
		"role":               "admin",
		"has_active_package": true,
		"allowed_menus": []string{
			"dashboard",
			"user_management",
			"application_management",
			"package_management",
			"report",
			"settings",
		},
	}
}

// GetUserMenuAccess - Menu akses untuk user (DI PERBAIKI)
func (s *UserService) GetUserMenuAccess(hasParent bool, packageInfo map[string]interface{}) map[string]interface{} {

	hasActivePackage := false
	packageName := ""

	if packageInfo != nil {
		if val, ok := packageInfo["has_active_package"].(bool); ok {
			hasActivePackage = val
		}
		if val, ok := packageInfo["package_name"].(string); ok {
			// 🔥 NORMALIZE: lowercase dan trim spasi
			packageName = strings.ToLower(strings.TrimSpace(val))
		}
	}

	menus := []string{}

	// 🔥 CHILD USER
	if hasParent {
		if hasActivePackage {
			menus = []string{
				"dashboard",
				"my_application",
				"profile",
			}
		} else {
			menus = []string{
				"dashboard",
				"profile",
			}
		}
	} else {
		// 🔥 USER UTAMA
		if hasActivePackage {

			// 🔥 CEK BERDASARKAN NORMALIZED NAME
			switch packageName {
			case "premium":
				menus = []string{
					"dashboard",
					"my_application",
					"add_application",
					"child_user_management",
					"report",
					"profile",
				}

			case "basic":
				menus = []string{
					"dashboard",
					"my_application",
					"add_application",
					"profile",
				}

			// 🔥 TRIAL APLIKASI PEMBAYARAN atau TRIAL saja
			case "trial aplikasi pembayaran", "trial aplikasi", "trial":
				menus = []string{
					"dashboard",
					"my_application",
					"add_application",
					"profile",
				}

			default:
				// 🔥 FALLBACK: Jika package name tidak dikenali
				// Cek apakah mengandung kata kunci
				if strings.Contains(packageName, "premium") {
					menus = []string{
						"dashboard",
						"my_application",
						"add_application",
						"child_user_management",
						"report",
						"profile",
					}
				} else if strings.Contains(packageName, "basic") {
					menus = []string{
						"dashboard",
						"my_application",
						"add_application",
						"profile",
					}
				} else if strings.Contains(packageName, "trial") {
					menus = []string{
						"dashboard",
						"my_application",
						"add_application",
						"profile",
					}
				} else {
					// Default jika tidak match sama sekali
					menus = []string{
						"dashboard",
						"profile",
					}
				}
			}

		} else {
			menus = []string{
				"dashboard",
				"profile",
			}
		}
	}

	return map[string]interface{}{
		"can_access_all":     false,
		"role":               "user",
		"has_parent":         hasParent,
		"has_active_package": hasActivePackage,
		"package_name":       packageName,
		"allowed_menus":      menus,
	}
}

// CheckUserActivePackage - Cek paket aktif user dari tabel vouchers (Status: Used)
// CheckUserActivePackage - Cek paket aktif user (cek parent juga)
func (s *UserService) CheckUserActivePackage(userID string) map[string]interface{} {
	var user model.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		fmt.Printf("[DEBUG] User %s: Tidak ditemukan\n", userID)
		return map[string]interface{}{
			"has_active_package": false,
			"package_name":       nil,
			"package_id":         nil,
			"expired_at":         nil,
			"status":             "inactive",
			"kode_voucher":       nil,
		}
	}

	// ✅ CEK VOUCHER USER ITU SENDIRI
	var voucher model.Voucher
	err := s.DB.Table("vouchers v").
		Where("v.user_id = ?", userID).
		Where("v.status = ?", "used").
		Where("v.valid_until > ?", time.Now()).
		Order("v.created_at DESC").
		First(&voucher).Error

	if err == nil {
		// ✅ DITEMUKAN VOUCHER USER - PAKET AKTIF
		var pkg model.Packages
		s.DB.First(&pkg, "id = ?", voucher.PackagesID)

		fmt.Printf("[DEBUG] User %s: Paket aktif (voucher sendiri)\n", userID)

		return map[string]interface{}{
			"has_active_package": true,
			"package_id":         voucher.PackagesID,
			"package_name":       pkg.Name,
			"kode_voucher":       voucher.KodeVoucher,
			"valid_from":         voucher.ValidFrom,
			"expired_at":         voucher.ValidUntil,
			"status":             voucher.Status,
		}
	}

	// ❌ TIDAK ADA VOUCHER SENDIRI - CEK PARENT
	if user.ParentID != nil && *user.ParentID != "" {
		fmt.Printf("[DEBUG] User %s: Cek voucher parent %s\n", userID, *user.ParentID)

		var parentVoucher model.Voucher
		errParent := s.DB.Table("vouchers v").
			Where("v.user_id = ?", *user.ParentID).
			Where("v.status = ?", "used").
			Where("v.valid_until > ?", time.Now()).
			Order("v.created_at DESC").
			First(&parentVoucher).Error

		if errParent == nil {
			// ✅ DITEMUKAN VOUCHER PARENT - PAKET AKTIF (WARISAN)
			var pkg model.Packages
			s.DB.First(&pkg, "id = ?", parentVoucher.PackagesID)

			fmt.Printf("[DEBUG] User %s: Paket aktif (warisan dari parent)\n", userID)

			return map[string]interface{}{
				"has_active_package": true,
				"package_id":         parentVoucher.PackagesID,
				"package_name":       pkg.Name,
				"kode_voucher":       parentVoucher.KodeVoucher,
				"valid_from":         parentVoucher.ValidFrom,
				"expired_at":         parentVoucher.ValidUntil,
				"status":             parentVoucher.Status,
				"inherited_from":     *user.ParentID, // ✅ TAMBAHAN: Tandai ini warisan
			}
		}
	}

	// ❌ TIDAK ADA PAKET SAMA SEKALI
	fmt.Printf("[DEBUG] User %s: Tidak punya paket aktif\n", userID)

	return map[string]interface{}{
		"has_active_package": false,
		"package_name":       nil,
		"package_id":         nil,
		"expired_at":         nil,
		"status":             "inactive",
		"kode_voucher":       nil,
	}
}
