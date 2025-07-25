package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterKunjunganRoutes(rg *gin.RouterGroup) {
	kunjungan := rg.Group("/kunjungan")
	kunjungan.Use(middleware.AuthMiddleware())
	{
		kunjungan.POST("/", middleware.RequireRole("frontliner", "admin"), controllers.CreateKunjungan)
		kunjungan.GET("/", middleware.RequireRole("frontliner", "dokter", "farmasi", "admin"), controllers.GetKunjunganList)
		kunjungan.PUT("/:id/status", middleware.RequireRole("dokter", "farmasi", "admin"), controllers.UpdateKunjunganStatus)
	}
}
