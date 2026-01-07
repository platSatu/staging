package service

import (
	"backend_go/internal/model"
	"path/filepath"

	"gorm.io/gorm"
)

type TicketUserService struct {
	db *gorm.DB
}

func NewTicketUserService(db *gorm.DB) *TicketUserService {
	return &TicketUserService{db: db}
}

// GetTicketQrcodesByUserID mengambil semua ticket qrcodes berdasarkan UserID
// func (s *TicketUserService) GetTicketQrcodesByUserID(userID string) ([]model.TicketQrcode, error) {
// 	var tickets []model.TicketQrcode
// 	err := s.db.Where("user_id = ?", userID).Find(&tickets).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return tickets, nil
// }

// func (s *TicketUserService) GetTicketQrcodesByUserID(userID string) ([]model.TicketQrcode, error) {
// 	var tickets []model.TicketQrcode
// 	err := s.db.Where("user_id = ?", userID).Find(&tickets).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Ubah path lokal menjadi URL publik
// 	for i := range tickets {
// 		if tickets[i].DirectoryQrcode != "" {
// 			// Ambil nama file saja
// 			fileName := filepath.Base(tickets[i].DirectoryQrcode)
// 			// Buat URL publik
// 			tickets[i].DirectoryQrcode = "/qrcodes/" + fileName
// 		}
// 	}

//		return tickets, nil
//	}
func (s *TicketUserService) GetTicketQrcodesByUserID(userID string) ([]model.TicketQrcode, error) {
	var tickets []model.TicketQrcode
	err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	// Ubah path lokal menjadi URL publik
	for i := range tickets {
		if tickets[i].DirectoryQrcode != "" {
			// Ambil nama file saja
			fileName := filepath.Base(tickets[i].DirectoryQrcode)
			// Buat URL publik
			tickets[i].DirectoryQrcode = "/qrcodes/" + fileName
		}
	}

	return tickets, nil
}
