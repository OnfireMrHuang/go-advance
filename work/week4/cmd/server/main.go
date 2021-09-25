package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"gitlab.mypaas.com.cn/dmp/gopkg/logging"
	"gopkg.in/ini.v1"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/conf"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/server"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/app"
)

var (
	flagConf     string
	flagLogLevel string
)

func init() {
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf conf.ini")
	flag.StringVar(&flagLogLevel, "log_level", "error", "log level, eg: -log_level debug")
}

func newCronApp(logger *logrus.Logger, cs *server.CronServer) *app.App {
	return app.New(
		app.Name("crontab"),
		app.Logger(logger),
		app.Server(
			cs,
		),
	)
}

func main() {
	flag.Parse()
	log := logging.NewLogger(flagLogLevel)

	var bc conf.Bootstrap
	err := ini.MapTo(&bc, flagConf+"/config.ini")
	if err != nil {
		panic(err)
	}
	instance, cleanup, err := initApp(&bc, context.Background(), log)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// 运行app
	if err := instance.Run(); err != nil {
		panic(err)
	}
}
