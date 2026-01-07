package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRegistrasiUlangRoutes(r *gin.Engine, db *gorm.DB) {
	// Service untuk auth (WAJIB)
	userService := service.NewUserService(db)

	// Service untuk business logic
	registrasiUlangService := service.NewRegistrasiUlangService(db)
	registrasiUlangController := controller.NewRegistrasiUlangController(registrasiUlangService)

	registrasiUlangGroup := r.Group("/registrasi-ulang")
	{
		registrasiUlangGroup.Use(
			middleware.AuthMiddleware(userService),     // ✅ FIX: Gunakan auth middleware untuk keamanan
			middleware.RoleMiddleware("admin", "user"), // Asumsi role admin atau user bisa akses; sesuaikan jika perlu
		)

		registrasiUlangGroup.POST("", registrasiUlangController.RegistrasiUlang)
	}

	// Group baru untuk ticket categories (atau bisa digabung jika diinginkan)
	ticketCategoriesGroup := r.Group("/ticket-categories")
	{
		ticketCategoriesGroup.Use(
			middleware.AuthMiddleware(userService),     // Gunakan auth middleware
			middleware.RoleMiddleware("admin", "user"), // Role yang sama, sesuaikan jika perlu (misalnya, hanya admin)
		)

		ticketCategoriesGroup.GET("", registrasiUlangController.GetAllTicketKategori)
	}
}
