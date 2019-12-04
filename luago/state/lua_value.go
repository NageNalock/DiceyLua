package state

import . "DiceyLua/luago/api"

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