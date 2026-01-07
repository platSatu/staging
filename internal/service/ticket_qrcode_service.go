package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TicketQrcodeService struct {
	DB *gorm.DB
}

func NewTicketQrcodeService(db *gorm.DB) *TicketQrcodeService {
	return &TicketQrcodeService{DB: db}
}

// CREATE
func (s *TicketQrcodeService) CreateTicketQrcode(ticketQrcode *model.TicketQrcode) error {
	if ticketQrcode.ID == "" {
		ticketQrcode.ID = fmt.Sprintf("%s", gorm.Expr("UUID()")) // Atau gunakan uuid.New().String() jika perlu import
	}

	return s.DB.Create(ticketQrcode).Error
}

// READ ALL
func (s *TicketQrcodeService) GetAllTicketQrcodes() ([]model.TicketQrcode, error) {
	var ticketQrcodes []model.TicketQrcode
	result := s.DB.Order("created_at DESC").Find(&ticketQrcodes)
	return ticketQrcodes, result.Error
}

// READ BY ID
func (s *TicketQrcodeService) GetTicketQrcodeByID(id string) (*model.TicketQrcode, error) {
	var ticketQrcode model.TicketQrcode
	result := s.DB.First(&ticketQrcode, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketQrcode, result.Error
}

// UPDATE
func (s *TicketQrcodeService) UpdateTicketQrcode(ticketQrcode *model.TicketQrcode) error {
	if ticketQrcode.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketQrcode model.TicketQrcode
	if err := s.DB.First(&oldTicketQrcode, "id = ?", ticketQrcode.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket qrcode not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketQrcode.Quantity != 0 && ticketQrcode.Quantity != oldTicketQrcode.Quantity {
		updateData["quantity"] = ticketQrcode.Quantity
	}
	if ticketQrcode.Price != 0 && ticketQrcode.Price != oldTicketQrcode.Price {
		updateData["price"] = ticketQrcode.Price
	}
	if ticketQrcode.OrderID != "" && ticketQrcode.OrderID != oldTicketQrcode.OrderID {
		updateData["order_id"] = ticketQrcode.OrderID
	}
	if ticketQrcode.KodeBooking != "" && ticketQrcode.KodeBooking != oldTicketQrcode.KodeBooking {
		updateData["kode_booking"] = ticketQrcode.KodeBooking
	}
	if ticketQrcode.Qrcode != "" && ticketQrcode.Qrcode != oldTicketQrcode.Qrcode {
		updateData["qrcode"] = ticketQrcode.Qrcode
	}
	if ticketQrcode.DirectoryQrcode != "" && ticketQrcode.DirectoryQrcode != oldTicketQrcode.DirectoryQrcode {
		updateData["directory_qrcode"] = ticketQrcode.DirectoryQrcode
	}
	if ticketQrcode.PaymentStatus != "" && ticketQrcode.PaymentStatus != oldTicketQrcode.PaymentStatus {
		updateData["payment_status"] = ticketQrcode.PaymentStatus
	}
	if ticketQrcode.PaymentMethod != "" && ticketQrcode.PaymentMethod != oldTicketQrcode.PaymentMethod {
		updateData["payment_method"] = ticketQrcode.PaymentMethod
	}
	if !ticketQrcode.PaymentDate.IsZero() && ticketQrcode.PaymentDate != oldTicketQrcode.PaymentDate {
		updateData["payment_date"] = ticketQrcode.PaymentDate
	}
	if ticketQrcode.IsScanned != oldTicketQrcode.IsScanned {
		updateData["is_scanned"] = ticketQrcode.IsScanned
	}
	if !ticketQrcode.DateScanned.IsZero() && ticketQrcode.DateScanned != oldTicketQrcode.DateScanned {
		updateData["date_scanned"] = ticketQrcode.DateScanned
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketQrcode{}).Where("id = ?", ticketQrcode.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketQrcodeService) DeleteTicketQrcode(id string) error {
	return s.DB.Delete(&model.TicketQrcode{}, "id = ?", id).Error
}
