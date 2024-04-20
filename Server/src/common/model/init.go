package model

import (
	"common/model/adminlogin"
	"common/model/stadium"
	"common/model/userlogin"
	"common/model/game"
)

func InitModel() {
	adminlogin.Init()
	userlogin.Init()
	stadium.InitRecord()
	stadium.InitStadium()
	game.InitGame()
	game.InitGameRecord()
}
