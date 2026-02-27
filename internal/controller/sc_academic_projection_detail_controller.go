package controller

import (
	"backend_go/internal/model"
	"backend_go/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ScAcademicProjectionDetailController struct {
	Service *service.ScAcademicProjectionDetailService
}

func NewScAcademicProjectionDetailController(s *service.ScAcademicProjectionDetailService) *ScAcademicProjectionDetailController {
	return &ScAcademicProjectionDetailController{Service: s}
}

// CreateScAcademicProjectionDetail (single)
func (sc *ScAcademicProjectionDetailController) CreateScAcademicProjectionDetail(c *gin.Context) {
	var req struct {
		SubjectTypeID   *string  `json:"subject_type_id"`
		No              *int     `json:"no"`
		SubjectNameID   *string  `json:"subject_name_id"`
		Status          *string  `json:"status"`
		IssueDate       *string  `json:"issue_date"` // Format: YYYY-MM-DD
		PtDate          *string  `json:"pt_date"`    // Format: YYYY-MM-DD
		PtScore         *float64 `json:"pt_score"`
		AlphabetPtScore *string  `json:"alphabet_pt_score"`
		EndDate         *string  `json:"end_date"` // Format: YYYY-MM-DD
		Paces           *int     `json:"paces"`
		PrevPace        *int     `json:"prev_pace"`
		NextPace        *int     `json:"next_pace"`
		Pages           *int     `json:"pages"`
		OrderListID     *string  `json:"order_list_id"`
		Order           *bool    `json:"order"`
		ProductID       *string  `json:"product_id"`
		AcademicYearID  *string  `json:"academic_year_id"`
		IsProcessed     *bool    `json:"is_processed"`
		OrderNote       *string  `json:"order_note"`
		SubjectID       *string  `json:"subject_id"`
		AssignmentID    *string  `json:"assignment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Parse dates
	var issueDate, ptDate, endDate *time.Time
	if req.IssueDate != nil && *req.IssueDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.IssueDate); err == nil {
			issueDate = &parsed
		}
	}
	if req.PtDate != nil && *req.PtDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.PtDate); err == nil {
			ptDate = &parsed
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &parsed
		}
	}

	// ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	detail := &model.ScAcademicProjectionDetail{
		UserID:          userID.(string),
		SubjectTypeID:   req.SubjectTypeID,
		No:              req.No,
		SubjectNameID:   req.SubjectNameID,
		Status:          req.Status,
		IssueDate:       issueDate,
		PtDate:          ptDate,
		PtScore:         req.PtScore,
		AlphabetPtScore: req.AlphabetPtScore,
		EndDate:         endDate,
		Paces:           req.Paces,
		PrevPace:        req.PrevPace,
		NextPace:        req.NextPace,
		Pages:           req.Pages,
		OrderListID:     req.OrderListID,
		Order:           req.Order,
		ProductID:       req.ProductID,
		AcademicYearID:  req.AcademicYearID,
		IsProcessed:     req.IsProcessed,
		OrderNote:       req.OrderNote,
		SubjectID:       req.SubjectID,
		AssignmentID:    req.AssignmentID,
	}

	if err := sc.Service.CreateScAcademicProjectionDetail(detail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    detail,
	})
}

// CreateMultipleScAcademicProjectionDetail (untuk foreach saat create academic projection)
func (sc *ScAcademicProjectionDetailController) CreateMultipleScAcademicProjectionDetail(c *gin.Context) {
	var req []struct {
		SubjectTypeID   *string  `json:"subject_type_id"`
		No              *int     `json:"no"`
		SubjectNameID   *string  `json:"subject_name_id"`
		Status          *string  `json:"status"`
		IssueDate       *string  `json:"issue_date"` // Format: YYYY-MM-DD
		PtDate          *string  `json:"pt_date"`    // Format: YYYY-MM-DD
		PtScore         *float64 `json:"pt_score"`
		AlphabetPtScore *string  `json:"alphabet_pt_score"`
		EndDate         *string  `json:"end_date"` // Format: YYYY-MM-DD
		Paces           *int     `json:"paces"`
		PrevPace        *int     `json:"prev_pace"`
		NextPace        *int     `json:"next_pace"`
		Pages           *int     `json:"pages"`
		OrderListID     *string  `json:"order_list_id"`
		Order           *bool    `json:"order"`
		ProductID       *string  `json:"product_id"`
		AcademicYearID  *string  `json:"academic_year_id"`
		IsProcessed     *bool    `json:"is_processed"`
		OrderNote       *string  `json:"order_note"`
		SubjectID       *string  `json:"subject_id"`
		AssignmentID    *string  `json:"assignment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// ambil user_id dari token
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	var details []*model.ScAcademicProjectionDetail
	for _, r := range req {
		// Parse dates
		var issueDate, ptDate, endDate *time.Time
		if r.IssueDate != nil && *r.IssueDate != "" {
			if parsed, err := time.Parse("2006-01-02", *r.IssueDate); err == nil {
				issueDate = &parsed
			}
		}
		if r.PtDate != nil && *r.PtDate != "" {
			if parsed, err := time.Parse("2006-01-02", *r.PtDate); err == nil {
				ptDate = &parsed
			}
		}
		if r.EndDate != nil && *r.EndDate != "" {
			if parsed, err := time.Parse("2006-01-02", *r.EndDate); err == nil {
				endDate = &parsed
			}
		}

		detail := &model.ScAcademicProjectionDetail{
			UserID:          userID.(string),
			SubjectTypeID:   r.SubjectTypeID,
			No:              r.No,
			SubjectNameID:   r.SubjectNameID,
			Status:          r.Status,
			IssueDate:       issueDate,
			PtDate:          ptDate,
			PtScore:         r.PtScore,
			AlphabetPtScore: r.AlphabetPtScore,
			EndDate:         endDate,
			Paces:           r.Paces,
			PrevPace:        r.PrevPace,
			NextPace:        r.NextPace,
			Pages:           r.Pages,
			OrderListID:     r.OrderListID,
			Order:           r.Order,
			ProductID:       r.ProductID,
			AcademicYearID:  r.AcademicYearID,
			IsProcessed:     r.IsProcessed,
			OrderNote:       r.OrderNote,
			SubjectID:       r.SubjectID,
			AssignmentID:    r.AssignmentID,
		}
		details = append(details, detail)
	}

	if err := sc.Service.CreateMultipleScAcademicProjectionDetail(details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    details,
	})
}

