package models

// TahunDistinct is a distinct set of (tahun_id, tahun_name) from
// webapi.datacontent.
type TahunDistinct struct {
	TahunID   string `gorm:"column:tahun_id" json:"tahun_id"`
	TahunName string `gorm:"column:tahun_name" json:"tahun_name"`
}

func (TahunDistinct) TableName() string { return "webapi.datacontent" }
