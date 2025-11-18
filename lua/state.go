package lua

import (
	"runtime"
	"unsafe"

	"github.com/CompeyDev/lei/ffi"
)

type LuaOptions struct {
	InitMemoryState *MemoryState
	CollectGarbage  bool
	IsSafe          bool
}

type Lua struct {
	inner *StateWithMemory
}

func (l *Lua) Memory() *MemoryState {
	return l.inner.MemState()
}

func (l *Lua) CreateTable() LuaTable {
	state := l.inner.luaState
	ffi.NewTable(state)

	return LuaTable{
		vm:    l,
		index: int(ffi.GetTop(state)),
	}
}

func (l *Lua) CreateString(str string) LuaString {
	state := l.inner.luaState

	ffi.PushString(state, str)
	index := ffi.GetTop(state)
	return LuaString{vm: l, index: int(index)}
}

func (l *Lua) Close() {
	l.inner.Close()
}

func (l *Lua) state() *ffi.LuaState {
	return l.inner.luaState
}

func New() *Lua {
	return NewWith(StdLibALLSAFE, LuaOptions{
		CollectGarbage: true,
		IsSafe:         true,
	})
}

func NewWith(libs StdLib, options LuaOptions) *Lua {
	if libs.Contains(StdLibPACKAGE) {
		// TODO: disable c modules for package lib
	}

	state := newStateWithAllocator(options.InitMemoryState)
	if state == nil {
		panic("Failed to create Lua state")
	}

	ffi.RequireLib(state.luaState, "_G", unsafe.Pointer(ffi.BaseOpener()), true)
	ffi.Pop(state.luaState, 1)

	// TODO: luau jit stuff

	type Library struct {
		lib  StdLib
		name string
	}

	luaLibs := map[Library]unsafe.Pointer{
		{StdLibCOROUTINE, ffi.LUA_COLIBNAME}:  unsafe.Pointer(ffi.CoroutineOpener()),
		{StdLibTABLE, ffi.LUA_TABLIBNAME}:     unsafe.Pointer(ffi.TableOpener()),
		{StdLibOS, ffi.LUA_OSLIBNAME}:         unsafe.Pointer(ffi.OsOpener()),
		{StdLibSTRING, ffi.LUA_STRLIBNAME}:    unsafe.Pointer(ffi.StringOpener()),
		{StdLibUTF8, ffi.LUA_UTF8LIBNAME}:     unsafe.Pointer(ffi.Utf8Opener()),
		{StdLibBIT, ffi.LUA_BITLIBNAME}:       unsafe.Pointer(ffi.Bit32Opener()),
		{StdLibBUFFER, ffi.LUA_BUFFERLIBNAME}: unsafe.Pointer(ffi.BufferOpener()),
		// TODO: vector lib
		{StdLibMATH, ffi.LUA_MATHLIBNAME}: unsafe.Pointer(ffi.MathOpener()),
		{StdLibBUFFER, ffi.LUA_DBLIBNAME}: unsafe.Pointer(ffi.DebugOpener()),
		// TODO: package lib
	}

	for library, open := range luaLibs {
		// FIXME: check safety here maybe?

		if libs.Contains(library.lib) {
			ffi.RequireLib(state.luaState, library.name, unsafe.Pointer(open), true)
		}
	}

	lua := &Lua{inner: state}
	runtime.SetFinalizer(lua, func(l *Lua) {
		if options.CollectGarbage {
			ffi.LuaGc(l.state(), ffi.LUA_GCCOLLECT, 0)
		}

		l.Close()
	})

	return lua
}
