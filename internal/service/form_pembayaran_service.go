package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FormPembayaranService struct {
	DB *gorm.DB
}

func NewFormPembayaranService(db *gorm.DB) *FormPembayaranService {
	return &FormPembayaranService{DB: db}
}

func isValidDate(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// CREATE
func (s *FormPembayaranService) CreateFormPembayaran(f *model.FormPembayaran) error {
	if f.ID == "" {
		f.ID = uuid.New().String()
	}

	if f.NamaForm == "" {
		return fmt.Errorf("nama_form is required")
	}

	if f.Jumlah <= 0 {
		return fmt.Errorf("jumlah must be > 0")
	}

	if f.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if !isValidDate(f.TanggalMulai) {
		return fmt.Errorf("tanggal_mulai is required and must be in YYYY-MM-DD format")
	}

	if f.TanggalSelesai != nil && !isValidDate(*f.TanggalSelesai) {
		return fmt.Errorf("tanggal_selesai must be in YYYY-MM-DD format if provided")
	}

	// Validasi FK kategori_id
	var kategori model.KategoriPembayaran
	if err := s.DB.First(&kategori, "id = ?", f.KategoriID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("kategori_id not found")
		}
		return err
	}

	// Validasi FK denda_id
	var denda model.AturanDenda
	if err := s.DB.First(&denda, "id = ?", f.DendaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("denda_id not found")
		}
		return err
	}

	// Validasi cicilan
	if f.BolehCicilan && (f.JumlahCicilan == nil || *f.JumlahCicilan <= 0) {
		return fmt.Errorf("jumlah_cicilan is required and must be > 0 if boleh_cicilan is true")
	}

	// Set default status jika kosong
	if f.Status == "" {
		f.Status = "active"
	}

	return s.DB.Create(f).Error
}

// READ ALL
func (s *FormPembayaranService) GetAllFormPembayaran() ([]model.FormPembayaran, error) {
	var forms []model.FormPembayaran
	result := s.DB.Preload("Kategori").Preload("Denda").Find(&forms)
	return forms, result.Error
}

// READ BY ID
func (s *FormPembayaranService) GetFormPembayaranByID(id string) (*model.FormPembayaran, error) {
	var f model.FormPembayaran
	result := s.DB.Preload("Kategori").Preload("Denda").First(&f, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &f, result.Error
}

// UPDATE
func (s *FormPembayaranService) UpdateFormPembayaran(f *model.FormPembayaran) error {
	if f.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var old model.FormPembayaran
	if err := s.DB.First(&old, "id = ?", f.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("form pembayaran not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if f.NamaForm != "" && f.NamaForm != old.NamaForm {
		updateData["nama_form"] = f.NamaForm
	}

	if f.Jumlah > 0 && f.Jumlah != old.Jumlah {
		updateData["jumlah"] = f.Jumlah
	}

	if f.UserID != "" && f.UserID != old.UserID {
		updateData["user_id"] = f.UserID
	}

	if f.KategoriID != "" && f.KategoriID != old.KategoriID {
		var kategori model.KategoriPembayaran
		if err := s.DB.First(&kategori, "id = ?", f.KategoriID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("kategori_id not found")
			}
			return err
		}
		updateData["kategori_id"] = f.KategoriID
	}

	if f.DendaID != "" && f.DendaID != old.DendaID {
		var denda model.AturanDenda
		if err := s.DB.First(&denda, "id = ?", f.DendaID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("denda_id not found")
			}
			return err
		}
		updateData["denda_id"] = f.DendaID
	}

	if f.BolehCicilan != old.BolehCicilan {
		updateData["boleh_cicilan"] = f.BolehCicilan
	}

	if f.JumlahCicilan != nil {
		updateData["jumlah_cicilan"] = f.JumlahCicilan
	}

	if f.Status != "" && f.Status != old.Status {
		updateData["status"] = f.Status
	}

	if f.TanggalMulai != "" && f.TanggalMulai != old.TanggalMulai {
		if !isValidDate(f.TanggalMulai) {
			return fmt.Errorf("tanggal_mulai must be in YYYY-MM-DD format")
		}
		updateData["tanggal_mulai"] = f.TanggalMulai
	}

	if f.TanggalSelesai != nil && (old.TanggalSelesai == nil || *f.TanggalSelesai != *old.TanggalSelesai) {
		if !isValidDate(*f.TanggalSelesai) {
			return fmt.Errorf("tanggal_selesai must be in YYYY-MM-DD format")
		}
		updateData["tanggal_selesai"] = f.TanggalSelesai
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.FormPembayaran{}).Where("id = ?", f.ID).Updates(updateData).Error
}

// DELETE
func (s *FormPembayaranService) DeleteFormPembayaran(id string) error {
	return s.DB.Delete(&model.FormPembayaran{}, "id = ?", id).Error
}
