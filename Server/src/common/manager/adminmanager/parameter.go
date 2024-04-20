package adminmanager

type AddAdminParameter struct {
	Name     string
	Password string
	Level    uint
}

type DeleteAdminParameter struct {
	ID uint
}

type UpdateAdminParameter struct {
	ID    uint
	Level uint
}

type AdminListParameter struct {
	Name    string
	Level   uint
	Page    int
	PerPage int
}

type Admin struct {
	ID    uint
	Name  string
	Level int
}

type Info struct {
	Name  string
	Level uint
}