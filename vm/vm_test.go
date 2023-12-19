package vm

import (
	"fmt"
	"luago/chunk"
	"os"
	"os/exec"
	"testing"
)

var helloworld = []byte{
	0x1B, 0x4C, 0x75, 0x61, 0x53, 0x00, 0x19, 0x93, 0x0D, 0x0A, 0x1A, 0x0A, 0x04, 0x08, 0x04, 0x08, // .LuaS...........
	0x08, 0x78, 0x56, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x28, 0x77, // .xV...........(w
	0x40, 0x01, 0x11, 0x40, 0x68, 0x65, 0x6C, 0x6C, 0x6F, 0x5F, 0x77, 0x6F, 0x72, 0x6C, 0x64, 0x2E, // @..@hello_world.
	0x6C, 0x75, 0x61, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x04, 0x00, // lua.............
	0x00, 0x00, 0x06, 0x00, 0x40, 0x00, 0x41, 0x40, 0x00, 0x00, 0x24, 0x40, 0x00, 0x01, 0x26, 0x00, // ....@.A@..$@..&.
	0x80, 0x00, 0x02, 0x00, 0x00, 0x00, 0x04, 0x06, 0x70, 0x72, 0x69, 0x6E, 0x74, 0x04, 0x0E, 0x48, // ........print..H
	0x65, 0x6C, 0x6C, 0x6F, 0x2C, 0x20, 0x57, 0x6F, 0x72, 0x6C, 0x64, 0x21, 0x01, 0x00, 0x00, 0x00, // ello, World!....
	0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, // ................
	0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, // ................
	0x00, 0x00, 0x05, 0x5F, 0x45, 0x4E, 0x56, // ..._ENV
}

func TestLuaVM(t *testing.T) {
	proto := chunk.Undump(helloworld)
	vm := NewState()
	vm.Load(proto)
	vm.Run()
}

func TestBinaryChunk(t *testing.T) {
	proto := chunk.Undump(helloworld)
	list(proto)

}

func list(f *chunk.Prototype) {
	chunk.PrintHeader(f)
	printCode(f)
	chunk.PrintDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

func printCode(f *chunk.Prototype) {
	fmt.Printf("\tpc\tline\topname\ta\tb\tc\n")
	for pc, c := range f.Code {
		i := Instruction(c)

		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}

		operands := ""
		switch i.OpMode() {
		case IABC:
			a, b, c := i.ABC()
			operands = fmt.Sprintf("%d\t%d\t%d", a, b, c)
		case IABx:
			a, b := i.ABx()
			operands = fmt.Sprintf("%d\t%d\t", a, b)
		case IAsBx:
			a, b := i.AsBx()
			operands = fmt.Sprintf("%d\t%d\t", a, b)
		case IAx:
			a := i.Ax()
			operands = fmt.Sprintf("%d\t", a)
		}

		fmt.Printf("\t%d\t[%s]\t%s\t%s\n", pc, line, i.OpName(), operands)
	}
}
func loadLuaScript(name, content string) (*State, error) {
	luaFile := "/tmp/" + name + ".lua"
	outFile := "/tmp/" + name + ".out"
	os.WriteFile(luaFile, []byte(content), os.ModePerm)
	if err := exec.Command("../lua-5.3.6/src/luac", "-o", outFile, luaFile).Run(); err != nil {
		return nil, err
	}
	if buf, err := os.ReadFile(outFile); err != nil {
		return nil, err
	} else {
		proto := chunk.Undump(buf)
		vm := NewState()
		vm.Load(proto)
		return vm, nil
	}
}
func runLuaScript(name, content string) error {
	if vm, err := loadLuaScript(name, content); err != nil {
		return err
	} else {
		vm.Run()
		return nil
	}
}
func TestLuaScript(t *testing.T) {
	script := ``
	if err := runLuaScript("test", script); err != nil {
		t.Fatal(err)
	}
}
func TestArith(t *testing.T) {
	nums := []interface{}{1, 1.0, 1.1, 2, 2.0, 2.2}
	operators := []string{"+", "-", "*", "/", "//", "%", "^"}
	for _, op := range operators {
		for _, a := range nums {
			for _, b := range nums {
				arith := fmt.Sprintf("%v%s%v", a, op, b)
				script := fmt.Sprintf("print(\"%v%s%v=\", %s)", a, op, b, arith)
				if err := runLuaScript("arith-"+fmt.Sprintf("%v-%v", a, b), script); err != nil {
					t.Fatal(script, err)
				}
			}
		}
	}
}
func TestFor(t *testing.T) {
	script := `
local sum = 0
for i = 1, 100 do
  if i % 2 == 0 then
    sum = sum + i
  end
end
`
	if err := runLuaScript("for", script); err != nil {
		t.Fatal(err)
	}
}
func TestTable(t *testing.T) {
	script := `
local t = {"a", "b", "c"}
t[2] = "B"
t["foo"] = "Bar"
local s = t[3] .. t[2] .. t[1] .. t["foo"] .. #t
`
	if err := runLuaScript("table", script); err != nil {
		t.Fatal(err)
	}
}
func TestClosure(t *testing.T) {
	script := `
local function max(...)
  local args = {...}
  local val, idx
  for i = 1, #args do
    if val == nil or args[i] > val then
	  val, idx = args[i], i
	end
  end
  return val, idx
end
local v1 = max(3, 9, 7, 128, 35)
assert(v1 == 128)
local v2, i2 = max(3, 9, 7, 128, 35)
assert(v2 == 128 and i2 == 4)
local v3, i3 = max(max(3, 9, 7, 128, 35))
assert(v3 == 128 and i3 == 1)
local t = {max(3,9,7,128,35)}
assert(t[1]==128 and t[2]==4)
`
	if err := runLuaScript("closure", script); err != nil {
		t.Fatal(err)
	}
}

