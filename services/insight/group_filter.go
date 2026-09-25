package services

import "fmt"

// turtahunGroupFilter builds an extra SQL condition that restricts `turtahun_id`
// to the members of `turtahun_group_id` (from webapi.master_tahun_turunan) for
// the given domain. It returns "" when turtahunGroupID is empty, so the param
// stays optional.
//
// The condition references positional placeholders $N-1 (domain) and $N
// (group), and appends both values to args. Callers must pass the same args
// slice to the query they built the clause for.
func turtahunGroupFilter(domainID, turtahunGroupID string, args *[]interface{}) string {
	if turtahunGroupID == "" {
		return ""
	}
	*args = append(*args, domainID, turtahunGroupID)
	domainPos := len(*args) - 1
	groupPos := len(*args)
	return fmt.Sprintf(` AND turtahun_id IN (
		SELECT turtahun_id FROM webapi.master_tahun_turunan
		WHERE domain_id = $%d AND group_turth_id = $%d
	)`, domainPos, groupPos)
}
