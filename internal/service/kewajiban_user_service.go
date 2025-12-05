// package service

// import (
// 	"backend_go/internal/model"
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"time"

// 	"github.com/google/uuid"
// 	"gorm.io/gorm"
// )

// type KewajibanUserService struct {
// 	DB *gorm.DB
// }

// func NewKewajibanUserService(db *gorm.DB) *KewajibanUserService {
// 	return &KewajibanUserService{DB: db}
// }

// // Helper untuk validasi format tanggal (YYYY-MM-DD)
// func (s *KewajibanUserService) validateDateFormat(dateStr string) bool {
// 	if dateStr == "" {
// 		return false
// 	}
// 	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)
// 	if !matched {
// 		return false
// 	}
// 	_, err := time.Parse("2006-01-02", dateStr)
// 	return err == nil
// }

// // Validasi status
// func (s *KewajibanUserService) validateStatus(status string) error {
// 	if status != "active" && status != "lunas" {
// 		return fmt.Errorf("status must be 'active' or 'lunas'")
// 	}
// 	return nil
// }

// // CREATE KewajibanUser dari FormPembayaran
// func (s *KewajibanUserService) CreateKewajibanUser(kewajiban *model.KewajibanUser) error {
// 	if kewajiban.ID == "" {
// 		kewajiban.ID = uuid.New().String()
// 	}

// 	// Set default status jika kosong
// 	if kewajiban.Status == "" {
// 		kewajiban.Status = "active"
// 	}

// 	// Set default jumlah_sisa sama dengan jumlah_total jika kosong atau 0
// 	if kewajiban.JumlahSisa == 0 {
// 		kewajiban.JumlahSisa = kewajiban.JumlahTotal
// 	}

// 	// Validasi status
// 	if err := s.validateStatus(kewajiban.Status); err != nil {
// 		return err
// 	}

// 	// Validasi JumlahTotal
// 	if kewajiban.JumlahTotal <= 0 {
// 		return fmt.Errorf("jumlah_total must be > 0")
// 	}

// 	// Validasi JumlahSisa
// 	if kewajiban.JumlahSisa < 0 || kewajiban.JumlahSisa > kewajiban.JumlahTotal {
// 		return fmt.Errorf("jumlah_sisa must be >= 0 and <= jumlah_total")
// 	}

// 	// Validasi FK user_id
// 	var user model.User
// 	if err := s.DB.First(&user, "id = ?", kewajiban.UserID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("user_id not found")
// 		}
// 		return err
// 	}

// 	// Validasi FK form_id
// 	var form model.FormPembayaran
// 	if err := s.DB.Preload("Denda").First(&form, "id = ?", kewajiban.FormID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("form_id not found")
// 		}
// 		return err
// 	}

// 	// Validasi FK kategori_id
// 	var kategori model.KategoriPembayaran
// 	if err := s.DB.First(&kategori, "id = ?", kewajiban.KategoriID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("kategori_id not found")
// 		}
// 		return err
// 	}

// 	// Validasi FK denda_id
// 	var denda model.AturanDenda
// 	if err := s.DB.First(&denda, "id = ?", kewajiban.DendaID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("denda_id not found")
// 		}
// 		return err
// 	}

// 	// Validasi TanggalMulai (string, harus valid dan tidak kosong)
// 	if !s.validateDateFormat(kewajiban.TanggalMulai) {
// 		return fmt.Errorf("tanggal_mulai is required and must be in YYYY-MM-DD format")
// 	}

// 	// Validasi TanggalJatuhTempo (opsional, tapi jika ada harus valid)
// 	if kewajiban.TanggalJatuhTempo != nil && !s.validateDateFormat(*kewajiban.TanggalJatuhTempo) {
// 		return fmt.Errorf("tanggal_jatuh_tempo must be in YYYY-MM-DD format if provided")
// 	}

// 	// Simpan KewajibanUser
// 	if err := s.DB.Create(kewajiban).Error; err != nil {
// 		return err
// 	}

