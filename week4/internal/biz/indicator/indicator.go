package indicator

import (
	"context"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz"
)

type Data struct {
	Data [][]string
	// 其他元信息...
	Module string
	Object string
}

// 定义一个指标接口，各个模块实现各自己的指标
type Indicator interface {
	Pull(ctx context.Context, node *biz.Node, object string, day string) (*Data, error)
}
