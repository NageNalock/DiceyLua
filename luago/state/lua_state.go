package state

type luaState struct {
	stack *luaStack
}

func New() *luaState {
	return &luaState{stack:newLuaStack(20)}  // 栈的初始大小设为 20(拍脑袋
}