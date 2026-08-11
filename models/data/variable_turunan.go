package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// VariableTurunan maps webapi.master_variable_turunan.
type VariableTurunan struct {
	Id         int64  `gorm:"column:id" json:"id"`
	DomainID   string `gorm:"column:domain_id" json:"domain_id"`
	VarID      string `gorm:"column:var_id" json:"var_id"`
	VarName    string `gorm:"column:var_name" json:"var_name"`
	TurvarID   string `gorm:"column:turvar_id" json:"turvar_id"`
	TurvarName string `gorm:"column:turvar_name" json:"turvar_name"`
	GetAt      HelperTime.Time   `gorm:"column:get_at" json:"get_at"`
}

func (VariableTurunan) TableName() string { return "webapi.master_variable_turunan" }
