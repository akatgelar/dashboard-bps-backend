package services

import (
	DB "akatgelar/dashboard-bps-backend/database"
	"akatgelar/dashboard-bps-backend/models"
	"database/sql"
	"fmt"
	"math"
	"net/http"
	"sort"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

const (
	perWilayahMaxCluster = 5
	perWilayahMinCluster = 2
)

// PerWilayahValue is one region row in /insight/per-wilayah.
type PerWilayahValue struct {
	VervarID         string   `json:"vervar_id"`
	VervarName       string   `json:"vervar_name"`
	DataContentID    string   `json:"datacontent_id"`
	DataContentValue *float64 `json:"datacontent_value"`
	Unit             string   `json:"unit"`
	Value            *float64 `json:"value"`
	Color            string   `json:"color"`
}

// PerWilayahRange is one cluster range in /insight/per-wilayah.
type PerWilayahRange struct {
	From         *float64 `json:"from"`
	To           *float64 `json:"to"`
	Color        string   `json:"color"`
	TotalCluster int      `json:"total_cluster"`
}

// PerWilayahData holds the grouped payload of /insight/per-wilayah.
type PerWilayahData struct {
	Value []PerWilayahValue `json:"value"`
	Range []PerWilayahRange `json:"range"`
}

// PerWilayahMetadata holds the metadata for /insight/per-wilayah.
type PerWilayahMetadata struct {
	LastUpdateData     string `json:"last_update_data"`
	LastUpdatePipeline string `json:"last_update_pipeline"`
}

// lerpColor interpolates a color between start and end hex colors by fraction f
// (0..1). start is the low color (#FFFFFF) and end the high color (#4856ff).
func lerpColor(f float64) string {
	const (
		rs, gs, bs = 197, 178, 255 // #c5b2ff
		re, ge, be = 19, 49, 160   // #1331a0
	)
	if f < 0 {
		f = 0
	}
	if f > 1 {
		f = 1
	}
	r := int(math.Round(rs + (re-rs)*f))
	g := int(math.Round(gs + (ge-gs)*f))
	b := int(math.Round(bs + (be-bs)*f))
	return fmt.Sprintf("#%02X%02X%02X", r, g, b)
}

// GetPerWilayahData godoc
// @Summary      Insight per wilayah
// @Description  Get data per region with percentile-based color clusters from webapi.datacontent
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        domain_id    query string true  "domain_id"
// @Param        var_id       query string true  "var_id"
// @Param        turvar_id    query string true  "turvar_id"
// @Param        tahun_id     query string true  "tahun_id"
// @Param        turtahun_id  query string true  "turtahun_id"
// @Success      200  {object}  models.BaseResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/per-wilayah [get]
func GetPerWilayahData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	turvarID := c.Query("turvar_id")
	tahunID := c.Query("tahun_id")
	turtahunID := c.Query("turtahun_id")

	if domainID == "" || varID == "" || turvarID == "" || tahunID == "" || turtahunID == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, turvar_id, tahun_id, turtahun_id",
		})
		return
	}

	rows, err := DB.DB_SQL_POSTGRES.Query(
		`SELECT vervar_id, vervar_name, datacontent_id, datacontent_value, unit
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5
		 ORDER BY vervar_name ASC`,
		domainID, varID, turvarID, tahunID, turtahunID,
	)
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}
	defer rows.Close()

	var regions []PerWilayahValue
	var values []float64
	for rows.Next() {
		var vervarID, vervarName, datacontentID, unit sql.NullString
		var value sql.NullFloat64
		if err := rows.Scan(&vervarID, &vervarName, &datacontentID, &value, &unit); err != nil {
			sentry.CaptureException(err)
			c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
			return
		}
		var vPtr *float64
		if value.Valid {
			v := value.Float64
			vPtr = &v
			values = append(values, v)
		}
		regions = append(regions, PerWilayahValue{
			VervarID:         vervarID.String,
			VervarName:       vervarName.String,
			DataContentID:    datacontentID.String,
			DataContentValue: vPtr,
			Unit:             unit.String,
		})
	}
	if err := rows.Err(); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	// build clusters from percentiles of datacontent_value
	cluster := buildClusters(values)

	// apply value + color to each region
	for i := range regions {
		regions[i].Value = regions[i].DataContentValue
		if regions[i].DataContentValue != nil {
			regions[i].Color = cluster.colorFor(*regions[i].DataContentValue)
		} else {
			regions[i].Color = ""
		}
	}

	// metadata
	var meta PerWilayahMetadata
	if err := DB.DB_SQL_POSTGRES.QueryRow(
		`SELECT
		   to_char(MAX(last_updated_at), 'YYYY-MM-DD HH24:MI:SS'),
		   to_char(MAX(get_at), 'YYYY-MM-DD HH24:MI:SS')
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5`,
		domainID, varID, turvarID, tahunID, turtahunID,
	).Scan(&meta.LastUpdateData, &meta.LastUpdatePipeline); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	data := PerWilayahData{
		Value: regions,
		Range: cluster.ranges,
	}
	c.JSON(http.StatusOK, models.BaseResponse{
		Status:   true,
		Message:  "Get data success",
		Data:     data,
		Metadata: meta,
	})
}

type wilayahCluster struct {
	ranges []PerWilayahRange
}

// buildClusters produces between 2 and 5 percentile-based ranges from the given
// values, partitioning the sorted values into equal-sized quantile clusters so
// every value maps to exactly one cluster.
func buildClusters(values []float64) *wilayahCluster {
	c := &wilayahCluster{}
	if len(values) == 0 {
		return c
	}
	sorted := append([]float64(nil), values...)
	sort.Float64s(sorted)

	// number of clusters: based on distinct values, clamped to [2,5]
	unique := 1
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1] {
			unique++
		}
	}
	n := unique
	if n < perWilayahMinCluster {
		n = perWilayahMinCluster
	}
	if n > perWilayahMaxCluster {
		n = perWilayahMaxCluster
	}

	ranges := make([]PerWilayahRange, n)
	// each cluster takes an equal quantile share (rank-based partition)
	start := 0
	for k := 0; k < n; k++ {
		end := int(float64(len(sorted)) * float64(k+1) / float64(n))
		if end > len(sorted) {
			end = len(sorted)
		}
		if end <= start {
			end = start + 1
		}
		if end > len(sorted) {
			end = len(sorted)
		}

		frac := 0.0
		if n > 1 {
			frac = float64(k) / float64(n-1)
		}
		from, to := sorted[start], sorted[end-1]
		color := lerpColor(frac)
		ranges[k] = PerWilayahRange{
			From:         &from,
			To:           &to,
			Color:        color,
			TotalCluster: end - start,
		}
		start = end
	}
	c.ranges = ranges
	return c
}

// colorFor returns the cluster color for a given value.
func (c *wilayahCluster) colorFor(v float64) string {
	for i := range c.ranges {
		r := &c.ranges[i]
		if r.From != nil && r.To != nil && v >= *r.From && v <= *r.To {
			return r.Color
		}
	}
	return ""
}
