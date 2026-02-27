package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Panggil route modul user langsung ke root, tanpa prefix /api
	InitUserRoutes(r, db)

	InitAuthRoutes(r, db)

	InitCategoryAplikasiRoutes(r, db)

	InitEventRoutes(r, db)

	InitJenisTiketRoutes(r, db)

	InitTiketRoutes(r, db)

	InitEventUserRoutes(r, db)

	InitPaymentUserRoutes(r, db)

	InitPaymentCategoryRoutes(r, db)

	InitPaymentFormRoutes(r, db)

	//InitPaymentStudentSettingsRoutes(r, db)

	InitPaymentInvoiceRoutes(r, db) // Diperbaiki dari InitPaymentInvoicesRoutes

	//InitPaymentInvoicesRoutes(r, db)

	InitPaymentInstallmentsRoutes(r, db)

	InitPaymentPenaltiesRoutes(r, db)

	InitPaymentPaymentsRoutes(r, db)

	InitPaymentPenaltySettingsRoutes(r, db)

	InitKategoriPembayaranRoutes(r, db)

	InitAturanDendaRoutes(r, db)

	InitFormPembayaranRoutes(r, db)

	InitKewajibanUserRoutes(r, db)

	InitTransaksiRoutes(r, db)

	InitCicilanUserRoutes(r, db)

	InitLaporanRoutes(r, db)

	InitCategoryPackagesRoutes(r, db) // Tambahkan ini

	InitPackagesRoutes(r, db)

	InitProfileRoutes(r, db)

	InitVoucherRoutes(r, db)

	InitDepositRoutes(r, db)

	InitTicketEventRoutes(r, db)

	InitTicketRegisterRoutes(r, db)

	InitTicketKategoriRoutes(r, db)

	InitTicketJenisQuantityRoutes(r, db)

	InitTicketMethodPembayaranRoutes(r, db)

	InitTicketQrcodeRoutes(r, db)

	InitTicketVoucherRoutes(r, db)

	InitTicketResellerSettingRoutes(r, db)

	InitTicketHistoryRoutes(r, db)

	InitTicketEmailKategoryRoutes(r, db)

	InitTicketTemplateRoutes(r, db)

	InitTicketBlastRoutes(r, db)

	InitTicketFeeSettingRoutes(r, db)

	InitPurchaseRoutes(r, db)

	InitRegistrasiUlangRoutes(r, db)

	InitTicketUserRoutes(r, db)

	InitScHabitsRoutes(r, db)

	InitScSubHabitsRoutes(r, db)

	InitScAlphabetProgressRoutes(r, db)

	InitScAlphabetProgressYearlyRoutes(r, db)

	InitScSubjectListSemesterRoutes(r, db)

	InitScSubjectListSemesterYearlyRoutes(r, db)

	InitScGradeRoutes(r, db)

	InitScTemplateHabitsRoutes(r, db)

	InitScAcademicProjectionRoutes(r, db)

	InitScAcademicProjectionDetailRoutes(r, db)

	InitScStudentRoutes(r, db)

	InitScAcademicYearRoutes(r, db)

	InitScLearningCenterRoutes(r, db)

	InitScSubjectTypeGeneralRoutes(r, db)

	InitScTeacherRoutes(r, db)

	InitScSubjectTypeGeneralSubjectRoutes(r, db)

	InitSCUserRoutes(r, db)

	InitTypeUserAplikasiRoutes(r, db)

	// Tambah ini di routes.go Anda:

	InitScPembelianBukuRoutes(r, db)

	InitXenditRoutes(r, db)

	// nanti bisa ditambahkan route lain, misal product, order, dsb
	// InitProductRoutes(r, db)
}
