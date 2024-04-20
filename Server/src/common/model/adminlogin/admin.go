package adminlogin

import(
	"common/connection"

	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	AdminName	string `gorm:"not null;index"`
	Password	string `gorm:"not null"`
	Level		uint `gorm:"bot null"`
}

func Init() {
	connection.GetDB().AutoMigrate(&Admin{})
}