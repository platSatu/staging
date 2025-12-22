package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid" // Tambahkan import untuk generate UUID
	"gorm.io/gorm"
)

type TicketVoucherService struct {
	DB *gorm.DB
}

func NewTicketVoucherService(db *gorm.DB) *TicketVoucherService {
	return &TicketVoucherService{DB: db}
}

// CREATE
func (s *TicketVoucherService) CreateTicketVoucher(ticketVoucher *model.TicketVoucher) error {
	// Generate UUID untuk ID jika kosong (karena primary key auto-generated)
	if ticketVoucher.ID == "" {
		ticketVoucher.ID = uuid.New().String()
	}

	// Set default status jika kosong
	if ticketVoucher.Status == "" {
		ticketVoucher.Status = "active"
	}

	return s.DB.Create(ticketVoucher).Error
}

// READ ALL
func (s *TicketVoucherService) GetAllTicketVouchers() ([]model.TicketVoucher, error) {
	var ticketVouchers []model.TicketVoucher
	result := s.DB.Find(&ticketVouchers)
	return ticketVouchers, result.Error
}

// READ BY ID
func (s *TicketVoucherService) GetTicketVoucherByID(id string) (*model.TicketVoucher, error) {
	var ticketVoucher model.TicketVoucher
	result := s.DB.First(&ticketVoucher, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ticketVoucher, result.Error
}

// UPDATE
func (s *TicketVoucherService) UpdateTicketVoucher(ticketVoucher *model.TicketVoucher) error {
	if ticketVoucher.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var old model.TicketVoucher
	if err := s.DB.First(&old, "id = ?", ticketVoucher.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("ticket voucher not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	// Update string
	if ticketVoucher.Nama != "" && ticketVoucher.Nama != old.Nama {
		updateData["nama"] = ticketVoucher.Nama
	}
	if ticketVoucher.KodeVoucher != "" && ticketVoucher.KodeVoucher != old.KodeVoucher {
		updateData["kode_voucher"] = ticketVoucher.KodeVoucher
	}

	// Update time
	if !ticketVoucher.TanggalMulai.IsZero() && !ticketVoucher.TanggalMulai.Equal(old.TanggalMulai) {
		updateData["tanggal_mulai"] = ticketVoucher.TanggalMulai
	}
	if !ticketVoucher.TanggalExpired.IsZero() && !ticketVoucher.TanggalExpired.Equal(old.TanggalExpired) {
		updateData["tanggal_expired"] = ticketVoucher.TanggalExpired
	}

	// Update numeric
	if ticketVoucher.HargaFlat != 0 && ticketVoucher.HargaFlat != old.HargaFlat {
		updateData["harga_flat"] = ticketVoucher.HargaFlat
	}
	if ticketVoucher.Persentase != 0 && ticketVoucher.Persentase != old.Persentase {
		updateData["persentase"] = ticketVoucher.Persentase
	}

	// Update integer
	if ticketVoucher.Quota != 0 && ticketVoucher.Quota != old.Quota {
		updateData["quota"] = ticketVoucher.Quota
	}
	if ticketVoucher.Terpakai != 0 && ticketVoucher.Terpakai != old.Terpakai {
		updateData["terpakai"] = ticketVoucher.Terpakai
	}
	if ticketVoucher.Sisa != 0 && ticketVoucher.Sisa != old.Sisa {
		updateData["sisa"] = ticketVoucher.Sisa
	}

	// Update status
	if ticketVoucher.Status != "" && ticketVoucher.Status != old.Status {
		updateData["status"] = ticketVoucher.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.TicketVoucher{}).
		Where("id = ?", ticketVoucher.ID).
		Updates(updateData).Error
}

// DELETE
func (s *TicketVoucherService) DeleteTicketVoucher(id string) error {
	return s.DB.Delete(&model.TicketVoucher{}, "id = ?", id).Error
}

// CheckVoucherAvailability - Fungsi untuk mengecek ketersediaan voucher berdasarkan kondisi:
// - Status voucher adalah "active"
// - Kode voucher sama dengan yang dicari
// - Quota (sisa) di atas nol (> 0)
// - Tanggal expired belum melebihi tanggal sekarang (tanggal_expired >= now)
// Fungsi ini hanya membaca data tanpa mengubah quota, sehingga bisa dicek berkali-kali selama kondisi terpenuhi.
// Kolom id sebagai primary key tidak digunakan dalam cek ini, hanya kode_voucher.
func (s *TicketVoucherService) CheckVoucherAvailability(kodeVoucher string) (*model.TicketVoucher, error) {
	var voucher model.TicketVoucher
	err := s.DB.Where("kode_voucher = ? AND status = ?", kodeVoucher, "active").First(&voucher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("voucher tidak ditemukan atau tidak aktif")
		}
		return nil, err
	}

	// Cek quota saja
	if voucher.Quota <= 0 {
		return nil, fmt.Errorf("voucher sudah habis")
	}

	now := time.Now()
	if now.Before(voucher.TanggalMulai) {
		return nil, fmt.Errorf("voucher belum aktif")
	}
	if now.After(voucher.TanggalExpired) {
		return nil, fmt.Errorf("voucher sudah expired")
	}

	return &voucher, nil
}
