package stadiummanager

type StadiumParameter struct {
	Name     string
	Category string
	Start string
	End   string
}

type StadiumListParameter struct {
	Name     string
	Category string
	Page     int
	PerPage  int
}

type UpdateParameter struct {
	ID    uint
	Name  string
	Start string
	End   string
}
