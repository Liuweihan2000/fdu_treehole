//+build wireinject

package di

import (
	"GoProject/fudan_bbs/controller"
	"GoProject/fudan_bbs/internal/dao"
	"github.com/google/wire"
)

//go:generate wire
func InitApp() (*App, error) {
	panic(wire.Build(dao.Provider, controller.NewServer, NewApp))
}
