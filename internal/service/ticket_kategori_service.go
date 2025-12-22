package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid" // Tambahkan import untuk UUID
	"gorm.io/gorm"
)

type TicketKategoriService struct {
	DB *gorm.DB
}

func NewTicketKategoriService(db *gorm.DB) *TicketKategoriService {
	return &TicketKategoriService{DB: db}
}

// CREATE
func (s *TicketKategoriService) CreateTicketKategori(ticketKategori *model.TicketKategori) error {
	if ticketKategori.ID == "" {
		ticketKategori.ID = uuid.New().String() // Generate UUID baru
	}

	if ticketKategori.Status == "" {
		ticketKategori.Status = "active"
	}

	return s.DB.Create(ticketKategori).Error
}

// READ ALL
func (s *TicketKategoriService) GetAllTicketKategoris() ([]model.TicketKategori, error) {
	var ticketKategoris []model.TicketKategori
	result := s.DB.Find(&ticketKategoris)
	return ticketKategoris, result.Error
}

// READ BY ID
func (s *TicketKategoriService) GetTicketKategoriByID(id string) (*model.TicketKategori, error) {
	var ticketKategori model.TicketKategori
	result := s.DB.First(&ticketKategori, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketKategori, result.Error
}

// UPDATE
func (s *TicketKategoriService) UpdateTicketKategori(ticketKategori *model.TicketKategori) error {
	if ticketKategori.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTicketKategori model.TicketKategori
	if err := s.DB.First(&oldTicketKategori, "id = ?", ticketKategori.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket kategori not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if ticketKategori.Nama != "" && ticketKategori.Nama != oldTicketKategori.Nama {
		updateData["nama"] = ticketKategori.Nama
	}
	if ticketKategori.Image != "" && ticketKategori.Image != oldTicketKategori.Image {
		updateData["image"] = ticketKategori.Image // Tambahkan logika update untuk image
	}
	if ticketKategori.StokAwal != 0 && ticketKategori.StokAwal != oldTicketKategori.StokAwal {
		updateData["stok_awal"] = ticketKategori.StokAwal
	}
	// Biarkan terjual nullable: update jika berbeda, termasuk 0
	if ticketKategori.Terjual != oldTicketKategori.Terjual {
		updateData["terjual"] = ticketKategori.Terjual
	}
	// Biarkan sisa nullable: update jika berbeda, termasuk 0
	if ticketKategori.Sisa != oldTicketKategori.Sisa {
		updateData["sisa"] = ticketKategori.Sisa
	}
	if ticketKategori.Harga != 0 && ticketKategori.Harga != oldTicketKategori.Harga {
		updateData["harga"] = ticketKategori.Harga
	}

	// Tambahkan logika update untuk min_quantity
	if ticketKategori.MinQuantity != 0 && ticketKategori.MinQuantity != oldTicketKategori.MinQuantity {
		updateData["min_quantity"] = ticketKategori.MinQuantity
	}

	// Tambahkan logika update untuk max_quantity
	if ticketKategori.MaxQuantity != 0 && ticketKategori.MaxQuantity != oldTicketKategori.MaxQuantity {
		updateData["max_quantity"] = ticketKategori.MaxQuantity
	}

	if ticketKategori.Status != "" && ticketKategori.Status != oldTicketKategori.Status {
		updateData["status"] = ticketKategori.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketKategori{}).Where("id = ?", ticketKategori.ID).Updates(updateData).Error
}

// DELETE
func (s *TicketKategoriService) DeleteTicketKategori(id string) error {
	return s.DB.Delete(&model.TicketKategori{}, "id = ?", id).Error
}
