package models

// TahunTurunanDistinct is a distinct set of (turtahun_id, turtahun_name) from
// webapi.datacontent.
type TahunTurunanDistinct struct {
	TurtahunID   string `gorm:"column:turtahun_id" json:"turtahun_id"`
	TurtahunName string `gorm:"column:turtahun_name" json:"turtahun_name"`
}

func (TahunTurunanDistinct) TableName() string { return "webapi.datacontent" }
