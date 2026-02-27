package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type ScStudentService struct {
	DB *gorm.DB
}

func NewScStudentService(db *gorm.DB) *ScStudentService {
	return &ScStudentService{DB: db}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CREATE
func (s *ScStudentService) CreateScStudent(student *model.ScStudent) error {
	if student.ID == "" {
		student.ID = uuid.New().String()
	}

	if student.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", student.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(student).Error
}

// READ ALL
func (s *ScStudentService) GetAllScStudent() ([]model.ScStudent, error) {
	var students []model.ScStudent
	result := s.DB.Find(&students)
	return students, result.Error
}

// READ BY ID
func (s *ScStudentService) GetScStudentByID(id string) (*model.ScStudent, error) {
	var student model.ScStudent
	result := s.DB.First(&student, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &student, result.Error
}

// UPDATE
func (s *ScStudentService) UpdateScStudent(student *model.ScStudent) error {
	if student.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldStudent model.ScStudent
	if err := s.DB.First(&oldStudent, "id = ?", student.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_student not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if student.UserID != "" && student.UserID != oldStudent.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", student.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = student.UserID
	}

	if student.Name != nil && (oldStudent.Name == nil || *student.Name != *oldStudent.Name) {
		updateData["name"] = student.Name
	}

	if student.Address != nil && (oldStudent.Address == nil || *student.Address != *oldStudent.Address) {
		updateData["address"] = student.Address
	}

	if student.Tin != nil && (oldStudent.Tin == nil || *student.Tin != *oldStudent.Tin) {
		updateData["tin"] = student.Tin
	}

	if student.Tags != nil && (oldStudent.Tags == nil || *student.Tags != *oldStudent.Tags) {
		updateData["tags"] = student.Tags
	}

	if student.StudentType != nil && (oldStudent.StudentType == nil || *student.StudentType != *oldStudent.StudentType) {
		updateData["student_type"] = student.StudentType
	}

	if student.LcID != nil && (oldStudent.LcID == nil || *student.LcID != *oldStudent.LcID) {
		updateData["lc_id"] = student.LcID
	}

	if student.LevelID != nil && (oldStudent.LevelID == nil || *student.LevelID != *oldStudent.LevelID) {
		updateData["level_id"] = student.LevelID
	}

	if student.Phone != nil && (oldStudent.Phone == nil || *student.Phone != *oldStudent.Phone) {
		updateData["phone"] = student.Phone
	}

	if student.Mobile != nil && (oldStudent.Mobile == nil || *student.Mobile != *oldStudent.Mobile) {
		updateData["mobile"] = student.Mobile
	}

	if student.Email != nil && (oldStudent.Email == nil || *student.Email != *oldStudent.Email) {
		updateData["email"] = student.Email
	}

	if student.Language != nil && (oldStudent.Language == nil || *student.Language != *oldStudent.Language) {
		updateData["language"] = student.Language
	}

	if student.StudentStatus != nil && (oldStudent.StudentStatus == nil || *student.StudentStatus != *oldStudent.StudentStatus) {
		updateData["student_status"] = student.StudentStatus
	}

	if student.PartnerType != nil && (oldStudent.PartnerType == nil || *student.PartnerType != *oldStudent.PartnerType) {
		updateData["partner_type"] = student.PartnerType
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScStudent{}).Where("id = ?", student.ID).Updates(updateData).Error
}

// DELETE
func (s *ScStudentService) DeleteScStudent(id string) error {
	return s.DB.Delete(&model.ScStudent{}, "id = ?", id).Error
}

func (s *ScStudentService) GetAllScStudentByUserID(userID string) ([]model.ScStudent, error) {
	var list []model.ScStudent
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// CreateUserAndStudent creates a new user and student simultaneously
func (s *ScStudentService) CreateUserAndStudent(parentID string, userData *model.User, studentData *model.ScStudent) (*model.User, *model.ScStudent, error) {
	// Validasi parent_id ada (user yang login)
	var parentUser model.User
	if err := s.DB.First(&parentUser, "id = ?", parentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, fmt.Errorf("parent user not found")
		}
		return nil, nil, err
	}

	// Set parent_id untuk user baru (child) dengan ID user yang login
	userData.ParentID = &parentID

	// Generate UUID untuk user jika kosong
	if userData.ID == "" {
		userData.ID = uuid.New().String()
	}

	// Set default status dan role jika kosong
	if userData.Status == "" {
		userData.Status = "active"
	}
	if userData.Role == "" {
		userData.Role = "user"
	}

	// Hash password sebelum simpan
	hashedPassword, err := hashPassword(userData.Password)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %v", err)
	}
	userData.Password = hashedPassword

	// Validasi email unik
	var existingUser model.User
	if err := s.DB.Where("email = ?", userData.Email).First(&existingUser).Error; err == nil {
		return nil, nil, fmt.Errorf("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	// Simpan user baru (child)
	if err := s.DB.Create(userData).Error; err != nil {
		return nil, nil, err
	}

	// Untuk student, user_id diisi dengan parentID (user yang login)
	studentData.UserID = parentID

	// Generate UUID untuk student jika kosong
	if studentData.ID == "" {
		studentData.ID = uuid.New().String()
	}

	// Simpan student baru
	if err := s.DB.Create(studentData).Error; err != nil {
		return nil, nil, err
	}

	return userData, studentData, nil
}
