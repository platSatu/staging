package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScHabitsService struct {
	DB *gorm.DB
}

func NewScHabitsService(db *gorm.DB) *ScHabitsService {
	return &ScHabitsService{DB: db}
}

// CREATE
func (s *ScHabitsService) CreateScHabits(habit *model.ScHabits) error {
	if habit.ID == "" {
		habit.ID = uuid.New().String()
	}

	if habit.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	if habit.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", habit.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(habit).Error
}

// READ ALL
func (s *ScHabitsService) GetAllScHabits() ([]model.ScHabits, error) {
	var habits []model.ScHabits
	result := s.DB.Find(&habits)
	return habits, result.Error
}

// READ BY ID
func (s *ScHabitsService) GetScHabitsByID(id string) (*model.ScHabits, error) {
	var habit model.ScHabits
	result := s.DB.First(&habit, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &habit, result.Error
}

// UPDATE
func (s *ScHabitsService) UpdateScHabits(habit *model.ScHabits) error {
	if habit.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldHabit model.ScHabits
	if err := s.DB.First(&oldHabit, "id = ?", habit.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_habits not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if habit.UserID != "" && habit.UserID != oldHabit.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", habit.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = habit.UserID
	}

	if habit.Subject != "" && habit.Subject != oldHabit.Subject {
		updateData["subject"] = habit.Subject
	}

	if habit.Description != nil {
		updateData["description"] = habit.Description
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScHabits{}).Where("id = ?", habit.ID).Updates(updateData).Error
}

// DELETE
func (s *ScHabitsService) DeleteScHabits(id string) error {
	return s.DB.Delete(&model.ScHabits{}, "id = ?", id).Error
}

func (s *ScHabitsService) GetAllScHabitsByUserID(userID string) ([]model.ScHabits, error) {
	var list []model.ScHabits
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}