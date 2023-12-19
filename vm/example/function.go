package example

import "luago/vm/api"

// ExampleGoFunc 一个可以在lua虚拟机调用的golang函数实现示例
func ExampleGoFuncAdd(state api.State, args ...interface{}) []interface{} {
	a, b := args[0].(int), args[1].(int)
	return []interface{}{a + b}
}

var _ api.GoFunc = api.GoFunc(ExampleGoFuncAdd)
