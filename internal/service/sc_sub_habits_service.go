package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScSubHabitsService struct {
	DB *gorm.DB
}

func NewScSubHabitsService(db *gorm.DB) *ScSubHabitsService {
	return &ScSubHabitsService{DB: db}
}

// CREATE
func (s *ScSubHabitsService) CreateScSubHabits(subHabit *model.ScSubHabits) error {
	if subHabit.ID == "" {
		subHabit.ID = uuid.New().String()
	}

	if subHabit.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if subHabit.HabitsID == "" {
		return fmt.Errorf("habits_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", subHabit.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	// Validasi FK habits_id
	var habit model.ScHabits
	if err := s.DB.First(&habit, "id = ?", subHabit.HabitsID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("habits_id not found")
		}
		return err
	}

	return s.DB.Create(subHabit).Error
}

// READ ALL
func (s *ScSubHabitsService) GetAllScSubHabits() ([]model.ScSubHabits, error) {
	var subHabits []model.ScSubHabits
	result := s.DB.Find(&subHabits)
	return subHabits, result.Error
}

// READ BY ID
func (s *ScSubHabitsService) GetScSubHabitsByID(id string) (*model.ScSubHabits, error) {
	var subHabit model.ScSubHabits
	result := s.DB.First(&subHabit, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subHabit, result.Error
}

// UPDATE
func (s *ScSubHabitsService) UpdateScSubHabits(subHabit *model.ScSubHabits) error {
	if subHabit.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldSubHabit model.ScSubHabits
	if err := s.DB.First(&oldSubHabit, "id = ?", subHabit.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_sub_habits not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if subHabit.UserID != "" && subHabit.UserID != oldSubHabit.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", subHabit.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = subHabit.UserID
	}

	if subHabit.HabitsID != "" && subHabit.HabitsID != oldSubHabit.HabitsID {
		// Validasi FK habits_id
		var habit model.ScHabits
		if err := s.DB.First(&habit, "id = ?", subHabit.HabitsID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("habits_id not found")
			}
			return err
		}
		updateData["habits_id"] = subHabit.HabitsID
	}

	if subHabit.Subject != nil {
		updateData["subject"] = subHabit.Subject
	}

	if subHabit.Description != nil {
		updateData["description"] = subHabit.Description
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScSubHabits{}).Where("id = ?", subHabit.ID).Updates(updateData).Error
}

// DELETE
func (s *ScSubHabitsService) DeleteScSubHabits(id string) error {
	return s.DB.Delete(&model.ScSubHabits{}, "id = ?", id).Error
}

func (s *ScSubHabitsService) GetAllScSubHabitsByUserID(userID string) ([]model.ScSubHabits, error) {
	var list []model.ScSubHabits
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (s *ScSubHabitsService) GetAllScSubHabitsByHabitsID(habitsID string) ([]model.ScSubHabits, error) {
	var list []model.ScSubHabits
	err := s.DB.Where("habits_id = ?", habitsID).Find(&list).Error
	return list, err
}