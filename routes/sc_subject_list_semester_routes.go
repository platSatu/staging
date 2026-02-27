package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScSubjectListSemesterRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	subjectService := service.NewScSubjectListSemesterService(db)
	subjectController := controller.NewScSubjectListSemesterController(subjectService)

	subjectGroup := r.Group("/sc-subject-list-semester")
	{
		// Auth middleware menggunakan userService (bukan subjectService)
		subjectGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		subjectGroup.POST("", subjectController.CreateScSubjectListSemester)
		subjectGroup.GET("", subjectController.GetAllScSubjectListSemester)
		subjectGroup.GET("/:id", subjectController.GetScSubjectListSemesterByID)
		subjectGroup.PUT("/:id", subjectController.UpdateScSubjectListSemester)
		subjectGroup.DELETE("/:id", subjectController.DeleteScSubjectListSemester)
	}
}