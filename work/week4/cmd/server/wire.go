// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz/collector"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/conf"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/data"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/server"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/service"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/app"
)

// initApp init kratos application.
func initApp(*conf.Bootstrap, context.Context, *logrus.Logger) (*app.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, collector.ProviderSet, data.ProviderSet, newCronApp))
}
