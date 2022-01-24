package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
)

const (
	PushRef byte = iota
	PushInt32
	PushFloat32
	AddI32
	AddF32
	// Prints to stdout whatever value is "top"
	PrintTop
	// This is always last
	End
)

type TtalkVm struct {
	ops []byte
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
		fmt.Printf("Header is wrong %x", tvm.ops[0:4])
		panic("")
	}
	// Start past the header
	ops := tvm.ops[5:]
	ptr := 0
	for {
		switch ops[ptr] {
		case PushInt32:
			ptrI32 := unsafe.Pointer(&ops[ptr+1])
			tvm.Push(*(*int32)(ptrI32))
			ptr += 5 // 1 + 4 bytes for this op

		case AddI32:
			right := tvm.Pop().(int32)
			left := tvm.Pop().(int32)
			tvm.Push(left + right)
			ptr += 1

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
	// then pushes two 32 bit signed ints to stack, then adds
	progCode := bytes.NewBuffer([]byte{0x74, 0x61, 0x6c, 0x6b, 0})
	progCode.WriteByte(PushInt32)
	binary.Write(progCode, binary.LittleEndian, int32(101))
	progCode.WriteByte(PushInt32)
	binary.Write(progCode, binary.LittleEndian, int32(202))
	progCode.Write([]byte{AddI32, End})

	vm := NewVMFromByteCode(progCode.Bytes())
	vm.Interpret()
	// Prints 303
	fmt.Println(vm.Top())
}