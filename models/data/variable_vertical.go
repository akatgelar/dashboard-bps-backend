package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// VariableVertical maps webapi.master_variable_vertical.
type VariableVertical struct {
	Id         int64  `gorm:"column:id" json:"id"`
	DomainID   string `gorm:"column:domain_id" json:"domain_id"`
	VarID      string `gorm:"column:var_id" json:"var_id"`
	VarName    string `gorm:"column:var_name" json:"var_name"`
	VervarID   string `gorm:"column:vervar_id" json:"vervar_id"`
	VervarName string `gorm:"column:vervar_name" json:"vervar_name"`
	GetAt      HelperTime.Time   `gorm:"column:get_at" json:"get_at"`
}

func (VariableVertical) TableName() string { return "webapi.master_variable_vertical" }
