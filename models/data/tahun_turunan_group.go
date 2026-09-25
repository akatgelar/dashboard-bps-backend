package models

// TahunTurunanGroup is the period group (group_turth) of a tahun turunan,
// resolved by joining webapi.datacontent with webapi.master_tahun_turunan.
type TahunTurunanGroup struct {
	TurtahunGroupID   string `gorm:"column:turtahun_group_id" json:"turtahun_group_id"`
	TurtahunGroupName string `gorm:"column:turtahun_group_name" json:"turtahun_group_name"`
}

func (TahunTurunanGroup) TableName() string { return "webapi.datacontent" }
