package gamemanager

type CreateParameter struct {
	Number    string // 学号
	UserName  string
	Title     string // 比赛标题
	Date      string
	StadiumID uint
	Stadium   string // 体育场馆名称
	Start     string
	End       string
	Remark    string
	Maximum   uint // 最大人数
}


type CancelParameter struct {
	GameID       uint
	BookRecordID uint
}

type ApplyParameter struct {
	Number   string // 学号
	UserName string
	GameID   uint
}

type CancelApplyParameter struct {
	RecordID uint
}

type AcceptApplyParameter struct {
	RecordID uint
}

type RefuseApplyParameter struct {
	RecordID uint
}

type GameListParameter struct {
	Number   string
	UserName string
	Date     string
	Stadium  string
	Title    string
	CanApply bool
	Applyer  string
	Page     int
	PerPage  int
}

type ApplyRecordParameter struct {
	Number     string
	UserName   string
	Date       string
	Stadium    string
	Title      string
	HostNumber string
	Page       int
	PerPage    int
}
