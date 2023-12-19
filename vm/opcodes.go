package vm

import (
	"fmt"
)

/* OpCode */
const (
	OP_MOVE = iota
	OP_LOADK
	OP_LOADKX
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_MOD
	OP_POW
	OP_DIV
	OP_IDIV
	OP_BAND
	OP_BOR
	OP_BXOR
	OP_SHL
	OP_SHR
	OP_UNM
	OP_BNOT
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
	LEN_OPCODE // opcodes length
)

type opcode struct {
	testFlag byte // operator is a test (next instruction must be a jump)
	setAFlag byte // instruction set register A
	argBMode byte // B arg mode
	argCMode byte // C arg mode
	opMode   byte // op mode
	name     string
	action   func(i Instruction, vm *State)
}

var opcodes = [LEN_OPCODE]opcode{
	/*T A    B       C     mode         name       action */
	{0, 1, OpArgR, OpArgN, IABC /* */, "MOVE    ", opMove},     // R(A) := R(B)
	{0, 1, OpArgK, OpArgN, IABx /* */, "LOADK   ", opLoadK},    // R(A) := Kst(Bx)
	{0, 1, OpArgN, OpArgN, IABx /* */, "LOADKX  ", opLoadKx},   // R(A) := Kst(extra arg)
	{0, 1, OpArgU, OpArgU, IABC /* */, "LOADBOOL", opLoadBool}, // R(A) := (bool)B; if (C) pc++
	{0, 1, OpArgU, OpArgN, IABC /* */, "LOADNIL ", opLoadNil},  // R(A), R(A+1), ..., R(A+B) := nil
	{0, 1, OpArgU, OpArgN, IABC /* */, "GETUPVAL", opGetUpval}, // R(A) := UpValue[B]
	{0, 1, OpArgU, OpArgK, IABC /* */, "GETTABUP", opGetTabup}, // R(A) := UpValue[B][RK(C)]
	{0, 1, OpArgR, OpArgK, IABC /* */, "GETTABLE", opGetTable}, // R(A) := R(B)[RK(C)]
	{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABUP", opSetTabup}, // UpValue[A][RK(B)] := RK(C)
	{0, 0, OpArgU, OpArgN, IABC /* */, "SETUPVAL", opSetUpval}, // UpValue[B] := R(A)
	{0, 0, OpArgK, OpArgK, IABC /* */, "SETTABLE", opSetTable}, // R(A)[RK(B)] := RK(C)
	{0, 1, OpArgU, OpArgU, IABC /* */, "NEWTABLE", opNewTable}, // R(A) := {} (size = B,C)
	{0, 1, OpArgR, OpArgK, IABC /* */, "SELF    ", opSelf},     // R(A+1) := R(B); R(A) := R(B)[RK(C)]
	{0, 1, OpArgK, OpArgK, IABC /* */, "ADD     ", opAdd},      // R(A) := RK(B) + RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "SUB     ", opSub},      // R(A) := RK(B) - RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "MUL     ", opMul},      // R(A) := RK(B) * RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "MOD     ", opMod},      // R(A) := RK(B) % RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "POW     ", opPow},      // R(A) := RK(B) ^ RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "DIV     ", opDiv},      // R(A) := RK(B) / RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "IDIV    ", opIdiv},     // R(A) := RK(B) // RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "BAND    ", opBand},     // R(A) := RK(B) & RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "BOR     ", opBor},      // R(A) := RK(B) | RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "BXOR    ", opBxor},     // R(A) := RK(B) ~ RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "SHL     ", opShl},      // R(A) := RK(B) << RK(C)
	{0, 1, OpArgK, OpArgK, IABC /* */, "SHR     ", opShr},      // R(A) := RK(B) >> RK(C)
	{0, 1, OpArgR, OpArgN, IABC /* */, "UNM     ", opUnm},      // R(A) := -R(B)
	{0, 1, OpArgR, OpArgN, IABC /* */, "BNOT    ", opBnot},     // R(A) := ~R(B)
	{0, 1, OpArgR, OpArgN, IABC /* */, "NOT     ", opNot},      // R(A) := not R(B)
	{0, 1, OpArgR, OpArgN, IABC /* */, "LEN     ", opLen},      // R(A) := length of R(B)
	{0, 1, OpArgR, OpArgR, IABC /* */, "CONCAT  ", opConcat},   // R(A) := R(B).. ... ..R(C)
	{0, 0, OpArgR, OpArgN, IAsBx /**/, "JMP     ", opJmp},      // pc+=sBx; if (A) close all upvalues >= R(A - 1)
	{1, 0, OpArgK, OpArgK, IABC /* */, "EQ      ", opEq},       // if ((RK(B) == RK(C)) ~= A) then pc++
	{1, 0, OpArgK, OpArgK, IABC /* */, "LT      ", opLt},       // if ((RK(B) <  RK(C)) ~= A) then pc++
	{1, 0, OpArgK, OpArgK, IABC /* */, "LE      ", opLe},       // if ((RK(B) <= RK(C)) ~= A) then pc++
	{1, 0, OpArgN, OpArgU, IABC /* */, "TEST    ", opTest},     // if not (R(A) <=> C) then pc++
	{1, 1, OpArgR, OpArgU, IABC /* */, "TESTSET ", opTestSet},  // if (R(B) <=> C) then R(A) := R(B) else pc++
	{0, 1, OpArgU, OpArgU, IABC /* */, "CALL    ", opCall},     // R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
	{0, 1, OpArgU, OpArgU, IABC /* */, "TAILCALL", opTailcall}, // return R(A)(R(A+1), ... ,R(A+B-1))
	{0, 0, OpArgU, OpArgN, IABC /* */, "RETURN  ", opReturn},   // return R(A), ... ,R(A+B-2)
	{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORLOOP ", opForLoop},  // R(A)+=R(A+2); if R(A) <?= R(A+1) then { pc+=sBx; R(A+3)=R(A) }
	{0, 1, OpArgR, OpArgN, IAsBx /**/, "FORPREP ", opForPrep},  // R(A)-=R(A+2); pc+=sBx
	{0, 0, OpArgN, OpArgU, IABC /* */, "TFORCALL", opTforCall}, // R(A+3), ... ,R(A+2+C) := R(A)(R(A+1), R(A+2));
	{0, 1, OpArgR, OpArgN, IAsBx /**/, "TFORLOOP", opTforLoop}, // if R(A+1) ~= nil then { R(A)=R(A+1); pc += sBx }
	{0, 0, OpArgU, OpArgU, IABC /* */, "SETLIST ", opSetList},  // R(A)[(C-1)*FPF+i] := R(A+i), 1 <= i <= B
	{0, 1, OpArgU, OpArgN, IABx /* */, "CLOSURE ", opClosure},  // R(A) := closure(KPROTO[Bx])
	{0, 1, OpArgU, OpArgN, IABC /* */, "VARARG  ", opVararg},   // R(A), R(A+1), ..., R(A+B-2) = vararg
	{0, 0, OpArgU, OpArgU, IAx /*  */, "EXTRAARG", opExtraArg}, // extra (larger) argument for previous opcode
}

// move 寄存器b的值拷贝到寄存器a
func opMove(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	vm.stack.slots[a] = vm.stack.slots[b]
}

// loadK 加载常量表bx处数据到寄存器a (bx只有18bit 最多加载262143个常量情况 否则需要使用loadKx函数)
func opLoadK(i Instruction, vm *State) {
	a, bx := i.ABx()
	vm.stack.slots[a] = vm.stack.c.proto.Constants[bx]
}

// loadKx 扩展loadK函数加载下一个EXTRAARG(iAx模式)指令中ax(26bit 最多加载67108864个常量)常量到寄存器a
func opLoadKx(i Instruction, vm *State) {
	a, _ := i.ABx()
	ax := vm.Fetch().Ax()
	vm.stack.slots[a] = vm.stack.c.proto.Constants[ax]
}

// loadBool 寄存器a设置bool值(b非零为true) 如果c非零则跳过下一条指令
func opLoadBool(i Instruction, vm *State) {
	a, b, c := i.ABC()
	vm.stack.slots[a] = b != 0
	if c != 0 {
		vm.stack.pc++
	}
}

// loadNil 寄存器a开始设置b个nil
func opLoadNil(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	for j := a + b; j >= a; j-- {
		if j < len(vm.stack.slots) {
			vm.stack.slots[j] = nil
		} else {
			vm.stack.slots = append(vm.stack.slots, nil)
		}
	}
}

// getUpval 寄存器a 设置为 upvalue b
func opGetUpval(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	vm.stack.slots[a] = *vm.stack.c.upvals[b].val
}

func opGetTabup(i Instruction, vm *State) {
	a, b, c := i.ABC()
	k := argK(vm, c)
	if t, ok := (*vm.stack.c.upvals[b].val).(LuaTable); ok {
		vm.stack.slots[a] = t.Get(k) //XXX:考虑是否需要调用元方法
	} else {
		panic(fmt.Sprintf("op getTabup r[%d] not table", b))
	}
}
func opGetTable(i Instruction, vm *State) {
	a, b, c := i.ABC()
	k := argK(vm, c)
	if v, ok := vm.getTable(vm.stack.slots[b], k); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Sprintf("op getTable r[%d] not table", b))
	}
}
func opSetTabup(i Instruction, vm *State) {
	a, b, c := i.ABC()
	k := argK(vm, b)
	v := argK(vm, c)
	if t, ok := (*vm.stack.c.upvals[a].val).(LuaTable); ok {
		t.Put(k, v) //XXX:考虑是否需要调用元方法
	} else {
		panic(fmt.Sprintf("op setTabup r[%d] not table", a))
	}
}

