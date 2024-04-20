package game

import (
	"net/http"

	"common/manager/gamemanager"
	"webserver/controller/utils"

	"github.com/gin-gonic/gin"
)

// 返回比赛活动列表 近期
func GetGameList(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query GameListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	games, total, err := gamemanager.GameList(&gamemanager.GameListParameter{
		UserName: query.UserName,
		Title: query.Title,
		Date:     query.Date,
		Stadium:  query.Stadium,
		CanApply: query.CanApply,
		Page:     query.Page,
		PerPage:  query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: games,
		Total: total,
	})
}

// 返回过期比赛活动列表
func GetOutdatedGameList(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query GameListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	games, total, err := gamemanager.OutdatedGameList(&gamemanager.GameListParameter{
		UserName: query.UserName,
		Date:     query.Date,
		Stadium:  query.Stadium,
		Title:    query.Title,
		CanApply: query.CanApply,
		Page:     query.Page,
		PerPage:  query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: games,
		Total: total,
	})
}

// 返回用户举办的比赛 未过期
func GetHostGameList(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query GameListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	games, total, err := gamemanager.GameList(&gamemanager.GameListParameter{
		Number:  number,
		Date:    query.Date,
		Title:   query.Title,
		Page:    query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: games,
		Total: total,
	})
}

// 返回用户举办已过期比赛
func GetOutdatedHostGameList(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query GameListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	games, total, err := gamemanager.OutdatedGameList(&gamemanager.GameListParameter{
		Number:  number,
		Date:    query.Date,
		Title:   query.Title,
		Page:    query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: games,
		Total: total,
	})
}

// 创建比赛活动
func CreateGame(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	username, exists := c.Get("name")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	name := username.(string)

	var query CreateGameQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = gamemanager.CreateGame(&gamemanager.CreateParameter{
		Number: number,
		UserName: name,
		Title: query.Title,
		Date: query.Date,
		StadiumID: query.StadiumID,
		Stadium: query.Stadium,
		Start: query.Start,
		End: query.End,
		Remark: query.Remark,
		Maximum: query.Maximum,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "创建比赛活动失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "创建成功", nil)
}

// 报名比赛活动
func ApplyGame(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	username, exists := c.Get("name")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	name := username.(string)
	var query ApplyGameQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = gamemanager.ApplyGame(&gamemanager.ApplyParameter{
		Number: number,
		UserName: name,
		GameID: query.ID,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "报名失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "报名成功", nil)
}

// 返回报名记录
func GetApplyRecord(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query RecordListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := gamemanager.RecentRecords(&gamemanager.ApplyRecordParameter{
		Number: number,
		UserName: query.UserName,
		Date: query.Date,
		Title: query.Title,
		Stadium: query.Stadium,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: records,
		Total: total,
	})
}

// 返回已过期报名记录
func OutdatedApplyRecord(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query RecordListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := gamemanager.OutdatedRecord(&gamemanager.ApplyRecordParameter{
		Number: number,
		UserName: query.UserName,
		Date: query.Date,
		Title: query.Title,
		Stadium: query.Stadium,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: records,
		Total: total,
	})
}

// 取消报名
func CancelApply(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query CancelApplyQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = gamemanager.CancelApplyGame(&gamemanager.CancelApplyParameter{RecordID: query.ID})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "取消失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "取消成功", nil)
}

// 返回近期申请记录
func GetApplication(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query ApplicationQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := gamemanager.RecentRecords(&gamemanager.ApplyRecordParameter{
		HostNumber: number,
		Date: query.Date,
		Title: query.Title,
		UserName: query.UserName,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: records,
		Total: total,
	})
}

// 返回过期的申请记录
func GetOutdatedApplication(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query ApplicationQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := gamemanager.OutdatedRecord(&gamemanager.ApplyRecordParameter{
		HostNumber: number,
		Date: query.Date,
		Title: query.Title,
		UserName: query.UserName,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: records,
		Total: total,
	})
}

// 接受报名
func AcceptApply(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query AcceptApplyQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = gamemanager.AcceptApply(&gamemanager.AcceptApplyParameter{
		RecordID: query.ID,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "通过失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "通过成功", nil)
}

// 拒绝报名
func RefuseApply(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query AcceptApplyQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = gamemanager.RefuseApply(&gamemanager.RefuseApplyParameter{
		RecordID: query.ID,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "拒绝失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "拒绝成功", nil)
}

// 取消比赛
func CancelGame(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query CancelGameQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = gamemanager.CancelGame(&gamemanager.CancelParameter{
		GameID: query.ID,
		BookRecordID: query.BookReocrdID,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "取消失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "取消成功", nil)
}