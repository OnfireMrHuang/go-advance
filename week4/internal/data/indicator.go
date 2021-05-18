package data

import (
	"context"
	"github.com/sirupsen/logrus"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz"
)

type Indicator struct {
	data *Data
	log  *logrus.Logger
}

func NewIndicator(log *logrus.Logger, data *Data) *Indicator {
	return &Indicator{
		data: data,
		log:  log,
	}
}

// PO定义
type Node struct {
	Name      string `json:"name" db:"name"`
	Path      string `json:"path" db:"path"`
	ApiKey    string `json:"api_key" db:"api_key"`
	ApiSecret string `json:"api_secret" db:"api_secret"`
}

type IndicatorItem struct {
	Module string `json:"module" db:"module"`
	Object string `json:"object" db:"object"`
}

func (c *Indicator) AllNode(ctx context.Context) ([]*biz.Node, error) {
	db, err := c.data.conn(ctx)
	if err != nil {
		return nil, err
	}
	sess := db.GetSession()
	querySql := "select name,path,api_key,api_secret from t_node"
	var nodes []Node
	err = sess.SelectContext(ctx, &nodes, querySql)
	if err != nil {
		return nil, err
	}
	var result []*biz.Node
	for _, node := range nodes {
		result = append(result, &biz.Node{
			Name:      node.Name,
			Path:      node.Path,
			ApiKey:    node.ApiKey,
			ApiSecret: node.ApiSecret,
		})
	}
	return result, nil
}

func (c *Indicator) Node(ctx context.Context, name string) (*biz.Node, error) {
	db, err := c.data.conn(ctx)
	if err != nil {
		return nil, err
	}
	sess := db.GetSession()
	querySql := "select name,path,api_key,api_secret from t_node where name=?"
	var node Node
	err = sess.SelectContext(ctx, &node, querySql, name)
	if err != nil {
		return nil, err
	}
	result := biz.Node{
		Name:      node.Name,
		Path:      node.Path,
		ApiKey:    node.ApiKey,
		ApiSecret: node.ApiSecret,
	}
	return &result, nil
}

func (c *Indicator) AllIndicatorItem(ctx context.Context) ([]*biz.IndicatorItem, error) {
	db, err := c.data.conn(ctx)
	if err != nil {
		return nil, err
	}
	sess := db.GetSession()
	querySql := "select module,object from t_indicators"
	var items []IndicatorItem
	err = sess.SelectContext(ctx, &items, querySql)
	if err != nil {
		return nil, err
	}
	var result []*biz.IndicatorItem
	for _, item := range items {
		result = append(result, &biz.IndicatorItem{
			Module: item.Module,
			Object: item.Object,
		})
	}
	return result, nil
}
