package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TicketHistoryService struct {
	DB *gorm.DB
}

func NewTicketHistoryService(db *gorm.DB) *TicketHistoryService {
	return &TicketHistoryService{DB: db}
}

// CREATE
func (s *TicketHistoryService) CreateTicketHistory(ticketHistory *model.TicketHistory) error {
	if ticketHistory.ID == "" {
		ticketHistory.ID = fmt.Sprintf("%s", gorm.Expr("UUID()")) // Atau gunakan uuid.New().String() jika perlu import
	}

	return s.DB.Create(ticketHistory).Error
}

// READ ALL
func (s *TicketHistoryService) GetAllTicketHistories() ([]model.TicketHistory, error) {
	var ticketHistories []model.TicketHistory
	result := s.DB.Find(&ticketHistories)
	return ticketHistories, result.Error
}

// READ BY ID
func (s *TicketHistoryService) GetTicketHistoryByID(id string) (*model.TicketHistory, error) {
	var ticketHistory model.TicketHistory
	result := s.DB.First(&ticketHistory, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketHistory, result.Error
}

// UPDATE
func (s *TicketHistoryService) UpdateTicketHistory(ticketHistory *model.TicketHistory) error {
	if ticketHistory.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketHistory model.TicketHistory
	if err := s.DB.First(&oldTicketHistory, "id = ?", ticketHistory.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket history not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketHistory.Qrcode != "" && ticketHistory.Qrcode != oldTicketHistory.Qrcode {
		updateData["qrcode"] = ticketHistory.Qrcode
	}
	if !ticketHistory.ScannedAt.IsZero() && ticketHistory.ScannedAt != oldTicketHistory.ScannedAt {
		updateData["scanned_at"] = ticketHistory.ScannedAt
	}
	if ticketHistory.ScannedByDevice != "" && ticketHistory.ScannedByDevice != oldTicketHistory.ScannedByDevice {
		updateData["scanned_by_device"] = ticketHistory.ScannedByDevice
	}
	if ticketHistory.IPAddress != "" && ticketHistory.IPAddress != oldTicketHistory.IPAddress {
		updateData["ip_address"] = ticketHistory.IPAddress
	}
	if ticketHistory.Browser != "" && ticketHistory.Browser != oldTicketHistory.Browser {
		updateData["browser"] = ticketHistory.Browser
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketHistory{}).Where("id = ?", ticketHistory.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketHistoryService) DeleteTicketHistory(id string) error {
	return s.DB.Delete(&model.TicketHistory{}, "id = ?", id).Error
}