package binchunk

import (
	"encoding/binary"
	"math"
)

/***
解析二进制 chunk
 */

type reader struct {
	data []byte  // 要解析的 chunk 数据
}

/*
读取一个字节
 */
func (self *reader) readByte() byte  {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

/*
使用小端的方式读取一个 cint
 */
func (self *reader) readUnit32() uint32 {
	u := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return u
}

/*
使用小端的方式读取一个 size_t
 */
func (self *reader) readUint64() uint64 {
	u := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return u
}

/*
读取一个 Lua 整型
 */
func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

/*
读取一个 Lua 浮点数
*/
func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

func (self *reader) readString() string {
	size := uint(self.readByte())  // 首位表示长度

	if size == 0{
		// null
		return ""
	}
	if size == 0xFF {
		// 长字符串, 长度在第二位
		size = uint(self.readUint64())
	}

	// 默认为短字符串, 去掉 size 位
	bytes := self.readBytes(size - 1)
	return string(bytes)
}

/*
检查头部
 */
func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}

	if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}

	if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}

	if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch")
	}

    header := self.readByte()
    if header != LUAC_VERSION {
    	panic("version mismatch!")
	}

	if header != LUAC_FORMAT {
		panic("format mismatch!")
	}

	if header != CINT_SIZE {
		panic("int size mismatch!")
	}

	if header != CSZIET_SIZE {
		panic("size_t size mismatch!")
	}

	if header != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}

	if header != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}

	if header != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}
}
