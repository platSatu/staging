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
	result := s.DB.Find(&users)
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

// GetUserProfile ambil profile user berdasarkan userID
func (s *UserService) GetUserProfile(userID string) (*model.User, error) {
	return s.GetUserByID(userID)
}
