// internal/service/sc_user_service.go

package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SCUserService struct {
	DB *gorm.DB
}

func NewSCUserService(db *gorm.DB) *SCUserService {
	return &SCUserService{DB: db}
}

// CREATE - Membuat SC user baru dengan ParentID dari user yang login
func (s *SCUserService) CreateSCUser(scUser *model.User, parentID string, aplikasiID string) error {
	log.Printf("[DEBUG] Starting CreateSCUser for email: %s", scUser.Email)

	// Validasi parentID tidak kosong
	if strings.TrimSpace(parentID) == "" {
		return fmt.Errorf("parent ID is required")
	}

	// Validasi parent user exists
	var parentUser model.User
	if err := s.DB.First(&parentUser, "id = ?", parentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("parent user not found")
		}
		return fmt.Errorf("error finding parent user: %v", err)
	}

	// Validasi required fields
	if strings.TrimSpace(scUser.Email) == "" {
		return fmt.Errorf("email is required")
	}
	if strings.TrimSpace(scUser.Password) == "" {
		return fmt.Errorf("password is required")
	}
	if strings.TrimSpace(scUser.FullName) == "" {
		return fmt.Errorf("full name is required")
	}

	// Generate ID jika kosong
	if strings.TrimSpace(scUser.ID) == "" {
		scUser.ID = uuid.New().String()
		log.Printf("[DEBUG] Generated new ID: %s", scUser.ID)
	}

	// Set default status
	if strings.TrimSpace(scUser.Status) == "" {
		scUser.Status = "active"
	}

	// Set default role
	if strings.TrimSpace(scUser.Role) == "" {
		scUser.Role = "user"
	}

	// Set ParentID (user yang login)
	scUser.ParentID = &parentID

	// Set EmailVerifiedAt ke sekarang
	now := time.Now()
	scUser.EmailVerifiedAt = &now

	// Generate username jika kosong
	if strings.TrimSpace(scUser.Username) == "" && strings.TrimSpace(scUser.FullName) != "" {
		username, err := s.generateUniqueUsername(scUser.FullName)
		if err != nil {
			return err
		}
		scUser.Username = username
		log.Printf("[DEBUG] Generated username: %s", scUser.Username)
	}

	// Generate KodeReferal jika kosong
	if scUser.KodeReferal == nil || strings.TrimSpace(*scUser.KodeReferal) == "" {
		code, err := s.generateKodeReferal(scUser.Username)
		if err != nil {
			return err
		}
		scUser.KodeReferal = &code
		log.Printf("[DEBUG] Generated KodeReferal: %s", code)
	}

	// Hash password
	log.Printf("[DEBUG] Hashing password for: %s", scUser.Email)
	hashed, err := bcrypt.GenerateFromPassword([]byte(scUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("[ERROR] Failed to hash password: %v", err)
		return fmt.Errorf("failed to hash password: %v", err)
	}

	// Validasi panjang hashed password
	if len(hashed) > 255 {
		log.Printf("[ERROR] Hashed password length %d exceeds database field limit (255)", len(hashed))
		return fmt.Errorf("hashed password length exceeds database field limit")
	}

	scUser.Password = string(hashed)
	log.Printf("[DEBUG] Password hashed successfully")

	// Cek duplicate email
	log.Printf("[DEBUG] Checking duplicate email: %s", scUser.Email)
	var count int64
	err = s.DB.Model(&model.User{}).Where("email = ?", scUser.Email).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %v", err)
	}
	if count > 0 {
		log.Printf("[ERROR] Duplicate email: %s", scUser.Email)
		return fmt.Errorf("user dengan email sudah ada")
	}

	// Create user
	log.Printf("[DEBUG] Creating SC user: %s", scUser.Email)
	err = s.DB.Create(scUser).Error
	if err != nil {
		log.Printf("[ERROR] Failed to create SC user: %v", err)
		return fmt.Errorf("failed to create sc user: %v", err)
	}

	log.Printf("[DEBUG] SC user created successfully: %s", scUser.ID)

	// =====================================================
	// JIKA ROLE ADALAH STUDENT DAN aplikasiID TIDAK KOSONG,
	// BUAT TypeUserAplikasi
	// =====================================================
	if scUser.Role == "student" && aplikasiID != "" {
		log.Printf("[DEBUG] Creating TypeUserAplikasi for student: %s", scUser.ID)
		log.Printf("[DEBUG] AplikasiID: %s", aplikasiID)
		log.Printf("[DEBUG] ParentID: %s", parentID)

		typeUserAplikasi := &model.TypeUserAplikasi{
			ID:         uuid.New().String(),
			UserID:     scUser.ID,
			ParentID:   parentID,
			AplikasiID: aplikasiID,
			Status:     "active",
		}

		log.Printf("[DEBUG] TypeUserAplikasi data: %+v", typeUserAplikasi)

		if err := s.DB.Create(typeUserAplikasi).Error; err != nil {
			log.Printf("[ERROR] Failed to create TypeUserAplikasi: %v", err)
			// NOTE: Tidak return error karena user sudah berhasil dibuat
			// Jika ingin rollback, gunakan transaction
		} else {
			log.Printf("[DEBUG] TypeUserAplikasi created successfully: %s", typeUserAplikasi.ID)
		}
	}

	return nil
}

