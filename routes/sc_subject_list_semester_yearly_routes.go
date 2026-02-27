package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScSubjectListSemesterYearlyRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	subjectService := service.NewScSubjectListSemesterYearlyService(db)
	subjectController := controller.NewScSubjectListSemesterYearlyController(subjectService)

	subjectGroup := r.Group("/sc-subject-list-semester-yearly")
	{
		// Auth middleware menggunakan userService (bukan subjectService)
		subjectGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		subjectGroup.POST("", subjectController.CreateScSubjectListSemesterYearly)
		subjectGroup.GET("", subjectController.GetAllScSubjectListSemesterYearly)
		subjectGroup.GET("/:id", subjectController.GetScSubjectListSemesterYearlyByID)
		subjectGroup.PUT("/:id", subjectController.UpdateScSubjectListSemesterYearly)
		subjectGroup.DELETE("/:id", subjectController.DeleteScSubjectListSemesterYearly)
	}
}