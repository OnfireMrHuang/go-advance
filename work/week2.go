package work

import (
	"errors"
	"log"
)

//1. 我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
// 这种错误表示没有查询出结果行，在工程中我一般都是直接handle住，返回成功同时返回空结果集；如果是一些特殊场景一定要返回有效结果集合，这个时候会
// 返回wrap error，但是也一定会备注函数说明。

// sentinel error
var ErrNoRows  = errors.New("sql.ErrNoRows")

func mock() ([]string,error) {
	return []string{}, ErrNoRows
}


func Week2Demo() ([]string,error) {
	result,err := mock()
	if err == ErrNoRows {
		// handle住
		log.Printf("%v",err)
		return []string{},nil
	}
	return result,err
}