// READ ALL - Ambil semua SC user
func (s *SCUserService) GetAllSCUsers() ([]model.User, error) {
	var users []model.User
	result := s.DB.Where("parent_id IS NOT NULL").Order("created_at DESC").Find(&users)
	return users, result.Error
}

// READ BY ID - Ambil SC user by ID
func (s *SCUserService) GetSCUserByID(id string) (*model.User, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id is required")
	}

	var user model.User
	result := s.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// READ BY PARENT ID - Ambil semua SC user dari parent tertentu
func (s *SCUserService) GetSCUsersByParentID(parentID string) ([]model.User, error) {
	if strings.TrimSpace(parentID) == "" {
		return nil, fmt.Errorf("parent ID is required")
	}

	var users []model.User
	result := s.DB.Where("parent_id = ?", parentID).Order("created_at DESC").Find(&users)
	return users, result.Error
}

// UPDATE - Update SC user
func (s *SCUserService) UpdateSCUser(scUser *model.User) error {
	// Pastikan ID tidak kosong
	if strings.TrimSpace(scUser.ID) == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldUser model.User
	if err := s.DB.First(&oldUser, "id = ?", scUser.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc user not found")
		}
		return err
	}

	// Validasi user ini adalah SC user (punya parent)
	if oldUser.ParentID == nil {
		return fmt.Errorf("user is not a sc user")
	}

	updateData := map[string]interface{}{}

	// Update FullName
	if scUser.FullName != "" && scUser.FullName != oldUser.FullName {
		updateData["full_name"] = scUser.FullName
	}

	// Update Username jika berbeda (username boleh duplikat)
	if scUser.Username != "" && scUser.Username != oldUser.Username {
		updateData["username"] = scUser.Username
	}

	// Update Password jika diisi
	if scUser.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(scUser.Password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %v", err)
		}
		if len(hashed) > 255 {
			return fmt.Errorf("hashed password length exceeds database field limit")
		}
		updateData["password"] = string(hashed)
	}

	// Update Status
	if scUser.Status != "" && scUser.Status != oldUser.Status {
		updateData["status"] = scUser.Status
	}

	// Update Role
	if scUser.Role != "" && scUser.Role != oldUser.Role {
		updateData["role"] = scUser.Role
	}

	// Update KodeReferal (kode referal boleh duplikat)
	if scUser.KodeReferal != nil && *scUser.KodeReferal != "" {
		updateData["kode_referal"] = *scUser.KodeReferal
	}

	// Update Saldo
	if scUser.Saldo != oldUser.Saldo {
		updateData["saldo"] = scUser.Saldo
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.User{}).Where("id = ?", scUser.ID).Updates(updateData).Error
}

// DELETE - Hapus SC user (permanent, tanpa soft delete)
func (s *SCUserService) DeleteSCUser(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("id is required")
	}

	var user model.User
	if err := s.DB.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc user not found")
		}
		return err
	}

	// Validasi user ini adalah SC user
	if user.ParentID == nil {
		return fmt.Errorf("user is not a sc user")
	}

	// Hapus TypeUserAplikasi terkait terlebih dahulu
	log.Printf("[DEBUG] Deleting TypeUserAplikasi for user: %s", id)
	if err := s.DB.Where("user_id = ?", id).Delete(&model.TypeUserAplikasi{}).Error; err != nil {
		log.Printf("[ERROR] Failed to delete TypeUserAplikasi: %v", err)
	}

	// Delete permanent user
	return s.DB.Delete(&user).Error
}

// generateUniqueUsername - Generate username unik
func (s *SCUserService) generateUniqueUsername(fullName string) (string, error) {
	baseUsername := strings.ReplaceAll(strings.ToLower(fullName), " ", "")
	baseUsername = strings.ReplaceAll(baseUsername, "'", "")
	baseUsername = strings.ReplaceAll(baseUsername, "\"", "")

	// Hapus karakter non-alphanumeric
	var result strings.Builder
	for _, r := range baseUsername {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}
	baseUsername = result.String()

	if baseUsername == "" {
		baseUsername = "user"
	}

	username := baseUsername
	counter := 1

	for {
		var count int64
		err := s.DB.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
		if err != nil {
			return "", fmt.Errorf("failed to check username uniqueness: %v", err)
		}
		if count == 0 {
			break
		}
		username = fmt.Sprintf("%s%d", baseUsername, counter)
		counter++

		// Safety check
		if counter > 1000 {
			return "", fmt.Errorf("failed to generate unique username after 1000 attempts")
		}
	}

	return username, nil
}

// generateKodeReferal - Generate kode referal unik
func (s *SCUserService) generateKodeReferal(username string) (string, error) {
	prefix := strings.ToUpper(username)
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}

	for i := 0; i < 100; i++ {
		suffix := fmt.Sprintf("%06d", uuid.New().ID()%1000000)
		code := prefix + suffix

		var count int64
		s.DB.Model(&model.User{}).Where("kode_referal = ?", code).Count(&count)
		if count == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("gagal generate kode referal unik setelah 100 percobaan")
}
