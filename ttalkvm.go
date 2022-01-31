package main

import (
	"fmt"

	//"log"

	bld "github.com/Racinettee/ttalkvm/pkg/codebuild"
	talk "github.com/Racinettee/ttalkvm/pkg/vm"
)

func main() {
	// Program begins with magic header "talk\0"
	progCode := bld.NewModuleBuffer()
	msg := []byte("Hello World!")
	// The next header is data, begining with how much total data there is
	bld.WriteU32(progCode, uint32(len(msg)))
	// The data sits linearly together
	progCode.Write(msg)
	// Next is the code - pushes two 32 bit signed ints to stack, then adds
	progCode.WriteByte(talk.PushInt32)
	bld.WriteI32(progCode, 101)
	progCode.WriteByte(talk.PushInt32)
	bld.WriteI32(progCode, 202)
	progCode.Write([]byte{
		talk.AddI32, talk.PrintTop, talk.PopTop, talk.PushNil, talk.PushString,
	})
	// Now we're adding a string reference to be printed
	bld.WriteU32(progCode, 0)
	bld.WriteU32(progCode, uint32(len(msg)))
	// Print the string
	progCode.Write([]byte{talk.PrintTop, talk.NativeCall, talk.CallD})
	addFuncPos := progCode.Len() + 5 // 4 bytes (for loc + 1 for end to reach the last func)
	bld.WriteU32(progCode, uint32(addFuncPos))
	progCode.WriteByte(talk.End)
	progCode.WriteByte(talk.PushInt32)
	bld.WriteI32(progCode, 1)
	progCode.WriteByte(talk.PushInt32)
	bld.WriteI32(progCode, 2)
	progCode.WriteByte(talk.PushInt32)
	bld.WriteI32(progCode, 3)
	progCode.Write([]byte{talk.AddI32, talk.AddI32, talk.Return, byte(1)})

	vm := talk.NewVMFromByteCode(progCode.Bytes())

	vm.AddFunction(string(msg), func(vm *talk.TtalkVm, reciever interface{}) int {
		fmt.Println("Message from Go")
		return 0
	})
	// Prints 303 then Hello World! then Message from Go
	vm.Interpret()
}
