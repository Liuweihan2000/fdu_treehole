//+build wireinject

package di

import (
	"GoProject/fudan_bbs/controller"
	"GoProject/fudan_bbs/internal/dal"
	"github.com/google/wire"
)

//go:generate wire
func InitApp() (*App, error) {
	panic(wire.Build(dal.Provider, controller.NewServer, NewApp))
}
