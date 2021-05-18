package app

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

// App is an application components lifecycle manager
type App struct {
	opts   options
	ctx    context.Context
	sigs   []os.Signal
	cancel func()
}

// 创建一个应用来管理程序的生命周期
func New(opts ...Option) *App {
	options := options{
		log: logrus.New(),
	}
	if id, err := uuid.NewUUID(); err == nil {
		options.id = id.String()
	}
	for _, o := range opts {
		o(&options)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		opts:   options,
		ctx:    ctx,
		cancel: cancel,
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
}

// Run executes all OnStart hooks registered with the application's Lifecycle.
func (a *App) Run() error {
	// 使用errgroup来管理应用
	g, ctx := errgroup.WithContext(a.ctx)

	// 启动服务
	for _, s := range a.opts.server {
		g.Go(func() error {
			return s.Start(ctx)
		})
	}

	// 监听
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.sigs...)
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return a.Stop()
			}
		}
	})

	if err := g.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
