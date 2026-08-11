package routes

import (
	ServiceInsight "akatgelar/dashboard-bps-backend/services/insight"

	"github.com/gin-gonic/gin"
)

func RouteInsight(g *gin.RouterGroup) {

	// insight/big-number
	g.GET("/insight/big-number", func(c *gin.Context) {
		ServiceInsight.GetBigNumberData(c)
	})

	// insight/per-tahun
	g.GET("/insight/per-tahun", func(c *gin.Context) {
		ServiceInsight.GetPerTahunData(c)
	})

	// insight/per-wilayah
	g.GET("/insight/per-wilayah", func(c *gin.Context) {
		ServiceInsight.GetPerWilayahData(c)
	})

	// insight/per-wilayah-turvar
	g.GET("/insight/per-wilayah-turvar", func(c *gin.Context) {
		ServiceInsight.GetPerWilayahTurvarData(c)
	})

}
