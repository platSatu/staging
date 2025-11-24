package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Panggil route modul user langsung ke root, tanpa prefix /api
	InitUserRoutes(r, db)

	InitAuthRoutes(r, db)

	InitCategoryAplikasiRoutes(r, db)

	InitEventRoutes(r, db)

	InitJenisTiketRoutes(r, db)

	InitTiketRoutes(r, db)

	InitEventUserRoutes(r, db)

	// nanti bisa ditambahkan route lain, misal product, order, dsb
	// InitProductRoutes(r, db)
}
