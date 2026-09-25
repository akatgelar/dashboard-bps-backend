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
}

func (Variable) TableName() string { return "webapi.master_variable" }
