package main

import (
	"backend_go/config"
	"backend_go/helper"
	"backend_go/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	r := gin.Default()

	r.Use(helper.CorsMiddleware())

	db := config.InitDB()

	routes.SetupRoutes(r, db)

	r.Run(":7070")
}
