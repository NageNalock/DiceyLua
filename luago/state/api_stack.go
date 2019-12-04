package state

/**
类似 Peek()
 */
func (self *luaState) GetTop() int {
	return self.stack.top
}

func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}

func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true  // 忽略失败的情况
}

func (self *luaState) Pop(n int) {
	for i := 0;i < n;i++ {
		self.stack.pop()
	}
}

func (self *luaState) Copy(fromIdx, toIdx int) {
	self.stack.set(toIdx, self.stack.get(fromIdx))
}

/**
将指定位置的元素推入栈顶
 */
func (self *luaState) PushValue(idx int) {
	self.stack.push(self.stack.get(idx))
}

/**
PushValue() 的反操作, 将栈顶元素弹出放入指定位置
 */
func (self *luaState) Replace(idx int) {
	self.stack.set(idx, self.stack.pop())
}

/**
将栈顶元素弹出插入到指定位置
 */
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.Pop(1)
}

/**
很关键,
将 [idx, top]去年内的值朝栈顶旋转 n 个位置, 如果 n 为负数则朝下旋转

ex:
栈中元素为 [a, b, c, d, e]
Rotate(2, 1) 后为 [a, e, b, c, d]
Rotate(2, -1) 后为 [a, c, d, e, b]
 */
func (self *luaState) Rotate(idx, n int) {
	t := self.stack.top - 1
	p := self.stack.absIndex(idx)

	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}

	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

/**
设置栈顶值为指定索引的值
 */
func (self *luaState) SetTop(idx int) {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("stack underFlow!")
	}

	n := self.stack.top - newTop
	if n > 0 {
		for i := 0;i < n;i++ {
			self.stack.pop()
		}
	} else {
		for i := 0;i > n;i-- {
			self.stack.push(nil)
		}
	}
}