// setUpval 将寄存器a设置到 upvalue b
func opSetUpval(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	*vm.stack.c.upvals[b].val = vm.stack.slots[a]
}
func opSetTable(i Instruction, vm *State) {
	a, b, c := i.ABC()
	k := argK(vm, b)
	v := argK(vm, c)
	if ok := vm.setTable(vm.stack.slots[a], k, v); !ok {
		panic(fmt.Sprintf("op setTable r[%d] not table", a))
	}
}
func opNewTable(i Instruction, vm *State) {
	a, b, c := i.ABC()
	vm.stack.slots[a] = newLuaTable(Fb2int(b), Fb2int(c))
}
func opSelf(i Instruction, vm *State) {
	a, b, c := i.ABC()
	self := vm.stack.slots[b]
	vm.stack.slots[a+1] = self
	vm.stack.slots[a], _ = vm.getTable(self, argK(vm, c))
}

func opAdd(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.add(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to add a '%T' with a '%T'", a, b))
	}
}

func opSub(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.sub(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to sub a '%T' with a '%T'", a, b))
	}
}

func opMul(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.mul(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to mul a '%T' with a '%T'", a, b))
	}
}

func opMod(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.mod(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to mod a '%T' with a '%T'", a, b))
	}
}

func opPow(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.pow(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to pow a '%T' with a '%T'", a, b))
	}
}

