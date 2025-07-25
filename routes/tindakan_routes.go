package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterTindakanRoutes(rg *gin.RouterGroup) {
	tindakan := rg.Group("/tindakan")
	tindakan.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		tindakan.POST("/", controllers.CreateTindakan)
		tindakan.GET("/", controllers.GetAllTindakan)
		tindakan.PUT("/:id", controllers.UpdateTindakan)
		tindakan.DELETE("/:id", controllers.DeleteTindakan)
	}
}
