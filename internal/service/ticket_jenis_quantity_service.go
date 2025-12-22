package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"log" // Tambahkan import untuk logging

	"github.com/google/uuid" // Tambahkan import untuk generate UUID
	"gorm.io/gorm"
)

type TicketJenisQuantityService struct {
	DB *gorm.DB
}

func NewTicketJenisQuantityService(db *gorm.DB) *TicketJenisQuantityService {
	return &TicketJenisQuantityService{DB: db}
}

// CREATE
func (s *TicketJenisQuantityService) CreateTicketJenisQuantity(ticketJenisQuantity *model.TicketJenisQuantity) error {
	log.Printf("Starting CreateTicketJenisQuantity with ticketJenisQuantity: %+v", ticketJenisQuantity)

	// Generate ID otomatis jika kosong
	if ticketJenisQuantity.ID == "" {
		ticketJenisQuantity.ID = uuid.New().String()
		log.Printf("Generated ID: %s", ticketJenisQuantity.ID)
	}

	// UserID harus di-set dari session login di controller, tapi jika kosong, warning
	if ticketJenisQuantity.UserID == "" {
		log.Printf("Warning: UserID is empty, should be set from logged-in user")
	}

	// Set status default jika kosong
	if ticketJenisQuantity.Status == "" {
		ticketJenisQuantity.Status = "active"
		log.Printf("Set default Status: %s", ticketJenisQuantity.Status)
	}

	log.Printf("Calling DB.Create for ticketJenisQuantity: %+v", ticketJenisQuantity)
	err := s.DB.Create(ticketJenisQuantity).Error
	if err != nil {
		log.Printf("Error in DB.Create: %v", err)
		return err
	}
	log.Printf("Successfully created ticketJenisQuantity: %+v", ticketJenisQuantity)
	return nil
}

// READ ALL
func (s *TicketJenisQuantityService) GetAllTicketJenisQuantities() ([]model.TicketJenisQuantity, error) {
	log.Printf("Starting GetAllTicketJenisQuantities")
	var ticketJenisQuantities []model.TicketJenisQuantity
	result := s.DB.Find(&ticketJenisQuantities)
	if result.Error != nil {
		log.Printf("Error in DB.Find for all ticketJenisQuantities: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Successfully retrieved %d ticketJenisQuantities", len(ticketJenisQuantities))
	return ticketJenisQuantities, nil
}

// READ BY ID
func (s *TicketJenisQuantityService) GetTicketJenisQuantityByID(id string) (*model.TicketJenisQuantity, error) {
	log.Printf("Starting GetTicketJenisQuantityByID with id: %s", id)
	var ticketJenisQuantity model.TicketJenisQuantity
	result := s.DB.First(&ticketJenisQuantity, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("TicketJenisQuantity not found for id: %s", id)
		return nil, nil
	}
	if result.Error != nil {
		log.Printf("Error in DB.First for id %s: %v", id, result.Error)
		return nil, result.Error
	}
	log.Printf("Successfully retrieved ticketJenisQuantity: %+v", ticketJenisQuantity)
	return &ticketJenisQuantity, nil
}

// UPDATE
func (s *TicketJenisQuantityService) UpdateTicketJenisQuantity(ticketJenisQuantity *model.TicketJenisQuantity) error {
	log.Printf("Starting UpdateTicketJenisQuantity with ticketJenisQuantity: %+v", ticketJenisQuantity)

	if ticketJenisQuantity.ID == "" {
		log.Printf("Error: ID is required for update")
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketJenisQuantity model.TicketJenisQuantity
	log.Printf("Fetching old ticketJenisQuantity for id: %s", ticketJenisQuantity.ID)
	if err := s.DB.First(&oldTicketJenisQuantity, "id = ?", ticketJenisQuantity.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Ticket jenis quantity not found for id: %s", ticketJenisQuantity.ID)
			return fmt.Errorf("ticket jenis quantity not found")
		}
		log.Printf("Error fetching old ticketJenisQuantity: %v", err)
		return err
	}
	log.Printf("Old ticketJenisQuantity: %+v", oldTicketJenisQuantity)

	updateData := map[string]interface{}{}

	if ticketJenisQuantity.Name != "" && ticketJenisQuantity.Name != oldTicketJenisQuantity.Name {
		updateData["name"] = ticketJenisQuantity.Name
		log.Printf("Updating name to: %s", ticketJenisQuantity.Name)
	}
	if ticketJenisQuantity.Status != "" && ticketJenisQuantity.Status != oldTicketJenisQuantity.Status {
		updateData["status"] = ticketJenisQuantity.Status
		log.Printf("Updating status to: %s", ticketJenisQuantity.Status)
	}

	if len(updateData) == 0 {
		log.Printf("No fields to update for id: %s", ticketJenisQuantity.ID)
		return nil // Tidak ada yang diupdate
	}

	log.Printf("Update data: %+v", updateData)
	log.Printf("Calling DB.Updates for id: %s", ticketJenisQuantity.ID)
	err := s.DB.Model(&model.TicketJenisQuantity{}).Where("id = ?", ticketJenisQuantity.ID).Updates(updateData).Error
	if err != nil {
		log.Printf("Error in DB.Updates: %v", err)
		return err
	}
	log.Printf("Successfully updated ticketJenisQuantity for id: %s", ticketJenisQuantity.ID)
	return nil
}

// DELETE
func (s *TicketJenisQuantityService) DeleteTicketJenisQuantity(id string) error {
	log.Printf("Starting DeleteTicketJenisQuantity with id: %s", id)
	err := s.DB.Delete(&model.TicketJenisQuantity{}, "id = ?", id).Error
	if err != nil {
		log.Printf("Error in DB.Delete for id %s: %v", id, err)
		return err
	}
	log.Printf("Successfully deleted ticketJenisQuantity for id: %s", id)
	return nil
}
