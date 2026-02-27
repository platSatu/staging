package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type DepositService struct {
	DB *gorm.DB
}

func NewDepositService(db *gorm.DB) *DepositService {
	return &DepositService{DB: db}
}

// generateUniqueNumber generates a unique number with prefix, ensuring no duplicates in DB
func (s *DepositService) generateUniqueNumber(prefix string, length int, field string) (string, error) {
	for {
		// Generate random number
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(999999) + 1 // 1 to 999999 for 6 digits, adjust as needed
		var generated string
		if length == 6 {
			generated = fmt.Sprintf("%s%06d", prefix, number)
		} else {
			generated = fmt.Sprintf("%s%d", prefix, number)
		}

		// Check for uniqueness in DB
		var count int64
		query := s.DB.Model(&model.Deposit{})
		if field == "no_invoice" {
			query = query.Where("no_invoice = ?", generated)
		} else if field == "order_id" {
			query = query.Where("order_id = ?", generated)
		}
		query.Count(&count)
		if count == 0 {
			return generated, nil
		}
		// If duplicate, loop again
	}
}

// GenerateNoInvoice generates a unique no_invoice (e.g., PS123456)
func (s *DepositService) GenerateNoInvoice() (string, error) {
	return s.generateUniqueNumber("PS", 0, "no_invoice") // Length 0 means variable, but example is 6 digits
}

// GenerateOrderID generates a unique order_id (e.g., ORD-123456)
func (s *DepositService) GenerateOrderID() (string, error) {
	return s.generateUniqueNumber("ORD-", 6, "order_id")
}

// GetLastSaldoByUserID gets the last saldo for a user
func (s *DepositService) GetLastSaldoByUserID(userID string) (float64, error) {
	var lastDeposit model.Deposit
	result := s.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&lastDeposit)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return 0, nil // Jika belum ada deposit, saldo 0
	}
	if result.Error != nil {
		return 0, result.Error
	}
	return lastDeposit.Saldo, nil
}

// CREATE
func (s *DepositService) CreateDeposit(deposit *model.Deposit) error {
	// Validasi kredit minimal 10.000 jika kredit > 0
	if deposit.Kredit > 0 && deposit.Kredit < 10000 {
		return fmt.Errorf("kredit must be at least 10,000 if positive")
	}
	if deposit.Debit < 0 {
		return fmt.Errorf("debit cannot be negative")
	}

	// Ambil saldo terakhir user
	lastSaldo, err := s.GetLastSaldoByUserID(deposit.UserID)
	if err != nil {
		return fmt.Errorf("failed to get last saldo: %v", err)
	}

	// Hitung saldo baru: saldo terakhir + kredit - debit
	newSaldo := lastSaldo + deposit.Kredit - deposit.Debit
	deposit.Saldo = newSaldo

	// Set transaksi_tanggal jika belum ada
	if deposit.TransaksiTanggal.IsZero() {
		deposit.TransaksiTanggal = time.Now()
	}

	// Generate unique no_invoice and order_id
	noInvoice, err := s.GenerateNoInvoice()
	if err != nil {
		return fmt.Errorf("failed to generate no_invoice: %v", err)
	}
	orderID, err := s.GenerateOrderID()
	if err != nil {
		return fmt.Errorf("failed to generate order_id: %v", err)
	}
	deposit.NoInvoice = noInvoice
	deposit.OrderID = orderID

	// Simpan ke database
	if err := s.DB.Create(deposit).Error; err != nil {
		return fmt.Errorf("failed to create deposit: %v", err)
	}

	return nil
}

// READ ALL
func (s *DepositService) GetAllDeposits() ([]model.Deposit, error) {
	var deposits []model.Deposit
	result := s.DB.Find(&deposits)
	return deposits, result.Error
}

// READ BY ID
func (s *DepositService) GetDepositByID(id string) (*model.Deposit, error) {
	var deposit model.Deposit
	result := s.DB.First(&deposit, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &deposit, result.Error
}

// READ BY USER ID
func (s *DepositService) GetDepositsByUserID(userID string) ([]model.Deposit, error) {
	var deposits []model.Deposit
	result := s.DB.Where("user_id = ?", userID).Find(&deposits)
	return deposits, result.Error
}

// UPDATE
func (s *DepositService) UpdateDeposit(deposit *model.Deposit) error {
	if deposit.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldDeposit model.Deposit
	if err := s.DB.First(&oldDeposit, "id = ?", deposit.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("deposit not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if deposit.Debit != oldDeposit.Debit {
		updateData["debit"] = deposit.Debit
	}

	if deposit.Kredit != oldDeposit.Kredit {
		updateData["kredit"] = deposit.Kredit
	}

	if deposit.Saldo != oldDeposit.Saldo {
		updateData["saldo"] = deposit.Saldo
	}

	if !deposit.TransaksiTanggal.IsZero() && deposit.TransaksiTanggal != oldDeposit.TransaksiTanggal {
		updateData["transaksi_tanggal"] = deposit.TransaksiTanggal
	}

	if deposit.TransaksiStatus != "" && deposit.TransaksiStatus != oldDeposit.TransaksiStatus {
		updateData["transaksi_status"] = deposit.TransaksiStatus
	}

	if deposit.TransaksiMethod != "" && deposit.TransaksiMethod != oldDeposit.TransaksiMethod {
		updateData["transaksi_method"] = deposit.TransaksiMethod
	}

	if deposit.Keterangan != "" && deposit.Keterangan != oldDeposit.Keterangan {
		updateData["keterangan"] = deposit.Keterangan
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.Deposit{}).Where("id = ?", deposit.ID).Updates(updateData).Error
}

// DELETE
func (s *DepositService) DeleteDeposit(id string) error {
	return s.DB.Delete(&model.Deposit{}, "id = ?", id).Error
}
