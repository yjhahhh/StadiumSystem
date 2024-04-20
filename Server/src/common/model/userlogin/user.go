package userlogin

import(

	"common/connection"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Number		string `gorm:"not null;unique;index"`
	Password	string `gorm:"not null"`
	Name		string `gorm:"not null"`
	Department	string `gorm:"not null"`
	Major		string `gorm:"not null"`
	Class		uint `gorm:"not null"`
}

func Init() {
	err := connection.GetDB().AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}
}