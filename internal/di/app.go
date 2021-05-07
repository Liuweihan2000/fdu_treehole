package di

import (
	"GoProject/fudan_bbs/controller"
	d "GoProject/fudan_bbs/internal/dal"
	"github.com/gin-gonic/gin"
)

type App struct {
	dao  d.DaoInterface
	http *gin.Engine
}

func NewApp(dao d.DaoInterface, h *gin.Engine) (app *App, err error) {
	app = &App{
		dao:  dao,
		http: h,
	}
	controller.DaoInstance = dao
	return
}
