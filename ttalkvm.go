package main

import (
	"fmt"

	//"log"

	talk "github.com/Racinettee/ttalkvm/pkg/vm"
	bld "github.com/Racinettee/ttalkvm/pkg/codebuild"
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
	progCode.Write([]byte{talk.PrintTop, talk.NativeCall, talk.End})

	vm := talk.NewVMFromByteCode(progCode.Bytes())

	vm.AddFunction(string(msg), func(vm *talk.TtalkVm, reciever interface{}) int {
		fmt.Println("Message from Go")
		return 0
	})
	// Prints 303 then Hello World! then Message from Go
	vm.Interpret()
}
