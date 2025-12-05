package service

import (
	"fmt"

	"gorm.io/gorm"
)

type LaporanService struct {
	DB *gorm.DB
}

func NewLaporanService(db *gorm.DB) *LaporanService {
	return &LaporanService{DB: db}
}

type LaporanKategoriResponse struct {
	KategoriID   string              `json:"kategori_id"`
	NamaKategori string              `json:"nama_kategori"`
	JumlahUser   int64               `json:"jumlah_user"`
	TotalTagihan float64             `json:"total_tagihan"`
	TotalBayar   float64             `json:"total_bayar"`
	TotalSisa    float64             `json:"total_sisa"`
	Detail       []LaporanUserDetail `json:"detail,omitempty"`
}

type LaporanUserDetail struct {
	UserID      string        `json:"user_id"`
	NamaUser    string        `json:"nama_user"`
	JumlahTotal float64       `json:"jumlah_total"`
	JumlahBayar float64       `json:"jumlah_bayar"`
	JumlahSisa  float64       `json:"jumlah_sisa"`
	Status      string        `json:"status"`
	History     []HistoryItem `json:"history"`
}

type HistoryItem struct {
	FormNama    string  `json:"form_nama"`
	JumlahTotal float64 `json:"jumlah_total"`
	JumlahBayar float64 `json:"jumlah_bayar"`
	JumlahSisa  float64 `json:"jumlah_sisa"`
	Status      string  `json:"status"`
	Denda       float64 `json:"denda"`
	AdaCicilan  bool    `json:"ada_cicilan"` // Tambahan: true jika jumlah_bayar > 0
}

func (s *LaporanService) GetLaporanPerKategori(userID string) ([]LaporanKategoriResponse, error) {
	// Ambil semua anak berdasarkan partner_id (parent)
	var childIDs []string
	if err := s.DB.Raw(`
		SELECT id 
		FROM users 
		WHERE parent_id = ?
	`, userID).Scan(&childIDs).Error; err != nil {
		return nil, fmt.Errorf("error fetch childIDs: %w", err)
	}

	if len(childIDs) == 0 {
		return []LaporanKategoriResponse{}, nil
	}

	// --- 1. Summary per kategori (tampilkan semua kategori, bahkan yang kosong) ---
	var kategoriList []struct {
		KategoriID   string
		NamaKategori string
		JumlahUser   int64
		TotalTagihan float64
		TotalBayar   float64
		TotalSisa    float64
	}

	summaryQuery := `
	SELECT
		kp.id AS kategori_id,
		kp.nama_kategori AS nama_kategori,
		COALESCE((
			SELECT COUNT(DISTINCT ku.parent_id)
			FROM kewajiban_user ku
			WHERE ku.kategori_id = kp.id AND ku.parent_id IN (?)
		), 0) AS jumlah_user,
		COALESCE((
			SELECT SUM(ku.jumlah_total)
			FROM kewajiban_user ku
			WHERE ku.kategori_id = kp.id AND ku.parent_id IN (?)
		), 0) AS total_tagihan,
		COALESCE((
			SELECT SUM(
				(SELECT COALESCE(SUM(cu.jumlah_cicilan), 0)
				FROM cicilan_user cu
				WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id)
			)
			FROM kewajiban_user ku
			WHERE ku.kategori_id = kp.id AND ku.parent_id IN (?)
		), 0) AS total_bayar,
		COALESCE((
			SELECT SUM(
				ku.jumlah_total - COALESCE((
					SELECT SUM(cu.jumlah_cicilan)
					FROM cicilan_user cu
					WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id
				), 0)
			)
			FROM kewajiban_user ku
			WHERE ku.kategori_id = kp.id AND ku.parent_id IN (?)
		), 0) AS total_sisa
	FROM kategori_pembayaran kp
	ORDER BY kp.nama_kategori ASC
	`

	if err := s.DB.Raw(summaryQuery, childIDs, childIDs, childIDs, childIDs).Scan(&kategoriList).Error; err != nil {
		return nil, err
	}

	// --- 2. Detail per user per kategori ---
	var finalResponse []LaporanKategoriResponse

	for _, k := range kategoriList {
		var detailList []struct {
			UserID      string
			NamaUser    string
			JumlahTotal float64
			JumlahBayar float64
			JumlahSisa  float64
			Status      string
		}

		detailQuery := `
		SELECT
			ku.parent_id AS user_id,
			u.username AS nama_user,
			ku.jumlah_total,
			COALESCE((SELECT SUM(cu.jumlah_cicilan)
				FROM cicilan_user cu
				WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id),0) AS jumlah_bayar,
			ku.jumlah_total -
			COALESCE((SELECT SUM(cu.jumlah_cicilan)
				FROM cicilan_user cu
				WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id),0) AS jumlah_sisa,
			ku.status
		FROM kewajiban_user ku
		JOIN users u ON u.id = ku.parent_id
		WHERE ku.kategori_id = ? AND ku.parent_id IN (?)
		ORDER BY u.username ASC
		`

		if err := s.DB.Raw(detailQuery, k.KategoriID, childIDs).Scan(&detailList).Error; err != nil {
			return nil, err
		}

		// --- 3. History tiap user ---
		var laporanDetail []LaporanUserDetail

		for _, d := range detailList {
			var historyRows []HistoryItem

			historyQuery := `
			SELECT
				f.nama_form AS form_nama,
				ku.jumlah_total AS jumlah_total,
				COALESCE((SELECT SUM(cu.jumlah_cicilan)
					FROM cicilan_user cu
					WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id),0) AS jumlah_bayar,
				ku.jumlah_total -
				COALESCE((SELECT SUM(cu.jumlah_cicilan)
					FROM cicilan_user cu
					WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id),0) AS jumlah_sisa,
				ku.status,
				COALESCE((SELECT SUM(cu.denda)
					FROM cicilan_user cu
					WHERE cu.kewajiban_id = ku.id AND cu.parent_id = ku.parent_id),0) AS denda
			FROM kewajiban_user ku
			JOIN form_pembayaran f ON f.id = ku.form_id
			WHERE ku.parent_id = ?
			ORDER BY f.nama_form ASC
			`

			if err := s.DB.Raw(historyQuery, d.UserID).Scan(&historyRows).Error; err != nil {
				return nil, err
			}

			// Tambahkan indikator ada cicilan
			for i := range historyRows {
				historyRows[i].AdaCicilan = historyRows[i].JumlahBayar > 0
			}

			laporanDetail = append(laporanDetail, LaporanUserDetail{
				UserID:      d.UserID,
				NamaUser:    d.NamaUser,
				JumlahTotal: d.JumlahTotal,
				JumlahBayar: d.JumlahBayar,
				JumlahSisa:  d.JumlahSisa,
				Status:      d.Status,
				History:     historyRows,
			})
		}

		finalResponse = append(finalResponse, LaporanKategoriResponse{
			KategoriID:   k.KategoriID,
			NamaKategori: k.NamaKategori,
			JumlahUser:   k.JumlahUser,
			TotalTagihan: k.TotalTagihan,
			TotalBayar:   k.TotalBayar,
			TotalSisa:    k.TotalSisa,
			Detail:       laporanDetail,
		})
	}

	return finalResponse, nil
}