func opDiv(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.div(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to div a '%T' with a '%T'", a, b))
	}

}

func opIdiv(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if v, ok := vm.idiv(argK(vm, b), argK(vm, c)); ok {
		vm.stack.slots[a] = v
	} else {
		panic(fmt.Errorf("attempt to idiv a '%T' with a '%T'", a, b))
	}
}

func opBand(i Instruction, vm *State) {
	a, b, c := i.ABC()
	x := argK(vm, b)
	y := argK(vm, c)

	if ix, ok := x.(int); ok {
		if iy, ok := y.(int); ok {
			vm.stack.slots[a] = ix & iy //integer
			return
		}
	}

	panic(fmt.Sprintf("op band %T %T", x, y))
}

func opBor(i Instruction, vm *State) {
	a, b, c := i.ABC()
	x := argK(vm, b)
	y := argK(vm, c)

	if ix, ok := x.(int); ok {
		if iy, ok := y.(int); ok {
			vm.stack.slots[a] = ix | iy //integer
			return
		}
	}

	panic(fmt.Sprintf("op band %T %T", x, y))
}

func opBxor(i Instruction, vm *State) {
	a, b, c := i.ABC()
	x := argK(vm, b)
	y := argK(vm, c)

	if ix, ok := x.(int); ok {
		if iy, ok := y.(int); ok {
			vm.stack.slots[a] = ix ^ iy //integer
			return
		}
	}

	panic(fmt.Sprintf("op bxor %T %T", x, y))
}

