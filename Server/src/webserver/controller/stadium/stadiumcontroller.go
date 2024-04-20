package stadium

import (
	"common/manager/stadiummanager"
	"fmt"
	"net/http"
	"webserver/controller/utils"
	"webserver/jwt"

	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/sessions"
)

// 获取体育场馆列表
func StadiumList(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "user获取错误", nil)
		return
	}
	var query StadiumListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	data, total := stadiummanager.StadiumList(stadiummanager.StadiumListParameter{
		Name:     query.Name,
		Category: query.Categoty,
		Page:     query.Page,
		PerPage:  query.PerPage,
	})
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "OK", &utils.DataList{
		Data:  data,
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
	var query AddStadiumQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = stadiummanager.AddStadium(&stadiummanager.StadiumParameter{
		Name: query.Name,
		Category: query.Categoty,
		Start: query.Start,
		End: query.End,
	})
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
	var query UpadateStadiumQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取参数错误", nil)
		return
	}
	err = stadiummanager.UpdateStadium(&stadiummanager.UpdateParameter{
		ID: query.ID,
		Name: query.Name,
		Start: query.Start,
		End: query.End,
	})
	if err != nil {
		fmt.Println(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "更新失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "更新成功", nil)
	}
}
