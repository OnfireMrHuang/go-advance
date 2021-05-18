package indicator

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/http"
	"net/url"
)

const DataFlowName = "dataflow"

// 实体对象
type DataFlow struct {
	log    *logrus.Logger
	client *http.Caller
}

func NewDataFlow(log *logrus.Logger) *DataFlow {
	return &DataFlow{
		log:    log,
		client: http.NewCaller(),
	}
}

func (d *DataFlow) Pull(ctx context.Context, node *biz.Node, object string, day string) (*Data, error) {
	switch object {
	case ExecStatObject:
		return d.pullExecStat(ctx, node, day)
	}
	return nil, errors.New("不支持的指标内容")
}

func (d *DataFlow) pullExecStat(ctx context.Context, node *biz.Node, day string) (*Data, error) {
	params := url.Values{}
	params.Add("module", DataFlowName)
	params.Add("object", ExecStatObject)
	params.Add("day", day)
	token := ""

	d.log.Debugf("path %v module %v object %v\n", node.Path, DataFlowName, ExecStatObject)

	resp, err := d.client.Get(ctx, node.Path, params, token)
	if err != nil {
		return nil, err
	}
	response := struct {
		Code int        `json:"code"`
		Msg  string     `json:"msg"`
		Data []ExecStat `json:"data"`
	}{}
	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, errors.New("pull response code " + cast.ToString(response.Code) + " msg " + response.Msg)
	}
	// 将结果转换为[][]string类型
	var data Data
	for _, stat := range response.Data {
		var columns []string
		columns = append(columns, cast.ToString(stat.EnvCode))
		columns = append(columns, cast.ToString(stat.EnvName))
		columns = append(columns, cast.ToString(stat.ProjectCode))
		columns = append(columns, cast.ToString(stat.ProjectName))
		columns = append(columns, cast.ToString(stat.Num))
		columns = append(columns, cast.ToString(stat.FailNum))
		columns = append(columns, cast.ToString(stat.Count))
		columns = append(columns, cast.ToString(stat.FailCount))
		columns = append(columns, stat.Begin)
		columns = append(columns, stat.End)
		data.Data = append(data.Data, columns)
	}
	data.Module = DataFlowName
	data.Object = ExecStatObject
	return &data, nil
}
