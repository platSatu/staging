package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TicketMethodPembayaranService struct {
	DB *gorm.DB
}

func NewTicketMethodPembayaranService(db *gorm.DB) *TicketMethodPembayaranService {
	return &TicketMethodPembayaranService{DB: db}
}

// CREATE
func (s *TicketMethodPembayaranService) CreateTicketMethodPembayaran(ticketMethodPembayaran *model.TicketMethodPembayaran) error {
	if ticketMethodPembayaran.ID == "" {
		ticketMethodPembayaran.ID = fmt.Sprintf("%s", gorm.Expr("UUID()")) // Atau gunakan uuid.New().String() jika perlu import
	}

	if ticketMethodPembayaran.Status == "" {
		ticketMethodPembayaran.Status = "active"
	}

	return s.DB.Create(ticketMethodPembayaran).Error
}

// READ ALL
func (s *TicketMethodPembayaranService) GetAllTicketMethodPembayarans() ([]model.TicketMethodPembayaran, error) {
	var ticketMethodPembayarans []model.TicketMethodPembayaran
	result := s.DB.Find(&ticketMethodPembayarans)
	return ticketMethodPembayarans, result.Error
}

// READ BY ID
func (s *TicketMethodPembayaranService) GetTicketMethodPembayaranByID(id string) (*model.TicketMethodPembayaran, error) {
	var ticketMethodPembayaran model.TicketMethodPembayaran
	result := s.DB.First(&ticketMethodPembayaran, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketMethodPembayaran, result.Error
}

// UPDATE
func (s *TicketMethodPembayaranService) UpdateTicketMethodPembayaran(ticketMethodPembayaran *model.TicketMethodPembayaran) error {
	if ticketMethodPembayaran.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketMethodPembayaran model.TicketMethodPembayaran
	if err := s.DB.First(&oldTicketMethodPembayaran, "id = ?", ticketMethodPembayaran.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket method pembayaran not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketMethodPembayaran.Nama != "" && ticketMethodPembayaran.Nama != oldTicketMethodPembayaran.Nama {
		updateData["nama"] = ticketMethodPembayaran.Nama
	}
	if ticketMethodPembayaran.Status != "" && ticketMethodPembayaran.Status != oldTicketMethodPembayaran.Status {
		updateData["status"] = ticketMethodPembayaran.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketMethodPembayaran{}).Where("id = ?", ticketMethodPembayaran.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketMethodPembayaranService) DeleteTicketMethodPembayaran(id string) error {
	return s.DB.Delete(&model.TicketMethodPembayaran{}, "id = ?", id).Error
}