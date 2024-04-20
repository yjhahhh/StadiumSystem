package stadium

import (
	"common/connection"

	"gorm.io/gorm"
)


type Stadium struct {
	gorm.Model
	Name     string `gorm:"not null;unique"`
	Category string `gorm:"not null"`
}

func InitStadium() {
	err := connection.GetDB().AutoMigrate(&Stadium{})
	if err != nil {
		panic(err)
	}
}