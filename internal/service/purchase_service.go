package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/skip2/go-qrcode"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseService struct {
	DB *gorm.DB
}

func NewPurchaseService(db *gorm.DB) *PurchaseService {
	return &PurchaseService{DB: db}
}

// Struct untuk request pembelian (sesuai kode React)
type PurchaseRequest struct {
	TicketJenisID string         `json:"ticket_jenis_id"`
	Items         []PurchaseItem `json:"items"`
	VoucherCode   *string        `json:"voucher_code,omitempty"`
}

type PurchaseItem struct {
	TicketKategoriID string `json:"ticket_kategori_id"`
	Quantity         int    `json:"quantity"`
	FeeID            string `json:"fee_id"`
}

// generateUniqueOrderID: Generate order_id unik (contoh: PS666666662025)
// Format: PS + 9 digit angka acak + tahun (4 digit)
// Cek duplikat di database dengan max attempts
func (s *PurchaseService) generateUniqueOrderID(tx *gorm.DB) (string, error) {
	year := strconv.Itoa(time.Now().Year())
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		randomNum := rand.Intn(900000000) + 100000000 // 9 digit angka (100000000 - 999999999)
		orderID := "PS" + strconv.Itoa(randomNum) + year

		var count int64
		if err := tx.Model(&model.TicketQrcode{}).Where("order_id = ?", orderID).Count(&count).Error; err != nil {
			return "", fmt.Errorf("error checking order_id uniqueness: %v", err)
		}
		if count == 0 {
			return orderID, nil
		}
	}
	return "", errors.New("failed to generate unique order_id after max attempts")
}

// generateUniqueKodeBooking: Generate kode_booking unik (8 digit angka, contoh: 77778888)
// Cek duplikat di database dengan max attempts
func (s *PurchaseService) generateUniqueKodeBooking(tx *gorm.DB) (string, error) {
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		randomNum := rand.Intn(90000000) + 10000000 // 8 digit angka (10000000 - 99999999)
		kodeBooking := strconv.Itoa(randomNum)

		var count int64
		if err := tx.Model(&model.TicketQrcode{}).Where("kode_booking = ?", kodeBooking).Count(&count).Error; err != nil {
			return "", fmt.Errorf("error checking kode_booking uniqueness: %v", err)
		}
		if count == 0 {
			return kodeBooking, nil
		}
	}
	return "", errors.New("failed to generate unique kode_booking after max attempts")
}

// generateUniqueQrcode: Generate qrcode unik (8 digit angka, contoh: 98989898)
// Cek duplikat di database dengan max attempts
func (s *PurchaseService) generateUniqueQrcode(tx *gorm.DB) (string, error) {
	maxAttempts := 100
	for i := 0; i < maxAttempts; i++ {
		randomNum := rand.Intn(90000000) + 10000000 // 8 digit angka (10000000 - 99999999)
		qrcode := strconv.Itoa(randomNum)

		var count int64
		if err := tx.Model(&model.TicketQrcode{}).Where("qrcode = ?", qrcode).Count(&count).Error; err != nil {
			return "", fmt.Errorf("error checking qrcode uniqueness: %v", err)
		}
		if count == 0 {
			return qrcode, nil
		}
	}
	return "", errors.New("failed to generate unique qrcode after max attempts")
}

// ProcessPurchase: Fungsi utama untuk menangani pembelian
// - Insert ke ticket_qrcode (satu baris per quantity, jadi jika quantity 10, 10 baris)
// - Update ticket_kategori (kurangi sisa, tambah terjual)

// func (s *PurchaseService) ProcessPurchase(req *PurchaseRequest, userID string) error {
// 	// Validasi input dasar
// 	if userID == "" {
// 		return errors.New("userID tidak boleh kosong")
// 	}
// 	if len(req.Items) == 0 {
// 		return errors.New("tidak ada item untuk dibeli")
// 	}

// 	// Mulai transaksi
// 	tx := s.DB.Begin()
// 	defer func() {
// 		if r := recover(); r != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// Validasi voucher jika ada (gunakan field KodeVoucher dari model)
// 	var voucher *model.TicketVoucher
// 	if req.VoucherCode != nil && *req.VoucherCode != "" {
// 		if err := tx.Where("kode_voucher = ?", *req.VoucherCode).First(&voucher).Error; err != nil {
// 			tx.Rollback()
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return errors.New("voucher tidak valid")
// 			}
// 			return fmt.Errorf("error validating voucher: %v", err)
// 		}
// 		if voucher.Sisa <= 0 {
// 			tx.Rollback()
// 			return errors.New("voucher sudah habis")
// 		}
// 	}

