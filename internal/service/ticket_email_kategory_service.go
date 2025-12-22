package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid" // Tambahkan import untuk generate UUID
	"gorm.io/gorm"
)

type TicketEmailKategoryService struct {
	DB *gorm.DB
}

func NewTicketEmailKategoryService(db *gorm.DB) *TicketEmailKategoryService {
	return &TicketEmailKategoryService{DB: db}
}

// CREATE
func (s *TicketEmailKategoryService) CreateTicketEmailKategory(ticketEmailKategory *model.TicketEmailKategory) error {
	// Generate UUID jika ID kosong untuk konsistensi
	if ticketEmailKategory.ID == "" {
		ticketEmailKategory.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if ticketEmailKategory.Status == "" {
		ticketEmailKategory.Status = "active"
	}

	// Pastikan ID tidak diubah setelah generate (logika konsistensi)
	// GORM akan menggunakan ID yang sudah di-set saat insert
	return s.DB.Create(ticketEmailKategory).Error
}

// READ ALL
func (s *TicketEmailKategoryService) GetAllTicketEmailKategories() ([]model.TicketEmailKategory, error) {
	var ticketEmailKategories []model.TicketEmailKategory
	result := s.DB.Find(&ticketEmailKategories)
	return ticketEmailKategories, result.Error
}

// READ BY ID
func (s *TicketEmailKategoryService) GetTicketEmailKategoryByID(id string) (*model.TicketEmailKategory, error) {
	var ticketEmailKategory model.TicketEmailKategory
	result := s.DB.First(&ticketEmailKategory, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketEmailKategory, result.Error
}

// UPDATE
func (s *TicketEmailKategoryService) UpdateTicketEmailKategory(ticketEmailKategory *model.TicketEmailKategory) error {
	if ticketEmailKategory.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketEmailKategory model.TicketEmailKategory
	if err := s.DB.First(&oldTicketEmailKategory, "id = ?", ticketEmailKategory.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket email kategory not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	// Hanya update field yang berubah dan tidak kosong
	if ticketEmailKategory.Nama != "" && ticketEmailKategory.Nama != oldTicketEmailKategory.Nama {
		updateData["nama"] = ticketEmailKategory.Nama
	}
	if ticketEmailKategory.Keterangan != "" && ticketEmailKategory.Keterangan != oldTicketEmailKategory.Keterangan {
		updateData["keterangan"] = ticketEmailKategory.Keterangan
	}
	if ticketEmailKategory.Status != "" && ticketEmailKategory.Status != oldTicketEmailKategory.Status {
		updateData["status"] = ticketEmailKategory.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	// Update hanya field yang diperlukan, ID tetap tidak diubah untuk konsistensi
	return s.DB.Model(&model.TicketEmailKategory{}).Where("id = ?", ticketEmailKategory.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketEmailKategoryService) DeleteTicketEmailKategory(id string) error {
	return s.DB.Delete(&model.TicketEmailKategory{}, "id = ?", id).Error
}
