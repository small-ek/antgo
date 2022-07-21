package test

import (
	"fmt"
	"github.com/small-ek/antgo/container/array"
	"testing"
)

func TestArray(t *testing.T) {
	// 创建并发安全的int类型数组
	a := array.New()

	// 添加数据项
	for i := 0; i < 100; i++ {
		go a.Append(i)
	}
	// 获取当前数组长度
	fmt.Println(a.Len())

	// 获取当前数据项列表
	fmt.Println(a.List())

	// 获取指定索引项
	fmt.Println(a.Get(0))

	// 在指定索引前插入数据项
	a.InsertAfter(9, 8888)
	// 在指定索引后插入数据项

	fmt.Println(a.List())

	// 修改指定索引的数据项
	a.Set(9, 100)
	fmt.Println(a.List())

	// 搜索数据项，返回搜索到的索引位置
	fmt.Println(a.Search(100))

	a.Clear()
	fmt.Println(a.List())
}