// 	// Proses setiap item
// 	for _, item := range req.Items {
// 		if item.Quantity <= 0 {
// 			tx.Rollback()
// 			return errors.New("quantity harus lebih dari 0")
// 		}

// 		// Ambil data kategori TERBARU dalam transaksi (reload untuk menghindari stale data)
// 		var kategori model.TicketKategori
// 		if err := tx.Where("id = ?", item.TicketKategoriID).First(&kategori).Error; err != nil {
// 			tx.Rollback()
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return fmt.Errorf("kategori tidak ditemukan: %s", item.TicketKategoriID)
// 			}
// 			return fmt.Errorf("error fetching kategori: %v", err)
// 		}

// 		// Debug log (hapus setelah test)
// 		fmt.Printf("Kategori %s: Sisa=%d, Quantity=%d\n", kategori.Nama, kategori.Sisa, item.Quantity)

// 		// Validasi stok
// 		if item.Quantity > kategori.Sisa {
// 			tx.Rollback()
// 			return fmt.Errorf("stok tidak cukup untuk kategori %s", kategori.Nama)
// 		}

// 		// Update kategori: kurangi sisa, tambah terjual (total quantity)
// 		if err := tx.Model(&model.TicketKategori{}).Where("id = ?", item.TicketKategoriID).
// 			Updates(map[string]interface{}{
// 				"terjual": gorm.Expr("terjual + ?", item.Quantity),
// 				"sisa":    gorm.Expr("sisa - ?", item.Quantity),
// 			}).Error; err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("error updating kategori: %v", err)
// 		}

// 		// Insert ke ticket_qrcode: satu baris per quantity (jika quantity 10, 10 baris dengan quantity 1 masing-masing)
// 		for i := 0; i < item.Quantity; i++ {
// 			// Generate unik order_id, kode_booking, qrcode per baris
// 			orderID, err := s.generateUniqueOrderID(tx)
// 			if err != nil {
// 				tx.Rollback()
// 				return fmt.Errorf("gagal generate order_id: %v", err)
// 			}
// 			kodeBooking, err := s.generateUniqueKodeBooking(tx)
// 			if err != nil {
// 				tx.Rollback()
// 				return fmt.Errorf("gagal generate kode_booking: %v", err)
// 			}
// 			qrcodeStr, err := s.generateUniqueQrcode(tx)
// 			if err != nil {
// 				tx.Rollback()
// 				return fmt.Errorf("gagal generate qrcode: %v", err)
// 			}

// 			// Insert satu baris dengan quantity 1
// 			qrcode := model.TicketQrcode{
// 				ID:               uuid.New().String(),
// 				UserID:           userID, // UserID dari ticket register
// 				TicketKategoriID: item.TicketKategoriID,
// 				TicketJenisID:    req.TicketJenisID,
// 				Quantity:         1, // Quantity per baris adalah 1
// 				Price:            kategori.Harga,
// 				OrderID:          orderID,
// 				KodeBooking:      kodeBooking,
// 				Qrcode:           qrcodeStr,
// 				PaymentFeeID:     item.FeeID,
// 				PaymentStatus:    "pending",
// 				CreatedAt:        time.Now(),
// 				UpdatedAt:        time.Now(),
// 			}
// 			if voucher != nil {
// 				qrcode.TicketVoucherID = voucher.ID
// 			}
// 			if err := tx.Create(&qrcode).Error; err != nil {
// 				tx.Rollback()
// 				return fmt.Errorf("error inserting ticket_qrcode: %v", err)
// 			}
// 		}
// 	}

// 	// Update voucher jika digunakan (tambah terpakai sebanyak 1 per transaksi)
// 	if voucher != nil {
// 		if err := tx.Model(&model.TicketVoucher{}).Where("id = ?", voucher.ID).
// 			Update("terpakai", gorm.Expr("terpakai + ?", 1)).Error; err != nil {
// 			tx.Rollback()
// 			return fmt.Errorf("error updating voucher: %v", err)
// 		}
// 	}

// 	// Commit transaksi jika semua berhasil
// 	if err := tx.Commit().Error; err != nil {
// 		return fmt.Errorf("error committing transaction: %v", err)
// 	}

//		return nil
//	}

