package vm

// Operation code definitions and reference
const (
	PushRef byte = iota
	// command 1 byte + 4 bytes representing i32le
	PushInt32
	PushFloat32
	// command 1 byte, 4 bytes for loc (u32le) within data, 4 bytes for len (u32le)
	// Pushes a string from the data section to the top of the stack
	PushString
	// make a copy of the top element 1 byte cmd
	PushTop
	// Push nil value to the stack
	PushNil
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

	// 1 byte command, puts a new table at the top of stack
	NewTable

	// Pops 2: 1st: string from stack, 2nd: receiver
	// Calls a function from Go itself
	NativeCall

	// This is always last
	End
)
