package admin

import (
	"common/connection"
	"common/manager/adminmanager"
	"common/model/adminlogin"
	
	"net/http"
	"webserver/jwt"
	"webserver/controller/utils"
	"common/log"
	
	"github.com/gin-gonic/gin"
)

const (
	LoginOK      = "登陆成功"
	LoginFail    = "账号或密码错误"
	RegisterOK   = "注册成功"
	RegisterFail = "注册失败"
)


func Login(c *gin.Context) {
	var loginQuery LoginQuery
	err := c.ShouldBind(&loginQuery)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, utils.BadParameter, nil)
		return
	}
	admin := adminlogin.Admin{}

	db := connection.GetDB()
	result := db.Where(&adminlogin.Admin{AdminName: loginQuery.Name}).First(&admin)
	if result.Error != nil {
		log.Error(result.Error)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, LoginFail, nil)
		return
	}
	if loginQuery.Password != admin.Password {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, LoginFail, nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, LoginOK, utils.DataWithToken{
		Token: jwt.GenerateToken(&jwt.UserClaims{
			Number: admin.AdminName,
			Role: jwt.AdminRole,
			Level: admin.Level,
		}),
	})
}

func Register(c *gin.Context) {
	var registerQuery RegisterQuery
	err := c.ShouldBind(&registerQuery)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, utils.BadParameter, nil)
		return
	}
	admin := adminlogin.Admin{
		AdminName: registerQuery.Name,
		Password: registerQuery.Password,
		Level: registerQuery.Level,
	}
	db := connection.GetDB()
	result := db.Create(&admin)
	if result.Error != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, RegisterFail, nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, RegisterOK, nil)
	}
}

// 返回管理员列表
func GetAdminList(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
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
	var query AdminListQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "参数获取错误", nil)
		return
	}
	ret, total, err := adminmanager.GetAdminList(&adminmanager.AdminListParameter{
		Name: query.Name,
		Level: query.Level,
		Page: query.Page,
		PerPage: query.PerPage,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取成功", &utils.DataList {
		Data: ret,
		Total: total,
	})
}

// 添加管理员
func AddAdmin(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	roleStr, exists := c.Get("role")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	role := roleStr.(uint)
	levelStr, exists := c.Get("level")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	level := levelStr.(uint)
	if role != jwt.AdminRole || level != adminmanager.SuperLevel {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query RegisterQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "参数获取错误", nil)
		return
	}
	err = adminmanager.AddAdmin(&adminmanager.AddAdminParameter{
		Name: query.Name,
		Password: query.Password,
		Level: query.Level,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "添加失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "添加成功", nil)
	}
}

// 删除管理员
func DeleteAdmin(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	roleStr, exists := c.Get("role")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	role := roleStr.(uint)
	levelStr, exists := c.Get("level")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	level := levelStr.(uint)
	if role != jwt.AdminRole || level != adminmanager.SuperLevel {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query DeleteAdminQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "参数获取错误", nil)
		return
	}
	err = adminmanager.DeleteAdmin(&adminmanager.DeleteAdminParameter{
		ID: query.ID,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "删除失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "删除成功", nil)
	}
}

// 更改管理员
func UpdateAdmin(c *gin.Context) {
	_, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	roleStr, exists := c.Get("role")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	role := roleStr.(uint)
	levelStr, exists := c.Get("level")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	level := levelStr.(uint)
	if role != jwt.AdminRole || level != adminmanager.SuperLevel {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "无权限", nil)
		return
	}
	var query UpdateAdminQuery
	err := c.ShouldBind(&query)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "参数获取错误", nil)
		return
	}
	if !adminmanager.CheckLevel(query.Level) {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "管理员权限设置无效", nil)
		return
	}
	err = adminmanager.UpdateAdmin(&adminmanager.UpdateAdminParameter{
		ID: query.ID,
		Level: query.Level,
	})
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "更改失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "更改成功", nil)
	}
}

// 管理员信息
func GetAdminInfo(c *gin.Context) {
	user, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.UnauthorizedCode, "获取user失败", nil)
		return
	}
	name := user.(string)
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
	info, err := adminmanager.GetInfo(name)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode,"获取成功", &info)
	}
}