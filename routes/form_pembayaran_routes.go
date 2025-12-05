package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitFormPembayaranRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware
	userService := service.NewUserService(db)

	formService := service.NewFormPembayaranService(db)
	formController := controller.NewFormPembayaranController(formService)

	formGroup := r.Group("/form-pembayaran")
	{
		// Auth middleware menggunakan userService
		formGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)

		formGroup.POST("", formController.CreateFormPembayaran)
		formGroup.GET("", formController.GetAllFormPembayaran)
		formGroup.GET("/:id", formController.GetFormPembayaranByID)
		formGroup.PUT("/:id", formController.UpdateFormPembayaran)
		formGroup.DELETE("/:id", formController.DeleteFormPembayaran)
	}
}
