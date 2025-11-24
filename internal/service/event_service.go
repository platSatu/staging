package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService struct {
	DB *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{DB: db}
}

// CREATE
func (s *EventService) CreateEvent(event *model.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	if event.Status == "" {
		event.Status = "active"
	}

	// Validasi user_id ada di tabel users (opsional, jika ingin foreign key constraint)
	var user model.User
	if err := s.DB.First(&user, "id = ?", event.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	return s.DB.Create(event).Error
}

// READ ALL
func (s *EventService) GetAllEvents() ([]model.Event, error) {
	var events []model.Event
	result := s.DB.Find(&events)
	return events, result.Error
}

// READ BY ID
func (s *EventService) GetEventByID(id string) (*model.Event, error) {
	var event model.Event
	result := s.DB.First(&event, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &event, result.Error
}

// READ BY USER ID
func (s *EventService) GetEventsByUserID(userID string) ([]model.Event, error) {
	var events []model.Event
	result := s.DB.Where("user_id = ?", userID).Find(&events)
	return events, result.Error
}

// UPDATE
func (s *EventService) UpdateEvent(event *model.Event) error {
	if event.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldEvent model.Event
	if err := s.DB.First(&oldEvent, "id = ?", event.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if event.Name != "" && event.Name != oldEvent.Name {
		updateData["name"] = event.Name
	}

	if event.Description != "" && event.Description != oldEvent.Description {
		updateData["description"] = event.Description
	}

	if !event.Date.IsZero() && event.Date != oldEvent.Date {
		updateData["date"] = event.Date
	}

	if event.Location != "" && event.Location != oldEvent.Location {
		updateData["location"] = event.Location
	}

	if event.Status != "" && event.Status != oldEvent.Status {
		updateData["status"] = event.Status
	}

	if event.UserID != "" && event.UserID != oldEvent.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", event.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = event.UserID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Event{}).Where("id = ?", event.ID).Updates(updateData).Error
}

// DELETE
func (s *EventService) DeleteEvent(id string) error {
	return s.DB.Delete(&model.Event{}, "id = ?", id).Error
}
