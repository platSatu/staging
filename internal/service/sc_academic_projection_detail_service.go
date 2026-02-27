package service

import (
	"backend_go/internal/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ScAcademicProjectionDetailService struct {
	DB *gorm.DB
}

func NewScAcademicProjectionDetailService(db *gorm.DB) *ScAcademicProjectionDetailService {
	return &ScAcademicProjectionDetailService{DB: db}
}

// CREATE
func (s *ScAcademicProjectionDetailService) CreateScAcademicProjectionDetail(detail *model.ScAcademicProjectionDetail) error {
	if detail.ID == "" {
		detail.ID = uuid.New().String()
	}

	if detail.UserID == "" {
		return fmt.Errorf("user_id is required")
	}

	// Validasi FK user_id
	var user model.User
	if err := s.DB.First(&user, "id = ?", detail.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user_id not found")
		}
		return err
	}

	return s.DB.Create(detail).Error
}

// CREATE MULTIPLE (untuk foreach saat create academic projection)
func (s *ScAcademicProjectionDetailService) CreateMultipleScAcademicProjectionDetail(details []*model.ScAcademicProjectionDetail) error {
	for _, detail := range details {
		if err := s.CreateScAcademicProjectionDetail(detail); err != nil {
			return err
		}
	}
	return nil
}

// READ ALL
func (s *ScAcademicProjectionDetailService) GetAllScAcademicProjectionDetail() ([]model.ScAcademicProjectionDetail, error) {
	var details []model.ScAcademicProjectionDetail
	result := s.DB.Find(&details)
	return details, result.Error
}

// READ BY ID
func (s *ScAcademicProjectionDetailService) GetScAcademicProjectionDetailByID(id string) (*model.ScAcademicProjectionDetail, error) {
	var detail model.ScAcademicProjectionDetail
	result := s.DB.First(&detail, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &detail, result.Error
}

// UPDATE
func (s *ScAcademicProjectionDetailService) UpdateScAcademicProjectionDetail(detail *model.ScAcademicProjectionDetail) error {
	if detail.ID == "" {
		return fmt.Errorf("ID is required for update")
	}

	var oldDetail model.ScAcademicProjectionDetail
	if err := s.DB.First(&oldDetail, "id = ?", detail.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("sc_academic_projection_detail not found")
		}
		return err
	}

	updateData := map[string]interface{}{}

	if detail.UserID != "" && detail.UserID != oldDetail.UserID {
		// Validasi FK user_id
		var user model.User
		if err := s.DB.First(&user, "id = ?", detail.UserID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("user_id not found")
			}
			return err
		}
		updateData["user_id"] = detail.UserID
	}

	if detail.SubjectTypeID != nil && (oldDetail.SubjectTypeID == nil || *detail.SubjectTypeID != *oldDetail.SubjectTypeID) {
		updateData["subject_type_id"] = detail.SubjectTypeID
	}

	if detail.No != nil && (oldDetail.No == nil || *detail.No != *oldDetail.No) {
		updateData["no"] = detail.No
	}

	if detail.SubjectNameID != nil && (oldDetail.SubjectNameID == nil || *detail.SubjectNameID != *oldDetail.SubjectNameID) {
		updateData["subject_name_id"] = detail.SubjectNameID
	}

	if detail.Status != nil && (oldDetail.Status == nil || *detail.Status != *oldDetail.Status) {
		updateData["status"] = detail.Status
	}

	if detail.IssueDate != nil && (oldDetail.IssueDate == nil || !detail.IssueDate.Equal(*oldDetail.IssueDate)) {
		updateData["issue_date"] = detail.IssueDate
	}

	if detail.PtDate != nil && (oldDetail.PtDate == nil || !detail.PtDate.Equal(*oldDetail.PtDate)) {
		updateData["pt_date"] = detail.PtDate
	}

	if detail.PtScore != nil && (oldDetail.PtScore == nil || *detail.PtScore != *oldDetail.PtScore) {
		updateData["pt_score"] = detail.PtScore
	}

	if detail.AlphabetPtScore != nil && (oldDetail.AlphabetPtScore == nil || *detail.AlphabetPtScore != *oldDetail.AlphabetPtScore) {
		updateData["alphabet_pt_score"] = detail.AlphabetPtScore
	}

	if detail.EndDate != nil && (oldDetail.EndDate == nil || !detail.EndDate.Equal(*oldDetail.EndDate)) {
		updateData["end_date"] = detail.EndDate
	}

	if detail.Paces != nil && (oldDetail.Paces == nil || *detail.Paces != *oldDetail.Paces) {
		updateData["paces"] = detail.Paces
	}

	if detail.PrevPace != nil && (oldDetail.PrevPace == nil || *detail.PrevPace != *oldDetail.PrevPace) {
		updateData["prev_pace"] = detail.PrevPace
	}

	if detail.NextPace != nil && (oldDetail.NextPace == nil || *detail.NextPace != *oldDetail.NextPace) {
		updateData["next_pace"] = detail.NextPace
	}

	if detail.Pages != nil && (oldDetail.Pages == nil || *detail.Pages != *oldDetail.Pages) {
		updateData["pages"] = detail.Pages
	}

	if detail.OrderListID != nil && (oldDetail.OrderListID == nil || *detail.OrderListID != *oldDetail.OrderListID) {
		updateData["order_list_id"] = detail.OrderListID
	}

	if detail.Order != nil && (oldDetail.Order == nil || *detail.Order != *oldDetail.Order) {
		updateData["order"] = detail.Order
	}

	if detail.ProductID != nil && (oldDetail.ProductID == nil || *detail.ProductID != *oldDetail.ProductID) {
		updateData["product_id"] = detail.ProductID
	}

	if detail.AcademicYearID != nil && (oldDetail.AcademicYearID == nil || *detail.AcademicYearID != *oldDetail.AcademicYearID) {
		updateData["academic_year_id"] = detail.AcademicYearID
	}

	if detail.IsProcessed != nil && (oldDetail.IsProcessed == nil || *detail.IsProcessed != *oldDetail.IsProcessed) {
		updateData["is_processed"] = detail.IsProcessed
	}

	if detail.OrderNote != nil && (oldDetail.OrderNote == nil || *detail.OrderNote != *oldDetail.OrderNote) {
		updateData["order_note"] = detail.OrderNote
	}

	if detail.SubjectID != nil && (oldDetail.SubjectID == nil || *detail.SubjectID != *oldDetail.SubjectID) {
		updateData["subject_id"] = detail.SubjectID
	}

	if detail.AssignmentID != nil && (oldDetail.AssignmentID == nil || *detail.AssignmentID != *oldDetail.AssignmentID) {
		updateData["assignment_id"] = detail.AssignmentID
	}

	if len(updateData) == 0 {
		return nil // Tidak ada yang diupdate
	}

	return s.DB.Model(&model.ScAcademicProjectionDetail{}).Where("id = ?", detail.ID).Updates(updateData).Error
}

// DELETE
func (s *ScAcademicProjectionDetailService) DeleteScAcademicProjectionDetail(id string) error {
	return s.DB.Delete(&model.ScAcademicProjectionDetail{}, "id = ?", id).Error
}

func (s *ScAcademicProjectionDetailService) GetAllScAcademicProjectionDetailByUserID(userID string) ([]model.ScAcademicProjectionDetail, error) {
	var list []model.ScAcademicProjectionDetail
	err := s.DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}
