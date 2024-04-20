package ordermanager

type IdleTimes struct {
	Date  string
	Start string
	End   string
}

type IdleTimeQuery struct {
	StadiumID uint
	Page      int
	PerPage   int
}