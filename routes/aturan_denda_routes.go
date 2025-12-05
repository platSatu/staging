package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitAturanDendaRoutes(r *gin.Engine, db *gorm.DB) {
	// Buat userService untuk middleware
	userService := service.NewUserService(db)

	aturanService := service.NewAturanDendaService(db)
	aturanController := controller.NewAturanDendaController(aturanService)

	aturanGroup := r.Group("/aturan-denda")
	{
		// Auth middleware menggunakan userService
		aturanGroup.Use(
			middleware.AuthMiddleware(userService),
			middleware.RoleMiddleware("admin", "user"),
		)

		aturanGroup.POST("", aturanController.CreateAturanDenda)
		aturanGroup.GET("", aturanController.GetAllAturanDenda)
		aturanGroup.GET("/:id", aturanController.GetAturanDendaByID)
		aturanGroup.PUT("/:id", aturanController.UpdateAturanDenda)
		aturanGroup.DELETE("/:id", aturanController.DeleteAturanDenda)
	}
}
