package di

import (
	"GoProject/fudan_bbs/controller"
	d "GoProject/fudan_bbs/internal/dao"
	"github.com/gin-gonic/gin"
)

type App struct {
	dao  d.Dao
	http *gin.Engine
}

func NewApp(dao d.Dao, h *gin.Engine) (app *App, err error) {
	app = &App{
		dao:  dao,
		http: h,
	}
	controller.DaoInstance = dao
	return
}
