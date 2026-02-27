package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScSubjectTypeGeneralSubjectRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	subjectTypeGeneralSubjectService := service.NewScSubjectTypeGeneralSubjectService(db)
	subjectTypeGeneralSubjectController := controller.NewScSubjectTypeGeneralSubjectController(subjectTypeGeneralSubjectService)

	subjectTypeGeneralSubjectGroup := r.Group("/sc-subject-type-general-subject")
	{
		// Auth middleware menggunakan userService (bukan subjectTypeGeneralSubjectService)
		subjectTypeGeneralSubjectGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		subjectTypeGeneralSubjectGroup.POST("", subjectTypeGeneralSubjectController.CreateScSubjectTypeGeneralSubject)
		subjectTypeGeneralSubjectGroup.GET("", subjectTypeGeneralSubjectController.GetAllScSubjectTypeGeneralSubject)
		subjectTypeGeneralSubjectGroup.GET("/:id", subjectTypeGeneralSubjectController.GetScSubjectTypeGeneralSubjectByID)
		subjectTypeGeneralSubjectGroup.PUT("/:id", subjectTypeGeneralSubjectController.UpdateScSubjectTypeGeneralSubject)
		subjectTypeGeneralSubjectGroup.DELETE("/:id", subjectTypeGeneralSubjectController.DeleteScSubjectTypeGeneralSubject)
	}
}
