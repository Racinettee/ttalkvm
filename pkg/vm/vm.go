package vm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"unsafe"
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
	dataHeaderOffset uint32 = 5
	dataDataOffset   uint32 = dataHeaderOffset + 4
)

type TtalkCallable = func(vm *TtalkVm, reciever interface{}) int
type TtalkTable = map[string]interface{}

type TtalkVm struct {
	ops           []byte
	stack         []interface{}
	stackSize     int
	nativeMethods map[string]TtalkCallable
}

func NewVMFromByteCode(code []byte) TtalkVm {
	return TtalkVm{
		ops:           code,
		stack:         make([]interface{}, 100),
		stackSize:     0,
		nativeMethods: make(map[string]TtalkCallable),
	}
}

// Objects
func (tvm *TtalkVm) NewTable() {
	tvm.Push(make(TtalkTable))
}

// Host functionality
func (tvm *TtalkVm) AddFunction(name string, fn TtalkCallable) {
	tvm.nativeMethods[name] = fn
}

// Stack operations
func (tvm *TtalkVm) Push(obj interface{}) {
	if tvm.stackSize+1 >= len(tvm.stack) {
		newStack := make([]interface{}, len(tvm.stack)+100)
		copy(newStack, tvm.stack)
		tvm.stack = newStack
	}
	tvm.stack[tvm.stackSize] = obj
	tvm.stackSize += 1
}

func (tvm *TtalkVm) Pop() (result interface{}) {
	result = tvm.stack[tvm.stackSize-1]
	tvm.stack[tvm.stackSize-1] = nil
	tvm.stackSize -= 1
	return
}

func (tvm *TtalkVm) Top() interface{} {
	return tvm.stack[tvm.stackSize-1]
}

func (tvm *TtalkVm) ShiftDown(atIndex int) {
	// At index will be overridden by elements "above it" in the stack
	if atIndex >= tvm.stackSize {
		return
	}
	// Starting at atIndex, and climbing the stack till reaching the top
	for ; atIndex < tvm.stackSize; atIndex += 1 {
		if atIndex+1 == tvm.stackSize {
			continue
		}
		tvm.stack[atIndex] = tvm.stack[atIndex+1]
	}
}

func (tvm *TtalkVm) Interpret() {
	// First ensure the header:
	if !bytes.Equal(tvm.ops[0:5], []byte{0x74, 0x61, 0x6c, 0x6b, 0}) {
		log.Printf("Header is wrong %x", tvm.ops[0:5])
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
			ptr += 5 // 1 + 4 (command, bytes for i32) - stak + 1

		case PushString:
			ptrLoc := dataDataOffset + binary.LittleEndian.Uint32(ops[ptr+1:])
			ptrSiz := binary.LittleEndian.Uint32(ops[ptr+5:])
			str := tvm.ops[ptrLoc : ptrLoc+ptrSiz]
			tvm.Push(string(str))
			ptr += 9 // 1 + 4 + 4 (command, location, length) - stak + 1

		case PushTop:
			tvm.Push(tvm.Top())
			ptr += 1 // 1 (command) - stak + 1

		case PushNil:
			tvm.Push(nil)
			ptr += 1 // 1 (command) - stak + 1

		case AddI32:
			right := tvm.Pop().(int32)
			left := tvm.Pop().(int32)
			tvm.Push(left + right)
			ptr += 1 // 1 (command) - stak + 1

		case PopTop:
			tvm.Pop()
			ptr += 1 // 1 (command) - stak - 1

		case PopN:
			numPop := int(binary.LittleEndian.Uint32(tvm.ops[ptr+1:]))
			for i := 0; i < numPop; i += 1 {
				tvm.stack[tvm.stackSize-(1+i)] = nil
			}
			tvm.stackSize -= numPop
			ptr += 5 // 1 + 4 (command, n elements as u32) - stak - N

		case PrintTop:
			fmt.Println(tvm.Top())
			ptr += 1 // 1 (command) - stak +/- 0

		case NewTable:
			tvm.NewTable()
			ptr += 1 // 1 (command) - stak + 1

		case NativeCall:
			methodStr := tvm.Pop().(string)
			reciever := tvm.Pop()

			tvm.nativeMethods[methodStr](tvm, reciever)
			ptr += 1 // 1 (command) - stak - 2 + N
		case End:
			return
		}
	}
}
