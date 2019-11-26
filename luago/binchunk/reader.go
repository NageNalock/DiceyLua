package binchunk

import (
	"encoding/binary"
	"fmt"
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
func (self *reader) readByte() byte {
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
func (self *reader) readUint32() uint32 {
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

	if size == 0 {
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

    if self.readByte() != LUAC_VERSION {
    	panic("version mismatch!")
	}

	if self.readByte() != LuacFormat {
		panic("format mismatch!")
	}

	if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}

	if self.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}

	if self.readByte() != CSZIET_SIZE {
		panic("size_t size mismatch!")
	}

	if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	}

	if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch!")
	}

	if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}

	if self.readLuaInteger() != LUAC_INT {
		panic("endianness mismatch!")
	}

	if self.readLuaNumber() != LUAC_NUM {
		panic("float format mismatch")
	}
}

/*
读取函数原型
 */
func (self *reader) readProto(parentSource string) *Prototype {
	source := self.readString()
	if source == "" {
		source = parentSource  // 只有主函数才有源文件名, 如果没有源文件名则不是主函数
	}
	return &Prototype{
		Source:          source,
		LineDefined:     self.readUint32(),    // 起行号
		LastLineDefined: self.readUint32(),    // 止行号
		NumParams:       self.readByte(),      // 固定参数个数
		IsVararg:        self.readByte(),      // 是否是 Vararg 函数
		MaxStackSize:    self.readByte(),      // 寄存器数量
		Code:            self.readCode(),      // 指令表
		Constants:       self.readConstants(), // 常量表
		Upvalues:        self.readUpvalues(),
		Protos:          self.readProtos(source),  // 子函数原型
		LineInfo:        self.readLineInfo(),  // 行号表
		LocVars:         self.readLocVars(),  // 局部变量表
		UpvalueNames:    self.readUpvalueNames(),
	}
}

func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32() // 每条指令 4 字节
	}
	return code
}

func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

func (self *reader) readConstant() interface{} {
	switch self.readByte() {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.readByte() != 0
	case TAG_INTEGER:
		return self.readLuaInteger()
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_SHORT_STR:
		return self.readString()
	case TAG_LONG_STR:
		return self.readString()
	default:
		panic("read constant corrupted!")
	}
}

func (self *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}

	return upvalues
}

func (self *reader) readProtos(parentSource string) []*Prototype {
	prototypes := make([]*Prototype, self.readUint32())
	for i := range prototypes {
		prototypes[i] = self.readProto(parentSource)
	}

	return prototypes
}

func (self *reader) readLineInfo() []uint32 {
	lienInfo := make([]uint32, self.readByte())
	for i := range lienInfo {
		lienInfo[i] = self.readUint32()
		fmt.Println(lienInfo[i])
		fmt.Println("-")
	}

	return lienInfo
}

func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPC: self.readUint32(),
			EndPC:   self.readUint32(),
		}
	}

	return locVars
}

func (self *reader) readUpvalueNames() []string {
	upvalueNames := make([]string, self.readUint32())
	for i := range upvalueNames {
		upvalueNames[i] = self.readString()
		fmt.Println(upvalueNames[i])
	}

	return upvalueNames
}