package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScAcademicProjectionService struct {
	DB *gorm.DB
}

func NewScAcademicProjectionService(db *gorm.DB) *ScAcademicProjectionService {
	return &ScAcademicProjectionService{DB: db}
}

// CREATE
func (s *ScAcademicProjectionService) CreateScAcademicProjection(projection *model.ScAcademicProjection) error {
	if projection.ID == "" {
		projection.ID = uuid.New().String()
	}

	if projection.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", projection.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(projection).Error
}

// READ ALL
func (s *ScAcademicProjectionService) GetAllScAcademicProjection() ([]model.ScAcademicProjection, error) {
	var projections []model.ScAcademicProjection
	result := s.DB.Find(&projections)
	return projections, result.Error
}

// READ BY ID
func (s *ScAcademicProjectionService) GetScAcademicProjectionByID(id string) (*model.ScAcademicProjection, error) {
	var projection model.ScAcademicProjection
	result := s.DB.First(&projection, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &projection, result.Error
}

// UPDATE
func (s *ScAcademicProjectionService) UpdateScAcademicProjection(projection *model.ScAcademicProjection) error {
	if projection.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldProjection model.ScAcademicProjection
	if err := s.DB.First(&oldProjection, "id = ?", projection.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_academic_projection not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if projection.UserID != "" && projection.UserID != oldProjection.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", projection.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = projection.UserID
	}

	if projection.StudentID != nil && (oldProjection.StudentID == nil || *projection.StudentID != *oldProjection.StudentID) {
		updateData["student_id"] = projection.StudentID
	}

	if projection.AcademicYear != nil && (oldProjection.AcademicYear == nil || *projection.AcademicYear != *oldProjection.AcademicYear) {
		updateData["academic_year"] = projection.AcademicYear
	}

	if projection.Semester != nil && (oldProjection.Semester == nil || *projection.Semester != *oldProjection.Semester) {
		updateData["semester"] = projection.Semester
	}

	if projection.QuarterID != nil && (oldProjection.QuarterID == nil || *projection.QuarterID != *oldProjection.QuarterID) {
		updateData["quarter_id"] = projection.QuarterID
	}

	if projection.IsSplitIso != nil && (oldProjection.IsSplitIso == nil || *projection.IsSplitIso != *oldProjection.IsSplitIso) {
		updateData["is_split_iso"] = projection.IsSplitIso
	}

	if projection.Level != nil && (oldProjection.Level == nil || *projection.Level != *oldProjection.Level) {
		updateData["level"] = projection.Level
	}

	if projection.LcID != nil && (oldProjection.LcID == nil || *projection.LcID != *oldProjection.LcID) {
		updateData["lc_id"] = projection.LcID
	}

	if projection.TotalSchool != nil && (oldProjection.TotalSchool == nil || *projection.TotalSchool != *oldProjection.TotalSchool) {
		updateData["total_school"] = projection.TotalSchool
	}

	if projection.TotalPages != nil && (oldProjection.TotalPages == nil || *projection.TotalPages != *oldProjection.TotalPages) {
		updateData["total_pages"] = projection.TotalPages
	}

	if projection.Status != nil && (oldProjection.Status == nil || *projection.Status != *oldProjection.Status) {
		updateData["status"] = projection.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScAcademicProjection{}).Where("id = ?", projection.ID).Updates(updateData).Error
}

// DELETE
func (s *ScAcademicProjectionService) DeleteScAcademicProjection(id string) error {
	return s.DB.Delete(&model.ScAcademicProjection{}, "id = ?", id).Error
}

func (s *ScAcademicProjectionService) GetAllScAcademicProjectionByUserID(userID string) ([]model.ScAcademicProjection, error) {
	var list []model.ScAcademicProjection
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

// GetScAcademicProjectionByStudentID returns academic projection based on student_id
func (s *ScAcademicProjectionService) GetScAcademicProjectionByStudentID(studentID string) ([]model.ScAcademicProjection, error) {
	var projections []model.ScAcademicProjection

	// Tampilkan semua projection dimana student_id = user yang login
	result := s.DB.Where("student_id = ?", studentID).
		Order("created_at DESC").
		Find(&projections)

	return projections, result.Error
}
