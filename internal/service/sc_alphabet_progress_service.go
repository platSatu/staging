package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScAlphabetProgressService struct {
	DB *gorm.DB
}

func NewScAlphabetProgressService(db *gorm.DB) *ScAlphabetProgressService {
	return &ScAlphabetProgressService{DB: db}
}

// CREATE
func (s *ScAlphabetProgressService) CreateScAlphabetProgress(progress *model.ScAlphabetProgress) error {
	if progress.ID == "" {
		progress.ID = uuid.New().String()
	}

	if progress.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", progress.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(progress).Error
}

// READ ALL
func (s *ScAlphabetProgressService) GetAllScAlphabetProgress() ([]model.ScAlphabetProgress, error) {
	var progresses []model.ScAlphabetProgress
	result := s.DB.Find(&progresses)
	return progresses, result.Error
}

// READ BY ID
func (s *ScAlphabetProgressService) GetScAlphabetProgressByID(id string) (*model.ScAlphabetProgress, error) {
	var progress model.ScAlphabetProgress
	result := s.DB.First(&progress, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &progress, result.Error
}

// UPDATE
func (s *ScAlphabetProgressService) UpdateScAlphabetProgress(progress *model.ScAlphabetProgress) error {
	if progress.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldProgress model.ScAlphabetProgress
	if err := s.DB.First(&oldProgress, "id = ?", progress.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_alphabet_progress not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if progress.UserID != "" && progress.UserID != oldProgress.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", progress.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = progress.UserID
	}

	if progress.Level != "" && progress.Level != oldProgress.Level {
		updateData["level"] = progress.Level
	}

	if progress.LVertical != "" && progress.LVertical != oldProgress.LVertical {
		updateData["l_vertical"] = progress.LVertical
	}

	if progress.LHorizontal != "" && progress.LHorizontal != oldProgress.LHorizontal {
		updateData["l_horizontal"] = progress.LHorizontal
	}

	if progress.Score != "" && progress.Score != oldProgress.Score {
		updateData["score"] = progress.Score
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScAlphabetProgress{}).Where("id = ?", progress.ID).Updates(updateData).Error
}

// DELETE
func (s *ScAlphabetProgressService) DeleteScAlphabetProgress(id string) error {
	return s.DB.Delete(&model.ScAlphabetProgress{}, "id = ?", id).Error
}

func (s *ScAlphabetProgressService) GetAllScAlphabetProgressByUserID(userID string) ([]model.ScAlphabetProgress, error) {
	var list []model.ScAlphabetProgress
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}