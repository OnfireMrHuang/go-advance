package biz

import (
	"context"
)

// DO定义
type Node struct {
	Name      string
	Path      string
	ApiKey    string
	ApiSecret string
}

type IndicatorItem struct {
	Module string
	Object string
}

// repo接口定义
type UserNode interface {
	AllNode(ctx context.Context) ([]*Node, error)
	Node(ctx context.Context, name string) (*Node, error)
}

type Indicators interface {
	AllIndicatorItem(ctx context.Context) ([]*IndicatorItem, error)
}
