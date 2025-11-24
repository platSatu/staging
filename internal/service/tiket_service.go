package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TiketService struct {
	DB *gorm.DB
}

func NewTiketService(db *gorm.DB) *TiketService {
	return &TiketService{DB: db}
}

// CREATE
func (s *TiketService) CreateTiket(tiket *model.Tiket) error {
	if tiket.ID == "" {
		tiket.ID = uuid.New().String()
	}

	if !tiket.IsScanned {
		tiket.IsScanned = false
	}

	// Validasi user_id jika ada
	if tiket.UserID != nil {
		var user model.User
		if err := s.DB.First(&user, "id = ?", *tiket.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
	}

	// Validasi event_id ada di tabel events
	var event model.Event
	if err := s.DB.First(&event, "id = ?", tiket.EventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event not found")
		}
		return err
	}

	// Validasi jenis_ticket_id ada di tabel jenis_tiket
	var jenisTiket model.JenisTiket
	if err := s.DB.First(&jenisTiket, "id = ?", tiket.JenisTicketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("jenis tiket not found")
		}
		return err
	}

	return s.DB.Create(tiket).Error
}

// READ ALL
func (s *TiketService) GetAllTiket() ([]model.Tiket, error) {
	var tikets []model.Tiket
	result := s.DB.Find(&tikets)
	return tikets, result.Error
}

// READ BY ID
func (s *TiketService) GetTiketByID(id string) (*model.Tiket, error) {
	var tiket model.Tiket
	result := s.DB.First(&tiket, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tiket, result.Error
}

// READ BY USER ID
func (s *TiketService) GetTiketByUserID(userID string) ([]model.Tiket, error) {
	var tikets []model.Tiket
	result := s.DB.Where("user_id = ?", userID).Find(&tikets)
	return tikets, result.Error
}

// READ BY EVENT ID
func (s *TiketService) GetTiketByEventID(eventID string) ([]model.Tiket, error) {
	var tikets []model.Tiket
	result := s.DB.Where("event_id = ?", eventID).Find(&tikets)
	return tikets, result.Error
}

// READ BY KODE BOOKING
func (s *TiketService) GetTiketByKodeBooking(kodeBooking string) (*model.Tiket, error) {
	var tiket model.Tiket
	result := s.DB.Where("kode_booking = ?", kodeBooking).First(&tiket)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &tiket, result.Error
}

// UPDATE
func (s *TiketService) UpdateTiket(tiket *model.Tiket) error {
	if tiket.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTiket model.Tiket
	if err := s.DB.First(&oldTiket, "id = ?", tiket.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("tiket not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if tiket.Name != "" && tiket.Name != oldTiket.Name {
		updateData["name"] = tiket.Name
	}

	if tiket.Email != "" && tiket.Email != oldTiket.Email {
		updateData["email"] = tiket.Email
	}

	if tiket.Phone != "" && tiket.Phone != oldTiket.Phone {
		updateData["phone"] = tiket.Phone
	}

	if tiket.OrderID != "" && tiket.OrderID != oldTiket.OrderID {
		updateData["order_id"] = tiket.OrderID
	}

	if tiket.Quantity != 0 && tiket.Quantity != oldTiket.Quantity {
		updateData["quantity"] = tiket.Quantity
	}

	if tiket.Price != 0 && tiket.Price != oldTiket.Price {
		updateData["price"] = tiket.Price
	}

	if tiket.Subtotal != 0 && tiket.Subtotal != oldTiket.Subtotal {
		updateData["subtotal"] = tiket.Subtotal
	}

	if tiket.Total != 0 && tiket.Total != oldTiket.Total {
		updateData["total"] = tiket.Total
	}

	if tiket.Qrcode != "" && tiket.Qrcode != oldTiket.Qrcode {
		updateData["qrcode"] = tiket.Qrcode
	}

	if tiket.QrcodeDirectory != "" && tiket.QrcodeDirectory != oldTiket.QrcodeDirectory {
		updateData["qrcode_directory"] = tiket.QrcodeDirectory
	}

	if tiket.PaymentStatus != "" && tiket.PaymentStatus != oldTiket.PaymentStatus {
		updateData["payment_status"] = tiket.PaymentStatus
	}

	if tiket.PaymentMethod != "" && tiket.PaymentMethod != oldTiket.PaymentMethod {
		updateData["payment_method"] = tiket.PaymentMethod
	}

	if tiket.PaymentDate != nil && (oldTiket.PaymentDate == nil || !tiket.PaymentDate.Equal(*oldTiket.PaymentDate)) {
		updateData["payment_date"] = tiket.PaymentDate
	}

	if tiket.IsScanned != oldTiket.IsScanned {
		updateData["is_scanned"] = tiket.IsScanned
	}

	if tiket.ScanDate != nil && (oldTiket.ScanDate == nil || !tiket.ScanDate.Equal(*oldTiket.ScanDate)) {
		updateData["scan_date"] = tiket.ScanDate
	}

	if tiket.KodeBooking != "" && tiket.KodeBooking != oldTiket.KodeBooking {
		updateData["kode_booking"] = tiket.KodeBooking
	}

	if tiket.UserID != nil && (oldTiket.UserID == nil || *tiket.UserID != *oldTiket.UserID) {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", *tiket.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = tiket.UserID
	}

	if tiket.EventID != "" && tiket.EventID != oldTiket.EventID {
		// Validasi event_id baru ada
		var event model.Event
		if err := s.DB.First(&event, "id = ?", tiket.EventID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("event not found")
			}
			return err
		}
		updateData["event_id"] = tiket.EventID
	}

	if tiket.JenisTicketID != "" && tiket.JenisTicketID != oldTiket.JenisTicketID {
		// Validasi jenis_ticket_id baru ada
		var jenisTiket model.JenisTiket
		if err := s.DB.First(&jenisTiket, "id = ?", tiket.JenisTicketID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("jenis tiket not found")
			}
			return err
		}
		updateData["jenis_ticket_id"] = tiket.JenisTicketID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Tiket{}).Where("id = ?", tiket.ID).Updates(updateData).Error
}

// DELETE
func (s *TiketService) DeleteTiket(id string) error {
	return s.DB.Delete(&model.Tiket{}, "id = ?", id).Error
}
