package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRekamMedisRoutes(rg *gin.RouterGroup) {
	rekam := rg.Group("/rekam-medis")
	rekam.Use(middleware.AuthMiddleware(), middleware.RequireRole("dokter"))
	{
		rekam.POST("/", controllers.CreateRekamMedis)
		rekam.GET("/:id", middleware.RequireRole("dokter", "admin"), controllers.GetRekamMedisByID)
		rekam.GET("/kunjungan/:id", middleware.RequireRole("dokter", "admin", "farmasi"), controllers.GetRekamMedisByKunjungan)
	}
}
