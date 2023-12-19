package vm

const LFIELDS_PER_FLUSH = 50

type LuaTable interface {
	Get(luaValue) luaValue
	Put(luaValue, luaValue)
	Meta() LuaTable
	SetMeta(LuaTable)
	Len() int
	Next(luaValue) luaValue
	INext(luaValue) luaValue
}

// lua table
// 官方实现将可以转换为integer的float也存放在arr中, 此处简化操作将float全部放map
type luaTable struct {
	arr    []luaValue            //number类型压缩存储
	_map   map[luaValue]luaValue //其他类型map存储
	meta   LuaTable              //元表
	keys   map[luaValue]luaValue //used by next() 使用map将luaTable所有key组成单向链表
	change bool                  //used by next()
}

var _ LuaTable = (*luaTable)(nil)

func GoMapToLuaTable(gomap map[interface{}]interface{}) LuaTable {
	return &luaTable{
		arr:  make([]luaValue, 0),
		_map: gomap,
	}
}

func newLuaTable(nArr, nRec int) *luaTable {
	return &luaTable{
		arr:  make([]luaValue, 0, nArr),
		_map: make(map[luaValue]luaValue, nRec),
	}
}

func (t *luaTable) Meta() LuaTable {
	if t.meta == nil {
		t.meta = newLuaTable(0, 0)
	}
	return t.meta
}
func (t *luaTable) SetMeta(meta LuaTable) {
	t.meta = meta
}

func (t luaTable) Get(key luaValue) luaValue {
	key = _tryToInteger(key) // it will try to convert float to integer
	if idx, ok := key.(int); ok {
		if idx >= 1 && idx <= len(t.arr) {
			return t.arr[idx-1]
		}
	}
	return t._map[key]
}

func (t *luaTable) Put(key, val luaValue) {
	if key == nil {
		panic("table index is nil")
	}
	t.change = true
	key = _tryToInteger(key)
	if idx, ok := key.(int); ok && idx >= 1 {
		arrLen := len(t.arr)
		if idx <= arrLen {
			t.arr[idx-1] = val
			if idx == arrLen && val == nil {
				t._shrinkArray()
			}
		}
		if idx == arrLen+1 {
			delete(t._map, key)
			if val != nil {
				t.arr = append(t.arr, val)
				t._expandArray()
			}
			return
		}
	}
	if val != nil {
		t._map[key] = val
	} else {
		delete(t._map, key)
	}
}

func (t luaTable) Len() int {
	return len(t.arr)
}

func (t *luaTable) Next(key luaValue) luaValue {
	if t.keys == nil || (key == nil && t.change) {
		t._initKeys()
		t.change = false
	}
	return t.keys[key]
}
func (t luaTable) INext(key luaValue) luaValue {
	if key == nil {
		return t.arr[0]
	} else {
		if k, ok := toNumber(key); ok {
			nextKey := int(k) + 1
			if nextKey < len(t.arr) {
				return t.arr[nextKey]
			}
		}
	}
	return nil
}

func (t *luaTable) _initKeys() {
	t.keys = make(map[luaValue]luaValue, len(t.arr)+len(t._map))
	var beforeKey luaValue
	for i, v := range t.arr {
		if v != nil { //是否需要判断 v != nil
			t.keys[beforeKey] = i + 1
			beforeKey = i + 1
		}
	}
	for k, v := range t._map {
		if v != nil {
			t.keys[beforeKey] = k
			beforeKey = k
		}
	}

}

func _tryToInteger(key luaValue) luaValue {
	if x, ok := key.(int); ok {
		return x
	}
	if x, ok := key.(float64); ok && x == float64(int(x)) {
		return int(x)
	}
	return key
}
func (t *luaTable) _shrinkArray() {
	for i := len(t.arr) - 1; i >= 0; i-- {
		if t.arr[i] == nil {
			t.arr = t.arr[0:i]
		}
	}
}
func (t *luaTable) _expandArray() {
	for idx := int(len(t.arr)) + 1; true; idx++ {
		if val, found := t._map[idx]; found {
			delete(t._map, idx)
			t.arr = append(t.arr, val)
		} else {
			break
		}
	}
}
