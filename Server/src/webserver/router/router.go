package router

import (
	// "net/http"
	"common/config"
	"common/log"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"webserver/controller/admin"
	"webserver/controller/game"
	"webserver/controller/order"
	"webserver/controller/record"
	"webserver/controller/stadium"
	"webserver/controller/user"
	"webserver/jwt"
)

var router *gin.Engine
var one sync.Once

func InitRouter() {
	gConf := config.GetGlobalConfig()
	if gConf == nil {
		panic("globalconfig is nil")
	}

	gin.SetMode(gConf.Mod)

	router = gin.New()
	router.Use(log.LogMiddleware())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     gConf.AllowOrigins,
		AllowMethods:     []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "AccessToken", "X-CSRF-Token", "Authorization", "Token"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	router.Use(jwt.JwtVerify)

	initUserRouter()
	initAdminRouter()
	initStadiunRouter()
	initOrderRouter()
	initRecordRouter()
	initGameRouter()
}

func initUserRouter() {

	userRouter := router.Group("/user")

	userRouter.POST("/login", user.Login)
	userRouter.POST("/register", user.Register)
	userRouter.GET("/userinfo", user.GetUserInfo)
}

func initAdminRouter() {
	adminRouter := router.Group("/admin")

	adminRouter.POST("/login", admin.Login)
	adminRouter.POST("/register", admin.Register)
	adminRouter.GET("/adminlist", admin.GetAdminList)
	adminRouter.POST("/addadmin", admin.AddAdmin)
	adminRouter.POST("/deleteadmin", admin.DeleteAdmin)
	adminRouter.POST("updateadmin", admin.UpdateAdmin)
	adminRouter.GET("/admininfo", admin.GetAdminInfo)
}

func initStadiunRouter() {

	router.GET("/api/stadiumlist", stadium.StadiumList)
	router.POST("/api/addstadium", stadium.AddStadium)
	router.POST("/api/deletestadium", stadium.DeleteStadium)
	router.POST("/api/updatestadium", stadium.UpdateStadium)
}

func initOrderRouter() {
	orderRouter := router.Group("/api/stadium")

	orderRouter.GET("/allowabletimes", order.StadiumAllowableTimes)
	orderRouter.POST("/order", order.Order)
}

func initRecordRouter() {
	recordRouter := router.Group("/api/record")
	recordRouter.GET("/recordlist", record.OrderRecords)
	recordRouter.GET("/outdated/recordlist", record.OutdatedOrderRecords)
	recordRouter.GET("/admin/recordlist", record.AllOrderRecords)
	recordRouter.POST("cancelorder", record.CancelOrder)
	recordRouter.POST("/admin/cancelorder", record.AdminCancelORder)
	recordRouter.GET("/admin/outdated/recordlist", record.AllOutdatedOrderRecords)
}

func initGameRouter() {
	gameRouter := router.Group("/api/game")
	gameRouter.POST("/creategame", game.CreateGame)
	gameRouter.GET("/gamelist", game.GetGameList)
	gameRouter.GET("/outdated/gamelist", game.GetOutdatedGameList)
	gameRouter.GET("/heldgamelist", game.GetHostGameList)
	gameRouter.GET("/outdated/heldgamelist", game.GetOutdatedHostGameList)
	gameRouter.POST("/applygame", game.ApplyGame)
	gameRouter.GET("/applyrecord", game.GetApplyRecord)
	gameRouter.POST("/cancelapply", game.CancelApply)
	gameRouter.GET("/outdated/applyrecord", game.OutdatedApplyRecord)
	gameRouter.GET("/application", game.GetApplication)
	gameRouter.GET("/outdated/application", game.GetOutdatedApplication)
	gameRouter.POST("/accept", game.AcceptApply)
	gameRouter.POST("/refuse", game.RefuseApply)
	gameRouter.POST("/cancelgame", game.CancelGame)
}

func Start() {
	one.Do(func() {
		router.Run()

	})

}
