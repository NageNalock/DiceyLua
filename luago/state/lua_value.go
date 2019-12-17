package state

import (
	. "DiceyLua/luago/api"
	"DiceyLua/luago/number"
)

type luaValue interface {
	// 空接口表示各种不同类型的 Lua 值 todo
}

func typeOf(val luaValue) LuaType {
	switch val.(type) {
		case nil: {
			return LUA_TNIL
		}
		case bool: {
			return LUA_TBOOLEAN
		}
		case int64: {
			return LUA_TNUMBER
		}
		case float64: {
			return LUA_TNUMBER
		}
		case string: {
			return LUA_TSTRING
		}
		default: {
			panic("todo lua type!")
		}
	}
}

func convertToFloat(val luaValue) (float64, bool) {
	switch x := val.(type) {
	case float64:
		return x, true
	case int64:
		return float64(x), true
	case string:
		return number.ParseFloat(x)
	default:
		return 0, false
	}
}

func convertToInteger(val luaValue) (int64, bool) {
	switch x := val.(type) {
	case int64:
		return x, true
	case float64:
		return number.FloatToInteger(x)
	case string:
		if i, ok := number.ParseInteger(x); ok {
			return i, true
		}

		if i, ok := number.ParseFloat(x); ok {
			return number.FloatToInteger(i)
		}

		return 0, false
	default:
		return 0, false
	}
}

