package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// Subject maps webapi.master_subject.
type Subject struct {
	Id       int64  `gorm:"column:id" json:"id"`
	DomainID string `gorm:"column:domain_id" json:"domain_id"`
	SubID    string `gorm:"column:sub_id" json:"sub_id"`
	SubName  string `gorm:"column:sub_name" json:"sub_name"`
	GetAt    HelperTime.Time   `gorm:"column:get_at" json:"get_at"`
}

func (Subject) TableName() string { return "webapi.master_subject" }
