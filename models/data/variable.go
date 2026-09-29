package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// Variable maps webapi.master_variable.
type Variable struct {
	Id       int64           `gorm:"column:id" json:"id"`
	DomainID string          `gorm:"column:domain_id" json:"domain_id"`
	VarID    string          `gorm:"column:var_id" json:"var_id"`
	VarName  string          `gorm:"column:var_name" json:"var_name"`
	SubID    string          `gorm:"column:sub_id" json:"sub_id"`
	SubName  string          `gorm:"column:sub_name" json:"sub_name"`
	Def      string          `gorm:"column:def" json:"def"`
	Notes    string          `gorm:"column:notes" json:"notes"`
	Unit     string          `gorm:"column:unit" json:"unit"`
	GetAt    HelperTime.Time `gorm:"column:get_at" json:"get_at"`

	// GroupID/GroupName are not columns of master_variable: they are filled from
	// master_variable_turunan (group_turvar_id / name_group_turvar) via a select
	// expression set by the variable endpoints.
	GroupID   string `gorm:"->;column:group_id" json:"group_id"`
	GroupName string `gorm:"->;column:group_name" json:"group_name"`
}

func (Variable) TableName() string { return "webapi.master_variable" }
