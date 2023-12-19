package vm

import (
	"fmt"
	"luago/vm/api"
)

var baseFuncs = map[string]api.GoFunc{
	"print":        basePrint,
	"assert":       baseAssert,
	"error":        baseError,
	"select":       baseSelect,
	"ipairs":       baseIpairs,
	"pairs":        basePairs,
	"next":         baseNext,
	"load":         baseLoad,
	"loadfile":     baseLoadfile,
	"dofile":       baseDofile,
	"pcall":        basePcall,
	"xpcall":       baseXpcall,
	"getmetatable": baseGetmetatable,
	"setmetatable": baseSetmetatable,
	"rawequal":     baseRawequal,
	"rawlen":       baseRawlen,
	"rawget":       baseRawget,
	"rawset":       baseRawset,
	"type":         baseType,
	"tostring":     baseTostring,
	"tonumber":     baseTonumber,
}
var tabFuncs = map[string]api.GoFunc{
	// "move":   nil,
	"insert": tabInsert,
	// "remove": nil,
	// "sort":   nil,
	// "concat": nil,
	// "pack":   nil,
	// "unpack": nil,
}

// OpenLibs 注册基础库函数到lua虚拟机
func OpenLibs(vm api.State) {
	for name, f := range baseFuncs {
		vm.Register(name, f)
	}
	for name, f := range tabFuncs {
		vm.Register(name, f)
	}
}

func basePrint(_ api.State, args ...interface{}) []interface{} {
	fmt.Println(args...)
	return nil
}
func baseAssert(_ api.State, args ...interface{}) []interface{} {
	msg := "input:1: assertion failed!"
	if len(args) > 1 {
		msg = fmt.Sprintf("%v", args[1])
	}
	if v := args[0]; v != nil {
		return []interface{}{v}
	} else {
		panic(msg)
	}
}
func baseError(_ api.State, args ...interface{}) []interface{}  { return nil } //TODO:std basefunc
func baseSelect(_ api.State, args ...interface{}) []interface{} { return nil } //TODO:std basefunc
func baseIpairs(_ api.State, args ...interface{}) []interface{} {
	return []interface{}{newGoClosure(_iNext), args[0], 0}
}
func basePairs(_ api.State, args ...interface{}) []interface{} {
	return []interface{}{newGoClosure(baseNext), args[0], nil}
}
func baseNext(_ api.State, args ...interface{}) []interface{} {
	_t, key := args[0], args[1] //argument #1 #2
	if t, ok := _t.(LuaTable); ok {
		if nextKey := t.Next(key); nextKey == nil {
			return []interface{}{nil}
		} else {
			return []interface{}{nextKey, t.Get(nextKey)}
		}
	} else {
		panic(fmt.Errorf("input:2: bad argument #1 to 'next' (table expected, got %T)", _t))
	}
}                                                                     //TODO:std basefunc
func baseLoad(_ api.State, args ...interface{}) []interface{}         { return nil } //TODO:std basefunc
func baseLoadfile(_ api.State, args ...interface{}) []interface{}     { return nil } //TODO:std basefunc
func baseDofile(_ api.State, args ...interface{}) []interface{}       { return nil } //TODO:std basefunc
func basePcall(_ api.State, args ...interface{}) []interface{}        { return nil } //TODO:std basefunc
func baseXpcall(_ api.State, args ...interface{}) []interface{}       { return nil } //TODO:std basefunc
func baseGetmetatable(_ api.State, args ...interface{}) []interface{} { return nil } //TODO:std basefunc
func baseSetmetatable(_ api.State, args ...interface{}) []interface{} {
	arg1, arg2 := args[0], args[1]
	if meta, ok := arg2.(LuaTable); ok {
		if t, ok := arg1.(LuaTable); ok {
			t.SetMeta(meta)
		} else {
			panic("TODO: 非LuaTable添加元方法")
		}
	} else {
		panic(fmt.Sprintf("input:2: bad argument #2 to 'setmetatable' (nil or table expected, got %T)", arg2))
	}
	return []interface{}{}
}
func baseRawequal(_ api.State, args ...interface{}) []interface{} { return nil } //TODO:std basefunc
func baseRawlen(_ api.State, args ...interface{}) []interface{}   { return nil } //TODO:std basefunc
func baseRawget(_ api.State, args ...interface{}) []interface{}   { return nil } //TODO:std basefunc
func baseRawset(_ api.State, args ...interface{}) []interface{}   { return nil } //TODO:std basefunc
func baseType(_ api.State, args ...interface{}) []interface{}     { return nil } //TODO:std basefunc
func baseTostring(_ api.State, args ...interface{}) []interface{} { return nil } //TODO:std basefunc
func baseTonumber(_ api.State, args ...interface{}) []interface{} { return nil } //TODO:std basefunc

func _iNext(_ api.State, args ...interface{}) []interface{} {
	_t, key := args[0], args[1] //argument #1 #2
	if t, ok := _t.(LuaTable); ok {
		if nextKey := t.INext(key); nextKey == nil {
			return []interface{}{nil}
		} else {
			return []interface{}{nextKey, t.Get(nextKey)}
		}
	} else {
		panic(fmt.Errorf("input:2: bad argument #1 to 'next' (table expected, got %v)", _t))
	}
}

func tabInsert(_ api.State, args ...interface{}) []interface{} {
	panic("TODO:")
}
