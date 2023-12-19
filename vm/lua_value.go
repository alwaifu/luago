package vm

import (
	"fmt"
	"strconv"
)

/* basic types */
const (
	LUA_TNONE = iota - 1 // -1
	LUA_TNIL
	LUA_TBOOLEAN
	LUA_TLIGHTUSERDATA
	LUA_TNUMBER
	LUA_TSTRING
	LUA_TTABLE
	LUA_TFUNCTION
	LUA_TUSERDATA
	LUA_TTHREAD
)

// lua 变量
type luaValue = interface{}

func toString(v luaValue) (string, bool) {
	switch x := v.(type) {
	case string:
		return x, true
	case int, int64, float64:
		return fmt.Sprintf("%v", x), true // TODO: check lua float to string
	default:
		return "", false
	}
}

func toBool(v luaValue) bool {
	switch x := v.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

func toNumber(v luaValue) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// ArgK 获取 OpArgK 类型操作数具体数据
func argK(vm *State, rk int) luaValue {
	if rk > 0xFF {
		return vm.stack.c.proto.Constants[rk&0xFF] // constant
	} else {
		return vm.stack.slots[rk] //register
	}
}

// luaValue to float
func convertToFloat(val luaValue) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case string:
		fv, err := strconv.ParseFloat(v, 64)
		return fv, err == nil
	default:
		return 0, false
	}
}
