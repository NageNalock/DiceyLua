package api

type LuaVM interface {
	LuaState
	PC() int  // 返回当前 PC
	AddPC(n int)  // 修改 PC(用于实现跳转命令)
	Fetch() uint32  // 取出当前指令, 并将 PC 转向下一条指令(主要用于虚拟机循环)
	GetConst(idx int)  // 将指定常量推入栈顶
	GetRK(rk int)  // 将指定常量或栈值推入栈顶
}