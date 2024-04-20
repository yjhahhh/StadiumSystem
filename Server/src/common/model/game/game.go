package game

import (
	"common/connection"

	"gorm.io/gorm"
)

type Game struct {
	gorm.Model
	Title        string `gorm:"not null;index"`
	Date         string `gorm:"not null;index:idx_date_stadium"`
	Stadium      string `gorm:"not null;index:idx_date_stadium"`
	StadiumID    uint `gorm:"not null"`
	BookRecordID uint `gorm:"not null;unique"`
	Start        string `gorm:"not null"`
	End          string `gorm:"not null"`
	UserNumber   string `gorm:"not null;index"`
	UserName     string `gorm:"not null;index"`
	Remark       string 
	Count        uint `gorm:"not null"`
	Maximum      uint `gorm:"not null"`
}


func InitGame() {
	err := connection.GetDB().AutoMigrate(&Game{})
	if err != nil {
		panic(err)
	}
}