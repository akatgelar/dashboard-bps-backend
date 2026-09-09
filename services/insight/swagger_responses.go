package services

// BigNumberResponse documents GET /insight/big-number.
type BigNumberResponse struct {
	Status   bool              `json:"status"`
	Message  string            `json:"message"`
	Data     BigNumberData     `json:"data"`
	Metadata BigNumberMetadata `json:"metadata"`
}

// PerTahunResponse documents GET /insight/per-tahun.
type PerTahunResponse struct {
	Status   bool             `json:"status"`
	Message  string           `json:"message"`
	Data     []PerTahunRow    `json:"data"`
	Metadata PerTahunMetadata `json:"metadata"`
}

// PerWilayahResponse documents GET /insight/per-wilayah.
type PerWilayahResponse struct {
	Status   bool               `json:"status"`
	Message  string             `json:"message"`
	Data     PerWilayahData     `json:"data"`
	Metadata PerWilayahMetadata `json:"metadata"`
}

// PerWilayahTurvarResponse documents GET /insight/per-wilayah-turvar.
type PerWilayahTurvarResponse struct {
	Status   bool                     `json:"status"`
	Message  string                   `json:"message"`
	Data     []PerWilayahTurvarValue  `json:"data"`
	Metadata PerWilayahTurvarMetadata `json:"metadata"`
}
