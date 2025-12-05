package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AturanDendaService struct {
	DB *gorm.DB
}

func NewAturanDendaService(db *gorm.DB) *AturanDendaService {
	return &AturanDendaService{DB: db}
}

// Validasi tipe_denda
func (s *AturanDendaService) validateTipeDenda(tipe string) error {
	validTypes := []string{"harian_persentase", "harian_flat", "flat_cap"}
	for _, v := range validTypes {
		if tipe == v {
			return nil
		}
	}
	return fmt.Errorf("tipe_denda must be one of: harian_persentase, harian_flat, flat_cap")
}

// Validasi status
func (s *AturanDendaService) validateStatus(status string) error {
	if status != "active" && status != "inactive" {
		return fmt.Errorf("status must be 'active' or 'inactive'")
	}
	return nil
}

// CREATE
func (s *AturanDendaService) CreateAturanDenda(aturan *model.AturanDenda) error {
	if aturan.ID == "" {
		aturan.ID = uuid.New().String()
	}

	if aturan.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	if err := s.validateTipeDenda(aturan.TipeDenda); err != nil {
		return err
	}

	if err := s.validateStatus(aturan.Status); err != nil {
		return err
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", aturan.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	// Validasi logika berdasarkan tipe
	switch aturan.TipeDenda {
	case "harian_persentase":
		if aturan.Persentase == nil || *aturan.Persentase <= 0 {
			return fmt.Errorf("persentase is required and must be > 0 for tipe_denda 'harian_persentase'")
		}
	case "harian_flat":
		if aturan.JumlahFlat == nil || *aturan.JumlahFlat <= 0 {
			return fmt.Errorf("jumlah_flat is required and must be > 0 for tipe_denda 'harian_flat'")
		}
	case "flat_cap":
		if aturan.JumlahFlat == nil || *aturan.JumlahFlat <= 0 {
			return fmt.Errorf("jumlah_flat is required and must be > 0 for tipe_denda 'flat_cap'")
		}
		if aturan.CapMaksimal == nil || *aturan.CapMaksimal <= 0 {
			return fmt.Errorf("cap_maksimal is required and must be > 0 for tipe_denda 'flat_cap'")
		}
	}

	// Set default status jika kosong
	if aturan.Status == "" {
		aturan.Status = "active"
	}

	return s.DB.Create(aturan).Error
}

// READ ALL
func (s *AturanDendaService) GetAllAturanDenda() ([]model.AturanDenda, error) {
	var aturans []model.AturanDenda
	result := s.DB.Find(&aturans)
	return aturans, result.Error
}

// READ BY ID
func (s *AturanDendaService) GetAturanDendaByID(id string) (*model.AturanDenda, error) {
	var aturan model.AturanDenda
	result := s.DB.First(&aturan, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &aturan, result.Error
}

// UPDATE
func (s *AturanDendaService) UpdateAturanDenda(aturan *model.AturanDenda) error {
	if aturan.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldAturan model.AturanDenda
	if err := s.DB.First(&oldAturan, "id = ?", aturan.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("aturan denda not found")
		}
		return err
	}

	if err := s.validateTipeDenda(aturan.TipeDenda); err != nil {
		return err
	}

	if err := s.validateStatus(aturan.Status); err != nil {
		return err
	}

	updateData := map[string]interface{}{}

	if aturan.UserID != "" && aturan.UserID != oldAturan.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", aturan.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = aturan.UserID
	}

	if aturan.TipeDenda != "" && aturan.TipeDenda != oldAturan.TipeDenda {
		updateData["tipe_denda"] = aturan.TipeDenda
	}

	if aturan.Persentase != nil {
		updateData["persentase"] = aturan.Persentase
	}

	if aturan.JumlahFlat != nil {
		updateData["jumlah_flat"] = aturan.JumlahFlat
	}

	if aturan.CapMaksimal != nil {
		updateData["cap_maksimal"] = aturan.CapMaksimal
	}

	if aturan.Catatan != nil {
		updateData["catatan"] = aturan.Catatan
	}

	if aturan.Status != "" && aturan.Status != oldAturan.Status {
		updateData["status"] = aturan.Status
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.AturanDenda{}).Where("id = ?", aturan.ID).Updates(updateData).Error
}

// DELETE
func (s *AturanDendaService) DeleteAturanDenda(id string) error {
	return s.DB.Delete(&model.AturanDenda{}, "id = ?", id).Error
}

func (s *AturanDendaService) GetAllAturanDendaByUserID(userID string) ([]model.AturanDenda, error) {
	var list []model.AturanDenda
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
