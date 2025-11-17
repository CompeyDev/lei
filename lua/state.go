package lua

import (
	"unsafe"

	"github.com/CompeyDev/lei/ffi"
)

type LuaOptions struct {
	collectGarbage bool
	isSafe         bool
}

type Lua struct {
	state   *ffi.LuaState
	options LuaOptions
}

func (l *Lua) RawState() *ffi.LuaState {
	return l.state
}

func (l *Lua) CreateTable() LuaTable {
	state := l.state
	ffi.NewTable(state)

	return LuaTable{
		lua:   l,
		index: int(ffi.GetTop(state)),
	}
}

func (l *Lua) CreateString(str string) LuaString {
	state := l.state

	ffi.PushString(state, str)
	index := ffi.GetTop(state)
	return LuaString{lua: l, index: int(index)}
}

func New() *Lua {
	return NewWith(StdLibALLSAFE, LuaOptions{
		collectGarbage: true,
		isSafe:         true,
	})
}

func NewWith(libs StdLib, options LuaOptions) *Lua {
	if libs.Contains(StdLibPACKAGE) {
		// TODO: disable c modules for package lib
	}

	state := newStateWithAllocator()
	if state == nil {
		panic("Failed to create Lua state")
	}

	ffi.RequireLib(state, "_G", unsafe.Pointer(ffi.BaseOpener()), true)
	ffi.Pop(state, 1)

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
			ffi.RequireLib(state, library.name, unsafe.Pointer(open), true)
		}
	}

	// TODO: set finalizer to collect garbage if collectGarbage = true
	return &Lua{
		state:   state,
		options: options,
	}
}
