package record

type RecordQuery struct {
	Date    string `form:"Date"`
	Stadium string `form:"Stadium"`
	Page    int `form:"page" binding:"required"`
	PerPage int `form:"perPage" binding:"required"`
}

type AllRecordQuery struct {
	Number  string `form:"UserNo"`
	Date    string `form:"Date"`
	Stadium string `form:"Stadium"`
	Page    int `form:"page" binding:"required"`
	PerPage int `form:"perPage" binding:"required"`
}

type CancleOrder struct {
	OrderID uint `form:"id" binding:"required"`
}
