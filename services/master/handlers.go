package services

import (
	"fmt"

	ModelData "akatgelar/dashboard-bps-backend/models/data"

	"github.com/gin-gonic/gin"
)

// GetDomainData godoc
// @Summary      List master domain
// @Description  Get list master domain
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(domain_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  DomainListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/domain [get]
func GetDomainData(c *gin.Context) {
	var dest []ModelData.Domain
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "domain_name", "domain_url", "get_at"},
		filterable:    []string{"id", "domain_id", "domain_name", "domain_url", "is_active"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetSubjectData godoc
// @Summary      List master subject
// @Description  Get list master subject
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(sub_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  SubjectListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/subjek [get]
func GetSubjectData(c *gin.Context) {
	var dest []ModelData.Subject
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "sub_id", "sub_name", "get_at"},
		filterable:    []string{"id", "domain_id", "sub_id", "sub_name"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetVariableData godoc
// @Summary      List master variable
// @Description  Get list master variable
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(var_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  VariableListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/variable [get]
func GetVariableData(c *gin.Context) {
	var dest []ModelData.Variable
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "var_id", "var_name", "sub_id", "sub_name", "unit", "get_at"},
		filterable:    []string{"id", "domain_id", "var_id", "var_name", "sub_id", "sub_name", "unit"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetVariableDistinctData godoc
// @Summary      List master variable that have variable_turunan
// @Description  Get list master variable that have variable_turunan
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(var_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  VariableDistinctListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/variable-distinct [get]
func GetVariableDistinctData(c *gin.Context) {
	var dest []ModelData.Variable
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "var_id", "var_name", "sub_id", "sub_name", "unit", "get_at"},
		filterable:    []string{"id", "domain_id", "var_id", "var_name", "sub_id", "sub_name", "unit"},
		lastUpdateCol: "get_at",
		extraWhere: fmt.Sprintf(
			`var_id IN (SELECT DISTINCT var_id FROM %s)`,
			ModelData.VariableTurunan{}.TableName(),
		),
	}, &dest)
}

// GetVariableTurunanData godoc
// @Summary      List master variable turunan
// @Description  Get list master variable turunan
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(turvar_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  VariableTurunanListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/variable-turunan [get]
func GetVariableTurunanData(c *gin.Context) {
	var dest []ModelData.VariableTurunan
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "var_id", "turvar_id", "turvar_name", "get_at"},
		filterable:    []string{"id", "domain_id", "var_id", "turvar_id", "turvar_name"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetVariableVerticalData godoc
// @Summary      List master variable vertical
// @Description  Get list master variable vertical
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(vervar_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  VariableVerticalListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/variable-vertical [get]
func GetVariableVerticalData(c *gin.Context) {
	var dest []ModelData.VariableVertical
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "var_id", "vervar_id", "vervar_name", "get_at"},
		filterable:    []string{"id", "domain_id", "var_id", "vervar_id", "vervar_name"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetTahunData godoc
// @Summary      List master tahun
// @Description  Get list master tahun
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(tahun_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  TahunListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/tahun [get]
func GetTahunData(c *gin.Context) {
	var dest []ModelData.Tahun
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "tahun_id", "tahun_name", "get_at"},
		filterable:    []string{"id", "domain_id", "tahun_id", "tahun_name", "is_active"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetTahunTurunanData godoc
// @Summary      List master tahun turunan
// @Description  Get list master tahun turunan
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(turtahun_name)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  TahunTurunanListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/tahun-turunan [get]
func GetTahunTurunanData(c *gin.Context) {
	var dest []ModelData.TahunTurunan
	fetchMaster(c, masterConfig{
		defaultSort:   "id",
		sortable:      []string{"id", "domain_id", "turtahun_id", "turtahun_name", "group_turth_id", "name_group_turth", "get_at"},
		filterable:    []string{"id", "domain_id", "turtahun_id", "turtahun_name", "group_turth_id", "name_group_turth"},
		lastUpdateCol: "get_at",
	}, &dest)
}

// GetDataContentData godoc
// @Summary      List datacontent
// @Description  Get list datacontent (payload data)
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(datacontent_value)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  DataContentListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /datacontent [get]
func GetDataContentData(c *gin.Context) {
	var dest []ModelData.DataContent
	fetchMaster(c, masterConfig{
		defaultSort:        "id",
		sortable:           []string{"id", "domain_id", "var_id", "unit", "sub_id", "turvar_id", "tahun_id", "turtahun_id", "vervar_id", "datacontent_id", "datacontent_value", "get_at"},
		filterable:         []string{"id", "domain_id", "var_id", "unit", "sub_id", "turvar_id", "tahun_id", "turtahun_id", "vervar_id", "datacontent_id", "datacontent_value"},
		lastUpdateCol:      "last_updated_at",
		lastUpdateFiltered: true,
	}, &dest)
}

// GetTahunTurunanDistinctData godoc
// @Summary      List distinct tahun turunan from datacontent
// @Description  SELECT DISTINCT turtahun_id, turtahun_name FROM webapi.datacontent
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(turtahun_id)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  TahunTurunanDistinctListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/tahun-turunan-distinct [get]
func GetTahunTurunanDistinctData(c *gin.Context) {
	fetchTahunTurunanDistinct(c, masterConfig{
		defaultSort:        "turtahun_id",
		sortable:           []string{"turtahun_id", "turtahun_name"},
		filterable:         []string{"turtahun_id", "turtahun_name", "domain_id", "var_id"},
		lastUpdateCol:      "last_updated_at",
		lastUpdateFiltered: true,
	})
}

// GetTahunDistinctData godoc
// @Summary      List distinct tahun from datacontent
// @Description  SELECT DISTINCT tahun_id, tahun_name FROM webapi.datacontent
// @Tags         master
// @Accept       json
// @Produce      json
// @Param        filter    query string false "JSON array of {field,operator,value}. Operator supports: eq, neq, gt, gte, lt, lte, like" example([{"field":"domain_id","operator":"eq","value":"0000"}])
// @Param        sort      query string false "sort column (single)" example(tahun_id)
// @Param        order     query string false "asc | desc" Enums(asc,desc) default(asc)
// @Param        per_page  query int    false "per page, default 20" example(20)
// @Param        page      query int    false "page, 1-based" example(1)
// @Success      200  {object}  TahunDistinctListResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /master/tahun-distinct [get]
func GetTahunDistinctData(c *gin.Context) {
	fetchTahunDistinct(c, masterConfig{
		defaultSort:        "tahun_id",
		sortable:           []string{"tahun_id", "tahun_name"},
		filterable:         []string{"tahun_id", "tahun_name", "domain_id", "var_id"},
		lastUpdateCol:      "last_updated_at",
		lastUpdateFiltered: true,
	})
}
