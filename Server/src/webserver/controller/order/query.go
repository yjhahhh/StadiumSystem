package order


type OrderQuery struct {
	ID     uint `form:"ID" binding:"required"`
	Name   string `form:"Name" binding:"required"`
	Date   string `form:"Date" binding:"required"`
    Start  string `form:"Start" binding:"required"`
    End    string `form:"End" binding:"required"`
}

type IdleTimeQuery struct {
	ID      uint `form:"ID" binding:"required"`
	Page    int `form:"page" binding:"required"`
	PerPage int `form:"perPage" binding:"required"`
}