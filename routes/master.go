package routes

import (
	ServiceMaster "akatgelar/dashboard-bps-backend/services/master"

	"github.com/gin-gonic/gin"
)

func RouteMaster(g *gin.RouterGroup) {

	// master/domain
	g.GET("/master/domain", func(c *gin.Context) {
		ServiceMaster.GetDomainData(c)
	})

	// master/subjek
	g.GET("/master/subjek", func(c *gin.Context) {
		ServiceMaster.GetSubjectData(c)
	})

	// master/variable
	g.GET("/master/variable", func(c *gin.Context) {
		ServiceMaster.GetVariableData(c)
	})

	// master/variable-distinct
	g.GET("/master/variable-distinct", func(c *gin.Context) {
		ServiceMaster.GetVariableDistinctData(c)
	})

	// master/variable_turunan
	g.GET("/master/variable-turunan", func(c *gin.Context) {
		ServiceMaster.GetVariableTurunanData(c)
	})

	// master/variable_vertical
	g.GET("/master/variable-vertical", func(c *gin.Context) {
		ServiceMaster.GetVariableVerticalData(c)
	})

	// master/tahun
	g.GET("/master/tahun", func(c *gin.Context) {
		ServiceMaster.GetTahunData(c)
	})

	// master/tahun_turunan
	g.GET("/master/tahun-turunan", func(c *gin.Context) {
		ServiceMaster.GetTahunTurunanData(c)
	})

	// master/tahun-turunan-distinct
	g.GET("/master/tahun-turunan-distinct", func(c *gin.Context) {
		ServiceMaster.GetTahunTurunanDistinctData(c)
	})

	// master/tahun-turunan-group
	g.GET("/master/tahun-turunan-group", func(c *gin.Context) {
		ServiceMaster.GetTahunTurunanGroupData(c)
	})

	// master/tahun-turunan-group-distinct
	g.GET("/master/tahun-turunan-group-distinct", func(c *gin.Context) {
		ServiceMaster.GetTahunTurunanGroupDistinctData(c)
	})

	// master/tahun-distinct
	g.GET("/master/tahun-distinct", func(c *gin.Context) {
		ServiceMaster.GetTahunDistinctData(c)
	})

	// datacontent
	g.GET("/datacontent", func(c *gin.Context) {
		ServiceMaster.GetDataContentData(c)
	})

}
