package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CicilanUserService struct {
	DB *gorm.DB
}

func NewCicilanUserService(db *gorm.DB) *CicilanUserService {
	return &CicilanUserService{DB: db}
}

// Validasi status
func (s *CicilanUserService) validateStatus(status string) error {
	validStatuses := []string{"belum", "telat", "lunas"}
	for _, v := range validStatuses {
		if status == v {
			return nil
		}
	}
	return fmt.Errorf("status must be one of: belum, telat, lunas")
}

// CREATE
func (s *CicilanUserService) CreateCicilanUser(cicilan *model.CicilanUser) error {
	if cicilan.ID == "" {
		cicilan.ID = uuid.New().String()
	}

	if cicilan.JumlahCicilan <= 0 {
		return fmt.Errorf("jumlah_cicilan must be > 0")
	}

	if cicilan.Denda < 0 {
		return fmt.Errorf("denda must be >= 0")
	}

	if err := s.validateStatus(cicilan.Status); err != nil {
		return err
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", cicilan.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	// Validasi FK kewajiban_id
	var kewajiban model.KewajibanUser
	if err := s.DB.First(&kewajiban, "id = ?", cicilan.KewajibanID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("kewajiban_id not found")
		}
		return err
	}

	// Validasi FK parent_id (jika disediakan)
	if cicilan.ParentID != "" {
		var parentCicilan model.CicilanUser
		if err := s.DB.First(&parentCicilan, "id = ?", cicilan.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent_id not found")
			}
			return err
		}
		// Opsional: Tambahkan logika bisnis, misalnya parent harus milik user yang sama
		if parentCicilan.UserID != cicilan.UserID {
			return fmt.Errorf("parent cicilan must belong to the same user")
		}
	}

	// Set default status jika kosong
	if cicilan.Status == "" {
		cicilan.Status = "belum"
	}

	return s.DB.Create(cicilan).Error
}

// READ ALL
func (s *CicilanUserService) GetAllCicilanUser() ([]model.CicilanUser, error) {
	var cicilans []model.CicilanUser
	result := s.DB.Preload("User").Preload("Kewajiban").Preload("Parent").Find(&cicilans) // Preload relasi termasuk Parent
	return cicilans, result.Error
}

// READ BY ID
func (s *CicilanUserService) GetCicilanUserByID(id string) (*model.CicilanUser, error) {
	var cicilan model.CicilanUser
	result := s.DB.Preload("User").Preload("Kewajiban").Preload("Parent").First(&cicilan, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &cicilan, result.Error
}

// UPDATE
func (s *CicilanUserService) UpdateCicilanUser(cicilan *model.CicilanUser) error {
	if cicilan.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldCicilan model.CicilanUser
	if err := s.DB.First(&oldCicilan, "id = ?", cicilan.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("cicilan user not found")
		}
		return err
	}

	if err := s.validateStatus(cicilan.Status); err != nil {
		return err
	}

	updateData := map[string]interface{}{}

	if cicilan.UserID != "" && cicilan.UserID != oldCicilan.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", cicilan.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = cicilan.UserID
	}

	if cicilan.KewajibanID != "" && cicilan.KewajibanID != oldCicilan.KewajibanID {
		// Validasi FK kewajiban_id
		var kewajiban model.KewajibanUser
		if err := s.DB.First(&kewajiban, "id = ?", cicilan.KewajibanID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("kewajiban_id not found")
			}
			return err
		}
		updateData["kewajiban_id"] = cicilan.KewajibanID
	}

	if cicilan.ParentID != "" && cicilan.ParentID != oldCicilan.ParentID {
		// Validasi FK parent_id
		var parentCicilan model.CicilanUser
		if err := s.DB.First(&parentCicilan, "id = ?", cicilan.ParentID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("parent_id not found")
			}
			return err
		}
		// Opsional: Tambahkan logika bisnis, misalnya parent harus milik user yang sama
		if parentCicilan.UserID != cicilan.UserID {
			return fmt.Errorf("parent cicilan must belong to the same user")
		}
		updateData["parent_id"] = cicilan.ParentID
	}

	if cicilan.JumlahCicilan > 0 && cicilan.JumlahCicilan != oldCicilan.JumlahCicilan {
		updateData["jumlah_cicilan"] = cicilan.JumlahCicilan
	}

	if !cicilan.JatuhTempo.IsZero() && cicilan.JatuhTempo != oldCicilan.JatuhTempo {
		updateData["jatuh_tempo"] = cicilan.JatuhTempo
	}

	if cicilan.Status != "" && cicilan.Status != oldCicilan.Status {
		updateData["status"] = cicilan.Status
	}

	if cicilan.Denda >= 0 && cicilan.Denda != oldCicilan.Denda {
		updateData["denda"] = cicilan.Denda
	}

	if cicilan.TanggalBayar != nil {
		updateData["tanggal_bayar"] = cicilan.TanggalBayar
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.CicilanUser{}).Where("id = ?", cicilan.ID).Updates(updateData).Error
}