// GetAllScAcademicProjectionDetail (hanya milik user login)
func (sc *ScAcademicProjectionDetailController) GetAllScAcademicProjectionDetail(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   "unauthorized",
		})
		return
	}

	details, err := sc.Service.GetAllScAcademicProjectionDetailByUserID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    details,
	})
}

// GetScAcademicProjectionDetailByID
func (sc *ScAcademicProjectionDetailController) GetScAcademicProjectionDetailByID(c *gin.Context) {
	id := c.Param("id")

	detail, err := sc.Service.GetScAcademicProjectionDetailByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	if detail == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_projection_detail not found",
		})
		return
	}

	// hanya pemilik detail boleh melihat
	userID, _ := c.Get("userID")
	if detail.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    detail,
	})
}

// UpdateScAcademicProjectionDetail
func (sc *ScAcademicProjectionDetailController) UpdateScAcademicProjectionDetail(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		SubjectTypeID   *string  `json:"subject_type_id"`
		No              *int     `json:"no"`
		SubjectNameID   *string  `json:"subject_name_id"`
		Status          *string  `json:"status"`
		IssueDate       *string  `json:"issue_date"` // Format: YYYY-MM-DD
		PtDate          *string  `json:"pt_date"`    // Format: YYYY-MM-DD
		PtScore         *float64 `json:"pt_score"`
		AlphabetPtScore *string  `json:"alphabet_pt_score"`
		EndDate         *string  `json:"end_date"` // Format: YYYY-MM-DD
		Paces           *int     `json:"paces"`
		PrevPace        *int     `json:"prev_pace"`
		NextPace        *int     `json:"next_pace"`
		Pages           *int     `json:"pages"`
		OrderListID     *string  `json:"order_list_id"`
		Order           *bool    `json:"order"`
		ProductID       *string  `json:"product_id"`
		AcademicYearID  *string  `json:"academic_year_id"`
		IsProcessed     *bool    `json:"is_processed"`
		OrderNote       *string  `json:"order_note"`
		SubjectID       *string  `json:"subject_id"`
		AssignmentID    *string  `json:"assignment_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Parse dates
	var issueDate, ptDate, endDate *time.Time
	if req.IssueDate != nil && *req.IssueDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.IssueDate); err == nil {
			issueDate = &parsed
		}
	}
	if req.PtDate != nil && *req.PtDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.PtDate); err == nil {
			ptDate = &parsed
		}
	}
	if req.EndDate != nil && *req.EndDate != "" {
		if parsed, err := time.Parse("2006-01-02", *req.EndDate); err == nil {
			endDate = &parsed
		}
	}

	// pastikan detail itu milik user login
	existing, err := sc.Service.GetScAcademicProjectionDetailByID(id)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_projection_detail not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if existing.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	updateData := &model.ScAcademicProjectionDetail{
		ID:              id,
		SubjectTypeID:   req.SubjectTypeID,
		No:              req.No,
		SubjectNameID:   req.SubjectNameID,
		Status:          req.Status,
		IssueDate:       issueDate,
		PtDate:          ptDate,
		PtScore:         req.PtScore,
		AlphabetPtScore: req.AlphabetPtScore,
		EndDate:         endDate,
		Paces:           req.Paces,
		PrevPace:        req.PrevPace,
		NextPace:        req.NextPace,
		Pages:           req.Pages,
		OrderListID:     req.OrderListID,
		Order:           req.Order,
		ProductID:       req.ProductID,
		AcademicYearID:  req.AcademicYearID,
		IsProcessed:     req.IsProcessed,
		OrderNote:       req.OrderNote,
		SubjectID:       req.SubjectID,
		AssignmentID:    req.AssignmentID,
	}

	if err := sc.Service.UpdateScAcademicProjectionDetail(updateData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	updated, _ := sc.Service.GetScAcademicProjectionDetailByID(id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    updated,
	})
}

// DeleteScAcademicProjectionDetail
func (sc *ScAcademicProjectionDetailController) DeleteScAcademicProjectionDetail(c *gin.Context) {
	id := c.Param("id")

	detail, err := sc.Service.GetScAcademicProjectionDetailByID(id)
	if err != nil || detail == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "sc_academic_projection_detail not found",
		})
		return
	}

	userID, _ := c.Get("userID")
	if detail.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "forbidden",
		})
		return
	}

	if err := sc.Service.DeleteScAcademicProjectionDetail(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "sc_academic_projection_detail deleted",
	})
}
