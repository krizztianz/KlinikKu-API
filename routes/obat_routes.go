package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterObatRoutes(rg *gin.RouterGroup) {
	obat := rg.Group("/obat")
	obat.Use(middleware.AuthMiddleware(), middleware.InjectUserToContext())
	{
		obat.GET("/", controllers.GetAllObat)
		obat.GET("/:id", controllers.GetObatByID)
		obat.POST("/", middleware.RequireRole("admin"), controllers.CreateObat)
		obat.PUT("/:id", middleware.RequireRole("admin"), controllers.UpdateObat)
		obat.DELETE("/:id", middleware.RequireRole("admin"), controllers.DeleteObat)
	}
}
