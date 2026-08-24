package services

import (
	"akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	ModelData "akatgelar/dashboard-bps-backend/models/data"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const masterDefaultPerPage = 20

// masterConfig defines how a master endpoint should query its source table.
type masterConfig struct {
	defaultSort        string   // default sort column
	sortable           []string // columns allowed in sort
	filterable         []string // columns allowed in filter
	lastUpdateCol      string   // column used for metadata.last_update_data
	lastUpdateFiltered bool     // whether the filter applies to the last_update computation
	extraWhere         string   // additional raw WHERE clause applied to data & count queries
}

func (cfg masterConfig) sortableSet() map[string]bool {
	set := make(map[string]bool, len(cfg.sortable))
	for _, col := range cfg.sortable {
		set[col] = true
	}
	return set
}

func (cfg masterConfig) filterableSet() map[string]bool {
	set := make(map[string]bool, len(cfg.filterable))
	for _, col := range cfg.filterable {
		set[col] = true
	}
	return set
}

func queryInt(c *gin.Context, name string, def int) int {
	raw := c.DefaultQuery(name, strconv.Itoa(def))
	v, err := strconv.Atoi(raw)
	if err != nil || v < 1 {
		return def
	}
	return v
}

// lastPipelineValue returns the max(get_at) across the entire source table for
// the given model, used to populate metadata.last_update_pipeline.
func lastPipelineValue(model interface{}) time.Time {
	return maxColumnValue(database.DB_POSTGRES.Model(model), `"get_at"`)
}

// lastDataValue returns the max of cfg.lastUpdateCol. For datacontent
// (lastUpdateFiltered) the request filter is applied to the aggregation.
func lastDataValue(c *gin.Context, model interface{}, cfg masterConfig) time.Time {
	q := database.DB_POSTGRES.Model(model)
	if cfg.lastUpdateFiltered {
		if !applyFilter(q, c, cfg) {
			return time.Time{}
		}
	}
	return maxColumnValue(q, cfg.lastUpdateCol)
}

func maxColumnValue(q *gorm.DB, column string) time.Time {
	if column == "" {
		return time.Time{}
	}
	var val *time.Time
	if err := q.Select(fmt.Sprintf(`COALESCE(MAX(%s), NULL)`, column)).Scan(&val).Error; err != nil || val == nil {
		return time.Time{}
	}
	return *val
}

// formatDateTime renders a time as "yyyy-MM-dd HH:mm:ss" (no T/Z). Returns an
// empty string for a zero value.
func formatDateTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// applyFilter parses the `filter` param (JSON array of {field, operator, value})
// and applies each operator (eq, neq, gt, gte, lt, lte, like) to the query.
// It returns false (after writing an error response) when parsing fails.
func applyFilter(q *gorm.DB, c *gin.Context, cfg masterConfig) bool {
	filterJSON := c.DefaultQuery("filter", "[]")
	if strings.TrimSpace(filterJSON) == "" {
		return true
	}

	var filters []models.BaseFilter
	if err := json.Unmarshal([]byte(filterJSON), &filters); err != nil {
		c.JSON(http.StatusBadRequest, models.BaseResponse{Status: false, Message: "Invalid filter format"})
		return false
	}

	filterable := cfg.filterableSet()
	for _, filter := range filters {
		field := strings.TrimSpace(filter.Field)
		if field == "" || !filterable[field] {
			continue
		}

		operator := strings.ToLower(filter.Operator)
		switch operator {
		case "eq":
			q.Where(fmt.Sprintf(`"%s" = ?`, field), filter.Value)
		case "neq":
			q.Where(fmt.Sprintf(`"%s" != ?`, field), filter.Value)
		case "gt":
			q.Where(fmt.Sprintf(`"%s" > ?`, field), filter.Value)
		case "gte":
			q.Where(fmt.Sprintf(`"%s" >= ?`, field), filter.Value)
		case "lt":
			q.Where(fmt.Sprintf(`"%s" < ?`, field), filter.Value)
		case "lte":
			q.Where(fmt.Sprintf(`"%s" <= ?`, field), filter.Value)
		case "like":
			if val, ok := filter.Value.(string); ok {
				q.Where(fmt.Sprintf(`"%s" ILIKE ?`, field), "%"+val+"%")
			}
		}
	}
	return true
}

func buildOrder(c *gin.Context, cfg masterConfig) (bool, string) {
	sortable := cfg.sortableSet()

	sortCol := strings.TrimSpace(c.DefaultQuery("sort", ""))
	dir := strings.ToUpper(strings.TrimSpace(c.DefaultQuery("order", "asc")))
	if dir != "ASC" && dir != "DESC" {
		dir = "ASC"
	}

	if sortCol == "" || !sortable[sortCol] {
		if cfg.defaultSort == "" {
			return false, ""
		}
		return true, fmt.Sprintf(`"%s" %s`, cfg.defaultSort, dir)
	}

	return true, fmt.Sprintf(`"%s" %s`, sortCol, dir)
}

// fetchMaster runs the shared read-only master query and writes the response.
// model must be a pointer to a slice (e.g. *[]Domain).
func fetchMaster(c *gin.Context, cfg masterConfig, model interface{}) {
	perPage := queryInt(c, "per_page", masterDefaultPerPage)
	page := queryInt(c, "page", 1)
	offset := (page - 1) * perPage

	q := database.DB_POSTGRES.Model(model)
	if !applyFilter(q, c, cfg) {
		return
	}
	if cfg.extraWhere != "" {
		q = q.Where(cfg.extraWhere)
	}

	var totalData int64
	if err := q.Count(&totalData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Failed to count records: " + err.Error()})
		return
	}

	totalPage := int(totalData) / perPage
	if int(totalData)%perPage > 0 {
		totalPage++
	}

	if hasOrder, orderClause := buildOrder(c, cfg); hasOrder {
		q = q.Order(orderClause)
	}

	if err := q.Limit(perPage).Offset(offset).Find(model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Failed to fetch data: " + err.Error()})
		return
	}

	// Special case: VariableTurunan with empty result returns default value
	rv := reflect.ValueOf(model)
	if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Slice {
		sliceVal := rv.Elem()
		if sliceVal.Len() == 0 && sliceVal.Type().Elem() == reflect.TypeOf(ModelData.VariableTurunan{}) {
			defaultData := []ModelData.VariableTurunan{
				{TurvarID: "0", TurvarName: "Tidak Ada"},
			}
			c.JSON(http.StatusOK, models.BaseResponse{
				Status:  true,
				Message: "Get data success",
				Data:     defaultData,
				Metadata: models.BaseMetadata{
					TotalData:      1,
					TotalPage:      1,
					PerPage:        perPage,
					Page:           page,
					LastUpdateData:     formatDateTime(lastDataValue(c, model, cfg)),
					LastUpdatePipeline: formatDateTime(lastPipelineValue(model)),
				},
			})
			return
		}
	}

	// metadata timestamp
	lastPipeline := lastPipelineValue(model)
	lastData := lastDataValue(c, model, cfg)

	metadata := models.BaseMetadata{
		TotalData:          int(totalData),
		TotalPage:          totalPage,
		PerPage:            perPage,
		Page:               page,
		LastUpdateData:     formatDateTime(lastData),
		LastUpdatePipeline: formatDateTime(lastPipeline),
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Status:   true,
		Message:  "Get data success",
		Data:     model,
		Metadata: metadata,
	})
}

