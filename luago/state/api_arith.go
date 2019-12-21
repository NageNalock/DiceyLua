package state

import (
	. "DiceyLua/luago/api"
	"DiceyLua/luago/number"
	"math"
)

/**
所有算数运算符的集合
 */

var (
	iadd = func(a, b int64) int64 {
		return a + b
	}

	fadd = func(a, b float64) float64 {
		return a + b
	}

	isub = func(a, b int64) int64 {
		return a - b
	}

	fsub = func(a, b float64) float64 {
		return a - b
	}

	imul = func(a, b int64) int64 {
		return a * b
	}

	fmul = func(a, b float64) float64 {
		return a * b
	}

	imod = number.IMod

	fmod = number.FMod

	pow = math.Pow

	div = func(a, b float64) float64 {
		return a / b
	}

	iidiv = number.IFloorDiv

	fidiv = number.FFloorDiv

	band = func(a, b int64) int64 {
		return a & b
	}

	bor = func(a, b int64) int64 {
		return a | b
	}

	bxor = func(a, b int64) int64 {
		return a ^ b
	}

	shl = number.ShiftLeft

	shr = number.ShiftRight

	iunm = func(a, _ int64) int64 {
		return -a
	}

	funm = func(a, _ float64) float64 {
		return -a
	}

	bnot = func(a, _ int64) int64 {
		return ^a
	}
)

type operator struct {
	integerFunc func(int64, int64) int64
	floatFunc func(float64, float64) float64
}

var operators = []operator{
	{
		integerFunc: iadd,
		floatFunc:   fadd,
	},
	{
		integerFunc: isub,
		floatFunc:   fsub,
	},
	{
		integerFunc: imul,
		floatFunc:   fmul,
	},
	{
		integerFunc: imod,
		floatFunc:   fmod,
	},
	{
		integerFunc: nil,
		floatFunc:   pow,
	},
	{
		integerFunc: nil,
		floatFunc:   div,
	},
	{
		integerFunc: iidiv,
		floatFunc:   fidiv,
	},
	{
		integerFunc: band,
		floatFunc:   nil,
	},
	{
		integerFunc: bor,
		floatFunc:   nil,
	},
	{
		integerFunc: bxor,
		floatFunc:   nil,
	},
	{
		integerFunc: shl,
		floatFunc:   nil,
	},
	{
		integerFunc: shr,
		floatFunc:   nil,
	},
	{
		integerFunc: iunm,
		floatFunc:   funm,
	},
	{
		integerFunc: bnot,
		floatFunc:   nil,
	},
}

func (self *luaState) Arith(op ArithOp) {
	var a, b luaValue
	b = self.stack.pop()
	if op != LUA_OPUNM && op != LUA_OPBNOT {
		a = self.stack.pop()
	} else {
		// 取反的二元操作符相同
		a = b
	}

	operator := operators[op]
	if result := _arith(a, b, operator); result != nil {
		self.stack.push(result)
	} else {
		panic("arithmetic error!")
	}
}

func _arith(a, b luaValue, op operator) luaValue {
	if op.floatFunc == nil {
		// 操作类型为位运算
		if x, ok := convertToInteger(a); ok {
			if y, ok := convertToInteger(b); ok {
				return op.integerFunc(x, y)
			}
		}
	} else {
		if op.integerFunc != nil {
			// 加减乘除取反取模
			if x, ok := a.(int64); ok {
				if y, ok := b.(int64); ok {
					return op.integerFunc(x, y)
				}
			}
		}
		if x, ok := convertToFloat(a); ok {
			if y, ok := convertToFloat(b); ok {
				return op.floatFunc(x, y)
			}
		}
	}

	return nil
}