package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFarmasiRoutes(rg *gin.RouterGroup) {
	farmasi := rg.Group("/farmasi")
	farmasi.Use(middleware.AuthMiddleware(), middleware.InjectUserToContext(), middleware.RequireRole("farmasi", "admin"))
	{
		farmasi.GET("/resep", controllers.GetPendingResep)                  // List resep belum ditebus
		farmasi.GET("/resep/:id", controllers.GetResepDetail)               // Detail resep
		farmasi.PATCH("/resep/:id/tebus", controllers.MarkResepAsCompleted) // Tandai sudah ditebus
	}
}
