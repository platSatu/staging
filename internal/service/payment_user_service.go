package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PaymentUserService struct {
	DB *gorm.DB
}

func NewPaymentUserService(db *gorm.DB) *PaymentUserService {
	return &PaymentUserService{DB: db}
}

// // CREATE
// func (s *PaymentUserService) CreatePaymentUser(user *model.PaymentUser) error {
// 	if user.ParentID == "" {
// 		return fmt.Errorf("parent_id is required")
// 	}

// 	if user.ID == "" {
// 		user.ID = uuid.New().String()
// 	}

// 	// Check duplicate email
// 	var count int64
// 	err := s.DB.Model(&model.PaymentUser{}).Where("email = ?", user.Email).Count(&count).Error
// 	if err != nil {
// 		return fmt.Errorf("failed to check email uniqueness: %v", err)
// 	}
// 	if count > 0 {
// 		return fmt.Errorf("user dengan email sudah ada")
// 	}

// 	// Check duplicate username
// 	err = s.DB.Model(&model.PaymentUser{}).Where("username = ?", user.Username).Count(&count).Error
// 	if err != nil {
// 		return fmt.Errorf("failed to check username uniqueness: %v", err)
// 	}
// 	if count > 0 {
// 		return fmt.Errorf("user dengan username sudah ada")
// 	}

// 	// Hash password
// 	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return fmt.Errorf("failed to hash password: %v", err)
// 	}
// 	user.Password = string(hashed)

//		err = s.DB.Create(user).Error
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
// CREATE
func (s *PaymentUserService) CreatePaymentUser(user *model.PaymentUser) error {
	if user.ParentID == "" {
		return fmt.Errorf("parent_id is required")
	}

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Check duplicate email
	var count int64
	err := s.DB.Model(&model.PaymentUser{}).Where("email = ?", user.Email).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check email uniqueness: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("user dengan email sudah ada")
	}

	// Check duplicate username
	err = s.DB.Model(&model.PaymentUser{}).Where("username = ?", user.Username).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check username uniqueness: %v", err)
	}
	if count > 0 {
		return fmt.Errorf("user dengan username sudah ada")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	user.Password = string(hashed)

	// Gunakan transaction untuk memastikan atomicity (kedua insert berhasil atau gagal bersama)
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert ke tabel users (PaymentUser)
	err = tx.Create(user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Otomatis insert ke type_user_aplikasi setelah create user berhasil
	// Hanya jika parent_id ada (sudah dicek di atas)
	typeUserAplikasi := &model.TypeUserAplikasi{
		ID:         uuid.New().String(),                    // Generate UUID baru untuk ID
		UserID:     user.ID,                                // Mengikuti user_id dari user yang baru dibuat
		ParentID:   user.ParentID,                          // Mengikuti parent_id dari user
		AplikasiID: "56d73587-11c2-45ab-939f-a0bf9fa9968d", // Aplikasi ID dari sistem
		Status:     "active",                               // Status active
	}

	// Insert ke type_user_aplikasi tanpa validasi FK (untuk menghindari error jika data belum ada)
	err = tx.Create(typeUserAplikasi).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create type_user_aplikasi: %v", err)
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

// READ ALL
func (s *PaymentUserService) GetAllPaymentUsers(parentID string) ([]model.PaymentUser, error) {
	var users []model.PaymentUser
	err := s.DB.Where("parent_id = ?", parentID).Find(&users).Error
	return users, err
}

// READ BY ID
func (s *PaymentUserService) GetPaymentUserByID(id string) (*model.PaymentUser, error) {
	var user model.PaymentUser
	result := s.DB.First(&user, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, result.Error
}

// UPDATE
func (s *PaymentUserService) UpdatePaymentUser(user *model.PaymentUser) error {
	if user.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldUser model.PaymentUser
	if err := s.DB.First(&oldUser, "id = ?", user.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("payment user not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if user.Username != "" && user.Username != oldUser.Username {
		var count int64
		s.DB.Model(&model.PaymentUser{}).Where("username = ? AND id != ?", user.Username, user.ID).Count(&count)
		if count > 0 {
			return fmt.Errorf("user dengan username sudah ada")
		}
		updateData["username"] = user.Username
	}

	if user.Email != "" && user.Email != oldUser.Email {
		var count int64
		s.DB.Model(&model.PaymentUser{}).Where("email = ? AND id != ?", user.Email, user.ID).Count(&count)
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
		updateData["password"] = string(hashed)
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.PaymentUser{}).Where("id = ?", user.ID).Updates(updateData).Error
}

// DELETE
func (s *PaymentUserService) DeletePaymentUser(id string) error {
	return s.DB.Delete(&model.PaymentUser{}, "id = ?", id).Error
}
