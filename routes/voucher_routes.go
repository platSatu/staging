package routes

import (
	"backend_go/internal/controller"
	"backend_go/internal/service"
	"backend_go/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitVoucherRoutes(r *gin.Engine, db *gorm.DB) {
	voucherService := service.NewVoucherService(db)
	voucherController := controller.NewVoucherController(voucherService)

	voucherGroup := r.Group("/vouchers")
	{
		voucherGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)

		voucherGroup.POST("", voucherController.CreateVoucher)
		voucherGroup.GET("", voucherController.GetAllVouchers)
		voucherGroup.GET("/:id", voucherController.GetVoucherByID)
		voucherGroup.GET("/user", voucherController.GetVouchersByUser) // Endpoint khusus untuk user login
		voucherGroup.PUT("/:id", voucherController.UpdateVoucher)
		voucherGroup.DELETE("/:id", voucherController.DeleteVoucher)
	}
}
