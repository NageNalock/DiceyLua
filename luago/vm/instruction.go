package vm

type Instruction uint32

/**
从指令中提取操作码数
 */
func (self Instruction) Opcode() int {
	return int(self & 0x3F)  // 3F = 111111, 32bit 的指令中, 低 6 为操作码, 剩余 25 位为操作数
}

/**
从 iABC 模式指令中提取参数
 */
func (self Instruction) ABC() (a, b, c int)  {
	a = int(self >> 6 & 0xFF)  // 6-13 为 A
	c = int(self >> 14 & 0x1FF)  // 14-22 为 C
	b = int(self >> 23 & 0x1FF)  // 23-32 为 B
	return
}

func (self Instruction) ABx() (a, bx int) {
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

const MAXARG_Bx  = 1 << 18 -1  // 2^18 - 1 = 262143, 最大能表示的无符号数
const MAXARG_sBx  = MAXARG_Bx >> 1 // 最大能表示的无符号数的一半
func (self Instruction) AsBx() (a, sbx int)  {
	a, bx := self.ABx()
	return a, bx - MAXARG_sBx  // 有符号数, 使用偏移二进制码
}

func (self Instruction) Ax() int {
	return int(self >> 6)
}

func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}