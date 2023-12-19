package example

import (
	"luago/chunk"
	"luago/vm"
	"os"
	"os/exec"
	"strings"
	"testing"

	gopherlua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

const LUA_SCRIPT = `local function fib(n)
if n < 2 then return n end
return fib(n - 2) + fib(n - 1)
end
print(fib(35))`

func BenchmarkRunScript(b *testing.B) {
	b.Run("my_lua_vm", BenchmarkMyRunScript)
	b.Run("gopher_lua_vm", BenchmarkGopherluaRunScript)
}
func BenchmarkMyRunScript(b *testing.B) {
	if vm, _, err := initVM(LUA_SCRIPT); err != nil {
		b.Fatal(err)
	} else {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			vm.Run()
		}
	}
}

func BenchmarkGopherluaRunScript(b *testing.B) {
	if _, vm, err := initVM(LUA_SCRIPT); err != nil {
		b.Fatal(err)
	} else {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			vm.Call(0, gopherlua.MultRet)
		}
	}
}

func initVM(script string) (vm1 *vm.State, vm2 *gopherlua.LState, err error) {
	luaFile := "/tmp/benchmark.lua"
	outFile := "/tmp/benchmark.out"
	os.WriteFile(luaFile, []byte(script), os.ModePerm)
	if err := exec.Command("../lua-5.3.6/src/luac", "-o", outFile, luaFile).Run(); err != nil {
		return vm1, vm2, err
	}
	if buf, err := os.ReadFile(outFile); err != nil {
		return vm1, vm2, err
	} else {
		proto := chunk.Undump(buf)
		vm1 = vm.NewState()
		vm1.Load(proto)
	}

	if chunk, err := parse.Parse(strings.NewReader(LUA_SCRIPT), "benchmark"); err != nil {
		return vm1, vm2, err
	} else {
		if proto, err := gopherlua.Compile(chunk, "benchmark"); err != nil {
			return vm1, vm2, err
		} else {
			vm2 = gopherlua.NewState()
			vm2.Push(vm2.NewFunctionFromProto(proto))
		}
	}
	return vm1, vm2, nil

}
