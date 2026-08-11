package routes

import (
	"akatgelar/dashboard-bps-backend/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	ServiceSwagger "akatgelar/dashboard-bps-backend/services/swagger"
)

func RouteIndex(g *gin.RouterGroup) {

	appPort := os.Getenv("APP_PORT")
	ginMode := os.Getenv("GIN_MODE")
	host := ""
	if ginMode == "debug" {
		host = "localhost:" + appPort
	} else if ginMode == "release" {
		host = "https://dashboard-bps.akatgelar.app/api/"
	}
	fmt.Print("host")
	fmt.Print(host)

	// route "/" -> swagger
	g.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, host + "/swagger/index.html")
	})

	// docs swagger served under /swagger/*
	ServiceSwagger.SwaggerInfo.Title = "Dashboard BPS API"
	ServiceSwagger.SwaggerInfo.Description = "This is a documentation for Dashboard BPS API."
	ServiceSwagger.SwaggerInfo.Version = "2.0"
	ServiceSwagger.SwaggerInfo.Host = host
	ServiceSwagger.SwaggerInfo.BasePath = "/"
	ServiceSwagger.SwaggerInfo.Schemes = []string{"http", "https"}
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// root
	g.GET("/is_alive", func(c *gin.Context) {
		c.JSON(
			200,
			models.BaseResponse{
				Status:  true,
				Message: "Welcome to Dashboard BPS API",
				Data:    gin.H{},
			},
		)
	})
}
