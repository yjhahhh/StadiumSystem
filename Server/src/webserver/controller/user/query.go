package user

type LoginQuery struct {
	Number   string `form:"number" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterQuery struct {
	Number     string `form:"number" binding:"required"`
	Password   string `form:"password" binding:"required"`
	Name       string `form:"name" binding:"required"`
	Department string `form:"department" binding:"required"`
	Major      string `form:"major" binding:"required"`
	Class      uint `form:"class" binding:"required"`
}

type UserInfo struct {
	Number     string
	Name       string
	Department string
	Major      string
	Class      uint
}