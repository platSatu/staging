package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransaksiService struct {
	DB *gorm.DB
}

func NewTransaksiService(db *gorm.DB) *TransaksiService {
	return &TransaksiService{DB: db}
}

// Validasi tipe_transaksi
func (s *TransaksiService) validateTipeTransaksi(tipe string) error {
	validTypes := []string{"pembayaran", "denda", "penyesuaian"}
	for _, v := range validTypes {
		if tipe == v {
			return nil
		}
	}
	return fmt.Errorf("tipe_transaksi must be one of: pembayaran, denda, penyesuaian")
}

// Validasi status_gateway
func (s *TransaksiService) validateStatusGateway(status *string) error {
	if status == nil {
		return nil
	}
	validStatuses := []string{"pending", "success", "failed"}
	for _, v := range validStatuses {
		if *status == v {
			return nil
		}
	}
	return fmt.Errorf("status_gateway must be one of: pending, success, failed")
}

// Hitung denda jika terlambat
func (s *TransaksiService) calculateDenda(kewajibanID string, tanggalBayar time.Time) (float64, error) {
	var kewajiban model.KewajibanUser
	if err := s.DB.Preload("Denda").First(&kewajiban, "id = ?", kewajibanID).Error; err != nil {
		return 0, err
	}

	if kewajiban.DendaID == "" || kewajiban.Denda.ID == "" {
		return 0, nil
	}

	denda := &kewajiban.Denda

	// Convert TanggalJatuhTempo dari *string ke time.Time
	if kewajiban.TanggalJatuhTempo == nil {
		return 0, nil
	}
	tglJatuhTempo, err := time.Parse("2006-01-02", *kewajiban.TanggalJatuhTempo)
	if err != nil {
		return 0, err
	}

	if !tanggalBayar.After(tglJatuhTempo) {
		return 0, nil
	}

	hariTerlambat := int(tanggalBayar.Sub(tglJatuhTempo).Hours() / 24)
	var totalDenda float64

	switch denda.TipeDenda {
	case "harian_persentase":
		if denda.Persentase != nil {
			totalDenda = kewajiban.JumlahSisa * (*denda.Persentase / 100) * float64(hariTerlambat)
			if denda.CapMaksimal != nil && totalDenda > *denda.CapMaksimal {
				totalDenda = *denda.CapMaksimal
			}
		}
	case "harian_flat":
		if denda.JumlahFlat != nil {
			totalDenda = *denda.JumlahFlat * float64(hariTerlambat)
			if denda.CapMaksimal != nil && totalDenda > *denda.CapMaksimal {
				totalDenda = *denda.CapMaksimal
			}
		}
	case "flat_cap":
		if denda.JumlahFlat != nil {
			totalDenda = *denda.JumlahFlat
			if denda.CapMaksimal != nil && totalDenda > *denda.CapMaksimal {
				totalDenda = *denda.CapMaksimal
			}
		}
	}

	return totalDenda, nil
}

// CREATE
func (s *TransaksiService) CreateTransaksi(transaksi *model.Transaksi) error {
	if transaksi.ID == "" {
		transaksi.ID = uuid.New().String()
	}

	if transaksi.Jumlah <= 0 {
		return fmt.Errorf("jumlah must be > 0")
	}

	if err := s.validateTipeTransaksi(transaksi.TipeTransaksi); err != nil {
		return err
	}

	if err := s.validateStatusGateway(transaksi.StatusGateway); err != nil {
		return err
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", transaksi.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	// Validasi FK kewajiban_id
	var kewajiban model.KewajibanUser
	if err := s.DB.First(&kewajiban, "id = ?", transaksi.KewajibanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("kewajiban_id not found")
		}
		return err
	}

	// Validasi parent_id jika ada
	if transaksi.ParentID != "" {
		var parent model.Transaksi
		if err := s.DB.First(&parent, "id = ?", transaksi.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent_id not found")
			}
			return err
		}
	}

	// Set tanggal sekarang jika kosong
	if transaksi.Tanggal.IsZero() {
		transaksi.Tanggal = time.Now()
	}

	// Hitung denda jika pembayaran terlambat
	if transaksi.TipeTransaksi == "pembayaran" {
		denda, err := s.calculateDenda(transaksi.KewajibanID, transaksi.Tanggal)
		if err != nil {
			return err
		}
		if denda > 0 {
			// Buat transaksi denda terpisah
			dendaTransaksi := model.Transaksi{
				ID:            uuid.New().String(),
				UserID:        transaksi.UserID,
				KewajibanID:   transaksi.KewajibanID,
				ParentID:      transaksi.ID,
				TipeTransaksi: "denda",
				Jumlah:        denda,
				Tanggal:       transaksi.Tanggal,
			}
			if err := s.DB.Create(&dendaTransaksi).Error; err != nil {
				return err
			}
		}
	}

	return s.DB.Create(transaksi).Error
}

