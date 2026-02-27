package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScSubjectTypeGeneralRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	subjectTypeGeneralService := service.NewScSubjectTypeGeneralService(db)
	subjectTypeGeneralController := controller.NewScSubjectTypeGeneralController(subjectTypeGeneralService)

	subjectTypeGeneralGroup := r.Group("/sc-subject-type-general")
	{
		// Auth middleware menggunakan userService (bukan subjectTypeGeneralService)
		subjectTypeGeneralGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		subjectTypeGeneralGroup.POST("", subjectTypeGeneralController.CreateScSubjectTypeGeneral)
		subjectTypeGeneralGroup.GET("", subjectTypeGeneralController.GetAllScSubjectTypeGeneral)
		subjectTypeGeneralGroup.GET("/:id", subjectTypeGeneralController.GetScSubjectTypeGeneralByID)
		subjectTypeGeneralGroup.PUT("/:id", subjectTypeGeneralController.UpdateScSubjectTypeGeneral)
		subjectTypeGeneralGroup.DELETE("/:id", subjectTypeGeneralController.DeleteScSubjectTypeGeneral)
	}
}
