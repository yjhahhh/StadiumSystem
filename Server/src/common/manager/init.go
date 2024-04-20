package manager

import (
	"common/manager/ordermanager"
	"common/manager/stadiummanager"
)

func InitManager() {
	ordermanager.InitTimetable()
	stadiummanager.Init()
}
