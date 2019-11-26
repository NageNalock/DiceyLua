package binchunk

type binaryChunk struct {
	header  // 头部
	// sizeUpvalues  // 主函数 upvalue 数量
	mainFunc *Prototype // 主函数原型
}

type header struct {
	signature [4]byte  // 魔数
	version byte
	format byte
	luacData [6]byte
	cintSize byte
	sizetSize byte
	instructionSize byte
	luaIntegerSize byte
	luaNumberSize byte
	luaInt int64
	luacNum float64
}

const (
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LuacFormat       = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSZIET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

type Prototype struct {
	Source string
	LineDefined uint32  // 起行号
	LastLineDefined uint32  // 止行号
	NumParams byte
	IsVararg byte  // 是否是 Var 参数函数
	MaxStackSize byte
	Code []uint32  // 指令表
	Constants []interface{}
	Upvalues []Upvalue
	Protos []*Prototype
	LineInfo []uint32  // 行号表
	LocVars []LocVar
	UpvalueNames []string
}

const (
	// 常量表的 Tag 值
	TAG_NIL = 0x00
	TAG_BOOLEAN = 0x01
	TAG_NUMBER = 0x03
	TAG_INTEGER = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR = 0x14
)

type Upvalue struct {
	Instack byte
	Idx byte
}

type LocVar struct {
	VarName string
	StartPC uint32   // 起索引
	EndPC uint32  // 止索引
}

func Undump(data []byte) *Prototype {
	// 用于解析二进制 chunk
	reader := &reader{data}
	reader.checkHeader()  // 校验头部
	reader.readByte()  // 跳过 Upvalue 数量
	return reader.readProto("")  // 读取函数原型
}