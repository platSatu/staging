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

	// Public route for cron job testing (bisa dipake dari browser/cron)
	r.GET("/vouchers/auto-expire", func(c *gin.Context) {
		affected, err := voucherService.AutoExpireVouchers()
		if err != nil {
			c.JSON(500, gin.H{
				"success": false,
				"message": "Error",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"success":  true,
			"message":  "Vouchers updated to expired",
			"affected": affected,
		})
	})

	voucherGroup := r.Group("/vouchers")
	{
		voucherGroup.Use(
			middleware.AuthMiddleware(service.NewUserService(db)), // Menggunakan UserService untuk auth
			middleware.RoleMiddleware("admin", "user"),
		)
		voucherGroup.POST("/redeem", voucherController.RedeemVoucher)
		voucherGroup.POST("/buy", voucherController.BuyPackage)
		voucherGroup.POST("", voucherController.CreateVoucher)
		voucherGroup.GET("", voucherController.GetAllVouchers)
		voucherGroup.GET("/:id", voucherController.GetVoucherByID)
		voucherGroup.GET("/user", voucherController.GetVouchersByUser) // Endpoint khusus untuk user login
		voucherGroup.PUT("/:id", voucherController.UpdateVoucher)
		voucherGroup.DELETE("/:id", voucherController.DeleteVoucher)
	}
}
