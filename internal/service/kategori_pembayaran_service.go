package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KategoriPembayaranService struct {
	DB *gorm.DB
}

func NewKategoriPembayaranService(db *gorm.DB) *KategoriPembayaranService {
	return &KategoriPembayaranService{DB: db}
}

// Validasi status (opsional, bisa ditambah enum jika perlu)
func (s *KategoriPembayaranService) validateStatus(status string) error {
	if status != "active" && status != "inactive" {
		return fmt.Errorf("status must be 'active' or 'inactive'")
	}
	return nil
}

// CREATE
func (s *KategoriPembayaranService) CreateKategoriPembayaran(kategori *model.KategoriPembayaran) error {
	if kategori.ID == "" {
		kategori.ID = uuid.New().String()
	}

	if kategori.NamaKategori == "" {
		return fmt.Errorf("nama_kategori is required")
	}

	if kategori.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if err := s.validateStatus(kategori.Status); err != nil {
		return err
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", kategori.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	// Set default status jika kosong
	if kategori.Status == "" {
		kategori.Status = "active"
	}

	return s.DB.Create(kategori).Error
}

// READ ALL
func (s *KategoriPembayaranService) GetAllKategoriPembayaran() ([]model.KategoriPembayaran, error) {
	var kategoris []model.KategoriPembayaran
	result := s.DB.Find(&kategoris)
	return kategoris, result.Error
}

// READ BY ID
func (s *KategoriPembayaranService) GetKategoriPembayaranByID(id string) (*model.KategoriPembayaran, error) {
	var kategori model.KategoriPembayaran
	result := s.DB.First(&kategori, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &kategori, result.Error
}

// UPDATE
func (s *KategoriPembayaranService) UpdateKategoriPembayaran(kategori *model.KategoriPembayaran) error {
	if kategori.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldKategori model.KategoriPembayaran
	if err := s.DB.First(&oldKategori, "id = ?", kategori.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("kategori pembayaran not found")
		}
		return err
	}

	if err := s.validateStatus(kategori.Status); err != nil {
		return err
	}

	updateData := map[string]interface{}{}

	if kategori.UserID != "" && kategori.UserID != oldKategori.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", kategori.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = kategori.UserID
	}

	if kategori.NamaKategori != "" && kategori.NamaKategori != oldKategori.NamaKategori {
		updateData["nama_kategori"] = kategori.NamaKategori
	}

	if kategori.Deskripsi != nil {
		updateData["deskripsi"] = kategori.Deskripsi
	}

	if kategori.Status != "" && kategori.Status != oldKategori.Status {
		updateData["status"] = kategori.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.KategoriPembayaran{}).Where("id = ?", kategori.ID).Updates(updateData).Error
}

// DELETE
func (s *KategoriPembayaranService) DeleteKategoriPembayaran(id string) error {
	return s.DB.Delete(&model.KategoriPembayaran{}, "id = ?", id).Error
}

func (s *KategoriPembayaranService) GetAllKategoriPembayaranByUserID(userID string) ([]model.KategoriPembayaran, error) {
	var list []model.KategoriPembayaran
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