func TestUpvalue1(t *testing.T) {
	script := `
x = 1
function g() print(x); x = 2 end
function f() local x = 3; g() end
f()      --> 1
print(x) --> 2`
	if err := runLuaScript("upvalue1", script); err != nil {
		t.Fatal(err)
	}
}

func TestUpvalue2(t *testing.T) {
	script := `
function newCounter ()
  local count = 0
  return function () -- anonymous function
    count = count + 1
    return count
  end
end

c1 = newCounter()
print(c1()) --> 1
print(c1()) --> 2

c2 = newCounter()
print(c2()) --> 1
print(c1()) --> 3
print(c2()) --> 2`
	if err := runLuaScript("upvalue2", script); err != nil {
		t.Fatal(err)
	}
}

func TestUpvalue3(t *testing.T) {
	script := `
local u,v,w = 1,2,3
local function f()
	print(u,v,w)
	u = 11
end
v = 22
local function g()
	print(u,v,w)
	v = 222
end
w = 33
f()          --> 1	22	33
g()          --> 11	22	33
print(u,v,w) --> 11	222	33`
	if err := runLuaScript("upvalue3", script); err != nil {
		t.Fatal(err)
	}
}

func TestMeta(t *testing.T) {
	script := `
local mt = {}
function vector(x,y)
  local v = {x=x, y=y}
  setmetatable(v,mt)
  return v
end
mt.__add = function(v1, v2)
  return vector(v1.x+v2.x, v1.y+v2.y)
end
mt.__sub = function(v1, v2)
  return vector(v1.x-v2.x, v1.y-v2.y)
end
mt.__mul = function(v1, n)
  return vector(v1.x*n, v1.y*n)
end
mt.__div = function(v1, n)
  return vector(v1.x/n, v1.y/n)
end
mt.__len = function(v)
  return (v.x*v.x + v.y*v.y) ^ 0.5
end
mt.__eq = function(v1, v2)
  return v1.x == v2.x and v1.y == v2.y
end
mt.__index = function(v, k)
  if k == "print" then
    return function()
	  print("[" .. v.x .. "," .. v.y .. "]")
	end
  end
end
mt.__call = function(v)
  print("[" .. v.x .. "," .. v.y .. "]")
end
v1 = vector(1,2); v1:print()
v2 = vector(3,4); v2:print()
v3 = v1*2;        v3:print()
v4 = v1+v2;       v4:print()
print(#v2)
print(v1 == v2)
print(v2 == vector(3,4))
v4()
`
	if err := runLuaScript("meta", script); err != nil {
		t.Fatal(err)
	}
}
func TestLoop(t *testing.T) {
	script := `
t = {a=1,b=2,c=3}
for k,v in pairs(t) do
  print(k,v)
end
t= {"a","b","c"}
for k,v in pairs(t) do
  print(k,v)
end
	`
	if err := runLuaScript("loop", script); err != nil {
		t.Fatal(err)
	}
}
func TestFib(t *testing.T) {
	script := `
local function fib(n)
    if n < 2 then return n end
    return fib(n - 2) + fib(n - 1)
end
print(fib(35))`
	if err := runLuaScript("fib", script); err != nil {
		t.Fatal(err)
	}
}

func TestGoCallLua(t *testing.T) {
	if vm, err := loadLuaScript("goCallLua", "function f(x,y) return x+y end"); err != nil {
		t.Fatal(err)
	} else {
		vm.Run()
		if results, err := vm.CallByParam("f", 1, 2); err != nil {
			t.Fatal(err)
		} else if results[0] != 3 {
			t.Fatal("1+2 != 3")
		}
	}
}
