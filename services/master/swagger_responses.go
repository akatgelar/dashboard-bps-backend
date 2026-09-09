package services

import (
	"akatgelar/dashboard-bps-backend/models"
	ModelData "akatgelar/dashboard-bps-backend/models/data"
)

// MasterListEnvelope is the shared shape of every master list endpoint:
// { status, message, data: [...], metadata: { total_data, ... } }.
type MasterListEnvelope struct {
	Status   bool                `json:"status"`
	Message  string              `json:"message"`
	Data     []ModelData.Domain  `json:"data"`
	Metadata models.BaseMetadata `json:"metadata"`
}

// DomainListResponse documents GET /master/domain.
type DomainListResponse struct {
	Status   bool                `json:"status"`
	Message  string              `json:"message"`
	Data     []ModelData.Domain  `json:"data"`
	Metadata models.BaseMetadata `json:"metadata"`
}

// SubjectListResponse documents GET /master/subjek.
type SubjectListResponse struct {
	Status   bool                `json:"status"`
	Message  string              `json:"message"`
	Data     []ModelData.Subject `json:"data"`
	Metadata models.BaseMetadata `json:"metadata"`
}

// VariableListResponse documents GET /master/variable.
type VariableListResponse struct {
	Status   bool                 `json:"status"`
	Message  string               `json:"message"`
	Data     []ModelData.Variable `json:"data"`
	Metadata models.BaseMetadata  `json:"metadata"`
}

// VariableDistinctListResponse documents GET /master/variable-distinct.
type VariableDistinctListResponse struct {
	Status   bool                 `json:"status"`
	Message  string               `json:"message"`
	Data     []ModelData.Variable `json:"data"`
	Metadata models.BaseMetadata  `json:"metadata"`
}

// VariableTurunanListResponse documents GET /master/variable-turunan.
type VariableTurunanListResponse struct {
	Status   bool                        `json:"status"`
	Message  string                      `json:"message"`
	Data     []ModelData.VariableTurunan `json:"data"`
	Metadata models.BaseMetadata         `json:"metadata"`
}

// VariableVerticalListResponse documents GET /master/variable-vertical.
type VariableVerticalListResponse struct {
	Status   bool                         `json:"status"`
	Message  string                       `json:"message"`
	Data     []ModelData.VariableVertical `json:"data"`
	Metadata models.BaseMetadata          `json:"metadata"`
}

// TahunListResponse documents GET /master/tahun.
type TahunListResponse struct {
	Status   bool                `json:"status"`
	Message  string              `json:"message"`
	Data     []ModelData.Tahun   `json:"data"`
	Metadata models.BaseMetadata `json:"metadata"`
}

// TahunTurunanListResponse documents GET /master/tahun-turunan.
type TahunTurunanListResponse struct {
	Status   bool                     `json:"status"`
	Message  string                   `json:"message"`
	Data     []ModelData.TahunTurunan `json:"data"`
	Metadata models.BaseMetadata      `json:"metadata"`
}

// TahunTurunanDistinctListResponse documents GET /master/tahun-turunan-distinct.
type TahunTurunanDistinctListResponse struct {
	Status   bool                             `json:"status"`
	Message  string                           `json:"message"`
	Data     []ModelData.TahunTurunanDistinct `json:"data"`
	Metadata models.BaseMetadata              `json:"metadata"`
}

// TahunDistinctListResponse documents GET /master/tahun-distinct.
type TahunDistinctListResponse struct {
	Status   bool                      `json:"status"`
	Message  string                    `json:"message"`
	Data     []ModelData.TahunDistinct `json:"data"`
	Metadata models.BaseMetadata       `json:"metadata"`
}

// DataContentListResponse documents GET /datacontent.
type DataContentListResponse struct {
	Status   bool                    `json:"status"`
	Message  string                  `json:"message"`
	Data     []ModelData.DataContent `json:"data"`
	Metadata models.BaseMetadata     `json:"metadata"`
}
