package game

// 查询比赛活动列表
type GameListQuery struct {
	UserName string `form:"username"`
	Title    string `form:"title"`
	Date     string `form:"date"`
	Stadium  string `form:"stadium"`
	CanApply bool `form:"canApply"`
	Page     int `form:"page" binding:"required"`
	PerPage  int `form:"perPage" binding:"required"`
}

// 查询申请记录列表
type RecordListQuery struct {
	UserName string `form:"username"`
	Title    string `form:"title"`
	Date     string `form:"date"`
	Stadium  string `form:"stadium"`
	Page     int `form:"page" binding:"required"`
	PerPage  int `form:"perPage" binding:"required"`
}

// 举办比赛活动
type CreateGameQuery struct {
	Title     string `form:"title" binding:"required"`
	Date      string `form:"date" binding:"required"`
	Stadium   string `form:"stadium" binding:"required"`
	StadiumID uint `form:"stadiumID" binding:"required"`
	Start     string `form:"start" binding:"required"`
	End       string `form:"end" binding:"required"`
	Remark    string `form:"remark"`
	Maximum   uint `form:"maximum" binding:"required"`
}

// 报名比赛
type ApplyGameQuery struct {
	ID uint `form:"ID" binding:"required"`
}

// 取消报名
type CancelApplyQuery struct {
	ID uint `form:"ID" binding:"required"`
}

// 申请记录
type ApplicationQuery struct {
	UserName string `form:"username"`
	Title    string `form:"title"`
	Date     string `form:"date"`
	Page     int `form:"page" binding:"required"`
	PerPage  int `form:"perPage" binding:"required"`
}

// 接受报名
type AcceptApplyQuery struct {
	ID uint `form:"ID" binding:"required"`
}

// 拒绝报名
type RefuseApplyQuery struct {
	ID uint `form:"ID" binding:"required"`
}

// 取消比赛
type CancelGameQuery struct {
	ID           uint `form:"ID" binding:"required"`
	BookReocrdID uint `form:"bookRecordID" binding:"required"`
}