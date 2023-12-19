package vm

import (
	"fmt"
	"luago/chunk"
	"luago/vm/api"
)

// state lua虚拟机运行期各种状态
type State struct {
	stack   *stackFrame //虚拟机栈
	opcodes [LEN_OPCODE]opcode

	global LuaTable              //全局变量
	meta   map[luaValue]luaValue //元表
}

var _ api.State = api.State(&State{})

// NewState new state
func NewState() *State {
	// 构造lua vm
	vm := &State{
		opcodes: opcodes,
		global:  newLuaTable(0, 0),
		meta:    make(map[luaValue]luaValue),
	}
	OpenLibs(vm) //注册基础库函数
	return vm
}

// load state load binary chunk
func (vm *State) Load(proto *chunk.Prototype) {
	// 构造主函数
	slots := make([]luaValue, proto.MaxStackSize, proto.MaxStackSize+20)
	mainStack := newStackFrame(vm, slots, nil, newClosure(proto), nil)
	// 初始化全局变量
	if len(proto.Upvalues) > 0 {
		var g luaValue = vm.global
		mainStack.c.upvals[0] = updateValue{&g}
	}
	for i, uvInfo := range proto.Upvalues {
		_ = i
		_ = uvInfo
	}
	vm.stack = mainStack
}

// Resister 实现golang函数注册到lua虚拟机
func (vm *State) Register(name string, f api.GoFunc) {
	vm.global.Put(name, newGoClosure(f))
}

// Fetch 获取下一条指令
func (vm *State) Fetch() Instruction {
	stack := vm.stack
	i := stack.c.proto.Code[stack.pc]
	stack.pc++
	return Instruction(i)
}

// pushStack 压入函数栈
func (vm *State) pushStack(stack *stackFrame) {
	stack.prev = vm.stack
	vm.stack = stack
}

// popStack 弹出函数栈
func (vm *State) popStack() *stackFrame {
	stack := vm.stack
	vm.stack = stack.prev
	// stack.prev = nil
	return stack
}

// metaField 从元表获取元数据
func (vm *State) metaField(val, field luaValue) luaValue {
	if t, ok := val.(LuaTable); ok {
		return t.Meta().Get(field)
	}
	return vm.meta[field]
}
func (vm *State) callMetaMethod(mmName string, a, b luaValue) (luaValue, bool) {
	if m := vm.metaField(a, mmName); m != nil {
		if mm, ok := m.(*closure); ok {
			// vm.Push(mm)
			return _callClosure(vm, mm, a, b)[0], true
		}
	}
	return nil, false
}
func (vm *State) getTable(t, field luaValue) (luaValue, bool) {
	if tb, ok := t.(LuaTable); ok {
		if v := tb.Get(field); v != nil {
			return v, true
		}
	}
	mf := vm.metaField(t, "__index")
	if mf != nil {
		switch x := mf.(type) {
		case luaTable:
			return x.Get(field), true //此处不会递归调用元方法
		case *closure:
			// vm.Push(x)
			return _callClosure(vm, x, t, field)[0], true
		}
	}
	return nil, false
}
func (vm *State) setTable(t, field, value luaValue) bool {
	mf := vm.metaField(t, "__newindex")
	if tb, ok := t.(LuaTable); ok && mf == nil {
		tb.Put(field, value)
		return true
	}
	if mf != nil {
		switch x := mf.(type) {
		case luaTable:
			return vm.setTable(x, field, value)
		case *closure:
			// vm.Push(x)
			_callClosure(vm, x, t, field, value)
			return true
		}
	}
	return false
}

func (vm *State) _stackLevel() int {
	level := 0
	curStack := vm.stack
	for curStack != nil {
		level++
		curStack = curStack.prev
	}
	return level
}
func (vm *State) _printStackLevel() {
	level := vm._stackLevel()
	fmt.Print(level)
	for i := 0; i < level; i++ {
		fmt.Print(" ")
	}
}

func (vm *State) Run() {
	for {
		// defer func() {
		// 	if p := recover(); p != nil {
		// 		f := vm.stack.c.proto
		// 		fmt.Printf("%v, %s<%d,%d>pc:%d, \n", p, f.Source, f.LineDefined, f.LastLineDefined, vm.stack.pc)
		// 		panic(p)
		// 	}
		// }()
		// vm._printStackLevel()
		// fmt.Printf("run process count: %d\n", vm.stack.pc)
		inst := vm.Fetch()
		code := inst.Opcode()
		vm.opcodes[code].action(inst, vm)
		if code == OP_RETURN {
			break
		}
	}
}

// CallByParam for go call lua function
func (vm *State) CallByParam(name string, args ...interface{}) ([]interface{}, error) {
	if _f := vm.global.Get(name); _f == nil {
		return []interface{}{}, fmt.Errorf("not a function: %s", name)
	} else {
		var results []luaValue
		callOk := false
		if f, ok := _f.(*closure); ok {
			results = _callClosure(vm, f, args...)
			callOk = true
		} else {
			if mf := vm.metaField(_f, "__call"); mf != nil {
				if f, ok := mf.(*closure); ok {
					results = _callClosure(vm, f, append([]luaValue{_f}, args...)...)
					callOk = true
				}
			}
		}
		if !callOk {
			return results, fmt.Errorf("attempt to call a %T value", _f)
		} else {
			return results, nil
		}
	}
}
