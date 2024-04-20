package game

import (
	"common/connection"

	"gorm.io/gorm"
)

type GameRecord struct {
	gorm.Model
	UserNumber string `gorm:"not null;index"` // 申请人
	UserName   string `gorm:"not null"`
	Stadium    string `gorm:"not null"`
	GameID     uint   `gorm:"not null;index"`
	Title      string `gorm:"not null"`
	Date       string `gorm:"not null"`
	Start      string `gorm:"not null"`
	End        string `gorm:"not null"`
	HostNumber string   `gorm:"not null;index"`
	HostName   string `gorm:"not null"`
	// 申请状态
	Status     string `gorm:"not null"`
}

func InitGameRecord() {
	err := connection.GetDB().AutoMigrate(&GameRecord{})
	if err != nil {
		panic(err)
	}
}