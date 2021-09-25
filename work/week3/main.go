package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	group,ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
			_, _ = fmt.Fprintln(writer, "hello week3!")
		})
		server := http.Server{
			Addr: ":8080",
			Handler: mux,
		}
		go func() {
			select {
				// 不管是收到信号还是监听失败，这里都会收到消息
				case <-ctx.Done():
					_ = server.Shutdown(context.Background())
			}
		}()
		return server.ListenAndServe()
	})
	group.Go(func() error {
		c := make(chan os.Signal)
		signal.Notify(c,syscall.SIGQUIT,syscall.SIGINT)
		//阻塞直至有信号传入
		s := <-c
		fmt.Println("get signal:", s)
		return errors.New("收到停止信号")
	})
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}

