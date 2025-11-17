package lua

/*
#cgo CFLAGS: -I/usr/lib/gcc/x86_64-pc-linux-gnu/15.2.1/include
#include <stdlib.h>

void* allocator(void* ud, void* ptr, size_t osize, size_t nsize);
*/
import "C"
import (
	"runtime"
	"unsafe"

	"github.com/CompeyDev/lei/ffi"
)

const SYS_MIN_ALIGN = unsafe.Sizeof(uintptr(0)) * 2

type MemoryState struct {
	usedMemory   int
	memoryLimit  int
	ignoreLimit  bool
	limitReached bool
}

func newMemoryState() *MemoryState {
	return &MemoryState{
		usedMemory:   0,
		memoryLimit:  0,
		ignoreLimit:  false,
		limitReached: false,
	}
}

func GetMemoryState(state *ffi.LuaState) *MemoryState {
	var memState unsafe.Pointer
	ffi.GetAllocF(state, &memState)

	if memState == nil {
		panic("Luau state has no allocator userdata")
	}

	return (*MemoryState)(memState)
}

func (m *MemoryState) UsedMemory() int {
	return m.usedMemory
}

func (m *MemoryState) MemoryLimit() int {
	return m.memoryLimit
}

func (m *MemoryState) SetMemoryLimit(limit int) int {
	prevLimit := m.memoryLimit
	m.memoryLimit = limit
	return prevLimit
}

func RelaxLimitWith(state *ffi.LuaState, f func()) {
	memState := GetMemoryState(state)
	if memState != nil {
		memState.ignoreLimit = true
		f()
		memState.ignoreLimit = false
	} else {
		f()
	}
}

func LimitReached(state *ffi.LuaState) bool {
	return GetMemoryState(state).limitReached
}

func newStateWithAllocator() *ffi.LuaState {
	memState := newMemoryState()
	state := ffi.NewState(C.allocator, unsafe.Pointer(memState))
	runtime.KeepAlive(memState)

	return state
}