func opShl(i Instruction, vm *State) {
	a, b, c := i.ABC()
	x := argK(vm, b)
	y := argK(vm, c)

	if ix, ok := x.(int); ok {
		if iy, ok := y.(int); ok {
			vm.stack.slots[a] = vm.shiftL(ix, iy) //integer
			return
		}
	}

	panic(fmt.Sprintf("op shl %T %T", x, y))
}

func opShr(i Instruction, vm *State) {
	a, b, c := i.ABC()
	x := argK(vm, b)
	y := argK(vm, c)

	if ix, ok := x.(int); ok {
		if iy, ok := y.(int); ok {
			vm.stack.slots[a] = vm.shiftR(ix, iy) //integer
			return
		}
	}

	panic(fmt.Sprintf("op shr %T %T", x, y))
}

func opUnm(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	x := argK(vm, b)

	if ix, ok := x.(int); ok {
		vm.stack.slots[a] = -ix //integer
		return
	}
	if fx, ok := convertToFloat(x); ok {
		vm.stack.slots[a] = -fx //float
		return
	}
	panic(fmt.Sprintf("op unm %T ", x))
}

func opBnot(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	x := argK(vm, b)

	if ix, ok := x.(int); ok {
		vm.stack.slots[a] = ^ix //integer
		return
	}
	panic(fmt.Sprintf("op bnot %T ", x))
}

func opNot(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	vm.stack.slots[a] = toBool(vm.stack.slots[b])
}

func opLen(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	val := vm.stack.slots[b]
	if s, ok := val.(string); ok {
		vm.stack.slots[a] = len(s)
	} else if v, ok := vm.callMetaMethod("__len", val, nil); ok {
		vm.stack.slots[a] = v
	} else if t, ok := val.(LuaTable); ok {
		vm.stack.slots[a] = t.Len()
	} else {
		panic(fmt.Sprintf("attempt to get length of a %T value", val))
	}
}

func opConcat(i Instruction, vm *State) {
	a, b, c := i.ABC()
	s := "" // TODO:字符串拼接是否可优化
	for ; b <= c; b++ {
		if v, ok := toString(vm.stack.slots[b]); ok {
			s += v
		} else if v, ok := vm.callMetaMethod("__concat", s, vm.stack.slots[b]); ok {
			s += v.(string) //panic if v is not string
		} else {
			panic(fmt.Sprintf("op concat %T", vm.stack.slots[b]))
		}
	}
	vm.stack.slots[a] = s
}

// jmp 程序计数器增加sbx
func opJmp(i Instruction, vm *State) {
	a, sBx := i.AsBx()
	vm.stack.pc += sBx
	if a != 0 { //close 大于等于 a-1 的 upvalue
		for i := range vm.stack.openuvs {
			if i >= a-1 {
				delete(vm.stack.openuvs, i)
			}
		}
	}
}

