package collector

import (
	"context"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz/indicator"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/data"
	"sync"
)

// 业务逻辑的组装层，类似 DDD 的 domain 层，data 类似 DDD 的 repo，repo 接口在这里定义，使用依赖倒置的原则。
var ProviderSet = wire.NewSet(NewCollector)

// 聚合根对象（特殊的实体）
type Collector struct {
	log  *logrus.Logger
	data *data.Indicator
}

// 单例
var c *Collector
var once sync.Once

func NewCollector(log *logrus.Logger, data *data.Indicator) *Collector {
	once.Do(func() {
		c = &Collector{
			log:  log,
			data: data,
		}
		// 初始化指标模块的加载
		initIndicatorList(log)

	})
	return c
}

// 管理各个指标模块的列表
var indicatorList map[string]indicator.Indicator

// 初始化指标模块的列表
func initIndicatorList(log *logrus.Logger) {
	indicatorList = map[string]indicator.Indicator{
		indicator.DataFlowName: indicator.NewDataFlow(log),
	}
}

// 通过模块名称获取模块
func getIndicator(name string) indicator.Indicator {
	instance, ok := indicatorList[name]
	if !ok {
		return nil
	}
	return instance
}

func (c *Collector) CollectIndicators(ctx context.Context, day string) ([]*indicator.Data, error) {
	var indicatordatas []*indicator.Data
	nodes, err := c.data.AllNode(ctx)
	if err != nil {
		return nil, err
	}
	indicators, err := c.data.AllIndicatorItem(ctx)
	if err != nil {
		return nil, err
	}
	// 遍历客户节点
	for _, node := range nodes {
		// 遍历指标
		for _, item := range indicators {
			indicator2 := getIndicator(item.Module)
			if indicator2 == nil {
				c.log.Warnf("没有找到%s模块,跳过\n", item.Module)
				continue
			}
			indicatorData, err := indicator2.Pull(ctx, node, item.Object, day)
			if err != nil {
				return nil, errors.WithMessagef(err, "pull node %s 's module:%s object:%s fail",
					node.Name, item.Module, item.Object)
			}
			indicatordatas = append(indicatordatas, indicatorData)
		}
	}
	return indicatordatas, nil
}
