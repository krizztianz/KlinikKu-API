package routes

import (
	"KlinikKu/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.LoginHandler)
			auth.POST("/refresh", controllers.RefreshHandler)
			auth.POST("/generate", controllers.GeneratePassword)
		}

		RegisterDokterRoutes(api)
		RegisterPasienRoutes(api)
		RegisterKunjunganRoutes(api)
		RegisterRekamMedisRoutes(api)
		RegisterSpesialisasiRoutes(api)
		RegisterTindakanRoutes(api)
		RegisterDiagnosaRoutes(api)
		RegisterResepRoutes(api)
		RegisterFarmasiRoutes(api)
		RegisterObatRoutes(api)
		RegisterUserRoutes(api)
	}
}