func opEq(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if vm.compareEq(argK(vm, b), argK(vm, c)) != (a != 0) {
		vm.stack.pc++
	}
}

func opLt(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if lt, err := vm.compareLt(argK(vm, b), argK(vm, c)); err != nil {
		panic(err)
	} else if lt != (a != 0) {
		vm.stack.pc++
	}
}

func opLe(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if le, err := vm.compareLe(argK(vm, b), argK(vm, c)); err != nil {
		panic(err)
	} else if le != (a != 0) {
		vm.stack.pc++
	}
}

func opTest(i Instruction, vm *State) {
	a, _, c := i.ABC()
	if toBool(vm.stack.slots[a]) != (c != 0) {
		vm.stack.pc++
	}
}

func opTestSet(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if toBool(vm.stack.slots[b]) == (c != 0) {
		vm.stack.slots[a] = vm.stack.slots[b]
	} else {
		vm.stack.pc++
	}
}

func opCall(i Instruction, vm *State) {
	a, b, c := i.ABC()
	_f := vm.stack.slots[a]
	var args []luaValue
	if b > 0 {
		args = vm.stack.slots[a+1 : a+b]
	} else {
		args = vm.stack.slots[a+1 : vm.stack.top]
	}
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
		panic(fmt.Sprintf("attempt to call a %T value", _f))
	}

	if c == 1 { // no results
		vm.stack.top = a
	} else if c > 1 { // 部分返回
		vm.stack.check(a + c) //实际只需要(a+c-2)
		for j := a; j < a+c-1; j++ {
			if j-a < len(results) {
				vm.stack.slots[j] = results[j-a]
			} else {
				vm.stack.slots[j] = nil
			}
		}
		vm.stack.top = a + c - 1
	} else { // c == 0 全部返回
		vm.stack.check(a + len(results)) //实际需要(a+len(results)-1)
		for j := a; j < a+len(results); j++ {
			vm.stack.slots[j] = results[j-a]
		}
		vm.stack.top = a + len(results)
	}
}
func opTailcall(i Instruction, vm *State) {
	opCall(i, vm) // TODO: 优化尾递归调用
}
func opReturn(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	if b == 1 { //no return
	} else if b > 1 {
		vm.stack.results = vm.stack.slots[a : a+b-1] //return (b-1) values
	} else {
		vm.stack.results = vm.stack.slots[a:vm.stack.top] //all return
	}
}
func opForLoop(i Instruction, vm *State) {
	a, sBx := i.AsBx()

	//a -> index 循环开始
	limit := vm.stack.slots[a+1] //a+1 -> limit 循环结束
	step := vm.stack.slots[a+2]  //a+2 -> step  循环步长
	//a+3 -> i     循环变量

	// index += step
	if ra, ok := vm.add(vm.stack.slots[a], step); ok {
		vm.stack.slots[a] = ra
	} else {
		panic(fmt.Errorf("op forloop, r[a] + r[a+2], attemp to add '%T' with '%T'", vm.stack.slots[a], step))
	}

	if _step, ok := toNumber(step); ok && _step >= 0 {
		if le, err := vm.compareLe(vm.stack.slots[a], limit); err != nil {
			panic("op forloop, r[a] <= r[a+1], " + err.Error())
		} else if le {
			vm.stack.pc += sBx
			vm.stack.slots[a+3] = vm.stack.slots[a] //change loop local variable
		}
	} else if ok && _step < 0 {
		if le, err := vm.compareLe(limit, vm.stack.slots[a]); err != nil {
			panic("op forloop, r[a+1] <= r[a], " + err.Error())
		} else if le {
			vm.stack.pc += sBx
			vm.stack.slots[a+3] = vm.stack.slots[a] //change loop local variable
		}
	} else {
		panic(fmt.Sprintf("bad 'for' step (number expected, got %T)", step))
	}
}
func opForPrep(i Instruction, vm *State) {
	a, sBx := i.AsBx()
	//a -> index 循环开始
	//a+1 -> limit 循环结束
	//a+2 -> step  循环步长
	//a+3 -> i     循环变量
	if ra, ok := vm.sub(vm.stack.slots[a], vm.stack.slots[a+2]); ok {
		vm.stack.slots[a] = ra
	} else {
		panic(fmt.Errorf("op forloop, r[a] - r[a+2], attemp to sub '%T' with '%T'", vm.stack.slots[a], vm.stack.slots[a+2]))
	}
	vm.stack.pc += sBx
}
func opTforCall(i Instruction, vm *State) {
	// R(A+3),...,R(A+2+C) := R(A)(R(A+1),R(A+2))
	a, _, c := i.ABC()
	_f := vm.stack.slots[a].(*closure)
	results := _callClosure(vm, _f, vm.stack.slots[a+1], vm.stack.slots[a+2])
	for i, rLen := a+3, len(results); i <= a+2+c; i++ {
		j := i - (a + 3)
		if j < rLen {
			vm.stack.slots[a+3+j] = results[j]
		} else {
			vm.stack.slots[a+3+j] = nil
		}
	}
}
func opTforLoop(i Instruction, vm *State) {
	// if R(A+1) ~= nil then { R(A)=R(A+1); pc+=sBx }
	a, sBx := i.AsBx()
	if vm.stack.slots[a+1] != nil {
		vm.stack.slots[a] = vm.stack.slots[a+1]
		vm.stack.pc += sBx
	}
}
func opSetList(i Instruction, vm *State) {
	a, b, c := i.ABC()
	if t, ok := vm.stack.slots[a].(LuaTable); ok {
		if c > 0 {
			c = c - 1
		} else {
			c = vm.Fetch().Ax()
		}
		if b == 0 { //从寄存器a开始所有数据
			b = vm.stack.top - a
		}
		idx := c * LFIELDS_PER_FLUSH
		for j := 1; j <= b; j++ {
			idx++
			t.Put(idx, vm.stack.slots[a+j])
		}
	} else {
		panic(fmt.Sprintf("op setList r[%d] not table", a))
	}
}
func opClosure(i Instruction, vm *State) {
	a, bx := i.ABx()
	subProto := vm.stack.c.proto.Protos[bx]
	c := newClosure(subProto)
	vm.stack.slots[a] = c
	// 更新upvalues
	for i, uv := range subProto.Upvalues {
		if uv.Instack != 0 { // 捕获当前函数的局部变量
			if openuv, ok := vm.stack.openuvs[int(uv.Idx)]; ok { //前面已经有opClosure把该寄存器放入openuv
				c.upvals[i] = openuv
			} else { //首次将寄存器i放入openuv
				c.upvals[i] = updateValue{&vm.stack.slots[uv.Idx]}
				vm.stack.openuvs[int(uv.Idx)] = c.upvals[i]
			}
		} else { // 父函数upval引用
			c.upvals[i] = vm.stack.c.upvals[uv.Idx]
		}
	}
}
func opVararg(i Instruction, vm *State) {
	a, b, _ := i.ABC()
	if b > 1 {
		vm.stack.check(a + b) //is it necessary check
		for j := a; j <= a+b-2; j++ {
			if j < len(vm.stack.varargs) {
				vm.stack.slots[j] = vm.stack.varargs[j-a]
			} else {
				vm.stack.slots[j] = nil
			}
		}
		vm.stack.top = a + b - 1
	} else {
		vm.stack.check(a + len(vm.stack.varargs)) //is it necessary check
		for j, v := range vm.stack.varargs {
			vm.stack.slots[a+j] = v
		}
		vm.stack.top = a + len(vm.stack.varargs)
	}
}
func opExtraArg(i Instruction, vm *State) { panic("TODO:") } //TODO:to implement op extraArg