func (s *LaporanService) GetDetailPerForm(formID string, parentID string) ([]HistoryItem, error) {
	// Fungsi ini tetap ada jika diperlukan untuk endpoint lain, tapi tidak digunakan di sini.
	var historyRows []struct {
		FormNama    string
		JumlahTotal float64
		JumlahBayar float64
		JumlahSisa  float64
		Status      string
		Denda       float64
	}

	historyQuery := `
	SELECT
		f.nama_form AS form_nama,
		ku.jumlah_total AS jumlah_total,
		IFNULL((
			SELECT SUM(cu.jumlah_cicilan)
			FROM cicilan_user cu
			WHERE cu.kewajiban_id = ku.id 
			AND cu.parent_id = ?
		), 0) AS jumlah_bayar,
		ku.jumlah_total - IFNULL((
			SELECT SUM(cu.jumlah_cicilan)
			FROM cicilan_user cu
			WHERE cu.kewajiban_id = ku.id 
			AND cu.parent_id = ?
		), 0) AS jumlah_sisa,
		ku.status,
		IFNULL((
			SELECT SUM(cu.denda)
			FROM cicilan_user cu
			WHERE cu.kewajiban_id = ku.id 
			AND cu.parent_id = ?
		), 0) AS denda
	FROM kewajiban_user ku
	JOIN form_pembayaran f ON f.id = ku.form_id
	WHERE ku.form_id = ? AND ku.parent_id = ?
	ORDER BY ku.tanggal_mulai ASC
	`

	if err := s.DB.Raw(historyQuery, parentID, parentID, parentID, formID, parentID).Scan(&historyRows).Error; err != nil {
		return nil, fmt.Errorf("error GetDetailPerForm form %s: %w", formID, err)
	}

	var history []HistoryItem
	for _, h := range historyRows {
		history = append(history, HistoryItem{
			FormNama:    h.FormNama,
			JumlahTotal: h.JumlahTotal,
			JumlahBayar: h.JumlahBayar,
			JumlahSisa:  h.JumlahSisa,
			Status:      h.Status,
			Denda:       h.Denda,
			AdaCicilan:  h.JumlahBayar > 0, // Tambahan
		})
	}

	return history, nil
}
