package models

type BaseMetadata struct {
	TotalData          int    `json:"total_data"`
	TotalPage          int    `json:"total_page"`
	PerPage            int    `json:"per_page"`
	Page               int    `json:"page"`
	LastUpdateData     string `json:"last_update_data"`
	LastUpdatePipeline string `json:"last_update_pipeline"`
}
