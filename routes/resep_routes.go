package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterResepRoutes(rg *gin.RouterGroup) {
	resep := rg.Group("/resep")
	resep.Use(middleware.AuthMiddleware(), middleware.RequireRole("dokter"))
	{
		resep.POST("/", controllers.CreateResep)
		resep.GET("/:id", middleware.RequireRole("dokter", "admin", "farmasi"), controllers.GetResepByID)
		resep.GET("/by-rekam-medis/:id", middleware.RequireRole("dokter", "admin", "farmasi"), controllers.GetResepByRekamMedisID)
	}
}
