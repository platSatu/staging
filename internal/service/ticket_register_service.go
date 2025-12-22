package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"log" // Tambahkan import untuk logging
	"regexp"
	"time"

	"github.com/google/uuid" // Tambahkan import untuk generate UUID
	"gorm.io/gorm"
)

type TicketRegisterService struct {
	DB *gorm.DB
}

func NewTicketRegisterService(db *gorm.DB) *TicketRegisterService {
	return &TicketRegisterService{DB: db}
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (s *TicketRegisterService) emailAlreadyExists(email string) (bool, error) {
	var existing model.TicketRegister
	err := s.DB.First(&existing, "email = ?", email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

/* ==========================
   Create Public Register
========================== */

func (s *TicketRegisterService) CreatePublicTicketRegister(
	ticketRegister *model.TicketRegister,
) map[string]string {

	errorsMap := make(map[string]string)

	/* -------- Validasi -------- */

	if len(ticketRegister.FullName) < 3 {
		errorsMap["full_name"] = "Nama lengkap minimal 3 huruf"
	}

	if !isValidEmail(ticketRegister.Email) {
		errorsMap["email"] = "Format email tidak valid"
	} else {
		exists, err := s.emailAlreadyExists(ticketRegister.Email)
		if err != nil {
			errorsMap["email"] = "Terjadi kesalahan database"
		} else if exists {
			errorsMap["email"] = "Email sudah terdaftar"
		}
	}

	if ticketRegister.Handphone != nil {
		h := *ticketRegister.Handphone
		if len(h) < 10 || len(h) > 14 {
			errorsMap["handphone"] = "Nomor handphone harus 10–14 digit"
		}
	}

	if len(errorsMap) > 0 {
		return errorsMap
	}

	/* -------- Generate Data -------- */

	now := time.Now()
	expiredAt := now.Add(24 * time.Hour)

	// ✅ INI YANG SEBELUMNYA HILANG
	ticketRegister.ID = uuid.NewString()

	ticketRegister.Status = "active"
	ticketRegister.UserID = uuid.NewString() // Diperbaiki: Generate UUID baru untuk UserID agar konsisten dengan flow purchase tanpa login
	ticketRegister.TanggalRegister = now
	ticketRegister.PurchaseToken = uuid.NewString()
	ticketRegister.TokenExpiredAt = &expiredAt
	ticketRegister.PurchaseLinkSentAt = &now

	/* -------- Save -------- */

	if err := s.DB.Create(ticketRegister).Error; err != nil {
		errorsMap["general"] = "Gagal menyimpan data"
		return errorsMap
	}

	return nil
}

func (s *TicketRegisterService) ValidatePurchaseToken(token string) (*model.TicketRegister, error) {
	var tr model.TicketRegister

	err := s.DB.
		Where("purchase_token = ?", token).
		First(&tr).Error

	if err != nil {
		return nil, errors.New("token tidak valid")
	}

	if tr.TokenExpiredAt == nil || time.Now().After(*tr.TokenExpiredAt) {
		return nil, errors.New("token sudah kadaluarsa")
	}

	return &tr, nil
}

// CREATE
func (s *TicketRegisterService) CreateTicketRegister(ticketRegister *model.TicketRegister) error {
	log.Printf("Starting CreateTicketRegister with ticketRegister: %+v", ticketRegister)

	// Generate ID otomatis jika kosong
	if ticketRegister.ID == "" {
		ticketRegister.ID = uuid.New().String()
		log.Printf("Generated ID: %s", ticketRegister.ID)
	}

	// Set tanggal register otomatis jika kosong
	if ticketRegister.TanggalRegister.IsZero() {
		ticketRegister.TanggalRegister = time.Now()
		log.Printf("Set TanggalRegister to: %v", ticketRegister.TanggalRegister)
	}

	// Set status default jika kosong
	if ticketRegister.Status == "" {
		ticketRegister.Status = "active"
		log.Printf("Set default Status: %s", ticketRegister.Status)
	}

	log.Printf("Calling DB.Create for ticketRegister: %+v", ticketRegister)
	err := s.DB.Create(ticketRegister).Error
	if err != nil {
		log.Printf("Error in DB.Create: %v", err)
		return err
	}
	log.Printf("Successfully created ticketRegister: %+v", ticketRegister)
	return nil
}

// READ ALL
func (s *TicketRegisterService) GetAllTicketRegisters() ([]model.TicketRegister, error) {
	log.Printf("Starting GetAllTicketRegisters")
	var ticketRegisters []model.TicketRegister
	result := s.DB.Find(&ticketRegisters)
	if result.Error != nil {
		log.Printf("Error in DB.Find for all ticketRegisters: %v", result.Error)
		return nil, result.Error
	}
	log.Printf("Successfully retrieved %d ticketRegisters", len(ticketRegisters))
	return ticketRegisters, nil
}

// READ BY ID
func (s *TicketRegisterService) GetTicketRegisterByID(id string) (*model.TicketRegister, error) {
	log.Printf("Starting GetTicketRegisterByID with id: %s", id)
	var ticketRegister model.TicketRegister
	result := s.DB.First(&ticketRegister, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Printf("TicketRegister not found for id: %s", id)
		return nil, nil
	}
	if result.Error != nil {
		log.Printf("Error in DB.First for id %s: %v", id, result.Error)
		return nil, result.Error
	}
	log.Printf("Successfully retrieved ticketRegister: %+v", ticketRegister)
	return &ticketRegister, nil
}

// UPDATE
func (s *TicketRegisterService) UpdateTicketRegister(ticketRegister *model.TicketRegister) error {
	log.Printf("Starting UpdateTicketRegister with ticketRegister: %+v", ticketRegister)

	if ticketRegister.ID == "" {
		log.Printf("Error: ID is required for update")
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketRegister model.TicketRegister
	log.Printf("Fetching old ticketRegister for id: %s", ticketRegister.ID)
	if err := s.DB.First(&oldTicketRegister, "id = ?", ticketRegister.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Ticket register not found for id: %s", ticketRegister.ID)
			return fmt.Errorf("ticket register not found")
		}
		log.Printf("Error fetching old ticketRegister: %v", err)
		return err
	}
	log.Printf("Old ticketRegister: %+v", oldTicketRegister)

	updateData := map[string]interface{}{}

	if ticketRegister.FullName != "" && ticketRegister.FullName != oldTicketRegister.FullName {
		updateData["full_name"] = ticketRegister.FullName
		log.Printf("Updating full_name to: %s", ticketRegister.FullName)
	}
	if ticketRegister.Email != "" && ticketRegister.Email != oldTicketRegister.Email {
		updateData["email"] = ticketRegister.Email
		log.Printf("Updating email to: %s", ticketRegister.Email)
	}
	if ticketRegister.Handphone != nil && (oldTicketRegister.Handphone == nil || *ticketRegister.Handphone != *oldTicketRegister.Handphone) {
		updateData["handphone"] = ticketRegister.Handphone
		log.Printf("Updating handphone to: %s", *ticketRegister.Handphone)
	}
	if ticketRegister.Status != "" && ticketRegister.Status != oldTicketRegister.Status {
		updateData["status"] = ticketRegister.Status
		log.Printf("Updating status to: %s", ticketRegister.Status)
	}
	// TanggalRegister biasanya tidak diupdate, karena itu tanggal registrasi awal
	// Jika perlu update, uncomment baris di bawah
	// if !ticketRegister.TanggalRegister.IsZero() && ticketRegister.TanggalRegister != oldTicketRegister.TanggalRegister {
	// 	updateData["tanggal_register"] = ticketRegister.TanggalRegister
	// 	log.Printf("Updating tanggal_register to: %v", ticketRegister.TanggalRegister)
	// }

	if len(updateData) == 0 {
		log.Printf("No fields to update for id: %s", ticketRegister.ID)
		return nil // Tidak ada yang diupdate
	}

	log.Printf("Update data: %+v", updateData)
	log.Printf("Calling DB.Updates for id: %s", ticketRegister.ID)
	err := s.DB.Model(&model.TicketRegister{}).Where("id = ?", ticketRegister.ID).Updates(updateData).Error
	if err != nil {
		log.Printf("Error in DB.Updates: %v", err)
		return err
	}
	log.Printf("Successfully updated ticketRegister for id: %s", ticketRegister.ID)
	return nil
}

// DELETE
func (s *TicketRegisterService) DeleteTicketRegister(id string) error {
	log.Printf("Starting DeleteTicketRegister with id: %s", id)
	err := s.DB.Delete(&model.TicketRegister{}, "id = ?", id).Error
	if err != nil {
		log.Printf("Error in DB.Delete for id %s: %v", id, err)
		return err
	}
	log.Printf("Successfully deleted ticketRegister for id: %s", id)
	return nil
}
