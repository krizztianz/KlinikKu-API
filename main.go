package main

import (
	"KlinikKu/config"
	"KlinikKu/controllers"
	_ "KlinikKu/docs"
	"KlinikKu/routes"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Gagal load config: %v", err)
	}
	defer cfg.DB.Close() // karena conn dibuka di config.go

	// Jalankan migration
	config.RunMigrations(cfg.DB)

	// Inject DB ke controller
	controllers.InitDB(cfg.DB)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetupRoutes(r)
	log.Printf("Server running on %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}
