package main

import (
	// "os"
	"os"
	"ytb/src/controller"
	post "ytb/src/controller/api/post"
	api "ytb/src/controller/api/user"
	myJwt "ytb/src/controller/jwt"
	"ytb/src/controller/user"
	mdw "ytb/src/middleware"

	"github.com/astaxie/beego/orm"
	// "github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// var  client *redis.Client

func init() {

	// register model
	orm.RegisterDriver("mysql", orm.DRMySQL)

	// set default database
	err := orm.RegisterDataBase("default", "mysql", "root:123456@/golang?charset=utf8", 30)

	if err != nil {
		glog.Fatal("Failed to connect database:", err)
	}

	// create table
	orm.RunSyncdb("default", false, true)

}

func main() {

	os.Setenv("JWT_SECRET", "LOOPS")
	os.Setenv("REFRESH_SECRET", "LOOPS_REFRESH")

	server := echo.New()

	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		MaxAge:       86400,
	}))
	isValidToken := mdw.IsValidToken

	server.Static("/static", "public")

	server.GET("/", controller.Home, isValidToken)
	//api routes
	group := server.Group("/api", isValidToken)

	group.Use(middleware.Logger())

	group.GET("/user/:id", api.APIGetUserById) //Get user by id

	group.GET("/user/follower/:id", api.APIGetFollowUserById) //Get follower by id

	group.GET("/user/following/:id", api.APIGetFollowUserById) //Get following by id

	group.POST("/user/follow/:id", api.FollowHandler) //Follow user

	group.POST("/user/un-follow/:id", api.UnFollowHandler) //Unfollow user

	group.GET("/user/search&q=:q", api.APISearchUser) //Search user

	group.POST("/user/change_pass/:id", api.ChangePass) //Change password

	group.POST("/user/update/:id", api.APIUpdateUserById) //Update user

	group.GET("/post/:id", post.GetPostByIdPost) //Get post by id

	group.GET("/post/user/:uid", post.GetPostByUid) //Get post by user id

	group.POST("/post/create", post.CreatePost) //Create post

	server.POST("/token/refresh", myJwt.RefreshToken) //Refresh token

	server.POST("/user/register", user.Register) //Register user

	server.POST("/user/login", user.Login) //Login user

	server.Logger.Fatal(server.Start(":5500"))
}
