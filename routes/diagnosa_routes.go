package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDiagnosaRoutes(rg *gin.RouterGroup) {
	diagnosa := rg.Group("/diagnosa")
	diagnosa.Use(middleware.AuthMiddleware(), middleware.InjectUserToContext(), middleware.RequireRole("admin"))
	{
		diagnosa.POST("/", controllers.CreateDiagnosa)
		diagnosa.GET("/", controllers.GetAllDiagnosa)
		diagnosa.PUT("/:id", controllers.UpdateDiagnosa)
		diagnosa.DELETE("/:id", controllers.DeleteDiagnosa)
	}
}
