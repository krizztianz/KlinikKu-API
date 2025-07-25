package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterSpesialisasiRoutes(rg *gin.RouterGroup) {
	sp := rg.Group("/spesialisasi")
	sp.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		sp.POST("/", controllers.CreateSpesialisasi)
		sp.GET("/", controllers.GetAllSpesialisasi)
		sp.PUT("/:id", controllers.UpdateSpesialisasi)
		sp.DELETE("/:id", controllers.DeleteSpesialisasi)
	}
}
