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

// BigYearData holds aggregated values for one year /insight/big-number.
type BigYearData struct {
	MaxValue    *float64 `json:"max_value"`
	MaxName     string   `json:"max_name"`
	MinValue    *float64 `json:"min_value"`
	MinName     string   `json:"min_name"`
	AvgValue    *float64 `json:"avg_value"`
	MedianValue *float64 `json:"median_value"`
	IndoValue   *float64 `json:"indo_value"`
	JabarValue  *float64 `json:"jabar_value"`
	KotaValue   *float64 `json:"kota_value"`
}

// BigPrevYearData holds aggregated values of the previous year plus the
// comparison percentages against the current year.
type BigPrevYearData struct {
	AvgValue     *float64 `json:"avg_value"`
	MedianValue  *float64 `json:"median_value"`
	IndoValue    *float64 `json:"indo_value"`
	JabarValue   *float64 `json:"jabar_value"`
	KotaValue    *float64 `json:"kota_value"`
	AvgPersen    *float64 `json:"avg_persen"`
	MedianPersen *float64 `json:"median_persen"`
	IndoPersen   *float64 `json:"indo_persen"`
	JabarPersen  *float64 `json:"jabar_persen"`
	KotaPersen   *float64 `json:"kota_persen"`
}

// BigNumberData is the data payload of /insight/big-number.
type BigNumberData struct {
	TahunSekarangID     string          `json:"tahun_sekarang_id"`
	TahunSekarangName   string          `json:"tahun_sekarang_name"`
	TahunSekarangData   BigYearData     `json:"tahun_sekarang_data"`
	TahunSebelumnyaID   string          `json:"tahun_sebelumnya_id"`
	TahunSebelumnyaName string          `json:"tahun_sebelumnya_name"`
	TahunSebelumnyaData BigPrevYearData `json:"tahun_sebelumnya_data"`
}

// BigNumberMetadata holds the metadata for /insight/big-number.
type BigNumberMetadata struct {
	LastUpdateData     string `json:"last_update_data"`
	LastUpdatePipeline string `json:"last_update_pipeline"`
}

const indonesiaFilter = "LOWER(vervar_name) <> 'indonesia'"

// jabarVervarID is the vervar_id of Jawa Barat (3200) in the dataset.
const jabarVervarID = "3200"

// qryYearData runs the raw aggregation SQL for a given tahun and returns its
// statistical fields. turtahunGroupID is optional: when set, rows are limited
// to the turtahun members of that group (AND-combined with turtahunID).
//
// Which rows feed max/min/avg/median depends on the domain level:
//   - "0000" (national)      : exclude the Indonesia total and the Jawa Barat row
//   - "3200" (Jawa Barat)    : exclude the Jawa Barat (province total) row
//   - otherwise (kab/kota)   : exclude the domain's own region row
//
// The benchmark values (indo_value, jabar_value, kota_value) are always looked
// up separately from those aggregates.
func qryYearData(domainID, varID, turvarID, tahunID, turtahunID, turtahunGroupID string) (BigYearData, error) {
	var d BigYearData
	var maxName, minName sql.NullString
	args := []interface{}{domainID, varID, turvarID, tahunID, turtahunID}
	gf := turtahunGroupFilter(domainID, turtahunGroupID, &args)

	exclusion := ""
	kotaCond := ""
	switch domainID {
	case "0000":
		exclusion = ` AND ` + indonesiaFilter + ` AND vervar_id <> '` + jabarVervarID + `'`
	case "3200":
		exclusion = ` AND vervar_id <> '` + jabarVervarID + `'`
	default:
		args = append(args, domainID)
		p := len(args)
		exclusion = fmt.Sprintf(` AND vervar_id <> $%d`, p)
		kotaCond = fmt.Sprintf(` AND vervar_id = $%d`, p)
	}

	// kota_value is only meaningful for a kabupaten/kota domain; for the other
	// levels it stays NULL.
	kotaExpr := "NULL::numeric"
	if kotaCond != "" {
		kotaExpr = `(SELECT datacontent_value FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` +
			kotaCond + gf + ` LIMIT 1)`
	}

	q := `
		SELECT
		  (SELECT MAX(datacontent_value) FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` + exclusion + gf + `) AS max_value,
		  (SELECT vervar_name FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` + exclusion + gf + `
		     ORDER BY datacontent_value DESC NULLS LAST LIMIT 1) AS max_name,
		  (SELECT MIN(datacontent_value) FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` + exclusion + gf + `) AS min_value,
		  (SELECT vervar_name FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` + exclusion + gf + `
		     ORDER BY datacontent_value ASC NULLS LAST LIMIT 1) AS min_name,
		  (SELECT AVG(datacontent_value) FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` + exclusion + gf + `) AS avg_value,
		  (SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY datacontent_value)
		     FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` + exclusion + gf + `) AS median_value,
		  (SELECT datacontent_value FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5
		       AND LOWER(vervar_name) = 'indonesia'` + gf + ` LIMIT 1) AS indo_value,
		  (SELECT datacontent_value FROM webapi.datacontent
		     WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5
		       AND vervar_id = '` + jabarVervarID + `'` + gf + ` LIMIT 1) AS jabar_value,
		  ` + kotaExpr + ` AS kota_value
	`
	err := DB.DB_SQL_POSTGRES.QueryRow(q, args...).Scan(
		&d.MaxValue, &maxName, &d.MinValue, &minName, &d.AvgValue, &d.MedianValue, &d.IndoValue, &d.JabarValue, &d.KotaValue,
	)
	d.MaxName = maxName.String
	d.MinName = minName.String
	return d, err
}