// 	// Jika boleh cicilan, generate CicilanUser otomatis
// 	if form.BolehCicilan && form.JumlahCicilan != nil && *form.JumlahCicilan > 0 {
// 		jumlahPerCicilan := kewajiban.JumlahTotal / float64(*form.JumlahCicilan)
// 		for i := 1; i <= *form.JumlahCicilan; i++ {
// 			cicilan := model.CicilanUser{
// 				ID:            uuid.New().String(),
// 				UserID:        kewajiban.UserID,
// 				KewajibanID:   kewajiban.ID,
// 				JumlahCicilan: jumlahPerCicilan,
// 				Status:        "belum",
// 				Denda:         0, // Default, bisa dihitung ketika terlambat bayar
// 			}
// 			// Set JatuhTempo cicilan (misal tiap bulan)
// 			tglMulai, _ := time.Parse("2006-01-02", kewajiban.TanggalMulai)
// 			cicilan.JatuhTempo = tglMulai.AddDate(0, i-1, 0)

// 			if err := s.DB.Create(&cicilan).Error; err != nil {
// 				return fmt.Errorf("failed to create cicilan_user: %v", err)
// 			}
// 		}
// 	}

// 	return nil
// }

// // READ ALL
// func (s *KewajibanUserService) GetAllKewajibanUser() ([]model.KewajibanUser, error) {
// 	var kewajibans []model.KewajibanUser
// 	result := s.DB.Preload("User").
// 		Preload("Form").
// 		Preload("Kategori").
// 		Preload("Denda").
// 		Find(&kewajibans)
// 	return kewajibans, result.Error
// }

// // READ BY ID
// func (s *KewajibanUserService) GetKewajibanUserByID(id string) (*model.KewajibanUser, error) {
// 	var kewajiban model.KewajibanUser
// 	result := s.DB.Preload("User").
// 		Preload("Form").
// 		Preload("Kategori").
// 		Preload("Denda").
// 		First(&kewajiban, "id = ?", id)
// 	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		return nil, nil
// 	}
// 	return &kewajiban, result.Error
// }

// // UPDATE
// func (s *KewajibanUserService) UpdateKewajibanUser(kewajiban *model.KewajibanUser) error {
// 	if kewajiban.ID == "" {
// 		return fmt.Errorf("ID is required for update")
// 	}

// 	var oldKewajiban model.KewajibanUser
// 	if err := s.DB.First(&oldKewajiban, "id = ?", kewajiban.ID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return fmt.Errorf("kewajiban user not found")
// 		}
// 		return err
// 	}

// 	// Set default status jika kosong
// 	if kewajiban.Status == "" {
// 		kewajiban.Status = oldKewajiban.Status
// 	}

// 	if err := s.validateStatus(kewajiban.Status); err != nil {
// 		return err
// 	}

// 	updateData := map[string]interface{}{}

// 	if kewajiban.UserID != "" && kewajiban.UserID != oldKewajiban.UserID {
// 		updateData["user_id"] = kewajiban.UserID
// 	}

// 	if kewajiban.FormID != "" && kewajiban.FormID != oldKewajiban.FormID {
// 		updateData["form_id"] = kewajiban.FormID
// 	}

// 	if kewajiban.KategoriID != "" && kewajiban.KategoriID != oldKewajiban.KategoriID {
// 		updateData["kategori_id"] = kewajiban.KategoriID
// 	}

// 	if kewajiban.DendaID != "" && kewajiban.DendaID != oldKewajiban.DendaID {
// 		updateData["denda_id"] = kewajiban.DendaID
// 	}

// 	if kewajiban.JumlahTotal > 0 && kewajiban.JumlahTotal != oldKewajiban.JumlahTotal {
// 		updateData["jumlah_total"] = kewajiban.JumlahTotal
// 	}

// 	if kewajiban.JumlahSisa >= 0 && kewajiban.JumlahSisa != oldKewajiban.JumlahSisa {
// 		updateData["jumlah_sisa"] = kewajiban.JumlahSisa
// 	}

// 	if kewajiban.Status != "" && kewajiban.Status != oldKewajiban.Status {
// 		updateData["status"] = kewajiban.Status
// 	}

// 	if kewajiban.TanggalMulai != "" && kewajiban.TanggalMulai != oldKewajiban.TanggalMulai {
// 		updateData["tanggal_mulai"] = kewajiban.TanggalMulai
// 	}

