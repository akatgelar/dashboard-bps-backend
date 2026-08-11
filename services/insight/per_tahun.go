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
// @Param        domain_id    query string true  "domain_id"
// @Param        var_id       query string true  "var_id"
// @Param        turvar_id    query string true  "turvar_id"
// @Param        turtahun_id  query string true  "turtahun_id"
// @Success      200  {object}  models.BaseResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/per-tahun [get]
func GetPerTahunData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	turvarID := c.Query("turvar_id")
	turtahunID := c.Query("turtahun_id")

	if domainID == "" || varID == "" || turvarID == "" || turtahunID == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, turvar_id, turtahun_id",
		})
		return
	}

	rows, err := DB.DB_SQL_POSTGRES.Query(
		`SELECT
		   vervar_id,
		   vervar_name,
		   tahun_id,
		   tahun_name,
		   datacontent_id,
		   datacontent_value
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND turtahun_id=$4
		 ORDER BY vervar_name ASC, tahun_id ASC`,
		domainID, varID, turvarID, turtahunID,
	)
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
	if err := DB.DB_SQL_POSTGRES.QueryRow(
		`SELECT
		   to_char(MAX(last_updated_at), 'YYYY-MM-DD HH24:MI:SS'),
		   to_char(MAX(get_at), 'YYYY-MM-DD HH24:MI:SS')
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND turtahun_id=$4`,
		domainID, varID, turvarID, turtahunID,
	).Scan(&meta.LastUpdateData, &meta.LastUpdatePipeline); err != nil {
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
