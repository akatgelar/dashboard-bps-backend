package routes

import (
	"akatgelar/dashboard-bps-backend/models"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	ServiceSwagger "akatgelar/dashboard-bps-backend/services/swagger"
)

func RouteIndex(g *gin.RouterGroup) {

	// root
	g.GET("/", func(c *gin.Context) {
		c.JSON(
			200,
			models.BaseResponse{
				Status:  true,
				Message: "Welcome to Dashboard BPS API",
				Data: gin.H{
					"swagger": "/docs/index.html",
				},
			},
		)
	})

	// docs swagger
	appPort := os.Getenv("APP_PORT")
	ServiceSwagger.SwaggerInfo.Title = "Swagger Dashboard BPS API"
	ServiceSwagger.SwaggerInfo.Description = "This is a documentation for Dashboard BPS API."
	ServiceSwagger.SwaggerInfo.Version = "1.0"
	ServiceSwagger.SwaggerInfo.Host = "localhost:" + appPort
	ServiceSwagger.SwaggerInfo.BasePath = "/"
	ServiceSwagger.SwaggerInfo.Schemes = []string{"http", "https"}
	g.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}
