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
	Compiler        *Compiler
}

type Lua struct {
	inner    *StateWithMemory
	compiler *Compiler
}

func (l *Lua) Execute(name string, input []byte) ([]LuaValue, error) {
	// TODO: create a load function which doesnt execute

	state := l.inner.luaState
	initialStack := ffi.GetTop(state) // Track initial stack size

	if !isBytecode(input) {
		bytecode, err := l.compiler.Compile(string(input))
		if err != nil {
			return nil, err
		}

		input = bytecode
	}

	loadResult := ffi.LuauLoad(state, name, input, uint64(len(input)), 0)
	loadErr := newLoadError(state, int(loadResult))

	if loadErr != nil {
		return nil, loadErr
	}

	execResult := ffi.Pcall(state, 0, -1, 0)
	execErr := newLoadError(state, int(execResult))

	if execErr != nil {
		return nil, execErr
	}

	stackNow := ffi.GetTop(state)
	resultsCount := stackNow - initialStack

	if resultsCount == 0 {
		return nil, nil
	}

	// TODO: contemplate whether to return LuaValues or go values
	results := make([]LuaValue, resultsCount)
	for i := range resultsCount {
		// The stack has grown by the number of returns of the chunk from the
		// initial value tracked at the beginning. We add one to that due to
		// Lua's 1-based indexing system
		stackIndex := int32(initialStack + i + 1)
		results[i] = intoLuaValue(l, stackIndex)
	}

	ffi.Pop(state, resultsCount)

	return results, nil
}

func (l *Lua) Memory() *MemoryState {
	return l.inner.MemState()
}

func (l *Lua) CreateTable() *LuaTable {
	state := l.inner.luaState

	ffi.NewTable(state)
	index := ffi.Ref(state, -1)

	t := &LuaTable{
		vm:    l,
		index: int(index),
	}

	return t
}

func (l *Lua) CreateString(str string) *LuaString {
	state := l.inner.luaState

	ffi.PushString(state, str)

	index := ffi.Ref(state, -1)
	ffi.RawGetI(state, ffi.LUA_REGISTRYINDEX, int32(index))

	ffi.Pop(state, 1)

	s := &LuaString{vm: l, index: int(index)}
	return s
}

func (l *Lua) SetCompiler(compiler *Compiler) {
	l.compiler = compiler
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
		Compiler:       DefaultCompiler(),
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

	compiler := options.Compiler
	if compiler == nil {
		compiler = DefaultCompiler()
	}

	lua := &Lua{inner: state, compiler: compiler}
	runtime.SetFinalizer(lua, func(l *Lua) {
		if options.CollectGarbage {
			ffi.LuaGc(l.state(), ffi.LUA_GCCOLLECT, 0)
		}

		l.Close()
	})

	return lua
}