// fetchDistinct serves a master-distinct endpoint: SELECT DISTINCT <props>
// FROM webapi.datacontent, with the same params (filter, sort, per_page, page).
func fetchDistinct(c *gin.Context, cfg masterConfig, props string, dest interface{}) {
	sourceTable := ModelData.DataContent{}.TableName()

	perPage := queryInt(c, "per_page", masterDefaultPerPage)
	page := queryInt(c, "page", 1)
	offset := (page - 1) * perPage

	// data: SELECT DISTINCT <props> FROM webapi.datacontent
	query := database.DB_POSTGRES.Table(sourceTable).Select("DISTINCT " + props)
	if !applyFilter(query, c, cfg) {
		return
	}
	if hasOrder, orderClause := buildOrder(c, cfg); hasOrder {
		query = query.Order(orderClause)
	}

	if err := query.Limit(perPage).Offset(offset).Scan(dest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Failed to fetch data: " + err.Error()})
		return
	}

	// count: SELECT COUNT(*) FROM (SELECT DISTINCT ... ) sub
	sub := database.DB_POSTGRES.Table(sourceTable).Select("DISTINCT " + props)
	if !applyFilter(sub, c, cfg) {
		return
	}
	var totalData int64
	if err := database.DB_POSTGRES.Table("(?) AS sub", sub).Count(&totalData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Failed to count records: " + err.Error()})
		return
	}

	totalPage := int(totalData) / perPage
	if int(totalData)%perPage > 0 {
		totalPage++
	}

	lastPipeline := lastPipelineValue(&ModelData.DataContent{})
	lastData := lastDataValue(c, &ModelData.DataContent{}, cfg)

	metadata := models.BaseMetadata{
		TotalData:          int(totalData),
		TotalPage:          totalPage,
		PerPage:            perPage,
		Page:               page,
		LastUpdateData:     formatDateTime(lastData),
		LastUpdatePipeline: formatDateTime(lastPipeline),
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Status:   true,
		Message:  "Get data success",
		Data:     dest,
		Metadata: metadata,
	})
}

// fetchTahunTurunanDistinct serves /master/tahun-turunan-distinct:
// SELECT DISTINCT turtahun_id, turtahun_name FROM webapi.datacontent.
func fetchTahunTurunanDistinct(c *gin.Context, cfg masterConfig) {
	var results []ModelData.TahunTurunanDistinct
	fetchDistinct(c, cfg, `"turtahun_id", "turtahun_name"`, &results)
}

// fetchTahunDistinct serves /master/tahun-distinct:
// SELECT DISTINCT tahun_id, tahun_name FROM webapi.datacontent.
func fetchTahunDistinct(c *gin.Context, cfg masterConfig) {
	var results []ModelData.TahunDistinct
	fetchDistinct(c, cfg, `"tahun_id", "tahun_name"`, &results)
}
