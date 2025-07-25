package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDokterRoutes(rg *gin.RouterGroup) {
	dokter := rg.Group("/dokter")
	dokter.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		dokter.POST("/", controllers.CreateDokter)
		dokter.GET("/", controllers.GetAllDokter)
		dokter.PUT("/:id", controllers.UpdateDokter)
		dokter.DELETE("/:id", controllers.DeleteDokter)

		dokter.POST("/:id/spesialisasi", middleware.RequireRole("admin"), controllers.AssignSpesialisasiToDokter)
		dokter.DELETE("/:id/spesialisasi/:spesialisasi_id", middleware.RequireRole("admin"), controllers.RemoveSpesialisasiFromDokter)
	}
}
