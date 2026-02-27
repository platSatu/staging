// service/voucher_service.go
package service

import (
	"backend_go/internal/model"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VoucherService struct {
	DB *gorm.DB
}

func NewVoucherService(db *gorm.DB) *VoucherService {
	return &VoucherService{DB: db}
}

// GetLastSaldoByUserID gets the last saldo for a user
func (s *VoucherService) GetLastSaldoByUserID(userID string) (float64, error) {
	var lastDeposit model.Deposit
	result := s.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&lastDeposit)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if result.Error != nil {
		return 0, result.Error
	}
	return lastDeposit.Saldo, nil
}

// BuyPackage handles purchasing a package using deposit
func (s *VoucherService) BuyPackage(userID string, packagesID string) (*model.Voucher, error) {
	// Validasi user_id
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Ambil package
	var pkg model.Packages
	if err := s.DB.First(&pkg, "id = ?", packagesID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("package not found")
		}
		return nil, fmt.Errorf("error fetching package: %v", err)
	}

	// Cek package aktif
	if pkg.Status != "active" {
		return nil, fmt.Errorf("package is not active")
	}

	// Ambil saldo terakhir user
	lastSaldo, err := s.GetLastSaldoByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching saldo: %v", err)
	}

	// Validasi saldo cukup
	if lastSaldo < pkg.Price {
		return nil, fmt.Errorf("insufficient balance. your balance: %.2f, package price: %.2f", lastSaldo, pkg.Price)
	}

	// Hitung saldo baru setelah debit
	newSaldo := lastSaldo - pkg.Price

	// Generate kode voucher
	kodeVoucher := s.generateUniqueKodeVoucher()

	// Hitung valid_from dan valid_until
	now := time.Now()
	validFrom := now
	validUntil := now.AddDate(0, 0, pkg.DurationDays)

	// Insert ke vouchers
	voucher := &model.Voucher{
		ID:          uuid.New().String(),
		UserID:      userID,
		PackagesID:  packagesID,
		KodeVoucher: kodeVoucher,
		Status:      "active",
		ValidFrom:   &validFrom,
		ValidUntil:  &validUntil,
	}

	if err := s.DB.Create(voucher).Error; err != nil {
		return nil, fmt.Errorf("failed to create voucher: %v", err)
	}

	// Insert ke deposits (debit) untuk mencatat transaksi
	noInvoice, err := s.generateUniqueNoInvoice()
	if err != nil {
		return nil, fmt.Errorf("failed to generate no_invoice: %v", err)
	}
	orderID, err := s.generateUniqueOrderID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate order_id: %v", err)
	}

	deposit := &model.Deposit{
		UserID:           userID,
		Debit:            pkg.Price,
		Kredit:           0,
		Saldo:            newSaldo,
		TransaksiTanggal: now,
		TransaksiStatus:  "completed",
		TransaksiMethod:  "buy_package",
		Keterangan:       fmt.Sprintf("Purchase package: %s", pkg.Name),
		NoInvoice:        noInvoice,
		OrderID:          orderID,
	}

	if err := s.DB.Create(deposit).Error; err != nil {
		return nil, fmt.Errorf("failed to create deposit transaction: %v", err)
	}

	return voucher, nil
}

// generateUniqueNoInvoice generates a unique no_invoice
func (s *VoucherService) generateUniqueNoInvoice() (string, error) {
	for {
		noInvoice := generateNoInvoice()
		var count int64
		s.DB.Model(&model.Deposit{}).Where("no_invoice = ?", noInvoice).Count(&count)
		if count == 0 {
			return noInvoice, nil
		}
	}
}

// generateUniqueOrderID generates a unique order_id
func (s *VoucherService) generateUniqueOrderID() (string, error) {
	for {
		orderID := generateOrderID()
		var count int64
		s.DB.Model(&model.Deposit{}).Where("order_id = ?", orderID).Count(&count)
		if count == 0 {
			return orderID, nil
		}
	}
}

