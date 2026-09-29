package services

import (
	DB "akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

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
// @Param        turvar_id            query string true  "turvar_id, single or array (e.g. 530 or [530,531])" example(530)
// @Param        tahun_id             query string false "optional, limit to one tahun_id" example(125)
// @Param        vervar_id            query string false "optional, vervar_id single or array (e.g. 1100 or [1100,1200])" example(1100)
// @Param        start_tahun_id       query string false "optional, inclusive lower bound of tahun_id (numeric)" example(110)
// @Param        end_tahun_id         query string false "optional, inclusive upper bound of tahun_id (numeric)" example(125)
// @Success      200  {object}  PerWaktuResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/per-waktu [get]
func GetPerWaktuData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	turvarIDs := parseIDList(c, "turvar_id")
	tahunID := c.Query("tahun_id")
	vervarIDs := parseIDList(c, "vervar_id")
	startTahunID := c.Query("start_tahun_id")
	endTahunID := c.Query("end_tahun_id")

	if domainID == "" || varID == "" || len(turvarIDs) == 0 {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, turvar_id",
		})
		return
	}

	// tahun_id, vervar_id, start_tahun_id & end_tahun_id are optional. The
	// tahun_id bounds must be numeric because tahun_id is compared as an
	// integer (so the range follows chronological order).
	for _, bound := range []struct{ name, value string }{
		{"start_tahun_id", startTahunID},
		{"end_tahun_id", endTahunID},
	} {
		if bound.value != "" {
			if _, err := strconv.Atoi(bound.value); err != nil {
				c.JSON(http.StatusBadRequest, models.BaseResponse{
					Status:  false,
					Message: bound.name + " must be numeric",
				})
				return
			}
		}
	}

	filter := ""
	dataArgs := []interface{}{domainID, varID}
	filter += ` AND ` + sqlInClause("turvar_id", turvarIDs, &dataArgs)
	if tahunID != "" {
		dataArgs = append(dataArgs, tahunID)
		filter += fmt.Sprintf(` AND tahun_id=$%d`, len(dataArgs))
	}
	if len(vervarIDs) > 0 {
		filter += ` AND ` + sqlInClause("vervar_id", vervarIDs, &dataArgs)
	}
	if startTahunID != "" {
		dataArgs = append(dataArgs, startTahunID)
		filter += fmt.Sprintf(` AND CAST(tahun_id AS integer) >= $%d`, len(dataArgs))
	}
	if endTahunID != "" {
		dataArgs = append(dataArgs, endTahunID)
		filter += fmt.Sprintf(` AND CAST(tahun_id AS integer) <= $%d`, len(dataArgs))
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
		 WHERE domain_id=$1 AND var_id=$2` + filter + `
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
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3` + filter
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
