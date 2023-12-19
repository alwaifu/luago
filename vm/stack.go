package vm

// lua 虚拟机栈帧
type stackFrame struct {
	vm      *State      // 引用state便于获取全局变量等信息
	slots   []luaValue  // 数据栈 底层存储(编译期可以直接确定局部变量需要使用的寄存器数量，子函数返回值临时占用需额外检查)
	top     int         // 下一个可以使用的数据栈位置(栈顶)
	prev    *stackFrame // 上一个函数栈帧
	c       *closure    // 闭包
	varargs []luaValue
	pc      int                 // 程序计数器
	openuvs map[int]updateValue // open状态的upvalue, key为寄存器索引

	results []luaValue //函数执行完的返回值
}

// newStackFrame new a stack frame
func newStackFrame(vm *State, slots []luaValue, prev *stackFrame, c *closure, vargs []luaValue) *stackFrame {
	return &stackFrame{
		vm:      vm,
		slots:   slots,
		top:     len(slots),
		prev:    prev,
		c:       c,
		varargs: vargs,
		pc:      0,
		openuvs: make(map[int]updateValue),
		results: nil,
	}
}

// check and prepare 'want'th slots,
func (stack *stackFrame) check(want int) {
	length := len(stack.slots)
	if want >= length {
		stack.slots = append(stack.slots, make([]luaValue, 1+want-length)...)
	}
}

type updateValue struct {
	val *luaValue
}
