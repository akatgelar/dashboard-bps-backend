package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// Tahun maps webapi.master_tahun.
type Tahun struct {
	Id        int64  `gorm:"column:id" json:"id"`
	DomainID  string `gorm:"column:domain_id" json:"domain_id"`
	TahunID   string `gorm:"column:tahun_id" json:"tahun_id"`
	TahunName string `gorm:"column:tahun_name" json:"tahun_name"`
	GetAt     HelperTime.Time   `gorm:"column:get_at" json:"get_at"`
	IsActive  bool   `gorm:"column:is_active" json:"is_active"`
}

func (Tahun) TableName() string { return "webapi.master_tahun" }