func (s *PurchaseService) ProcessPurchase(req *PurchaseRequest, userID string) error {
	if userID == "" {
		return errors.New("userID tidak boleh kosong")
	}
	if len(req.Items) == 0 {
		return errors.New("tidak ada item untuk dibeli")
	}

	tx := s.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Validasi voucher jika ada
	var voucher *model.TicketVoucher
	if req.VoucherCode != nil && *req.VoucherCode != "" {
		if err := tx.Where("kode_voucher = ?", *req.VoucherCode).First(&voucher).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("voucher tidak valid")
			}
			return fmt.Errorf("error validating voucher: %v", err)
		}
		if voucher.Sisa <= 0 {
			tx.Rollback()
			return errors.New("voucher sudah habis")
		}
	}

	// Folder penyimpanan QR - gunakan path absolut untuk menghindari masalah relatif
	wd, err := os.Getwd()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal mendapatkan working directory: %v", err)
	}
	qrFolder := filepath.Join(wd, "qrcode")
	if err := os.MkdirAll(qrFolder, os.ModePerm); err != nil {
		tx.Rollback()
		return fmt.Errorf("gagal membuat folder QR: %v", err)
	}
	log.Printf("QR folder: %s", qrFolder)

	for _, item := range req.Items {
		if item.Quantity <= 0 {
			tx.Rollback()
			return errors.New("quantity harus lebih dari 0")
		}

		var kategori model.TicketKategori
		if err := tx.Where("id = ?", item.TicketKategoriID).First(&kategori).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("kategori tidak ditemukan: %s", item.TicketKategoriID)
			}
			return fmt.Errorf("error fetching kategori: %v", err)
		}

		if item.Quantity > kategori.Sisa {
			tx.Rollback()
			return fmt.Errorf("stok tidak cukup untuk kategori %s", kategori.Nama)
		}

		if err := tx.Model(&model.TicketKategori{}).Where("id = ?", item.TicketKategoriID).
			Updates(map[string]interface{}{
				"terjual": gorm.Expr("terjual + ?", item.Quantity),
				"sisa":    gorm.Expr("sisa - ?", item.Quantity),
			}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating kategori: %v", err)
		}

		for i := 0; i < item.Quantity; i++ {
			// Generate unik
			orderID, err := s.generateUniqueOrderID(tx)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal generate order_id: %v", err)
			}
			kodeBooking, err := s.generateUniqueKodeBooking(tx)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal generate kode_booking: %v", err)
			}
			qrcodeStr, err := s.generateUniqueQrcode(tx)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("gagal generate qrcode: %v", err)
			}

			// Generate QR code file
			qrFileName := fmt.Sprintf("%s.png", qrcodeStr)
			qrFilePath := filepath.Join(qrFolder, qrFileName)
			log.Printf("Generating QR for %s at %s", qrcodeStr, qrFilePath)
			err = qrcode.WriteFile(qrcodeStr, qrcode.Medium, 256, qrFilePath)
			if err != nil {
				log.Printf("Error generating QR: %v", err)
				tx.Rollback()
				return fmt.Errorf("gagal generate QR file: %v", err)
			}
			log.Printf("QR generated successfully at %s", qrFilePath)

			// Insert ke database
			qrcodeModel := model.TicketQrcode{
				ID:               uuid.New().String(),
				UserID:           userID,
				TicketKategoriID: item.TicketKategoriID,
				TicketJenisID:    req.TicketJenisID,
				TicketEventID:    "f9fa68e6-471b-4974-afe9-fa0873c22007",
				Quantity:         1,
				Price:            kategori.Harga,
				OrderID:          orderID,
				KodeBooking:      kodeBooking,
				Qrcode:           qrcodeStr,
				DirectoryQrcode:  qrFilePath, // path file QR (sekarang absolut)
				PaymentFeeID:     item.FeeID,
				PaymentStatus:    "pending",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			if voucher != nil {
				qrcodeModel.TicketVoucherID = voucher.ID
			}

			if err := tx.Create(&qrcodeModel).Error; err != nil {
				tx.Rollback()
				return fmt.Errorf("error inserting ticket_qrcode: %v", err)
			}
		}
	}

	// Update voucher jika digunakan
	if voucher != nil {
		if err := tx.Model(&model.TicketVoucher{}).Where("id = ?", voucher.ID).
			Update("terpakai", gorm.Expr("terpakai + ?", 1)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("error updating voucher: %v", err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GetPurchaseByToken: Fungsi untuk mendapatkan data purchase berdasarkan token (order_id)
// Mengembalikan list TicketQrcode yang match dengan order_id
func (s *PurchaseService) GetPurchaseByToken(token string) ([]model.TicketQrcode, error) {
	if token == "" {
		return nil, errors.New("token tidak boleh kosong")
	}

	var purchases []model.TicketQrcode
	if err := s.DB.Where("order_id = ?", token).Find(&purchases).Error; err != nil {
		return nil, fmt.Errorf("error fetching purchases: %v", err)
	}

	if len(purchases) == 0 {
		return nil, errors.New("purchase tidak ditemukan")
	}

	return purchases, nil
}
