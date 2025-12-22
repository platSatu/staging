package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid" // Tambahkan import untuk generate UUID
	"gorm.io/gorm"
)

type TicketTemplateService struct {
	DB *gorm.DB
}

func NewTicketTemplateService(db *gorm.DB) *TicketTemplateService {
	return &TicketTemplateService{DB: db}
}

// CREATE
func (s *TicketTemplateService) CreateTicketTemplate(ticketTemplate *model.TicketTemplate) error {
	// Generate UUID jika ID kosong untuk konsistensi
	if ticketTemplate.ID == "" {
		ticketTemplate.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if ticketTemplate.Status == "" {
		ticketTemplate.Status = "active"
	}

	// Pastikan ID tidak diubah setelah generate (logika konsistensi)
	// GORM akan menggunakan ID yang sudah di-set saat insert
	return s.DB.Create(ticketTemplate).Error
}

// READ ALL
func (s *TicketTemplateService) GetAllTicketTemplates() ([]model.TicketTemplate, error) {
	var ticketTemplates []model.TicketTemplate
	result := s.DB.Find(&ticketTemplates)
	return ticketTemplates, result.Error
}

// READ BY ID
func (s *TicketTemplateService) GetTicketTemplateByID(id string) (*model.TicketTemplate, error) {
	var ticketTemplate model.TicketTemplate
	result := s.DB.First(&ticketTemplate, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketTemplate, result.Error
}

// UPDATE
func (s *TicketTemplateService) UpdateTicketTemplate(ticketTemplate *model.TicketTemplate) error {
	if ticketTemplate.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketTemplate model.TicketTemplate
	if err := s.DB.First(&oldTicketTemplate, "id = ?", ticketTemplate.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket template not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	// Hanya update field yang berubah dan tidak kosong
	if ticketTemplate.Template != "" && ticketTemplate.Template != oldTicketTemplate.Template {
		updateData["template"] = ticketTemplate.Template
	}
	if ticketTemplate.Status != "" && ticketTemplate.Status != oldTicketTemplate.Status {
		updateData["status"] = ticketTemplate.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	// Update hanya field yang diperlukan, ID tetap tidak diubah untuk konsistensi
	return s.DB.Model(&model.TicketTemplate{}).Where("id = ?", ticketTemplate.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketTemplateService) DeleteTicketTemplate(id string) error {
	return s.DB.Delete(&model.TicketTemplate{}, "id = ?", id).Error
}
