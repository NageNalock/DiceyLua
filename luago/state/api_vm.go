package state

func (self *luaState) PC() int {
	return self.pc
}

func (self *luaState) AddPC(n int) {
	self.pc += n
}

func (self *luaState) Fetch() uint32 {
	i := self.proto.Code[self.pc]
	self.pc++
	return i
}

func (self *luaState) GetConst(idx int) {
	c := self.proto.Constants[idx]
	self.stack.push(c)
}

// 实际参数是 iABC 模式 OpArgK 类型参数
func (self *luaState) GetRK(rk int) {
	if rk > 0xFF {  // 即为常量
		self.GetConst(rk & 0xFF)
	} else {
		self.PushValue(rk + 1)
	}
}