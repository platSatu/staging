package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid" // Tambahkan import untuk UUID
	"gorm.io/gorm"
)

type TicketFeeSettingService struct {
	DB *gorm.DB
}

func NewTicketFeeSettingService(db *gorm.DB) *TicketFeeSettingService {
	return &TicketFeeSettingService{DB: db}
}

// CREATE
func (s *TicketFeeSettingService) CreateTicketFeeSetting(ticketFeeSetting *model.TicketFeeSetting) error {
	if ticketFeeSetting.ID == "" {
		ticketFeeSetting.ID = uuid.New().String() // Generate UUID baru
	}

	if ticketFeeSetting.Status == "" {
		ticketFeeSetting.Status = "active"
	}

	return s.DB.Create(ticketFeeSetting).Error
}

// READ ALL
func (s *TicketFeeSettingService) GetAllTicketFeeSettings() ([]model.TicketFeeSetting, error) {
	var ticketFeeSettings []model.TicketFeeSetting
	result := s.DB.Find(&ticketFeeSettings)
	return ticketFeeSettings, result.Error
}

// READ BY ID
func (s *TicketFeeSettingService) GetTicketFeeSettingByID(id string) (*model.TicketFeeSetting, error) {
	var ticketFeeSetting model.TicketFeeSetting
	result := s.DB.First(&ticketFeeSetting, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketFeeSetting, result.Error
}

// UPDATE
func (s *TicketFeeSettingService) UpdateTicketFeeSetting(ticketFeeSetting *model.TicketFeeSetting) error {
	if ticketFeeSetting.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketFeeSetting model.TicketFeeSetting
	if err := s.DB.First(&oldTicketFeeSetting, "id = ?", ticketFeeSetting.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket fee setting not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketFeeSetting.Name != "" && ticketFeeSetting.Name != oldTicketFeeSetting.Name {
		updateData["name"] = ticketFeeSetting.Name
	}
	if ticketFeeSetting.JenisFee != "" && ticketFeeSetting.JenisFee != oldTicketFeeSetting.JenisFee {
		updateData["jenis_fee"] = ticketFeeSetting.JenisFee
	}
	if ticketFeeSetting.Nominal != 0 && ticketFeeSetting.Nominal != oldTicketFeeSetting.Nominal {
		updateData["nominal"] = ticketFeeSetting.Nominal
	}
	if ticketFeeSetting.Status != "" && ticketFeeSetting.Status != oldTicketFeeSetting.Status {
		updateData["status"] = ticketFeeSetting.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketFeeSetting{}).Where("id = ?", ticketFeeSetting.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketFeeSettingService) DeleteTicketFeeSetting(id string) error {
	return s.DB.Delete(&model.TicketFeeSetting{}, "id = ?", id).Error
}
