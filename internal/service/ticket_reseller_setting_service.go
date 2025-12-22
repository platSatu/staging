package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TicketResellerSettingService struct {
	DB *gorm.DB
}

func NewTicketResellerSettingService(db *gorm.DB) *TicketResellerSettingService {
	return &TicketResellerSettingService{DB: db}
}

// CREATE
func (s *TicketResellerSettingService) CreateTicketResellerSetting(ticketResellerSetting *model.TicketResellerSetting) error {
	if ticketResellerSetting.ID == "" {
		ticketResellerSetting.ID = fmt.Sprintf("%s", gorm.Expr("UUID()")) // Atau gunakan uuid.New().String() jika perlu import
	}

	if ticketResellerSetting.Status == "" {
		ticketResellerSetting.Status = "active"
	}

	return s.DB.Create(ticketResellerSetting).Error
}

// READ ALL
func (s *TicketResellerSettingService) GetAllTicketResellerSettings() ([]model.TicketResellerSetting, error) {
	var ticketResellerSettings []model.TicketResellerSetting
	result := s.DB.Find(&ticketResellerSettings)
	return ticketResellerSettings, result.Error
}

// READ BY ID
func (s *TicketResellerSettingService) GetTicketResellerSettingByID(id string) (*model.TicketResellerSetting, error) {
	var ticketResellerSetting model.TicketResellerSetting
	result := s.DB.First(&ticketResellerSetting, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketResellerSetting, result.Error
}

// UPDATE
func (s *TicketResellerSettingService) UpdateTicketResellerSetting(ticketResellerSetting *model.TicketResellerSetting) error {
	if ticketResellerSetting.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketResellerSetting model.TicketResellerSetting
	if err := s.DB.First(&oldTicketResellerSetting, "id = ?", ticketResellerSetting.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket reseller setting not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketResellerSetting.Slug != "" && ticketResellerSetting.Slug != oldTicketResellerSetting.Slug {
		updateData["slug"] = ticketResellerSetting.Slug
	}
	if ticketResellerSetting.Name != "" && ticketResellerSetting.Name != oldTicketResellerSetting.Name {
		updateData["name"] = ticketResellerSetting.Name
	}
	if ticketResellerSetting.Description != nil && (oldTicketResellerSetting.Description == nil || *ticketResellerSetting.Description != *oldTicketResellerSetting.Description) {
		updateData["description"] = ticketResellerSetting.Description
	}
	if ticketResellerSetting.IDReseller != "" && ticketResellerSetting.IDReseller != oldTicketResellerSetting.IDReseller {
		updateData["id_reseller"] = ticketResellerSetting.IDReseller
	}
	if ticketResellerSetting.Status != "" && ticketResellerSetting.Status != oldTicketResellerSetting.Status {
		updateData["status"] = ticketResellerSetting.Status
	}
	if ticketResellerSetting.MethodPembayaran != "" && ticketResellerSetting.MethodPembayaran != oldTicketResellerSetting.MethodPembayaran {
		updateData["method_pembayaran"] = ticketResellerSetting.MethodPembayaran
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketResellerSetting{}).Where("id = ?", ticketResellerSetting.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketResellerSettingService) DeleteTicketResellerSetting(id string) error {
	return s.DB.Delete(&model.TicketResellerSetting{}, "id = ?", id).Error
}