// DELETE
func (s *CicilanUserService) DeleteCicilanUser(id string) error {
	// Opsional: Periksa apakah ada cicilan anak yang merujuk ke ini
	var childCount int64
	if err := s.DB.Model(&model.CicilanUser{}).Where("parent_id = ?", id).Count(&childCount).Error; err != nil {
		return err
	}
	if childCount > 0 {
		return fmt.Errorf("cannot delete cicilan with child cicilans")
	}

	return s.DB.Delete(&model.CicilanUser{}, "id = ?", id).Error
}

func (s *CicilanUserService) GetAllCicilanUserByUserID(userID string) ([]model.CicilanUser, error) {
	var cicilans []model.CicilanUser
	err := s.DB.
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Kewajiban").
		Preload("Parent").
		Find(&cicilans).Error
	return cicilans, err
}

func (s *CicilanUserService) GetCicilanByParentID(userID string, parentID string) ([]model.CicilanUser, error) {

	fmt.Println("=== GetCicilanByParentID CALLED ===")
	fmt.Println("UserID   :", userID)
	fmt.Println("ParentID :", parentID)

	var list []model.CicilanUser

	tx := s.DB.
		Preload("Kewajiban").
		Preload("Parent").
		Where("user_id = ? AND parent_id = ?", userID, parentID).
		Order("jatuh_tempo ASC")

	fmt.Println("Executing Query:", tx.Statement.SQL.String())

	err := tx.Find(&list).Error
	if err != nil {
		fmt.Println("❌ ERROR in GetCicilanByParentID:", err.Error())
		return nil, err
	}

	if len(list) == 0 {
		fmt.Println("⚠ No cicilan found for parent:", parentID)
	}

	fmt.Println("Total Cicilan Found:", len(list))
	fmt.Println("===================================")

	return list, nil
}

func (s *CicilanUserService) GetParentSummary(userID string) ([]map[string]interface{}, error) {

	fmt.Println("=== GetParentSummary CALLED ===")
	fmt.Println("UserID:", userID)

	var cicilan []model.CicilanUser

	err := s.DB.
		Where("user_id = ?", userID).
		Find(&cicilan).Error

	if err != nil {
		fmt.Println("❌ ERROR in GetParentSummary:", err.Error())
		return nil, err
	}

	fmt.Println("Total Cicilan Rows:", len(cicilan))

	if len(cicilan) == 0 {
		fmt.Println("⚠ No cicilan found for user:", userID)
	}

	result := map[string]map[string]interface{}{}

	for _, c := range cicilan {
		fmt.Printf("Processing Cicilan: ID=%s ParentID=%s Jumlah=%.2f\n",
			c.ID, c.ParentID, c.JumlahCicilan)

		parent := c.ParentID

		if _, exists := result[parent]; !exists {
			result[parent] = map[string]interface{}{
				"parent_id":      parent,
				"total_tagihan":  0.0,
				"total_bayar":    0.0,
				"total_sisa":     0.0,
				"jumlah_cicilan": 0,
			}
		}

		result[parent]["jumlah_cicilan"] = result[parent]["jumlah_cicilan"].(int) + 1
		result[parent]["total_tagihan"] = result[parent]["total_tagihan"].(float64) + c.JumlahCicilan
		result[parent]["total_sisa"] = result[parent]["total_sisa"].(float64) + c.JumlahCicilan
	}

	fmt.Println("Final Parent Summary:", result)
	fmt.Println("===================================")

	// convert ke array
	output := []map[string]interface{}{}
	for _, v := range result {
		output = append(output, v)
	}

	return output, nil
}
