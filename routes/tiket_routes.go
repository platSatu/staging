package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitTiketRoutes(r *gin.Engine, db *gorm.DB) {
	tiketService := service.NewTiketService(db)
	tiketController := controller.NewTiketController(tiketService)

	tiketGroup := r.Group("/tiket")
	{
		tiketGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		tiketGroup.POST("", tiketController.CreateTiket)
		tiketGroup.GET("", tiketController.GetAllTiket)
		tiketGroup.GET("/:id", tiketController.GetTiketByID)
		tiketGroup.GET("/user", tiketController.GetTiketByUser)                         // Endpoint khusus untuk user login
		tiketGroup.GET("/event/:event_id", tiketController.GetTiketByEvent)             // Endpoint berdasarkan event_id
		tiketGroup.GET("/booking/:kode_booking", tiketController.GetTiketByKodeBooking) // Endpoint berdasarkan kode booking
		tiketGroup.PUT("/:id", tiketController.UpdateTiket)
		tiketGroup.DELETE("/:id", tiketController.DeleteTiket)
	}
}
