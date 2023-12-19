package vm

import (
	"fmt"
	"math"
)

// add 加
func (vm *State) add(a, b luaValue) (luaValue, bool) {
	if x, ok := a.(int); ok {
		if y, ok := b.(int); ok {
			return x + y, true
		}
	}
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x + y, true
		}
	}
	if v, ok := vm.callMetaMethod("__add", a, b); ok {
		return v, true
	}
	return nil, false
}

// sub 减
func (vm *State) sub(a, b luaValue) (luaValue, bool) {
	if x, ok := a.(int); ok {
		if y, ok := b.(int); ok {
			return x - y, true
		}
	}
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x - y, true
		}
	}
	if v, ok := vm.callMetaMethod("__sub", a, b); ok {
		return v, true
	}
	// return nil, fmt.Errorf("attempt to sub a '%T' with a '%T'", a, b)
	return nil, false
}

// mul 乘
func (vm *State) mul(a, b luaValue) (luaValue, bool) {
	if x, ok := a.(int); ok {
		if y, ok := b.(int); ok {
			return x * y, true
		}
	}
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x * y, true
		}
	}
	if v, ok := vm.callMetaMethod("__mul", a, b); ok {
		return v, true
	}
	// return nil, fmt.Errorf("attemp to mul a '%T' with a '%T'", a, b)
	return nil, false
}

// div 除
func (vm *State) div(a, b luaValue) (luaValue, bool) {
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x / y, true
		}
	}
	if v, ok := vm.callMetaMethod("__div", a, b); ok {
		return v, true
	}
	// return nil, fmt.Errorf("attemp to div a '%T' with a '%T'", a, b)
	return nil, false
}

// idiv 整除
func (vm *State) idiv(a, b luaValue) (luaValue, bool) {
	if x, ok := a.(int); ok {
		if y, ok := b.(int); ok {
			return _iFloorDiv(x, y), true
		}
	}
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return math.Floor(x / y), true
		}
	}
	if v, ok := vm.callMetaMethod("__idiv", a, b); ok {
		return v, true
	}
	// return nil, fmt.Errorf("attemp to idiv a '%T' with a '%T'", a, b)
	return nil, false
}

// mod 取模
func (vm *State) mod(a, b luaValue) (luaValue, bool) {
	if x, ok := a.(int); ok {
		if y, ok := b.(int); ok {
			return _iMod(x, y), true
		}
	}
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return x - math.Floor(x/y)*y, true
		}
	}
	if v, ok := vm.callMetaMethod("__mod", a, b); ok {
		return v, true
	}
	// return nil, fmt.Errorf("attemp to mod a '%T' with a '%T'", a, b)
	return nil, false
}

// pow 乘方
func (vm *State) pow(a, b luaValue) (luaValue, bool) {
	if x, ok := convertToFloat(a); ok {
		if y, ok := convertToFloat(b); ok {
			return math.Pow(x, y), true
		}
	}
	if v, ok := vm.callMetaMethod("__pow", a, b); ok {
		return v, true
	}
	// return nil, fmt.Errorf("attemp to pow a '%T' with a '%T'", a, b)
	return nil, false
}

// lua integer整除运算单独实现(golang整除运算是直接截断而非向下取整)
// 5//3  -> lua:1, golang:1
// -5//3 -> lua:-2 golang:-1
func _iFloorDiv(a, b int) int {
	if a > 0 && b > 0 || a < 0 && b < 0 || a%b == 0 {
		return a / b
	} else {
		return a/b - 1
	}
}

// lua integer取模运算需要单独实现
// 5%3  -> lua:2 golang:2
// -5%3 -> lua:1 golang:-2
func _iMod(a, b int) int {
	return a - _iFloorDiv(a, b)*b
}

// lua 左移运算需要单独实现
func (vm *State) shiftL(a, n int) int {
	if n >= 0 {
		return a << n
	} else {
		return vm.shiftR(a, -n)
	}
}

// lua 右移运算需要单独实现
// -1>>63 -> lua:1 golang:-1
func (vm *State) shiftR(a, n int) int {
	if n >= 0 {
		return int(uint(a) >> n)
	} else {
		return vm.shiftL(a, -n)
	}
}

func (vm *State) compareEq(a, b luaValue) bool {
	switch x := a.(type) {
	case nil:
		return b == nil
	case bool:
		y, ok := b.(bool)
		return ok && x == y
	case string:
		y, ok := b.(string)
		return ok && x == y
	case int:
		switch y := b.(type) {
		case int:
			return x == y
		case float64:
			return float64(x) == y
		default:
			return false
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x == y
		case int:
			return x == float64(y)
		default:
			return false
		}
	default:
		if a == b {
			return true
		} else {
			if v, ok := vm.callMetaMethod("__eq", a, b); ok {
				return toBool(v)
			} else {
				return false
			}
		}
	}
}

func (vm *State) compareLt(a, b luaValue) (bool, error) {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x < y, nil
		} else {
			return false, fmt.Errorf("attempt to compare string with %T", b)
		}
	case int:
		switch y := b.(type) {
		case int:
			return x < y, nil
		case float64:
			return float64(x) < y, nil
		default:
			return false, fmt.Errorf("attempt to compare number with %T", b)
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x < y, nil
		case int:
			return x < float64(y), nil
		default:
			return false, fmt.Errorf("attempt to compare number with %T", b)
		}
	default:
		if v, ok := vm.callMetaMethod("__lt", a, b); ok {
			return toBool(v), nil
		} else {
			return false, fmt.Errorf("attempt to compare %T with %T", a, b)
		}
	}
}

func (vm *State) compareLe(a, b luaValue) (bool, error) {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y, nil
		} else {
			return false, fmt.Errorf("attempt to compare string with %T", b)
		}
	case int:
		switch y := b.(type) {
		case int:
			return x <= y, nil
		case float64:
			return float64(x) <= y, nil
		default:
			return false, fmt.Errorf("attempt to compare number with %T", b)
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y, nil
		case int:
			return x <= float64(y), nil
		default:
			return false, fmt.Errorf("attempt to compare number with %T", b)
		}
	default:
		if v, ok := vm.callMetaMethod("__le", a, b); ok {
			return toBool(v), nil
		} else {
			return false, fmt.Errorf("attempt to compare %T with %T", a, b)
		}
	}

}
