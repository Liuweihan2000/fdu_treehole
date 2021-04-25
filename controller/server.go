package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

type Server struct {
	e *gin.Engine
}

func NewServer() (e *gin.Engine, err error) {
	// 启动定时任务
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ = range ticker.C {
			DaoInstance.LoadThreadsToRedis()
		}
	}()

	engine := gin.Default()
	InitRouter(engine)
	go func() {
		if e := engine.Run(":9999"); e != nil {
			fmt.Println(e)
		}
	}()
	return
}

func InitRouter(e *gin.Engine) {
	// 主页面
	e.GET("/", index)

	// 错误处理页面
	e.GET("/err", err)

	// 注销
	e.GET("/logout", logout)

	// 注册页面
	g1 := e.Group("/sign_up")
	{
		g1.GET("/", signUp)
		g1.POST("/", signUpAction)
	}

	// 登录页面
	g2 := e.Group("/login")
	{
		g2.GET("/", login)
		g2.POST("/", loginAction)
	}

	// 密码重置页面
	g3 := e.Group("/reset")
	{
		g3.GET("/", reset)
		g3.POST("/", resetAction)
	}

	// 设置管理员页面
	//g4 := e.Group("/set_admin")
	//{
	//	g4.GET("/", setAdmin)
	//	g4.POST("/", setAdminAction)
	//}

	g5 := e.Group("/threads")
	{
		g5.GET("/create", createThread)
		g5.POST("/create", createThreadAction)
		g5.GET("/read", readThread)
		g5.GET("/search", searchThread)
		g5.POST("/read", followThreadAction)
		g5.GET("read_follow", readFollowThread)
		// g5.DELETE("/delete", deleteThreadAction)
	}

	g6 := e.Group("/posts")
	{
		g6.POST("/create", createPostAction)
		// g6.DELETE("/delete", deletePostAction)
	}

}
