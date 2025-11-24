package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitJenisTiketRoutes(r *gin.Engine, db *gorm.DB) {
	jenisTiketService := service.NewJenisTiketService(db)
	jenisTiketController := controller.NewJenisTiketController(jenisTiketService)

	jenisTiketGroup := r.Group("/jenis_tiket")
	{
		jenisTiketGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		jenisTiketGroup.POST("", jenisTiketController.CreateJenisTiket)
		jenisTiketGroup.GET("", jenisTiketController.GetAllJenisTiket)
		jenisTiketGroup.GET("/:id", jenisTiketController.GetJenisTiketByID)
		jenisTiketGroup.GET("/user", jenisTiketController.GetJenisTiketByUser)             // Endpoint khusus untuk user login
		jenisTiketGroup.GET("/event/:event_id", jenisTiketController.GetJenisTiketByEvent) // Endpoint berdasarkan event_id
		jenisTiketGroup.PUT("/:id", jenisTiketController.UpdateJenisTiket)
		jenisTiketGroup.DELETE("/:id", jenisTiketController.DeleteJenisTiket)
	}
}
