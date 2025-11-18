package lua

/*
#cgo CFLAGS: -I/usr/lib/gcc/x86_64-pc-linux-gnu/15.2.1/include
#include <stdlib.h>
#include <stdint.h>

extern void* allocator(void* ud, void* ptr, size_t osize, size_t nsize);
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

func NewMemoryState() *MemoryState {
	return &MemoryState{
		usedMemory:   0,
		memoryLimit:  0,
		ignoreLimit:  false,
		limitReached: false,
	}
}

func (m *MemoryState) Used() int {
	return m.usedMemory
}

func (m *MemoryState) Limit() int {
	return m.memoryLimit
}

func (m *MemoryState) SetLimit(limit int) int {
	prevLimit := m.memoryLimit
	m.memoryLimit = limit
	return prevLimit
}

func RelaxLimitWith(state *ffi.LuaState, f func()) {
	memState := getMemoryState(state)
	if memState != nil {
		memState.ignoreLimit = true
		f()
		memState.ignoreLimit = false
	} else {
		f()
	}
}

func LimitReached(state *ffi.LuaState) bool {
	return getMemoryState(state).limitReached
}

func getMemoryState(state *ffi.LuaState) *MemoryState {
	var memState unsafe.Pointer
	ffi.GetAllocF(state, &memState)

	if memState == nil {
		panic("Lua state has no allocator userdata")
	}

	return (*MemoryState)(memState)
}

//export allocator
func allocator(ud, ptr unsafe.Pointer, osize, nsize C.size_t) unsafe.Pointer {
	memState := (*MemoryState)(ud)

	// Avoid GC of pointer for this call period
	runtime.KeepAlive(memState)
	memState.limitReached = false

	// Free memory
	if nsize == 0 {
		if ptr != nil {
			C.free(ptr)
			memState.usedMemory -= int(osize)
		}
		return nil
	}

	if nsize > C.size_t(^uint(0)>>1) {
		return nil
	}

	var memDiff int
	if ptr != nil {
		memDiff = int(nsize) - int(osize)
	} else {
		memDiff = int(nsize)
	}

	memLimit := memState.memoryLimit
	newUsedMemory := memState.usedMemory + memDiff
	if memLimit > 0 && newUsedMemory > memLimit && !memState.ignoreLimit {
		memState.limitReached = true
		panic("allocations exceeded set limit for memory")
	}
	memState.usedMemory = newUsedMemory

	var newPtr unsafe.Pointer
	if ptr == nil {
		newPtr = C.malloc(nsize)
		if newPtr == nil {
			panic("memory allocation failed")
		}
	} else {
		newPtr = C.realloc(ptr, nsize)
		if newPtr == nil {
			panic("memory reallocation failed")
		}
	}

	return newPtr
}

type StateWithMemory struct {
	luaState *ffi.LuaState
	memState *MemoryState
	pinner   *runtime.Pinner
}

func newStateWithAllocator(initState *MemoryState) *StateWithMemory {
	var memState *MemoryState
	if initState != nil {
		memState = initState
	} else {
		memState = NewMemoryState()
	}

	// Pin the memory state to prevent GC from moving it
	pinner := &runtime.Pinner{}
	pinner.Pin(memState)

	state := ffi.NewState(C.allocator, unsafe.Pointer(memState))

	return &StateWithMemory{
		luaState: state,
		memState: memState,
		pinner:   pinner,
	}
}

func (s *StateWithMemory) LuaState() *ffi.LuaState {
	return s.luaState
}

func (s *StateWithMemory) MemState() *MemoryState {
	return s.memState
}

func (s *StateWithMemory) Close() {
	if s.pinner != nil {
		s.pinner.Unpin()
	}

	ffi.LuaClose(s.luaState)
}