// generateNoInvoice helper: PS + 6 angka
func generateNoInvoice() string {
	const digits = "0123456789"
	noInvoice := "PS"
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		noInvoice += string(digits[num.Int64()])
	}
	return noInvoice
}

// generateOrderID helper: ORD- + 6 angka
func generateOrderID() string {
	const digits = "0123456789"
	orderID := "ORD-"
	for i := 0; i < 6; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		orderID += string(digits[num.Int64()])
	}
	return orderID
}

// CREATE (Update untuk generate kode_voucher, status, valid_from, valid_until berdasarkan packages_id)
func (s *VoucherService) CreateVoucher(userID, packagesID string) (*model.Voucher, error) {
	// Validasi packages_id ada di tabel packages
	var pkg model.Packages
	if err := s.DB.First(&pkg, "id = ?", packagesID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("package not found")
		}
		return nil, err
	}

	// Validasi user_id ada di tabel users
	var user model.User
	if err := s.DB.First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Generate kode_voucher unik (KSD + 8 angka)
	kodeVoucher := s.generateUniqueKodeVoucher()

	// Set valid_from = sekarang, valid_until = valid_from + duration_days
	now := time.Now()
	validFrom := &now
	validUntil := &time.Time{}
	*validUntil = now.AddDate(0, 0, pkg.DurationDays)

	voucher := &model.Voucher{
		ID:          uuid.New().String(),
		UserID:      userID,
		PackagesID:  packagesID,
		KodeVoucher: kodeVoucher,
		Status:      "active",
		ValidFrom:   validFrom,
		ValidUntil:  validUntil,
	}

	return voucher, s.DB.Create(voucher).Error
}

// generateUniqueKodeVoucher membuat kode voucher unik: KSD + 8 angka
func (s *VoucherService) generateUniqueKodeVoucher() string {
	for {
		kode := generateKodeVoucher()
		// Cek apakah kode sudah ada di database
		var count int64
		s.DB.Model(&model.Voucher{}).Where("kode_voucher = ?", kode).Count(&count)
		if count == 0 {
			return kode
		}
		// Jika duplicate, loop lagi untuk generate ulang
	}
}

// generateKodeVoucher helper: KSD + 8 angka
func generateKodeVoucher() string {
	const digits = "0123456789"
	kode := "KSD"

	// 8 angka
	for i := 0; i < 8; i++ {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		kode += string(digits[num.Int64()])
	}

	return kode
}

// READ ALL
func (s *VoucherService) GetAllVouchers() ([]model.Voucher, error) {
	var vouchers []model.Voucher
	result := s.DB.Find(&vouchers)
	return vouchers, result.Error
}

