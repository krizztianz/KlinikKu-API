package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPasienRoutes(rg *gin.RouterGroup) {
	adminPasien := rg.Group("/pasien")
	adminPasien.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		adminPasien.PUT("/:id", controllers.UpdatePasien)
		adminPasien.DELETE("/:id", controllers.DeletePasien)
	}

	sharedPasien := rg.Group("/pasien")
	sharedPasien.Use(middleware.AuthMiddleware(), middleware.RequireRole("frontliner", "admin"))
	{
		sharedPasien.GET("/", controllers.GetAllPasien)
		sharedPasien.GET("/search", controllers.SearchPasien)
		sharedPasien.POST("/", controllers.CreatePasien)
	}
}
