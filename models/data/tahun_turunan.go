package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// TahunTurunan maps webapi.master_tahun_turunan.
type TahunTurunan struct {
	Id             int64  `gorm:"column:id" json:"id"`
	DomainID       string `gorm:"column:domain_id" json:"domain_id"`
	TurtahunID     string `gorm:"column:turtahun_id" json:"turtahun_id"`
	TurtahunName   string `gorm:"column:turtahun_name" json:"turtahun_name"`
	GroupTurthID   string `gorm:"column:group_turth_id" json:"group_turth_id"`
	GroupTurthName string `gorm:"column:group_turth_name" json:"group_turth_name"`
	GetAt          HelperTime.Time   `gorm:"column:get_at" json:"get_at"`
}

func (TahunTurunan) TableName() string { return "webapi.master_tahun_turunan" }
