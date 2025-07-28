package routes

import (
	"KlinikKu/controllers"
	"KlinikKu/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users", middleware.AuthMiddleware(), middleware.InjectUserToContext(), middleware.RequireRole("admin"))

	users.POST("", controllers.CreateUser)
	users.GET("", controllers.GetAllUsers)
	users.GET("/:id", controllers.GetUserByID)
	users.PUT("/:id", controllers.UpdateUser)
	users.DELETE("/:id", controllers.DeleteUser)
}
