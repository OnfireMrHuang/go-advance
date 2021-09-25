package file

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"mingyuanyun.com/mic/bigdata-opcenter/pkg/zip"
	"os"
)

func WriteFile(data *[][]string, filePath string) error {
	// 首先判断参数
	if data == nil {
		return errors.New("压缩数据为空")
	}
	if filePath == "" {
		return errors.New("压缩文件路径为空")
	}
	f, err := file(filePath)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	err = w.WriteAll(*data)
	if err != nil {
		return errors.Wrap(err, "write file error")
	}
	w.Flush()
	return nil
}

func Compress(fin, fout string) error {
	err := zip.GzipFile(fin, fout)
	if err != nil {
		return errors.Wrap(err, "compress zip file fail")
	}
	return nil
}

func file(path string) (*os.File, error) {
	var f *os.File
	var err error
	exist := fileExist(path)
	if !exist {
		//err = os.MkdirAll(path, os.ModePerm)
		//if err != nil {
		//	return nil, errors.Wrap(err, "create file dir error")
		//}
		f, err = os.Create(path)
		if err != nil {
			return nil, errors.Wrap(err, "create file error")
		}
		return f, nil
	}
	f, err = os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file error")
	}
	return f, nil
}

func fileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}
