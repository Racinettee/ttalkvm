package util

import (
	"bytes"
	"encoding/binary"
	"io"
)

func NewModuleBuffer() *bytes.Buffer {
	return bytes.NewBuffer([]byte{0x74, 0x61, 0x6c, 0x6b, 0})
}

func WriteU32(w io.Writer, u32 uint32) {
	binary.Write(w, binary.LittleEndian, u32)
}

func WriteI32(w io.Writer, i32 int32) {
	binary.Write(w, binary.LittleEndian, i32)
}
