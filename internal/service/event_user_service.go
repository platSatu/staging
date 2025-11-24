package service

import (
	"backend_go/internal/model"
	"errors"

	"gorm.io/gorm"
)

type EventUserService struct {
	DB *gorm.DB
}

func NewEventUserService(db *gorm.DB) *EventUserService {
	return &EventUserService{DB: db}
}

func (s *EventUserService) CreateEventUser(input model.EventUser) (*model.EventUser, error) {
	var existing model.EventUser
	if err := s.DB.Where("email = ? AND user_id = ?", input.Email, input.UserID).First(&existing).Error; err == nil {
		return nil, errors.New("email sudah terdaftar untuk tenant ini")
	}

	if err := s.DB.Create(&input).Error; err != nil {
		return nil, err
	}
	return &input, nil
}