// READ ALL
func (s *TransaksiService) GetAllTransaksi() ([]model.Transaksi, error) {
	var transaksis []model.Transaksi
	result := s.DB.Preload("User").
		Preload("Kewajiban").
		Find(&transaksis)
	return transaksis, result.Error
}

// READ BY ID
func (s *TransaksiService) GetTransaksiByID(id string) (*model.Transaksi, error) {
	var transaksi model.Transaksi
	result := s.DB.Preload("User").
		Preload("Kewajiban").
		First(&transaksi, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &transaksi, result.Error
}

// UPDATE
func (s *TransaksiService) UpdateTransaksi(transaksi *model.Transaksi) error {
	if transaksi.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldTransaksi model.Transaksi
	if err := s.DB.First(&oldTransaksi, "id = ?", transaksi.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("transaksi not found")
		}
		return err
	}

	if err := s.validateTipeTransaksi(transaksi.TipeTransaksi); err != nil {
		return err
	}

	if err := s.validateStatusGateway(transaksi.StatusGateway); err != nil {
		return err
	}

	updateData := map[string]interface{}{}

	if transaksi.UserID != "" && transaksi.UserID != oldTransaksi.UserID {
		var user model.User
		if err := s.DB.First(&user, "id = ?", transaksi.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = transaksi.UserID
	}

	if transaksi.KewajibanID != "" && transaksi.KewajibanID != oldTransaksi.KewajibanID {
		var kewajiban model.KewajibanUser
		if err := s.DB.First(&kewajiban, "id = ?", transaksi.KewajibanID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("kewajiban_id not found")
			}
			return err
		}
		updateData["kewajiban_id"] = transaksi.KewajibanID
	}

	if transaksi.ParentID != "" && transaksi.ParentID != oldTransaksi.ParentID {
		var parent model.Transaksi
		if err := s.DB.First(&parent, "id = ?", transaksi.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent_id not found")
			}
			return err
		}
		updateData["parent_id"] = transaksi.ParentID
	}

	if transaksi.TipeTransaksi != "" && transaksi.TipeTransaksi != oldTransaksi.TipeTransaksi {
		updateData["tipe_transaksi"] = transaksi.TipeTransaksi
	}

	if transaksi.Jumlah > 0 && transaksi.Jumlah != oldTransaksi.Jumlah {
		updateData["jumlah"] = transaksi.Jumlah
	}

	if !transaksi.Tanggal.IsZero() && !transaksi.Tanggal.Equal(oldTransaksi.Tanggal) {
		updateData["tanggal"] = transaksi.Tanggal
	}

	if transaksi.MetodePembayaran != nil {
		updateData["metode_pembayaran"] = transaksi.MetodePembayaran
	}

	if transaksi.ReferenceID != nil {
		updateData["reference_id"] = transaksi.ReferenceID
	}

	if transaksi.StatusGateway != nil {
		updateData["status_gateway"] = transaksi.StatusGateway
	}

	if transaksi.Catatan != nil {
		updateData["catatan"] = transaksi.Catatan
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.Transaksi{}).Where("id = ?", transaksi.ID).Updates(updateData).Error
}

// DELETE
func (s *TransaksiService) DeleteTransaksi(id string) error {
	return s.DB.Delete(&model.Transaksi{}, "id = ?", id).Error
}

// GetTransaksiByUserID returns transaksi based on parent_id
func (s *TransaksiService) GetTransaksiByParentID(parentID string) ([]model.Transaksi, error) {
	var transaksis []model.Transaksi

	// Tampilkan semua transaksi dimana parent_id = user yang login
	result := s.DB.Preload("User").
		Preload("Kewajiban").
		Where("parent_id = ?", parentID).
		Order("created_at DESC").
		Find(&transaksis)

	return transaksis, result.Error
}
