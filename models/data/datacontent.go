package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// DataContent maps webapi.datacontent.
type DataContent struct {
	Id               int64    `gorm:"column:id" json:"id"`
	Domain           string   `gorm:"column:domain_id" json:"domain_id"`
	VarID            string   `gorm:"column:var_id" json:"var_id"`
	VarName          string   `gorm:"column:var_name" json:"var_name"`
	Unit             string   `gorm:"column:var_unit" json:"var_unit"`
	SubName          string   `gorm:"column:sub_name" json:"sub_name"`
	TurvarID         string   `gorm:"column:turvar_id" json:"turvar_id"`
	TurvarName       string   `gorm:"column:turvar_name" json:"turvar_name"`
	TahunID          string   `gorm:"column:tahun_id" json:"tahun_id"`
	TahunName        string   `gorm:"column:tahun_name" json:"tahun_name"`
	TurtahunID       string   `gorm:"column:turtahun_id" json:"turtahun_id"`
	TurtahunName     string   `gorm:"column:turtahun_name" json:"turtahun_name"`
	VervarID         string   `gorm:"column:vervar_id" json:"vervar_id"`
	VervarName       string   `gorm:"column:vervar_name" json:"vervar_name"`
	DataContentID    string   `gorm:"column:datacontent_id" json:"datacontent_id"`
	DataContentValue *float64 `gorm:"column:datacontent_value" json:"datacontent_value"`
	LastUpdatedAt    HelperTime.Time     `gorm:"column:last_updated_at" json:"last_updated_at"`
	GetAt            HelperTime.Time     `gorm:"column:get_at" json:"get_at"`
}

func (DataContent) TableName() string { return "webapi.datacontent" }
