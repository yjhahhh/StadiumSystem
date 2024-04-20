package record

import (
	"net/http"

	"common/log"
	"common/manager/ordermanager"
	"webserver/controller/utils"
	"webserver/jwt"

	"github.com/gin-gonic/gin"
)

// 返回用户预约记录
func OrderRecords(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query RecordQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := ordermanager.GetOrderRecord(&ordermanager.Paramter{
		Number:  number,
		Stadium: query.Stadium,
		Date:    query.Date,
		Page:    query.Page,
		PerPage: query.PerPage,
	}, false)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取预约记录失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取预约记录成功", &utils.DataList{
			Data:  records,
			Total: total,
		})
	}
}

// 返回用户预约记录 已过期
func OutdatedOrderRecords(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	number := user.(string)
	var query RecordQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := ordermanager.GetOrderRecord(&ordermanager.Paramter{
		Number:  number,
		Stadium: query.Stadium,
		Date:    query.Date,
		Page:    query.Page,
		PerPage: query.PerPage,
	}, true)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取预约记录失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取预约记录成功", &utils.DataList{
			Data:  records,
			Total: total,
		})
	}
}

// 返回预约记录  验证管理员
func AllOrderRecords(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	roleStr, exists := c.Get("role")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	role := roleStr.(uint)
	if role != jwt.AdminRole {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query AllRecordQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := ordermanager.GetOrderRecord(&ordermanager.Paramter{
		Number:  query.Number,
		Date:    query.Date,
		Stadium: query.Stadium,
		Page:    query.Page,
		PerPage: query.PerPage,
	}, false)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取预约记录失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取预约记录成功", &utils.DataList{
			Data:  records,
			Total: total,
		})
	}
}

// 返回预约记录  验证管理员
func AllOutdatedOrderRecords(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	roleStr, exists := c.Get("role")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	role := roleStr.(uint)
	if role != jwt.AdminRole {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query AllRecordQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	records, total, err := ordermanager.GetOrderRecord(&ordermanager.Paramter{
		Number:  query.Number,
		Date:    query.Date,
		Stadium: query.Stadium,
		Page:    query.Page,
		PerPage: query.PerPage,
	}, true)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取预约记录失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取预约记录成功", &utils.DataList{
			Data:  records,
			Total: total,
		})
	}
}

// 取消预约
func CancelOrder(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query CancleOrder
	err := c.ShouldBind(&query)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = ordermanager.CancelOrder(query.OrderID)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "取消预约失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "取消预约成功", nil)
	}
}

// 管理员取消预约
func AdminCancelORder(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	roleStr, exists := c.Get("role")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	role := roleStr.(uint)
	if role != jwt.AdminRole {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query CancleOrder
	err := c.ShouldBind(&query)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = ordermanager.CancelOrder(query.OrderID)
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "取消预约失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "取消预约成功", nil)
	}
}
