package services

import (
	DB "akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// PerWaktuEntry holds a single period (tahun + turtahun) value inside a vervar group.
type PerWaktuEntry struct {
	TahunID          string   `json:"tahun_id"`
	TahunName        string   `json:"tahun_name"`
	TurtahunID       string   `json:"turtahun_id"`
	TurtahunName     string   `json:"turtahun_name"`
	DataContentID    string   `json:"datacontent_id"`
	DataContentValue *float64 `json:"datacontent_value"`
}

// PerWaktuRow is a vervar group with its list of period values.
type PerWaktuRow struct {
	VervarID   string          `json:"vervar_id"`
	VervarName string          `json:"vervar_name"`
	Data       []PerWaktuEntry `json:"data"`
}

// PerWaktuMetadata holds the metadata for /insight/per-waktu.
type PerWaktuMetadata struct {
	LastUpdateData     string `json:"last_update_data"`
	LastUpdatePipeline string `json:"last_update_pipeline"`
}

// GetPerWaktuData godoc
// @Summary      Insight per waktu
// @Description  Get data per region (vervar) grouped by vervar_id across all periods (tahun + turtahun) from webapi.datacontent
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        domain_id            query string true  "domain_id" example(0000)
// @Param        var_id               query string true  "var_id" example(286)
// @Param        turvar_id            query string true  "turvar_id" example(530)
// @Param        tahun_id             query string false "optional, limit to one tahun_id" example(125)
// @Success      200  {object}  PerWaktuResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/per-waktu [get]
func GetPerWaktuData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	turvarID := c.Query("turvar_id")
	tahunID := c.Query("tahun_id")

	if domainID == "" || varID == "" || turvarID == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, turvar_id",
		})
		return
	}

	// tahun_id is optional: when set, only that year's periods are returned.
	tahunFilter := ""
	dataArgs := []interface{}{domainID, varID, turvarID}
	if tahunID != "" {
		dataArgs = append(dataArgs, tahunID)
		tahunFilter = fmt.Sprintf(` AND tahun_id=$%d`, len(dataArgs))
	}

	dataQuery := `SELECT
		   vervar_id,
		   vervar_name,
		   tahun_id,
		   tahun_name,
		   turtahun_id,
		   turtahun_name,
		   datacontent_id,
		   datacontent_value
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3` + tahunFilter + `
		 ORDER BY vervar_name ASC, CAST(tahun_id AS integer) ASC, CAST(turtahun_id AS integer) ASC`

	rows, err := DB.DB_SQL_POSTGRES.Query(dataQuery, dataArgs...)
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}
	defer rows.Close()

	var groups []PerWaktuRow
	index := map[string]int{}
	for rows.Next() {
		var vervarID, vervarName, tahunIDr, tahunName, turtahunIDr, turtahunName, datacontentID sql.NullString
		var value sql.NullFloat64
		if err := rows.Scan(&vervarID, &vervarName, &tahunIDr, &tahunName, &turtahunIDr, &turtahunName, &datacontentID, &value); err != nil {
			sentry.CaptureException(err)
			c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
			return
		}

		key := vervarID.String
		idx, ok := index[key]
		if !ok {
			groups = append(groups, PerWaktuRow{
				VervarID:   vervarID.String,
				VervarName: vervarName.String,
				Data:       []PerWaktuEntry{},
			})
			idx = len(groups) - 1
			index[key] = idx
		}

		var valuePtr *float64
		if value.Valid {
			v := value.Float64
			valuePtr = &v
		}
		groups[idx].Data = append(groups[idx].Data, PerWaktuEntry{
			TahunID:          tahunIDr.String,
			TahunName:        tahunName.String,
			TurtahunID:       turtahunIDr.String,
			TurtahunName:     turtahunName.String,
			DataContentID:    datacontentID.String,
			DataContentValue: valuePtr,
		})
	}
	if err := rows.Err(); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	var meta PerWaktuMetadata
	metaArgs := dataArgs
	metaQuery := `SELECT
		   COALESCE(to_char(MAX(last_updated_at), 'YYYY-MM-DD HH24:MI:SS'), ''),
		   COALESCE(to_char(MAX(get_at), 'YYYY-MM-DD HH24:MI:SS'), '')
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3` + tahunFilter
	if err := DB.DB_SQL_POSTGRES.QueryRow(metaQuery, metaArgs...).Scan(&meta.LastUpdateData, &meta.LastUpdatePipeline); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Status:   true,
		Message:  "Get data success",
		Data:     groups,
		Metadata: meta,
	})
}
