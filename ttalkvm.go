package main

import (
	"bytes"
	"encoding/binary"
	"fmt"

	//"log"

	talk "github.com/Racinettee/ttalkvm/pkg/vm"
)

func main() {
	// Program begins with magic header "talk\0"
	progCode := bytes.NewBuffer([]byte{0x74, 0x61, 0x6c, 0x6b, 0})
	msg := []byte("Hello World!")
	// The next header is data, begining with how much total data there is
	binary.Write(progCode, binary.LittleEndian, uint32(len(msg)))
	// The data sits linearly together
	progCode.Write(msg)
	// Next is the code - pushes two 32 bit signed ints to stack, then adds
	progCode.WriteByte(talk.PushInt32)
	binary.Write(progCode, binary.LittleEndian, int32(101))
	progCode.WriteByte(talk.PushInt32)
	binary.Write(progCode, binary.LittleEndian, int32(202))
	progCode.Write([]byte{talk.AddI32, talk.PrintTop, talk.PopTop, talk.PushNil, talk.PushString})
	// Now we're adding a string reference to be printed
	binary.Write(progCode, binary.LittleEndian, uint32(0))
	binary.Write(progCode, binary.LittleEndian, uint32(len(msg)))
	// Print the string
	progCode.Write([]byte{talk.PrintTop, talk.NativeCall, talk.End})

	vm := talk.NewVMFromByteCode(progCode.Bytes())

	vm.AddFunction(string(msg), func(vm *talk.TtalkVm, reciever interface{}) int {
		fmt.Println("Yolo!")
		return 0
	})
	// Prints 303 then Hello World!
	vm.Interpret()
}
