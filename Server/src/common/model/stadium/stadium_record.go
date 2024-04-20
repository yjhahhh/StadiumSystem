package stadium

import (
	"common/connection"

	"gorm.io/gorm"
)

type StadiumRecord struct {
	gorm.Model
	StadiumId uint `gorm:"not null;index:idx_stadium_date"`
	// 日期 格式 2001-01-01
	Date    string `gorm:"not null;index:idx_stadium_date"`
	Stadium string `gorm:"not null"`
	// 时间段 格式 10:00
	UserNo string `gorm:"not null;index"`
	Start  string `gorm:"not null"`
	End    string `gorm:"not null"`
}

func InitRecord() {
	err := connection.GetDB().AutoMigrate(&StadiumRecord{})
	if err != nil {
		panic(err)
	}
}
