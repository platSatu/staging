package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TicketBlastService struct {
	DB *gorm.DB
}

func NewTicketBlastService(db *gorm.DB) *TicketBlastService {
	return &TicketBlastService{DB: db}
}

// CREATE
func (s *TicketBlastService) CreateTicketBlast(ticketBlast *model.TicketBlast) error {
	if ticketBlast.ID == "" {
		ticketBlast.ID = fmt.Sprintf("%s", gorm.Expr("UUID()")) // Atau gunakan uuid.New().String() jika perlu import
	}

	if ticketBlast.Status == "" {
		ticketBlast.Status = "active"
	}

	return s.DB.Create(ticketBlast).Error
}

// READ ALL
func (s *TicketBlastService) GetAllTicketBlasts() ([]model.TicketBlast, error) {
	var ticketBlasts []model.TicketBlast
	result := s.DB.Find(&ticketBlasts)
	return ticketBlasts, result.Error
}

// READ BY ID
func (s *TicketBlastService) GetTicketBlastByID(id string) (*model.TicketBlast, error) {
	var ticketBlast model.TicketBlast
	result := s.DB.First(&ticketBlast, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketBlast, result.Error
}

// UPDATE
func (s *TicketBlastService) UpdateTicketBlast(ticketBlast *model.TicketBlast) error {
	if ticketBlast.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketBlast model.TicketBlast
	if err := s.DB.First(&oldTicketBlast, "id = ?", ticketBlast.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket blast not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketBlast.Nama != "" && ticketBlast.Nama != oldTicketBlast.Nama {
		updateData["nama"] = ticketBlast.Nama
	}
	if ticketBlast.Subject != "" && ticketBlast.Subject != oldTicketBlast.Subject {
		updateData["subject"] = ticketBlast.Subject
	}
	if ticketBlast.CcOrBcc != "" && ticketBlast.CcOrBcc != oldTicketBlast.CcOrBcc {
		updateData["cc_or_bcc"] = ticketBlast.CcOrBcc
	}
	if ticketBlast.StatusPengiriman != "" && ticketBlast.StatusPengiriman != oldTicketBlast.StatusPengiriman {
		updateData["status_pengiriman"] = ticketBlast.StatusPengiriman
	}
	if ticketBlast.Status != "" && ticketBlast.Status != oldTicketBlast.Status {
		updateData["status"] = ticketBlast.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketBlast{}).Where("id = ?", ticketBlast.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketBlastService) DeleteTicketBlast(id string) error {
	return s.DB.Delete(&model.TicketBlast{}, "id = ?", id).Error
}
