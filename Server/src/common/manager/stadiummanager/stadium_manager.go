package stadiummanager

import (
	"fmt"
	"strings"
	"sync"

	"common/connection"
	"common/model/stadium"
	"common/utils/timeutils"

	"gorm.io/gorm"
)

// type StadiumManager struct {
var (
	// CategoryList  []Category // 场馆分类列表
	stadiumList []stadium.Stadium // 体育场馆列表
	lock        sync.RWMutex
)

// }

// var g_stadiumManager *StadiumManager

func Init() {
	var err error
	stadiumList, err = getStadiumFromDB()
	if err != nil {
		panic(err)
	}
}

// 从数据库获取体育场列表
func getStadiumFromDB() ([]stadium.Stadium, error) {
	var stadiums []stadium.Stadium
	err := connection.GetDB().Find(&stadiums).Error
	if err != nil {
		return nil, err
	}
	return stadiums, nil
}

// 返回体育场馆列表
func StadiumList(query StadiumListParameter) ([]stadium.Stadium, int) {
	lock.RLock()
	defer lock.RUnlock()
	start := query.PerPage * (query.Page - 1)
	end := start + query.PerPage
	if query.Name == "" && query.Category == "" {
		if start >= len(stadiumList) {
			start = len(stadiumList)
		}
		if end > len(stadiumList) {
			end = len(stadiumList)
		}
		ret := make([]stadium.Stadium, end-start)
		copy(ret, stadiumList[start:end])
		return ret, len(stadiumList)
	}
	ret := make([]stadium.Stadium, 0, len(stadiumList))
	for i := range stadiumList {
		if query.Name != "" && strings.Contains(stadiumList[i].Name, query.Name) {
			ret = append(ret, stadiumList[i])
		} else if query.Category != "" && strings.Contains(stadiumList[i].Category, query.Category) {
			ret = append(ret, stadiumList[i])
		}
	}

	if len(ret) < end {
		end = len(ret)
	}
	return ret[start:end], len(ret)
}

// 添加体育场馆
func AddStadium(parameter *StadiumParameter) error {
	legal, err := timeutils.IsLegal(parameter.Start, parameter.End)
	if err != nil {
		return err
	}
	if !legal {
		return fmt.Errorf("时间不合法")
	}
	stadium := stadium.Stadium{
		Name:     parameter.Name,
		Category: parameter.Category,
		Start: parameter.Start,
		End: parameter.End,
	}
	result := connection.GetDB().Create(&stadium)
	if result.Error != nil {
		return result.Error
	}
	lock.Lock()
	defer lock.Unlock()
	stadiumList = append(stadiumList, stadium)
	return nil
}

// 删除体育场馆
func DeleteStadium(id uint) error {
	result := connection.GetDB().Unscoped().Delete(&stadium.Stadium{}, id)
	if result.Error != nil {
		return result.Error
	}
	lock.Lock()
	defer lock.Unlock()
	for i := range stadiumList {
		if stadiumList[i].ID == id {
			stadiumList = append(stadiumList[:i], stadiumList[i+1:]...)
			break
		}
	}
	return nil
}

// 更新体育场馆
func UpdateStadium(parameter *UpdateParameter) error {
	legal, err := timeutils.IsLegal(parameter.Start, parameter.End)
	if err != nil {
		return err
	}
	if !legal {
		return fmt.Errorf("时间不合法")
	}
	result := connection.GetDB().Model(&stadium.Stadium{Model: gorm.Model{ID: parameter.ID}}).Updates(&stadium.Stadium {
		Name: parameter.Name,
		Start: parameter.Start,
		End: parameter.End,
	})
	if result.Error != nil {
		return result.Error
	}
	lock.Lock()
	defer lock.Unlock()
	for i := range stadiumList {
		if stadiumList[i].ID == parameter.ID {
			stadiumList[i].Name = parameter.Name
			stadiumList[i].Start = parameter.Start
			stadiumList[i].End = parameter.End
			break
		}
	}
	return nil
}