func pct(cur, prev *float64) *float64 {
	if cur == nil || prev == nil || *prev == 0 {
		return nil
	}
	v := ((*cur - *prev) / *prev) * 100
	return &v
}

// GetBigNumberData godoc
// @Summary      Insight big number
// @Description  Get aggregated big number (current & previous year comparison) from webapi.datacontent
// @Tags         insight
// @Accept       json
// @Produce      json
// @Param        domain_id            query string true  "domain_id" example(0000)
// @Param        var_id               query string true  "var_id" example(286)
// @Param        turvar_id            query string true  "turvar_id" example(530)
// @Param        tahun_id             query string true  "tahun_id" example(125)
// @Param        turtahun_id          query string true  "turtahun_id" example(0)
// @Param        turtahun_group_id    query string false "optional, filter turtahun by group (group_turth_id)" example(0)
// @Success      200  {object}  BigNumberResponse
// @Failure      400  {object}  models.BaseResponse
// @Failure      500  {object}  models.BaseResponse
// @Router       /insight/big-number [get]
func GetBigNumberData(c *gin.Context) {
	domainID := c.Query("domain_id")
	varID := c.Query("var_id")
	turvarID := c.Query("turvar_id")
	tahunID := c.Query("tahun_id")
	turtahunID := c.Query("turtahun_id")
	turtahunGroupID := c.Query("turtahun_group_id")

	if domainID == "" || varID == "" || turvarID == "" || tahunID == "" || turtahunID == "" {
		c.JSON(http.StatusBadRequest, models.BaseResponse{
			Status:  false,
			Message: "Missing required params: domain_id, var_id, turvar_id, tahun_id, turtahun_id",
		})
		return
	}

	// resolve current tahun name
	var tahunSekarangName string
	if err := DB.DB_SQL_POSTGRES.QueryRow(
		`SELECT tahun_name FROM webapi.master_tahun WHERE domain_id=$1 AND tahun_id=$2`,
		domainID, tahunID,
	).Scan(&tahunSekarangName); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}
	data := BigNumberData{
		TahunSekarangID:   tahunID,
		TahunSekarangName: tahunSekarangName,
	}

	// resolve previous year (tahun_name - 1)
	tahunSekarangNum, err := strconv.Atoi(tahunSekarangName)
	var tahunSebelumnyaID, tahunSebelumnyaName string
	if err == nil {
		prevName := strconv.Itoa(tahunSekarangNum - 1)
		if err := DB.DB_SQL_POSTGRES.QueryRow(
			`SELECT tahun_id, tahun_name FROM webapi.master_tahun WHERE domain_id=$1 AND tahun_name=$2`,
			domainID, prevName,
		).Scan(&tahunSebelumnyaID, &tahunSebelumnyaName); err == nil {
			data.TahunSebelumnyaID = tahunSebelumnyaID
			data.TahunSebelumnyaName = tahunSebelumnyaName
		}
	}

	// aggregate for current year
	curData, err := qryYearData(domainID, varID, turvarID, tahunID, turtahunID, turtahunGroupID)
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}
	data.TahunSekarangData = curData

	// aggregate for previous year (if any) and compute percentages
	if data.TahunSebelumnyaID != "" {
		prevData, err := qryYearData(domainID, varID, turvarID, data.TahunSebelumnyaID, turtahunID, turtahunGroupID)
		if err == nil {
			data.TahunSebelumnyaData = BigPrevYearData{
				AvgValue:     prevData.AvgValue,
				MedianValue:  prevData.MedianValue,
				IndoValue:    prevData.IndoValue,
				JabarValue:   prevData.JabarValue,
				KotaValue:    prevData.KotaValue,
				AvgPersen:    pct(curData.AvgValue, prevData.AvgValue),
				MedianPersen: pct(curData.MedianValue, prevData.MedianValue),
				IndoPersen:   pct(curData.IndoValue, prevData.IndoValue),
				JabarPersen:  pct(curData.JabarValue, prevData.JabarValue),
				KotaPersen:   pct(curData.KotaValue, prevData.KotaValue),
			}
		}
	}

	// metadata from the current year rows
	var meta BigNumberMetadata
	metaArgs := []interface{}{domainID, varID, turvarID, tahunID, turtahunID}
	metaQuery := `SELECT
		   COALESCE(to_char(MAX(last_updated_at), 'YYYY-MM-DD HH24:MI:SS'), ''),
		   COALESCE(to_char(MAX(get_at), 'YYYY-MM-DD HH24:MI:SS'), '')
		 FROM webapi.datacontent
		 WHERE domain_id=$1 AND var_id=$2 AND turvar_id=$3 AND tahun_id=$4 AND turtahun_id=$5` +
		turtahunGroupFilter(domainID, turtahunGroupID, &metaArgs)
	if err := DB.DB_SQL_POSTGRES.QueryRow(metaQuery, metaArgs...).Scan(&meta.LastUpdateData, &meta.LastUpdatePipeline); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusInternalServerError, models.BaseResponse{Status: false, Message: "Internal server error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.BaseResponse{
		Status:   true,
		Message:  "Get data success",
		Data:     data,
		Metadata: meta,
	})
}
