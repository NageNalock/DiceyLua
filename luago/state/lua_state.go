package state

type luaState struct {
	stack *luaStack
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

func New() *luaState {
	return &luaState{stack:newLuaStack(20)}  // 栈的初始大小设为 20(拍脑袋
}