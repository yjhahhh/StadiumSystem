package admin

type LoginQuery struct {
	Name       string `form:"name" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type RegisterQuery struct {
	Name        string `form:"name" binding:"required"`
	Password  string `form:"password" binding:"required"`
	Level     uint `form:"level" binding:"required"`
}

type AdminListQuery struct {
	Page    int `form:"page" binding:"required"`
	PerPage int `form:"perPage" binding:"required"`
	Name    string `form:"Name"`
	Level   uint `form:"Level"`
}

type DeleteAdminQuery struct {
	ID uint `form:"ID" binding:"required"`
}

type UpdateAdminQuery struct {
	ID    uint `form:"ID" binding:"required"`
	Level uint `form:"Level" binding:"required"`
}