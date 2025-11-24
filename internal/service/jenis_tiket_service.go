package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JenisTiketService struct {
	DB *gorm.DB
}

func NewJenisTiketService(db *gorm.DB) *JenisTiketService {
	return &JenisTiketService{DB: db}
}

// CREATE
func (s *JenisTiketService) CreateJenisTiket(jenisTiket *model.JenisTiket) error {
	if jenisTiket.ID == "" {
		jenisTiket.ID = uuid.New().String()
	}

	if jenisTiket.Status == "" {
		jenisTiket.Status = "active"
	}

	if jenisTiket.Terjual == 0 {
		jenisTiket.Terjual = 0
	}

	// Validasi user_id ada di tabel users
	var user model.User
	if err := s.DB.First(&user, "id = ?", jenisTiket.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found")
		}
		return err
	}

	// Validasi event_id ada di tabel events
	var event model.Event
	if err := s.DB.First(&event, "id = ?", jenisTiket.EventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("event not found")
		}
		return err
	}

	return s.DB.Create(jenisTiket).Error
}

// READ ALL
func (s *JenisTiketService) GetAllJenisTiket() ([]model.JenisTiket, error) {
	var jenisTikets []model.JenisTiket
	result := s.DB.Find(&jenisTikets)
	return jenisTikets, result.Error
}

// READ BY ID
func (s *JenisTiketService) GetJenisTiketByID(id string) (*model.JenisTiket, error) {
	var jenisTiket model.JenisTiket
	result := s.DB.First(&jenisTiket, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &jenisTiket, result.Error
}

// READ BY USER ID
func (s *JenisTiketService) GetJenisTiketByUserID(userID string) ([]model.JenisTiket, error) {
	var jenisTikets []model.JenisTiket
	result := s.DB.Where("user_id = ?", userID).Find(&jenisTikets)
	return jenisTikets, result.Error
}

// READ BY EVENT ID
func (s *JenisTiketService) GetJenisTiketByEventID(eventID string) ([]model.JenisTiket, error) {
	var jenisTikets []model.JenisTiket
	result := s.DB.Where("event_id = ?", eventID).Find(&jenisTikets)
	return jenisTikets, result.Error
}

// UPDATE
func (s *JenisTiketService) UpdateJenisTiket(jenisTiket *model.JenisTiket) error {
	if jenisTiket.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldJenisTiket model.JenisTiket
	if err := s.DB.First(&oldJenisTiket, "id = ?", jenisTiket.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("jenis tiket not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if jenisTiket.Name != "" && jenisTiket.Name != oldJenisTiket.Name {
		updateData["name"] = jenisTiket.Name
	}

	if jenisTiket.Stok != 0 && jenisTiket.Stok != oldJenisTiket.Stok {
		updateData["stok"] = jenisTiket.Stok
	}

	if jenisTiket.Terjual != oldJenisTiket.Terjual {
		updateData["terjual"] = jenisTiket.Terjual
	}

	if jenisTiket.Sisa != 0 && jenisTiket.Sisa != oldJenisTiket.Sisa {
		updateData["sisa"] = jenisTiket.Sisa
	}

	if jenisTiket.Harga != 0 && jenisTiket.Harga != oldJenisTiket.Harga {
		updateData["harga"] = jenisTiket.Harga
	}

	if jenisTiket.Status != "" && jenisTiket.Status != oldJenisTiket.Status {
		updateData["status"] = jenisTiket.Status
	}

	if jenisTiket.UserID != "" && jenisTiket.UserID != oldJenisTiket.UserID {
		// Validasi user_id baru ada
		var user model.User
		if err := s.DB.First(&user, "id = ?", jenisTiket.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user not found")
			}
			return err
		}
		updateData["user_id"] = jenisTiket.UserID
	}

	if jenisTiket.EventID != "" && jenisTiket.EventID != oldJenisTiket.EventID {
		// Validasi event_id baru ada
		var event model.Event
		if err := s.DB.First(&event, "id = ?", jenisTiket.EventID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("event not found")
			}
			return err
		}
		updateData["event_id"] = jenisTiket.EventID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.JenisTiket{}).Where("id = ?", jenisTiket.ID).Updates(updateData).Error
}

// DELETE
func (s *JenisTiketService) DeleteJenisTiket(id string) error {
	return s.DB.Delete(&model.JenisTiket{}, "id = ?", id).Error
}
