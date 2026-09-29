package services

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// parseIDList reads a query param that may be supplied as a single value
// (turvar_id=1), a JSON array (turvar_id=[1,2,3]), a comma-separated list
// (turvar_id=1,2,3) or repeated params (?turvar_id=1&turvar_id=2). Empty
// values are dropped and duplicates removed; order is preserved.
func parseIDList(c *gin.Context, name string) []string {
	var out []string
	seen := map[string]bool{}
	add := func(v string) {
		v = strings.TrimSpace(v)
		if v == "" || seen[v] {
			return
		}
		seen[v] = true
		out = append(out, v)
	}

	for _, raw := range c.QueryArray(name) {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		if strings.HasPrefix(raw, "[") {
			var arr []interface{}
			if err := json.Unmarshal([]byte(raw), &arr); err == nil {
				for _, item := range arr {
					add(fmt.Sprint(item))
				}
				continue
			}
		}
		for _, part := range strings.Split(raw, ",") {
			add(part)
		}
	}
	return out
}

// sqlInClause builds `<column> IN ($n, $n+1, ...)` with one positional
// placeholder per value, appending the values to args.
func sqlInClause(column string, values []string, args *[]interface{}) string {
	placeholders := make([]string, len(values))
	for i, v := range values {
		*args = append(*args, v)
		placeholders[i] = fmt.Sprintf("$%d", len(*args))
	}
	return column + " IN (" + strings.Join(placeholders, ",") + ")"
}
