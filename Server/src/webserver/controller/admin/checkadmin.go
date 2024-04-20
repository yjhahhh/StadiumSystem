package admin

import (
	"common/connection"
	"common/model/adminlogin"
	"gorm.io/gorm"
)



func CheckIsAdmin(name string) (bool, uint, error) {
	var admin adminlogin.Admin
	result := connection.GetDB().Where(&adminlogin.Admin{AdminName: name}).First(&admin)
	exists := true
	if result.Error == gorm.ErrRecordNotFound {
		exists = false
	}
	return exists, admin.Level, result.Error
}