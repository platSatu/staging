package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitScAcademicYearRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware (konsisten dengan user_routes.go)
	userService := service.NewUserService(db)

	academicYearService := service.NewScAcademicYearService(db)
	academicYearController := controller.NewScAcademicYearController(academicYearService)

	academicYearGroup := r.Group("/sc-academic-year")
	{
		// Auth middleware menggunakan userService (bukan academicYearService)
		academicYearGroup.Use(
			middleware.AuthMiddleware(userService), // Diperbaiki: Gunakan userService seperti pada user
			middleware.RoleMiddleware("admin", "user"),
		)

		academicYearGroup.POST("", academicYearController.CreateScAcademicYear)
		academicYearGroup.GET("", academicYearController.GetAllScAcademicYear)
		academicYearGroup.GET("/:id", academicYearController.GetScAcademicYearByID)
		academicYearGroup.PUT("/:id", academicYearController.UpdateScAcademicYear)
		academicYearGroup.DELETE("/:id", academicYearController.DeleteScAcademicYear)
	}
}
