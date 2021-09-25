package service

import (
	"context"
	"fmt"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz/collector"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz/indicator"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/conf"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/file"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/oss"
	"sync"
	"time"
)

// 实现了 api 定义的服务层，类似 DDD 的 application 层，处理 DTO 到 biz 领域实体的转换(DTO -> DO)，同时协同各类 biz 交互，但是不应处理复杂逻辑
var ProviderSet = wire.NewSet(NewService)

type Service struct {
	conf *conf.Bootstrap
	c    *collector.Collector
	log  *logrus.Logger
}

// NewGreeterService new a greeter service.
func NewService(conf *conf.Bootstrap, logger *logrus.Logger, c *collector.Collector) *Service {
	return &Service{conf: conf, c: c, log: logger}
}

func (s *Service) Run(ctx context.Context) error {
	// 获取当天时间
	curDay := time.Now().Format("2006-01-02")
	indicatorDatas, err := s.c.CollectIndicators(ctx, curDay)
	if err != nil {
		return errors.WithMessagef(err, "收集指标时出错")
	}
	// 初始化oss链接
	ossClient, err := oss.NewOSSClient(s.conf.Oss.Endpoint, s.conf.Oss.AccessKey, s.conf.Oss.AccessSecret)
	if err != nil {
		return errors.WithMessagef(err, "创建oss客户端时出错")
	}
	// 开启多个协程处理数据，只是提高性能，整个函数还是同步函数
	// 注意这里协程的数量，如果数量太多后续需要优化为分批处理
	wg := sync.WaitGroup{}
	for _, indicatorData := range indicatorDatas {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := s.handleIndicatorData(ctx, indicatorData, ossClient)
			if err != nil {
				s.log.Errorf("(%)指标获取后在处理时出错: %v",
					indicatorData.Module+"_"+indicatorData.Object, err)
			}

		}()
	}
	wg.Wait()
	return nil
}

func (s *Service) handleIndicatorData(ctx context.Context, data *indicator.Data, ossClient *oss.OSS) error {

	csvFilePath := fmt.Sprintf("%s/%s_%s.csv", s.conf.LocalDir.Csv, data.Module, data.Object)
	// 生成本地文件
	err := file.WriteFile(&data.Data, csvFilePath)
	if err != nil {
		return errors.WithMessage(err, "写指标数据到文件失败")
	}
	compressedFileName := fmt.Sprintf("%s_%s.zip", data.Module, data.Object)
	compressedFilePath := fmt.Sprintf("%s/%s_%s.zip", s.conf.LocalDir.Zip, data.Module, data.Object)
	// 压缩文件
	err = file.Compress(csvFilePath, compressedFilePath)
	if err != nil {
		return errors.WithMessage(err, "压缩文件失败")
	}
	// 上传文件
	err = ossClient.PubObjectFromFile(ctx, s.conf.Oss.Bucket, compressedFilePath, compressedFileName)
	if err != nil {
		return errors.WithMessage(err, "上传文件时出错")
	}
	return nil
}
