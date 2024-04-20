package utils

import (
	"github.com/gin-gonic/gin"
)

type DataWithToken struct {
	Token string
}

func ResponseHandler(c *gin.Context, httpStatus int, status int,msg string, data interface{}) {
	c.JSON(httpStatus, gin.H{"status" : status, "msg" : msg, "data" : data})
}


type DataList struct {
	Data  interface{}
	Total int `json:"total"`
}