//go:debug x509negativeserial=1

// @title Dashboard BPS API
// @version 2.0
// @description REST API WebAPI BPS (master & insight endpoints) - schema webapi.
// @host localhost:8080
// @BasePath /
package main

import (
	"akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	"akatgelar/dashboard-bps-backend/routes"

	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// r.SetTrustedProxies([]string{"127.0.0.1"})

	// DB init
	database.ConnnectDatabasePostgres()

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// hello
	fmt.Println("Service running")
	appPort := os.Getenv("APP_PORT")

	// routes
	root := r.Group("/")
	{
		routes.RouteIndex(root.Group("/"))
		routes.RouteMaster(root.Group("/"))
		routes.RouteInsight(root.Group("/"))
	}

	// route not found -> JSON
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, models.BaseResponse{
			Status:   false,
			Message:  "Route Not Found",
			Data:     nil,
			Metadata: nil,
		})
	})

	// method not allowed -> JSON
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, models.BaseResponse{
			Status:   false,
			Message:  "Method Not Allowed",
			Data:     nil,
			Metadata: nil,
		})
	})

	// run app
	if err := r.Run(":" + appPort); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
