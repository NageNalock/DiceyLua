package state

import . "DiceyLua/luago/api"


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