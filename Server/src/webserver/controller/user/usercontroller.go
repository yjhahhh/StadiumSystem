package user

import (
	"net/http"
	
	"common/connection"
	"common/log"
	"common/model/userlogin"
	"webserver/jwt"
	"webserver/controller/utils"

	"github.com/gin-gonic/gin"
)

const (
	LoginOK      = "登录成功"
	LoginFail    = "用户名或密码错误"
	RegisterOK   = "注册成功"
	RegisterFail = "注册失败"
	LogoutOK     = "LogoutOK"
	LogoutFail   = "LogoutFail"
)



func Login(c *gin.Context) {
	var loginQuery LoginQuery
	err := c.ShouldBind(&loginQuery)
	if err != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, utils.BadParameter, nil)
		return
	}
	user := userlogin.User {}
	db := connection.GetDB()
	result := db.Where(&userlogin.User{Number: loginQuery.Number}).First(&user)
	if result.Error != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, LoginFail, nil)
		return
	}
	if loginQuery.Password != user.Password {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, LoginFail, nil)
		return
	}
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, LoginOK, utils.DataWithToken{
		Token: jwt.GenerateToken(&jwt.UserClaims{
			Number: user.Number,
			Name: user.Name,
			Role: jwt.UserRole,
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

	user := userlogin.User {
		Number: registerQuery.Number,
		Password: registerQuery.Password,
		Name: registerQuery.Name,
		Department: registerQuery.Department,
		Major: registerQuery.Major,
		Class: registerQuery.Class,
	}
	db := connection.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, RegisterFail, nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, RegisterOK, nil)
	}
}

func Logout(c *gin.Context) {
	
	utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, LogoutOK, nil)
}

// 获取用户信息
func GetUserInfo(c *gin.Context) {
	userId, exists := c.Get("number")
	if !exists {
		utils.ResponseHandler(c, http.StatusUnauthorized, utils.FailCode, "user获取错误", nil)
		return
	}
	number := userId.(string)

	var user userlogin.User
	err := connection.GetDB().Where(&userlogin.User{Number: number}).First(&user).Error
	if err != nil {
		log.Error(err)
		utils.ResponseHandler(c, http.StatusOK, utils.FailCode, "获取用户信息失败", nil)
	} else {
		utils.ResponseHandler(c, http.StatusOK, utils.SuccessCode, "获取用户信息成功", &UserInfo {
			Number: user.Number,
			Name: user.Name,
			Department: user.Department,
			Major: user.Major,
			Class: user.Class,
		})
	}
}