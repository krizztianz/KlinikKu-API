package main

import (
	"KlinikKu/config"
	"KlinikKu/controllers"
	"KlinikKu/routes"
	"log"

	"github.com/gin-gonic/gin"
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
	routes.SetupRoutes(r)
	log.Printf("Server running on %s", cfg.ServerPort)
	r.Run(cfg.ServerPort)
}
