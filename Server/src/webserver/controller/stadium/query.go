package stadium

type StadiumListQuery struct {
	Name     string `form:"Name"`
	Categoty string `form:"Category"`
	Page     int `form:"page" binding:"required"`
	PerPage  int `form:"perPage" binding:"required"`
}

type AddStadiumQuery struct {
	Name     string `form:"name" binding:"required"`
	Categoty string `form:"category" binding:"required"`
}

type DeleteStadiumQuery struct {
	ID uint `form:"id" binding:"required"`
}

type UpadateStadiumQuery struct {
	ID   uint `form:"ID" binding:"required"`
	Name string `form:"Name" binding:"required"`
}