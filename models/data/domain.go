package models

import (
	HelperTime "akatgelar/dashboard-bps-backend/helpers"
)

// Domain maps webapi.master_domain.
type Domain struct {
	Id         int64  `gorm:"column:id" json:"id"`
	DomainID   string `gorm:"column:domain_id" json:"domain_id"`
	DomainName string `gorm:"column:domain_name" json:"domain_name"`
	DomainURL  string `gorm:"column:domain_url" json:"domain_url"`
	GetAt      HelperTime.Time   `gorm:"column:get_at" json:"get_at"`
	IsActive   bool   `gorm:"column:is_active" json:"is_active"`
}

func (Domain) TableName() string { return "webapi.master_domain" }
