// internal/routes/type_user_aplikasi_routes.go

package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTypeUserAplikasiRoutes(r *gin.Engine, db *gorm.DB) {
	userService := service.NewUserService(db)
	typeUserAplikasiService := service.NewTypeUserAplikasiService(db)
	typeUserAplikasiController := controller.NewTypeUserAplikasiController(typeUserAplikasiService)

	typeUserAplikasiGroup := r.Group("/type-user-aplikasi")
	{
		// Auth middleware
		typeUserAplikasiGroup.Use(middleware.AuthMiddleware(userService))

		// Routes
		typeUserAplikasiGroup.POST("", typeUserAplikasiController.CreateTypeUserAplikasi)
		typeUserAplikasiGroup.GET("", typeUserAplikasiController.GetAllTypeUserAplikasi)
		typeUserAplikasiGroup.GET("/user/:user_id", typeUserAplikasiController.GetTypeUserAplikasiByUserID)
		typeUserAplikasiGroup.GET("/:id", typeUserAplikasiController.GetTypeUserAplikasiByID)
		typeUserAplikasiGroup.PUT("/:id", typeUserAplikasiController.UpdateTypeUserAplikasi)
		typeUserAplikasiGroup.DELETE("/:id", typeUserAplikasiController.DeleteTypeUserAplikasi)
	}
}
