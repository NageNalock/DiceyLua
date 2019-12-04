package main

import (
	. "DiceyLua/luago/api"
	"DiceyLua/luago/binchunk"
	. "DiceyLua/luago/state"
	. "DiceyLua/luago/vm"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	// todo 以后改成单测框架
	ls := New()

	ls.PushBoolean(true)
	printStack(ls)

	ls.PushInteger(10)
	printStack(ls)

	ls.PushNil()
	printStack(ls)

	ls.PushString("hello")
	printStack(ls)

	ls.PushValue(-4)
	printStack(ls)

	ls.Replace(3)
	printStack(ls)

	ls.SetTop(6)
	printStack(ls)

	ls.Remove(-3)
	printStack(ls)

	ls.SetTop(-5)
	printStack(ls)
}

func printStack(ls LuaState) {
	// todo 写入栈的方法
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default:
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func decompilerMain() {
	// 测试二进制 chunk 部分是否正常
	if len(os.Args) > 1{
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(err)
		}

		proto := binchunk.Undump(data)
		list(proto)
	}
}

func list(f *binchunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)

	// 子函数原型
	for _,p := range f.Protos {
		list(p)
	}
}

func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"  // main 函数的起行号为 0, 非 main 函数不包含源文件名以节省空间
	}

	varargFlag := ""
	if f.IsVararg > 0 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s:%d,%d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ", f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n", len(f.LocVars), len(f.Constants), len(f.Protos))
}

func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])  // 格式化行号
		}
		i := Instruction(c)

		fmt.Printf("\t%d\t[%s]\t%s \t", pc + 1, line, i.OpName())
		printOperands(i)
		fmt.Printf("\n")
	}
}

/**
打印非 OpArgN 类型的操作数
 */
func printNotOpArgN(arg int) {
	if arg > 0xFF {
		fmt.Printf(" %d", -1-arg&0xFF)  // 表示常量表索引按照负数输出
	} else {
		fmt.Printf(" %d", arg)
	}
}

func printOperands(i Instruction) {
	switch i.OpMode() {
	case IABC:
		a, b, c := i.ABC()
		fmt.Printf("%d", a)
		if i.BMode() != OpArgN {  // OpArgN 类型的参数不打印
			printNotOpArgN(b)
		}

		if i.CMode() != OpArgN {
			printNotOpArgN(c)
		}
	case IABx:
		a, bx := i.ABx()
		fmt.Printf("%d", a)
		if i.BMode() == OpArgK {
			fmt.Printf(" %d", -1-bx)  // 表示常量表索引按照负数输出
		} else if i.BMode() == OpArgU {
			fmt.Printf(" %d", bx)
		}
	case IAsBx:
		a, sbx := i.AsBx()
		fmt.Printf("%d %d", a, sbx)
	case IAx:
		ax := i.Ax()
		fmt.Printf("%d", -1-ax)
	default:
		fmt.Println("error operation mode: ", i.OpMode())
	}
}

/*
打印常量表, 局部变量表和 Upvalue 表
 */
func printDetail(f *binchunk.Prototype)  {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPC, locVar.EndPC)
	}

	fmt.Printf("upvalues (%d(:\n", len(f.Upvalues))
	for i, upvalue := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, upvalName(f, i), upvalue.Instack, upvalue.Idx)
	}
}

func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}

func upvalName(f *binchunk.Prototype, i int) string {
	if len(f.UpvalueNames) > 0 {
		return f.UpvalueNames[i]
	}
	return "-"
}