package services

import (
	DB "akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	"database/sql"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// PerTahunYearEntry holds a single year value inside a vervar group.
type PerTahunYearEntry struct {
	TahunID          string   `json:"tahun_id"`
	TahunName        string   `json:"tahun_name"`
	DataContentID    string   `json:"datacontent_id"`
	DataContentValue *float64 `json:"datacontent_value"`
}

// PerTahunRow is a vervar group with its list of year values.
type PerTahunRow struct {
	VervarID   string              `json:"vervar_id"`
	VervarName string              `json:"vervar_name"`
	Data       []PerTahunYearEntry `json:"data"`
}

// PerTahunMetadata holds the metadata for /insight/per-tahun.
type PerTahunMetadata struct {
	LastUpdateData     string `json:"last_update_data"`
	LastUpdatePipeline string `json:"last_update_pipeline"`
}

// GetPerTahunData godoc
// @Summary      Insight per tahun
// @Description  Get data per region (vervar) & per year from webapi.datacontent
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        domain_id            query string true  "domain_id" example(0000)
// @Param        var_id               query string true  "var_id" example(286)
// @Param        turvar_id            query string true  "turvar_id" example(530)
// @Param        turtahun_id          query string true  "turtahun_id" example(0)
// @Param        turtahun_group_id    query string false "optional, filter turtahun by group (group_turth_id)" example(1)
// @Success      200  {object}  PerTahunResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/per-tahun [get]
func GetPerTahunData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	turvarID := c.Query("turvar_id")
	turtahunID := c.Query("turtahun_id")
	turtahunGroupID := c.Query("turtahun_group_id")

	if domainID == "" || varID == "" || turvarID == "" || turtahunID == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, turvar_id, turtahun_id",
		})
		return
	}

	dataArgs := []interface{}{domainID, varID, turvarID, turtahunID}
	dataQuery := `SELECT
		   vervar_id,
		   vervar_name,
		   tahun_id,
		   tahun_name,
		   datacontent_id,
		   datacontent_value
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND turtahun_id=$4` +
		turtahunGroupFilter(domainID, turtahunGroupID, &dataArgs) +
		` ORDER BY vervar_name ASC, tahun_id ASC`

	rows, err := DB.DB_SQL_POSTGRES.Query(dataQuery, dataArgs...)
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}
	defer rows.Close()

	var groups []PerTahunRow
	index := map[string]int{}
	for rows.Next() {
		var vervarID, vervarName, tahunIDr, tahunName, datacontentID sql.NullString
		var value sql.NullFloat64
		if err := rows.Scan(&vervarID, &vervarName, &tahunIDr, &tahunName, &datacontentID, &value); err != nil {
			sentry.CaptureException(err)
			c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
			return
		}

		key := vervarID.String
		idx, ok := index[key]
		if !ok {
			groups = append(groups, PerTahunRow{
				VervarID:   vervarID.String,
				VervarName: vervarName.String,
				Data:       []PerTahunYearEntry{},
			})
			idx = len(groups) - 1
			index[key] = idx
		}

		var valuePtr *float64
		if value.Valid {
			v := value.Float64
			valuePtr = &v
		}
		groups[idx].Data = append(groups[idx].Data, PerTahunYearEntry{
			TahunID:          tahunIDr.String,
			TahunName:        tahunName.String,
			DataContentID:    datacontentID.String,
			DataContentValue: valuePtr,
		})
	}
	if err := rows.Err(); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	var meta PerTahunMetadata
	metaArgs := []interface{}{domainID, varID, turvarID, turtahunID}
	metaQuery := `SELECT
		   COALESCE(to_char(MAX(last_updated_at), 'YYYY-MM-DD HH24:MI:SS'), ''),
		   COALESCE(to_char(MAX(get_at), 'YYYY-MM-DD HH24:MI:SS'), '')
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND turtahun_id=$4` +
		turtahunGroupFilter(domainID, turtahunGroupID, &metaArgs)
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
