package api

import "luago/chunk"

// lua调用的go函数
type GoFunc func(State, ...interface{}) []interface{}

// 暴露给外部使用的虚拟机状态
type State interface {
	// Pop() interface{} //Pop a value from lua vm stack
	// Push(interface{})        //Push a value to lua vm stack
	Load(proto *chunk.Prototype)
	Register(string, GoFunc) //Register a Go function to lua vm
	Run()
	CallByParam(funcName string, args ...interface{}) ([]interface{}, error) //Call global function
}
