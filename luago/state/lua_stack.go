package state

type luaStack struct {
	slots []luaValue  // 值
	top int  // 栈顶索引
}

func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

/**
检查栈的空间是否还可以容纳至少 n 个值
当不满足时则扩容
 */
func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n;i++ {
		self.slots = append(self.slots, nil)
	}
}

func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("stack overflow!")
	}

	self.slots[self.top] = val
	self.top++
}

func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		// 栈从 1 开始
		panic("stack underflow")
	}

	self.top--
	value := self.slots[self.top]
	self.slots[self.top] = nil
	return value
}

/**
将索引转换成绝对索引
 */
func (self *luaStack) absIndex(idx int) int {
	if idx < 0 {
		idx = self.top + idx + 1
	}
	return idx
}

/**
判断索引是否有效
 */
func (self *luaStack) isValid(idx int) bool {
	absIdx := self.absIndex(idx)
	return absIdx > 0 && absIdx <= self.top
}

/**
根据索引从栈里取值
 */
func (self *luaStack) get(idx int) luaValue {
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top {
		return self.slots[absIdx - 1]
	}
	return nil
}

/**
写值
 */
func (self *luaStack) put(idx int, val luaValue) {
	absIndex := self.absIndex(idx)
	if absIndex > 0 && absIndex <= self.top {
		self.slots[absIndex - 1] = val
		return
	}

	panic("invalid index!")
}