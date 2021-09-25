package server

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/conf"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/data"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/service"
)

type CronServer struct {
	c    *conf.Bootstrap
	s    *service.Service
	data *data.Data
	log  *logrus.Logger
}

func NewCronServer(c *conf.Bootstrap, service *service.Service, data *data.Data, log *logrus.Logger) *CronServer {
	return &CronServer{
		c:    c,
		s:    service,
		data: data,
		log:  log,
	}
}

func (s *CronServer) Start(ctx context.Context) error {
	tasks, err := s.data.AllCrontab(ctx)
	if err != nil {
		return err
	}
	c := newWithSecond()
	for _, task := range tasks {
		if !task.IsEnable {
			continue
		}
		if task.Name == "collector" {
			_, err = c.AddFunc(task.Schedule, func() {
				s.log.Infof("start exec task %s \n", task.Name)
				err := s.s.Run(context.Background())
				if err != nil {
					s.log.Errorf("cron scheduler execute error: %v", err)
				}
			})
		}
	}
	c.Start()
	<-ctx.Done() // 阻塞等待上下文取消
	return nil
}

func newWithSecond() *cron.Cron {
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	return cron.New(cron.WithParser(secondParser), cron.WithChain())
}
