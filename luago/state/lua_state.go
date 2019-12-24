package state

import "DiceyLua/luago/binchunk"

type luaState struct {
	stack *luaStack

	proto *binchunk.Prototype  // 函数原型, 用于提取指令或常量
	pc int  // 程序计数器
}

/* TODO 补齐剩余方法 */
func (self *luaState) IsTable(idx int) bool {
	panic("implement me")
}

func (self *luaState) IsThread(idx int) bool {
	panic("implement me")
}

func (self *luaState) IsFunction(idx int) bool {
	panic("implement me")
}

func New(stackSize int, proto *binchunk.Prototype) *luaState {
	return &luaState{
		stack: newLuaStack(stackSize),  // 栈的初始大小设为 20(拍脑袋
		proto: proto,
		pc:    0,
	}
}