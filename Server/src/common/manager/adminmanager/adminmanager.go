package adminmanager

import (
	"common/connection"
	"common/model/adminlogin"
	"fmt"

	"gorm.io/gorm"
)

const (
	SuperLevel   = 1
	GeneralLevel = 2
)

func CheckLevel(level uint) bool {
	if level != SuperLevel && level != GeneralLevel {
		return false
	}
	return true
}

// 添加管理员
func AddAdmin(parameter *AddAdminParameter) error {
	if !CheckLevel(parameter.Level) {
		return fmt.Errorf("invalid level")
	}
	admin := adminlogin.Admin {
		AdminName: parameter.Name,
		Password: parameter.Password,
		Level: parameter.Level,
	}
	return connection.GetDB().Create(&admin).Error
}

// 删除管理员
func DeleteAdmin(parameter *DeleteAdminParameter) error {
	return connection.GetDB().Unscoped().Delete(&adminlogin.Admin{}, parameter.ID).Error
}

// 更改管理员权限
func UpdateAdmin(parameter *UpdateAdminParameter) error {
	return connection.GetDB().Model(&adminlogin.Admin{Model: gorm.Model{ID: parameter.ID}}).Update("Level", parameter.Level).Error
}

// 获取管理员列表
func GetAdminList(parameter *AdminListParameter) ([]Admin, int, error) {
	var admins []adminlogin.Admin
	err := connection.GetDB().Where(&adminlogin.Admin{
		AdminName: parameter.Name,
		Level: parameter.Level,
	}).Find(&admins).Error
	if err != nil {
		return nil, 0, err
	}
	start := parameter.PerPage * (parameter.Page - 1)
	end := start + parameter.PerPage
	if start >= len(admins) {
		start = len(admins)
	}
	if end > len(admins) {
		end = len(admins)
	}
	ret := make([]Admin, 0, end - start)
	for i := start; i < end; i++ {
		ret = append(ret, Admin{
			ID: admins[i].ID,
			Name: admins[i].AdminName,
			Level: int(admins[i].Level),
		})
	}
	return ret, len(admins), err
}

// 获取管理员信息
func GetInfo(name string) (Info, error) {
	var admin adminlogin.Admin
	err := connection.GetDB().Where(&adminlogin.Admin{AdminName: name}).First(&admin).Error
	info := Info {
		Name: admin.AdminName,
		Level: admin.Level,
	}
	return info, err
}