// 	if kewajiban.TanggalJatuhTempo != nil {
// 		updateData["tanggal_jatuh_tempo"] = kewajiban.TanggalJatuhTempo
// 	}

// 	if len(updateData) == 0 {
// 		return nil
// 	}

// 	return s.DB.Model(&model.KewajibanUser{}).Where("id = ?", kewajiban.ID).Updates(updateData).Error
// }

// // DELETE
// func (s *KewajibanUserService) DeleteKewajibanUser(id string) error {
// 	return s.DB.Delete(&model.KewajibanUser{}, "id = ?", id).Error
// }

//	func (s *KewajibanUserService) GetAllKewajibanUserByUserID(userID string) ([]model.KewajibanUser, error) {
//		var list []model.KewajibanUser
//		err := s.DB.Where("user_id = ?", userID).Find(&list).Error
//		return list, err
//	}
package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KewajibanUserService struct {
	DB *gorm.DB
}

func NewKewajibanUserService(db *gorm.DB) *KewajibanUserService {
	return &KewajibanUserService{DB: db}
}

// Helper untuk validasi format tanggal (YYYY-MM-DD)
func (s *KewajibanUserService) validateDateFormat(dateStr string) bool {
	if dateStr == "" {
		return false
	}
	matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, dateStr)
	if !matched {
		return false
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// Validasi status
func (s *KewajibanUserService) validateStatus(status string) error {
	if status != "active" && status != "lunas" {
		return fmt.Errorf("status must be 'active' or 'lunas'")
	}
	return nil
}

// ===================================================================================
// CREATE KewajibanUser (HANYA membutuhkan user_id dan form_id dari admin)
// ===================================================================================
// ===================================================================================
// CREATE KewajibanUser (AUTO-FILL semua field penting)
// ===================================================================================
func (s *KewajibanUserService) CreateKewajibanUser(kewajiban *model.KewajibanUser) error {

	// Generate ID jika kosong
	if kewajiban.ID == "" {
		kewajiban.ID = uuid.New().String()
	}

	// ============================================
	// VALIDASI USER
	// ============================================
	var user model.User
	if err := s.DB.First(&user, "id = ?", kewajiban.UserID).Error; err != nil {
		return fmt.Errorf("user_id not found")
	}

	// ============================================
	// VALIDASI FORM & AUTO-FILL
	// ============================================
	var form model.FormPembayaran
	if err := s.DB.First(&form, "id = ?", kewajiban.FormID).Error; err != nil {
		return fmt.Errorf("form_id not found")
	}

	// AUTO FILL
	kewajiban.KategoriID = form.KategoriID
	kewajiban.DendaID = form.DendaID
	kewajiban.JumlahTotal = form.Jumlah
	kewajiban.JumlahSisa = form.Jumlah
	kewajiban.TanggalMulai = form.TanggalMulai
	kewajiban.Status = "active"

	// ============================================
	// Generate tanggal_jatuh_tempo
	// ============================================
	var jatuhTempo time.Time
	if s.validateDateFormat(form.TanggalMulai) {
		startDate, _ := time.Parse("2006-01-02", form.TanggalMulai)
		jatuhTempo = startDate.AddDate(0, 1, 0)
	} else {
		jatuhTempo = time.Now().AddDate(0, 0, 30)
	}
	jt := jatuhTempo.Format("2006-01-02")
	kewajiban.TanggalJatuhTempo = &jt

	// ============================================
	// Simpan KewajibanUser
	// ============================================
	if err := s.DB.Create(kewajiban).Error; err != nil {
		return err
	}

	// ============================================
	// CICILAN OTOMATIS
	// ============================================
	if form.BolehCicilan && form.JumlahCicilan != nil && *form.JumlahCicilan > 0 {

		jumlahPerCicilan := kewajiban.JumlahTotal / float64(*form.JumlahCicilan)

		var tglMulai time.Time
		if s.validateDateFormat(kewajiban.TanggalMulai) {
			tglMulai, _ = time.Parse("2006-01-02", kewajiban.TanggalMulai)
		} else {
			tglMulai = time.Now()
		}

		for i := 0; i < *form.JumlahCicilan; i++ {
			jTempo := tglMulai.AddDate(0, i, 0)

			cicilan := model.CicilanUser{
				ID:            uuid.New().String(),
				UserID:        kewajiban.UserID,
				ParentID:      kewajiban.ParentID, // PAKAI ParentID dari KewajibanUser
				KewajibanID:   kewajiban.ID,
				JumlahCicilan: jumlahPerCicilan,
				Status:        "belum",
				Denda:         0,
				JatuhTempo:    jTempo,
			}

			if err := s.DB.Create(&cicilan).Error; err != nil {
				return fmt.Errorf("failed to create cicilan_user: %v", err)
			}
		}
	}

	// ============================================
	// TRANSAKSI OTOMATIS
	// ============================================
	transaksi := model.Transaksi{
		ID:               uuid.New().String(),
		UserID:           kewajiban.UserID,
		ParentID:         kewajiban.ParentID,
		KewajibanID:      kewajiban.ID,
		TipeTransaksi:    "penyesuaian", // default tipe
		Jumlah:           kewajiban.JumlahTotal,
		Tanggal:          time.Now(),
		StatusGateway:    nil,
		MetodePembayaran: nil,
		ReferenceID:      nil,
		Catatan:          nil,
	}

	if err := s.DB.Create(&transaksi).Error; err != nil {
		return fmt.Errorf("failed to create transaksi: %v", err)
	}

	return nil
}

// ===================================================================================
// READ ALL
// ===================================================================================
func (s *KewajibanUserService) GetAllKewajibanUser() ([]model.KewajibanUser, error) {
	var kewajibans []model.KewajibanUser
	result := s.DB.Preload("User").
		Preload("Form").
		Preload("Kategori").
		Preload("Denda").
		Find(&kewajibans)
	return kewajibans, result.Error
}

// ===================================================================================
// READ BY ID
// ===================================================================================
func (s *KewajibanUserService) GetKewajibanUserByID(id string) (*model.KewajibanUser, error) {
	var kewajiban model.KewajibanUser
	result := s.DB.Preload("User").
		Preload("Form").
		Preload("Kategori").
		Preload("Denda").
		First(&kewajiban, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &kewajiban, result.Error
}

// ===================================================================================
// UPDATE (Admin hanya boleh update beberapa field)
// ===================================================================================
func (s *KewajibanUserService) UpdateKewajibanUser(kewajiban *model.KewajibanUser) error {
	if kewajiban.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var old model.KewajibanUser
	if err := s.DB.First(&old, "id = ?", kewajiban.ID).Error; err != nil {
		return fmt.Errorf("kewajiban user not found")
	}

	updateData := map[string]interface{}{}

	// Admin tidak boleh update kategori_id, denda_id, jumlah_total → itu milik form

	// Boleh update status
	if kewajiban.Status != "" && kewajiban.Status != old.Status {
		if err := s.validateStatus(kewajiban.Status); err != nil {
			return err
		}
		updateData["status"] = kewajiban.Status
	}

	// Update jumlah_sisa jika diperlukan
	if kewajiban.JumlahSisa >= 0 && kewajiban.JumlahSisa != old.JumlahSisa {
		updateData["jumlah_sisa"] = kewajiban.JumlahSisa
	}

	// Update tanggal jatuh tempo
	if kewajiban.TanggalJatuhTempo != nil {
		updateData["tanggal_jatuh_tempo"] = kewajiban.TanggalJatuhTempo
	}

	if len(updateData) == 0 {
		return nil
	}

	return s.DB.Model(&model.KewajibanUser{}).Where("id = ?", kewajiban.ID).Updates(updateData).Error
}

// ===================================================================================
// DELETE
// ===================================================================================
func (s *KewajibanUserService) DeleteKewajibanUser(id string) error {
	return s.DB.Delete(&model.KewajibanUser{}, "id = ?", id).Error
}

// ===================================================================================
// GET ALL BY USER
// ===================================================================================
func (s *KewajibanUserService) GetAllKewajibanUserByUserID(userID string) ([]model.KewajibanUser, error) {
	var list []model.KewajibanUser
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
