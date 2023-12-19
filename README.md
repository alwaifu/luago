# luago
A lua vm implement by Go


## 笔记
[《自己动手实现Lua 虚拟机、编译器和标准库》](https://github.com/zxh0/luago-book)读书笔记：[blog.md](blog.md)

## How to use
只实现了虚拟机部分，暂未实现编译器，需要官方luac编译得到`binary chunk`后加载字节码运行
1. 解压 [lua-5.3.6.tar.gz](./lua-5.3.6.tar.gz)（来自(https://github.com/lua/lua/releases/tag/v5.3.6)）
2. 修改 `lua-5.3.6` `Makefile` 第七行 `PLAT= none` 为对应平台后执行`make`命令
3. 参考 [vm_test.go](./vm/vm_test.go) 执行 lua binary chunk
## How about performance
参考 [benchmark_fib_test.go](./example/benchmark_fib_test.go) 该玩具项目与 [gopherlua](https://github.com/yuin/gopher-lua) 执行斐波纳切性能对比  
参考 (https://github.com/yuin/gopher-lua/wiki/Benchmarks) 对比其他语言执行斐波纳切性能