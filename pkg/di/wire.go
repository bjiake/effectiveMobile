//go:build wireinject
// +build wireinject

package di

import (
	http "effectiveMobile/pkg/api"
	"effectiveMobile/pkg/api/handler"
	"effectiveMobile/pkg/config"
	"effectiveMobile/pkg/db"
	"effectiveMobile/pkg/service"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectToBD, service.NewService, handler.NewHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
