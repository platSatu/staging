package service

import (
	"errors"
	"log"
	"time"

	"backend_go/internal/model" // Ganti dengan path import model Anda

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RegistrasiUlangService struct {
	DB *gorm.DB
}

func NewRegistrasiUlangService(db *gorm.DB) *RegistrasiUlangService {
	return &RegistrasiUlangService{DB: db}
}

// Fungsi RegistrasiUlang yang diperbarui: Mengetahui kategori ticket dan input ke TicketHistory
// Ditambahkan transaksi untuk atomicity dan mencegah race condition pada traffic tinggi
func (s *RegistrasiUlangService) RegistrasiUlang(qrcode string, scannedByUser string, scannedByDevice string, ipAddress string, browser string) (kategori string, err error) {
	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var ticket model.TicketQrcode

	// Cari berdasarkan qrcode
	if err := tx.Where("qrcode = ?", qrcode).First(&ticket).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("qrcode tidak ada di sistem kami")
		}
		return "", err
	}

	// Cek jika sudah di-scan
	if ticket.IsScanned {
		tx.Rollback()
		return "", errors.New("qrcode sudah melakukan registrasi ulang")
	}

	// Dapatkan kategori ticket (asumsi TicketQrcode memiliki field TicketKategoriID yang link ke TicketKategori.ID)
	var ticketKategori model.TicketKategori
	if err := tx.Where("id = ?", ticket.TicketKategoriID).First(&ticketKategori).Error; err != nil { // Ganti ticket.TicketKategoriID dengan field yang sesuai jika berbeda
		tx.Rollback()
		return "", errors.New("kategori ticket tidak ditemukan")
	}
	kategori = ticketKategori.Nama // Ambil nama kategori

	// Update is_scanned dan date_scanned pada TicketQrcode
	now := time.Now()
	if err := tx.Model(&model.TicketQrcode{}).Where("id = ?", ticket.ID).
		Updates(map[string]interface{}{
			"is_scanned":   true,
			"date_scanned": now,
		}).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	// Input ke TicketHistory
	ticketHistory := model.TicketHistory{
		UserID:          ticket.UserID,   // Ambil dari ticket yang di-scan (user pemilik tiket)
		ParentID:        ticket.ParentID, // Ambil dari ticket yang di-scan (sesuai model)
		Qrcode:          qrcode,
		ScannedAt:       now,
		ScannedByUser:   scannedByUser,   // Siapa yang scan (user yang login, dari parameter)
		ScannedByDevice: scannedByDevice, // Device yang digunakan (dari parameter)
		IPAddress:       ipAddress,       // IP address (dari parameter)
		Browser:         browser,         // Browser (dari parameter)
	}
	// Pastikan ID di-set sebelum create (untuk mencegah error duplikasi primary key kosong)
	if ticketHistory.ID == "" {
		// Generate UUID manual (import "github.com/google/uuid" di file Anda)
		ticketHistory.ID = uuid.New().String()
	}
	if err := tx.Create(&ticketHistory).Error; err != nil {
		// Log error untuk debug
		log.Printf("Error creating TicketHistory: %v, TicketHistory ID: %s", err, ticketHistory.ID)
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return kategori, nil
}

// Fungsi baru: GetAllTicketKategori - Mengambil seluruh data ticket kategori
// Mengembalikan slice dari TicketKategori, dengan sorting berdasarkan CreatedAt (terbaru dulu)
func (s *RegistrasiUlangService) GetAllTicketKategori() ([]model.TicketKategori, error) {
	var categories []model.TicketKategori

	// Query untuk mengambil semua data, dengan sorting
	if err := s.DB.Order("created_at DESC").Find(&categories).Error; err != nil {
		return nil, err
	}

	// Tambahkan logika untuk menghitung sisa tiket yang belum discan
	// Asumsi: model.TicketKategori memiliki field Terjual (int, jumlah tiket terjual) dan Sisa (int, yang akan dihitung)
	// Hitung jumlah TicketQrcode yang sudah discan (is_scanned = true) per kategori, lalu kurangi dari Terjual
	for i := range categories {
		var scannedCount int64
		// Hitung jumlah tiket yang sudah discan (is_scanned = true) untuk kategori ini
		if err := s.DB.Model(&model.TicketQrcode{}).Where("ticket_kategori_id = ? AND is_scanned = ?", categories[i].ID, true).Count(&scannedCount).Error; err != nil {
			return nil, err
		}
		// Hitung sisa: terjual dikurangi dengan jumlah yang sudah discan
		categories[i].Sisa = categories[i].Terjual - int(scannedCount)
	}

	return categories, nil
}
