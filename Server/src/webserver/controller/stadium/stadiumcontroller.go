package stadium

import (
	"common/manager/stadiummanager"
	"net/http"
	"webserver/controller/utils"
	"webserver/jwt"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/sessions"
)

// 获取体育场馆列表
func StadiumList(c *gin.Context) {
	var query StadiumListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	data, total := stadiummanager.StadiumList(stadiummanager.StadiumQuery{
		Name: query.Name,
		Category: query.Categoty,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "OK", &utils.DataList {
		Data: data,
		Total: total,
	})
}

// 添加体育场馆
func AddStadium(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	role := c.GetInt("role")
	if role != jwt.AdminRole {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query AddStadiumQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = stadiummanager.AddStadium(query.Name, query.Categoty)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "添加失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "添加成功", nil)
}

// 删除体育场馆
func DeleteStadium(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	role := c.GetInt("role")
	if role != jwt.AdminRole {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query DeleteStadiumQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = stadiummanager.DeleteStadium(query.ID)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "删除失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "删除成功", nil)
	}
}

// 更新体育场馆
func UpdateStadium(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	role := c.GetInt("role")
	if role != jwt.AdminRole {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query UpadateStadiumQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = stadiummanager.UpdateStadium(query.ID, query.Name)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "更新失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "更新成功", nil)
	}
}