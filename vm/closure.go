package vm

import (
	"luago/chunk"
	"luago/vm/api"
)

// closure lua闭包结构 分为两类: lua函数 go函数
type closure struct {
	proto  *chunk.Prototype // lua函数原型
	goFunc api.GoFunc       // go函数

	upvals []updateValue // 构造closure时确定upvalue
}

func newClosure(proto *chunk.Prototype) *closure {
	return &closure{
		proto:  proto,
		upvals: make([]updateValue, len(proto.Upvalues)),
	}
}
func newGoClosure(f api.GoFunc) *closure {
	return &closure{
		goFunc: f, // TODO: 设置upvalue
	}
}

// callClosure, in 'vm' state call closure 'c' with params 'args'
func _callClosure(vm *State, c *closure, args ...luaValue) []luaValue {
	if c.proto != nil {
		// lua 函数
		f := c.proto
		// fmt.Printf("call %s<%d,%d>\n", f.Source, f.LineDefined, f.LastLineDefined)
		nParams := f.NumParams         //函数固定参数个数
		varargs := []luaValue{}        //新栈帧的变长参数
		top := f.MaxStackSize          //数据栈顶
		slots := make([]luaValue, top) //数据栈(寄存器) //XXX:初始slots比top大一些是否会有性能提升

		if len(args) > int(nParams) && f.IsVararg == 1 {
			varargs = args[nParams:] // 传入参数多于固定参数 记录到vararg
		}
		for j := 0; j < int(nParams) && j < len(args); j++ {
			slots[j] = args[j] // 固定参数放在寄存器列表开头
		}
		// 构造新的函数栈帧
		subStack := newStackFrame(vm, slots, vm.stack, c, varargs)

		vm.pushStack(subStack) // 压栈
		vm.Run()               // 运行新栈帧
		vm.popStack()          // 出栈

		return subStack.results
	} else {
		// go 函数
		f := c.goFunc
		return f(vm, args...)
	}
}
