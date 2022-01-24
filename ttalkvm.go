package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	//"log"
	"unsafe"

	//"github.com/Racinettee/ttalkvm/pkg/util"
)

/*
 * Ttalk binary spec:
 * Header:
 * First 5 bytes: talk\0 (0x74, 0x61, 0x6c, 0x6b, 0)
 * Header section data:
 * 4 Bytes: the length of the data in the module followed by the data

 * Code follows the data immediately
 * End must be included somewhere for the interpreter to stop
 */
const (
	PushRef byte = iota
	// command 1 byte + 4 bytes representing i32le
	PushInt32
	PushFloat32
	// command 1 byte, 4 bytes for loc (u32le) within data, 4 bytes for len (u32le)
	PushString
	// make a copy of the top element 1 byte cmd
	PushTop
	// command 1 byte, pops 2 i32le from stack, pushes result
	AddI32
	AddF32

	// Pops top
	// command 1 byte
	PopTop
	// Pops N from the stack
	// command 1 byte, 4 bytes (u32le)
	PopN

	// Prints to stdout whatever value is "top"
	PrintTop
	// This is always last
	End
)

const dataHeaderOffset uint32 = 5
const dataDataOffset uint32 = dataHeaderOffset + 4

type TtalkVm struct {
	ops []byte
	dat []byte
	stack []interface{}
	stackSize int
}

func NewVMFromByteCode(code []byte) TtalkVm {
	return TtalkVm {
		ops: code,
		stack: make([]interface{}, 100),
		stackSize: 0,
	}
}

func (tvm *TtalkVm) Push(obj interface{}) {
	if tvm.stackSize+1 >= len(tvm.stack) {
		newStack := make([]interface{}, len(tvm.stack) + 100)
		copy(newStack, tvm.stack)
		tvm.stack = newStack
	}
	tvm.stack[tvm.stackSize] = obj
	tvm.stackSize += 1
}

func (tvm *TtalkVm) Pop() (result interface{}) {
	result = tvm.stack[tvm.stackSize - 1]
	tvm.stack[tvm.stackSize - 1] = nil
	tvm.stackSize -= 1
	return
}

func (tvm *TtalkVm) Top() interface{} {
	return tvm.stack[tvm.stackSize - 1]
}

func (tvm *TtalkVm) Interpret() {
	// First ensure the header:
	if bytes.Compare(tvm.ops[0:5], []byte{0x74, 0x61, 0x6c, 0x6b, 0}) != 0 {
		fmt.Printf("Header is wrong %x", tvm.ops[0:5])
		panic("")
	}
	// Get the data (4 bytes describe length 6 bytes passing the header):
	dataLen := binary.LittleEndian.Uint32(tvm.ops[dataHeaderOffset:])
	//stackPointers := util.NewIntStack(100)
	// Start past the data header (offset of 9 total bytes so far, plus the length)
	ops := tvm.ops[9+dataLen:]
	ptr := 0
	for {
		switch ops[ptr] {
		case PushInt32:
			ptrI32 := binary.LittleEndian.Uint32(ops[ptr+1:])
			tvm.Push(*(*int32)(unsafe.Pointer(&ptrI32)))
			ptr += 5 // 1 + 4 bytes for this op

		case PushString:
			ptrLoc := dataDataOffset + binary.LittleEndian.Uint32(ops[ptr+1:])
			ptrSiz := binary.LittleEndian.Uint32(ops[ptr+5:])
			//log.Printf("Str ref: (%v, %v)\n", ptrLoc, ptrSiz)
			str := tvm.ops[ptrLoc:ptrLoc+ptrSiz]
			tvm.Push(string(str))
			ptr += 9

		case PushTop:
			tvm.Push(tvm.Top())
			ptr += 1

		case AddI32:
			right := tvm.Pop().(int32)
			left := tvm.Pop().(int32)
			tvm.Push(left + right)
			ptr += 1

		case PopTop:
			tvm.Pop()
			ptr += 1
		
		case PopN:
			numPop := int(binary.LittleEndian.Uint32(tvm.ops[ptr+1:]))
			for i := 0; i < numPop; i += 1 {
				tvm.stack[tvm.stackSize-(1+i)] = nil
			}
			tvm.stackSize -= numPop
			ptr += 5

		case PrintTop:
			fmt.Println(tvm.Top())
			ptr += 1

		case End:
			return
		}
	}
}

func main() {
	// Program begins with magic header "talk\0"
	progCode := bytes.NewBuffer([]byte{0x74, 0x61, 0x6c, 0x6b, 0})
	msg := []byte("Hello World!")
	// The next header is data, begining with how much total data there is
	binary.Write(progCode, binary.LittleEndian, uint32(len(msg)))
	// The data sits linearly together
	progCode.Write(msg)
	// Next is the code - pushes two 32 bit signed ints to stack, then adds
	progCode.WriteByte(PushInt32)
	binary.Write(progCode, binary.LittleEndian, int32(101))
	progCode.WriteByte(PushInt32)
	binary.Write(progCode, binary.LittleEndian, int32(202))
	progCode.Write([]byte{AddI32,PrintTop, PopTop, PushString})
	// Now we're adding a string reference to be printed
	binary.Write(progCode, binary.LittleEndian, uint32(0))
	binary.Write(progCode, binary.LittleEndian, uint32(len(msg)))
	// Print the string
	progCode.Write([]byte{PrintTop, End})

	vm := NewVMFromByteCode(progCode.Bytes())
	// Prints 303 then Hello World!
	vm.Interpret()
}