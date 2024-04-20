package order

import (
	"net/http"

	"common/log"
	"common/manager/ordermanager"
	"webserver/controller/utils"

	"github.com/gin-gonic/gin"
)

const (
	OrderOK = "OrderOK"
	OrderFail = "OrderFail"
)

// 返回可预约时段
func StadiumAllowableTimes(c *gin.Context) {
	var query IdleTimeQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "参数获取错误", nil)
		return
	}
	idleTimes, total := ordermanager.GetAllowableTimes(&ordermanager.IdleTimeQuery{
		StadiumID: query.ID,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "OK", &utils.DataList {
		Data: idleTimes,
		Total: total,
	})
}

// 预约
func Order(c *gin.Context) {
	var query OrderQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "参数获取错误", nil)
		return
	}
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	err = ordermanager.OrderStadium(&ordermanager.OrderParameter{
		StadiumID: query.ID,
		Number: number,
		Stadium: query.Name,
		Date: query.Date,
		Start: query.Start,
		End: query.End,
	})
	if err != nil {
		log.Errorf("Order fail err = %s\n", err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, OrderFail, nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, OrderOK, nil)
	}
}