package number

import "math"

/**
运算符
 */

/**
整除操作
 */
func IFloorDiv(a, b int64) int64 {
	if a > 0 && b > 0 || a < 0 && b < 0 || a % b == 0 {
		return a / b
	} else {  // 非整除时向下取整, 如果除数结果为负且不整除时则是向上取整, 此时与 Go 的整除规则相同
		return a / b - 1
	}
}

func FFloorDiv(a, b float64) float64 {
	return math.Floor(a / b)
}

/**
取模本身与整除是类似的
 */
func IMod(a, b int64) int64 {
	return a - IFloorDiv(a, b) * b
}

func FMod(a, b float64) float64 {
	return a - FFloorDiv(a, b) * b
}

/*
Go 的右移运算分有符号和无符号
有符号的话是空缺补 1, 无符号是补 0.
Lua 则都为无符号, 因此需要转换
 */
func ShiftRight(a, n int64) int64 {
	if n >= 0 {
		return int64(uint64(a) >> uint64(n))
	} else {
		return ShiftLeft(a, -n)
	}
}

func ShiftLeft(a, n int64) int64 {
	if n >= 0 {
		return a << uint64(n)
	} else {
		return ShiftRight(a, -n)
	}
}
