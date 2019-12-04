package state

import (
	. "DiceyLua/luago/api"
	"fmt"
)


/* 都是使用索引来访问栈, 但不改变栈的状态 */

func (self *luaState) TypeName(tp LuaType) string {
	switch tp {
	case LUA_TNONE:
		return "no value"
	case LUA_TNIL:
		return "nil"
	case LUA_TBOOLEAN:
		return "boolean"
	case LUA_TNUMBER:
		return "number"
	case LUA_TSTRING:
		return "string"
	case LUA_TTABLE:
		return "table"
	case LUA_TFUNCTION:
		return "function"
	case LUA_TTHREAD:
		return "thread"
	default:
		return "userdata"
	}
}

/**
根据索引类型返回索引的值
(实际是返回一个 int)
 */
func (self *luaState) Type(idx int) LuaType {
	if self.stack.isValid(idx) {
		val := self.stack.get(idx)
		return typeOf(val)
	}

	return LUA_TNONE
}

func (self *luaState) IsNone(idx int) bool {
	return self.Type(idx) == LUA_TNONE
}

func (self *luaState) IsNil(idx int) bool {
	return self.Type(idx) == LUA_TNIL
}

func (self *luaState) IsNoneOrNil(idx int) bool {
	return self.Type(idx) <= LUA_TNIL
}

func (self *luaState) isBoolean(idx int) bool {
	return self.Type(idx) == LUA_TBOOLEAN
}

func (self *luaState) isString(idx int) bool {
	t := self.Type(idx)
	return t == LUA_TSTRING || t == LUA_TNUMBER
}

func (self *luaState) isNumber(idx int) bool {
	_, ok := self.ToIntegerX(idx)
	return ok
}

func (self *luaState) IsInteger(idx int) bool {
	val := self.stack.get(idx)
	_, ok := val.(int64)
	return ok
}

/**
只有当取出为空或者 false 时为假
其他皆为真
 */
func (self *luaState) ToBoolean(idx int) bool {
	val := self.stack.get(idx)
	return convertToBoolean(val)
}

func convertToBoolean(val luaValue) bool {
	switch x := val.(type) {
	case nil:
		return false
	case bool:
		return x
	default:
		return true
	}
}

/**
取出数字, 如果取出的类型不为数字则会进行强制类型转换
如果不是数字且转换失败, 则直接返回 0
 */
func (self *luaState) ToNumber(idx int) float64 {
	val, _ := self.ToNumberX(idx)
	return val
}

/**
取出数字, 如果取出的类型不为数字则会进行强制类型转换
如果不是数字且转换失败, 则直接返回 0 和 是否转换成功的布尔值
 */
func (self *luaState) ToNumberX(idx int) (float64, bool) {
	val := self.stack.get(idx)
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

func (self *luaState) ToInteger(idx int) int64 {
	val, _ := self.ToIntegerX(idx)
	return val
}

func (self *luaState) ToIntegerX(idx int) (int64, bool) {
	val := self.stack.get(idx)
	i, ok := val.(int64)  // todo 类型转换判断
	return i, ok
}

func (self *luaState) ToString(idx int) string {
	s, _ := self.ToStringX(idx)
	return s
}

func (self *luaState) ToStringX(idx int) (string, bool) {
	val := self.stack.get(idx)
	switch x := val.(type) {
		case string: {
			return x, true
		}
		case int64, float64: {
			s := fmt.Sprintf("%v", x)
			self.stack.set(idx, s)
			return s, true
		}
		default:
			return "", false
	}
}