// READ BY ID
func (s *VoucherService) GetVoucherByID(id string) (*model.Voucher, error) {
	var voucher model.Voucher
	result := s.DB.First(&voucher, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &voucher, result.Error
}

// READ BY USER ID
func (s *VoucherService) GetVouchersByUserID(userID string) ([]model.Voucher, error) {
	var vouchers []model.Voucher
	result := s.DB.Where("user_id = ?", userID).Find(&vouchers)
	return vouchers, result.Error
}

// UPDATE
func (s *VoucherService) UpdateVoucher(voucher *model.Voucher) error {
	if voucher.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldVoucher model.Voucher
	if err := s.DB.First(&oldVoucher, "id = ?", voucher.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("voucher not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if voucher.KodeVoucher != "" && voucher.KodeVoucher != oldVoucher.KodeVoucher {
		// Cek duplicate kode_voucher baru
		var count int64
		s.DB.Model(&model.Voucher{}).Where("kode_voucher = ? AND id != ?", voucher.KodeVoucher, voucher.ID).Count(&count)
		if count > 0 {
			return fmt.Errorf("kode_voucher already exists")
		}
		updateData["kode_voucher"] = voucher.KodeVoucher
	}

	if voucher.Status != "" && voucher.Status != oldVoucher.Status {
		updateData["status"] = voucher.Status
	}

	if voucher.ValidFrom != nil && (oldVoucher.ValidFrom == nil || !voucher.ValidFrom.Equal(*oldVoucher.ValidFrom)) {
		updateData["valid_from"] = voucher.ValidFrom
	}

	if voucher.ValidUntil != nil && (oldVoucher.ValidUntil == nil || !voucher.ValidUntil.Equal(*oldVoucher.ValidUntil)) {
		updateData["valid_until"] = voucher.ValidUntil
	}

	if voucher.UserID != "" && voucher.UserID != oldVoucher.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", voucher.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = voucher.UserID
	}

	if voucher.PackagesID != "" && voucher.PackagesID != oldVoucher.PackagesID {
		// Validasi packages_id baru ada
		var pkg model.Packages
		if err := s.DB.First(&pkg, "id = ?", voucher.PackagesID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("package not found")
			}
			return err
		}
		updateData["packages_id"] = voucher.PackagesID
	}

	// Validasi valid_from dan valid_until jika diupdate
	if updateData["valid_from"] != nil || updateData["valid_until"] != nil {
		validFrom := voucher.ValidFrom
		if validFrom == nil && oldVoucher.ValidFrom != nil {
			validFrom = oldVoucher.ValidFrom
		}
		validUntil := voucher.ValidUntil
		if validUntil == nil && oldVoucher.ValidUntil != nil {
			validUntil = oldVoucher.ValidUntil
		}
		if validFrom != nil && validUntil != nil && validFrom.After(*validUntil) {
			return fmt.Errorf("valid_from cannot be after valid_until")
		}
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Voucher{}).Where("id = ?", voucher.ID).Updates(updateData).Error
}

// DELETE
func (s *VoucherService) DeleteVoucher(id string) error {
	return s.DB.Delete(&model.Voucher{}, "id = ?", id).Error
}

// RedeemVoucher handles voucher redemption by kode_voucher
func (s *VoucherService) RedeemVoucher(userID, kodeVoucher string) (*model.Voucher, error) {
	// Validasi input
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if kodeVoucher == "" {
		return nil, fmt.Errorf("kode_voucher is required")
	}

	// Cari voucher berdasarkan kode_voucher
	var voucher model.Voucher
	result := s.DB.Where("kode_voucher = ?", kodeVoucher).First(&voucher)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("voucher not found")
		}
		return nil, fmt.Errorf("error fetching voucher: %v", result.Error)
	}

	// Validasi voucher milik user yang login
	if voucher.UserID != userID {
		return nil, fmt.Errorf("voucher does not belong to this user")
	}

	// Validasi status voucher masih active
	if voucher.Status != "active" {
		return nil, fmt.Errorf("voucher is already used or expired")
	}

	// Validasi belum expired
	if voucher.ValidUntil != nil {
		now := time.Now()
		if now.After(*voucher.ValidUntil) {
			return nil, fmt.Errorf("voucher has expired")
		}
	}

	// Update status menjadi used
	voucher.Status = "used"
	if err := s.DB.Save(&voucher).Error; err != nil {
		return nil, fmt.Errorf("failed to update voucher status: %v", err)
	}

	return &voucher, nil
}

// AutoExpireVouchers - Update voucher yang sudah expired berdasarkan valid_until
// Dipanggil oleh cron job
func (s *VoucherService) AutoExpireVouchers() (int64, error) {
	now := time.Now()

	// Set waktu hari ini jam 00:00:00
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Update voucher yang statusnya "used" dan valid_until <= hari ini
	result := s.DB.Model(&model.Voucher{}).
		Where("status = ?", "used").
		Where("valid_until <= ?", today).
		Update("status", "expired")

	if result.Error != nil {
		return 0, result.Error
	}

	affected := result.RowsAffected
	if affected > 0 {
		fmt.Printf("[VOUCHER] Auto-expired %d vouchers\n", affected)
	}

	return affected, nil
}
