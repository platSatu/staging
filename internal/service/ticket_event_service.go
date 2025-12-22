package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"log" // Tambahkan import untuk logging

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type TicketEventService struct {
	DB *gorm.DB
}

func NewTicketEventService(db *gorm.DB) *TicketEventService {
	return &TicketEventService{DB: db}
}

// helper untuk generate slug dari nama (menggunakan library slug)
func generateTicketEventSlug(nama string) string {
	log.Printf("Generating slug for nama: %s", nama)
	slug := slug.Make(nama)
	log.Printf("Generated slug: %s", slug)
	return slug
}

// CREATE
func (s *TicketEventService) CreateTicketEvent(ticketEvent *model.TicketEvent) error {
	log.Printf("Starting CreateTicketEvent with ticketEvent: %+v", ticketEvent)

	if ticketEvent.Slug == "" && ticketEvent.Nama != "" {
		ticketEvent.Slug = generateTicketEventSlug(ticketEvent.Nama)
		log.Printf("Slug generated and set: %s", ticketEvent.Slug)
	}

	// ID dan Status akan di-handle oleh BeforeCreate di model
	log.Printf("Calling DB.Create for ticketEvent: %+v", ticketEvent)
	err := s.DB.Create(ticketEvent).Error
	if err != nil {
		log.Printf("Error in DB.Create: %v", err)
		return err
	}
	log.Printf("Successfully created ticketEvent: %+v", ticketEvent)
	return nil
}

// READ ALL
func (s *TicketEventService) GetAllTicketEvents() ([]model.TicketEvent, error) {
	log.Printf("Starting GetAllTicketEvents")
	var ticketEvents []model.TicketEvent
	result := s.DB.Find(&ticketEvents)
	if result.Error != nil {
		log.Printf("Error in DB.Find for all ticketEvents: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Successfully retrieved %d ticketEvents", len(ticketEvents))
	return ticketEvents, nil
}

// READ BY ID
func (s *TicketEventService) GetTicketEventByID(id string) (*model.TicketEvent, error) {
	log.Printf("Starting GetTicketEventByID with id: %s", id)
	var ticketEvent model.TicketEvent
	result := s.DB.First(&ticketEvent, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("TicketEvent not found for id: %s", id)
		return nil, nil
	}
	if result.Error != nil {
		log.Printf("Error in DB.First for id %s: %v", id, result.Error)
		return nil, result.Error
	}
	log.Printf("Successfully retrieved ticketEvent: %+v", ticketEvent)
	return &ticketEvent, nil
}

// UPDATE - Perbaiki untuk menerima ID sebagai parameter terpisah, karena ID tidak dikirim di body JSON dari frontend
func (s *TicketEventService) UpdateTicketEvent(id string, ticketEvent *model.TicketEvent) error {
	log.Printf("Starting UpdateTicketEvent with id: %s and ticketEvent: %+v", id, ticketEvent)

	if id == "" {
		log.Printf("Error: ID is required for update")
		return fmt.Errorf("ID is required for update")
	}

	// Set ID dari parameter URL ke struct (karena frontend tidak kirim ID di body)
	ticketEvent.ID = id
	log.Printf("Set ticketEvent.ID to: %s", id)

	var oldTicketEvent model.TicketEvent
	log.Printf("Fetching old ticketEvent for id: %s", id)
	if err := s.DB.First(&oldTicketEvent, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Ticket event not found for id: %s", id)
			return fmt.Errorf("ticket event not found")
		}
		log.Printf("Error fetching old ticketEvent: %v", err)
		return err
	}
	log.Printf("Old ticketEvent: %+v", oldTicketEvent)

	updateData := map[string]interface{}{}

	// Jika nama dikirim (tidak kosong), update nama dan slug otomatis dari nama
	if ticketEvent.Nama != "" {
		updateData["nama"] = ticketEvent.Nama
		updateData["slug"] = generateTicketEventSlug(ticketEvent.Nama)
		log.Printf("Updating nama to: %s and slug to: %s", ticketEvent.Nama, updateData["slug"])
	} else if ticketEvent.Slug != "" && ticketEvent.Slug != oldTicketEvent.Slug {
		// Update slug manual hanya jika nama tidak dikirim dan slug berbeda
		updateData["slug"] = ticketEvent.Slug
		log.Printf("Updating slug to: %s", ticketEvent.Slug)
	}
	if ticketEvent.Keterangan != "" && ticketEvent.Keterangan != oldTicketEvent.Keterangan {
		updateData["keterangan"] = ticketEvent.Keterangan
		log.Printf("Updating keterangan to: %s", ticketEvent.Keterangan)
	}
	if !ticketEvent.TanggalEvent.IsZero() && ticketEvent.TanggalEvent != oldTicketEvent.TanggalEvent {
		updateData["tanggal_event"] = ticketEvent.TanggalEvent
		log.Printf("Updating tanggal_event to: %v", ticketEvent.TanggalEvent)
	}
	if ticketEvent.Alamat != "" && ticketEvent.Alamat != oldTicketEvent.Alamat {
		updateData["alamat"] = ticketEvent.Alamat
		log.Printf("Updating alamat to: %s", ticketEvent.Alamat)
	}
	if ticketEvent.Status != "" && ticketEvent.Status != oldTicketEvent.Status {
		updateData["status"] = ticketEvent.Status
		log.Printf("Updating status to: %s", ticketEvent.Status)
	}

	if len(updateData) == 0 {
		log.Printf("No fields to update for id: %s", id)
		return nil // Tidak ada yang diupdate
	}

	log.Printf("Update data: %+v", updateData)
	log.Printf("Calling DB.Updates for id: %s", id)
	err := s.DB.Model(&model.TicketEvent{}).Where("id = ?", id).Updates(updateData).Error
	if err != nil {
		log.Printf("Error in DB.Updates: %v", err)
		return err
	}
	log.Printf("Successfully updated ticketEvent for id: %s", id)
	return nil
}

// DELETE
func (s *TicketEventService) DeleteTicketEvent(id string) error {
	log.Printf("Starting DeleteTicketEvent with id: %s", id)
	err := s.DB.Delete(&model.TicketEvent{}, "id = ?", id).Error
	if err != nil {
		log.Printf("Error in DB.Delete for id %s: %v", id, err)
		return err
	}
	log.Printf("Successfully deleted ticketEvent for id: %s", id)
	return nil
}
