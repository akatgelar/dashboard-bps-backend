package services

import (
	DB "akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	"database/sql"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

// PerWilayahTurvarEntry holds a single variable-turunan value inside a region.
type PerWilayahTurvarEntry struct {
	TurvarID         string   `json:"turvar_id"`
	TurvarName       string   `json:"turvar_name"`
	DataContentID    string   `json:"datacontent_id"`
	DataContentValue *float64 `json:"datacontent_value"`
}

// PerWilayahTurvarValue is a region (vervar) group with its turvar values.
type PerWilayahTurvarValue struct {
	VervarID   string                  `json:"vervar_id"`
	VervarName string                  `json:"vervar_name"`
	Data       []PerWilayahTurvarEntry `json:"data"`
}

// PerWilayahTurvarMetadata holds the metadata for /insight/per-wilayah-turvar.
type PerWilayahTurvarMetadata struct {
	LastUpdateData     string `json:"last_update_data"`
	LastUpdatePipeline string `json:"last_update_pipeline"`
}

// GetPerWilayahTurvarData godoc
// @Summary      Insight per wilayah - turvar
// @Description  Get data per region & per variable turunan from webapi.datacontent
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        domain_id    query string true  "domain_id"
// @Param        var_id       query string true  "var_id"
// @Param        tahun_id     query string true  "tahun_id"
// @Param        turtahun_id  query string true  "turtahun_id"
// @Success      200  {object}  models.BaseResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/per-wilayah-turvar [get]
func GetPerWilayahTurvarData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	tahunID := c.Query("tahun_id")
	turtahunID := c.Query("turtahun_id")

	if domainID == "" || varID == "" || tahunID == "" || turtahunID == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, tahun_id, turtahun_id",
		})
		return
	}

	rows, err := DB.DB_SQL_POSTGRES.Query(
		`SELECT vervar_id, vervar_name, turvar_id, turvar_name, datacontent_id, datacontent_value
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND tahun_id=$3 AND turtahun_id=$4
		 ORDER BY vervar_name ASC, turvar_name ASC`,
		domainID, varID, tahunID, turtahunID,
	)
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}
	defer rows.Close()

	var groups []PerWilayahTurvarValue
	index := map[string]int{}
	for rows.Next() {
		var vervarID, vervarName, turvarID, turvarName, datacontentID sql.NullString
		var value sql.NullFloat64
		if err := rows.Scan(&vervarID, &vervarName, &turvarID, &turvarName, &datacontentID, &value); err != nil {
			sentry.CaptureException(err)
			c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
			return
		}

		key := vervarID.String
		idx, ok := index[key]
		if !ok {
			groups = append(groups, PerWilayahTurvarValue{
				VervarID:   vervarID.String,
				VervarName: vervarName.String,
				Data:       []PerWilayahTurvarEntry{},
			})
			idx = len(groups) - 1
			index[key] = idx
		}

		var valuePtr *float64
		if value.Valid {
			v := value.Float64
			valuePtr = &v
		}
		groups[idx].Data = append(groups[idx].Data, PerWilayahTurvarEntry{
			TurvarID:         turvarID.String,
			TurvarName:       turvarName.String,
			DataContentID:    datacontentID.String,
			DataContentValue: valuePtr,
		})
	}
	if err := rows.Err(); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	var meta PerWilayahTurvarMetadata
	if err := DB.DB_SQL_POSTGRES.QueryRow(
		`SELECT
		   to_char(MAX(last_updated_at), 'YYYY-MM-DD HH24:MI:SS'),
		   to_char(MAX(get_at), 'YYYY-MM-DD HH24:MI:SS')
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND tahun_id=$3 AND turtahun_id=$4`,
		domainID, varID, tahunID, turtahunID,
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
