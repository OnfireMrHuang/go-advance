package file

import (
	"context"
	"fmt"
	"mingyuanyun.com/mic/bigdata-opcenter/internal/biz/indicator"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/file"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/oss"
	"testing"
)

func TestOssFile(t *testing.T) {
	ossClient, err := oss.NewOSSClient("https://oss-cn-shenzhen.aliyuncs.com",
		"",
		"")
	if err != nil {
		t.Fatalf("创建oss客户端时出错 %v ", err)
	}
	data := indicator.Data{
		Data: [][]string{
			{
				"test1",
				"test2",
				"test3",
				"test4",
			},
			{
				"demo1",
				"demo2",
				"demo3",
				"demo4",
			},
		},
		Module: "dataflow",
		Object: "exec_stat",
	}
	csvFilePath := fmt.Sprintf("%s/%s_%s.csv", "/tmp/csv", "dataflow", "exec_stat")
	// 生成本地文件
	err = file.WriteFile(&data.Data, csvFilePath)
	if err != nil {
		t.Fatalf("写指标数据到文件失败 %v \n", err)
	}
	compressedFileName := fmt.Sprintf("%s_%s.zip", data.Module, data.Object)
	compressedFilePath := fmt.Sprintf("%s/%s_%s.zip", "/tmp/zip", data.Module, data.Object)
	// 压缩文件
	err = file.Compress(csvFilePath, compressedFilePath)
	if err != nil {
		t.Fatalf("压缩文件失败 %v \n", err)
	}
	// 上传文件
	err = ossClient.PubObjectFromFile(context.Background(), "dmp-test", compressedFilePath, compressedFileName)
	if err != nil {
		t.Fatalf("上传文件时出错 %v \n", err)
	}
}
