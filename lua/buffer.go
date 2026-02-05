package lua

import (
	"unsafe"

	"github.com/CompeyDev/lei/ffi"
)

type LuaBuffer struct {
	vm    *Lua
	index int
	size  uint64
}

func (b *LuaBuffer) Read(offset uint64, count uint64) []byte {
	b.deref(b.vm)
	defer ffi.Pop(b.vm.state(), 1)

	if buf := ffi.ToBuffer(b.vm.state(), -1, &b.size); buf != nil && offset <= b.size {
		// Clamp to the size if the count exceeds it
		if offset+count > b.size {
			count = b.size - offset
		}

		// Copy data to Go owned byte array for safety
		data := make([]byte, count)
		slice := unsafe.Slice((*byte)(buf), b.size)
		copy(data, slice[offset:offset+count])

		return data
	}

	return nil
}

func (b *LuaBuffer) Write(offset uint64, data []byte) {
	if len(data) == 0 {
		return
	}

	b.deref(b.vm)
	defer ffi.Pop(b.vm.state(), 1)

	if buf := ffi.ToBuffer(b.vm.state(), -1, &b.size); buf != nil && offset <= b.size {
		// Truncate the data to buffer end if exceeding
		count := uint64(len(data))
		if offset+count > b.size {
			count = b.size - offset
		}

		dest := unsafe.Slice((*byte)(buf), b.size)
		copy(dest[offset:offset+count], data[:count])
	}
}

//
// LuaValue implementation
//

var _ LuaValue = (*LuaBuffer)(nil)

func (b *LuaBuffer) lua() *Lua { return b.vm }
func (b *LuaBuffer) ref() int  { return b.index }
func (b *LuaBuffer) deref(lua *Lua) int {
	return int(ffi.GetRef(lua.state(), int32(b.ref())))
}
