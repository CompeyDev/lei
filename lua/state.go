package lua

import (
	"runtime"
	"runtime/cgo"
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
	inner      *StateWithMemory
	compiler   *Compiler
	fnRegistry *functionRegistry
}

func (l *Lua) Load(name string, input []byte) (*LuaChunk, error) {
	chunk := &LuaChunk{vm: l, bytecode: input}
	if !isBytecode(input) {
		bytecode, err := l.compiler.Compile(string(input))
		if err != nil {
			return nil, err
		}

		chunk.bytecode = bytecode
	}

	return chunk, nil
}

func (l *Lua) Memory() *MemoryState {
	return l.inner.MemState()
}

func (l *Lua) CreateTable() *LuaTable {
	state := l.inner.luaState

	ffi.NewTable(state)
	index := ffi.Ref(state, -1)

	t := &LuaTable{vm: l, index: int(index)}
	runtime.SetFinalizer(t, valueUnrefer[*LuaTable](t.lua()))

	return t
}

func (l *Lua) CreateString(str string) *LuaString {
	state := l.inner.luaState

	ffi.PushString(state, str)
	index := ffi.Ref(state, -1)

	s := &LuaString{vm: l, index: int(index)}
	runtime.SetFinalizer(s, valueUnrefer[*LuaString](s.lua()))

	return s
}

func (l *Lua) CreateFunction(fn GoFunction) *LuaChunk {
	state := l.state()

	entry := l.fnRegistry.register(fn)
	handle := cgo.NewHandle(entry)

	ud := (*uintptr)(ffi.NewUserdataDtor(state, uint64(unsafe.Sizeof(uintptr(0))), registryTrampolineDtor))
	*ud = uintptr(handle)

	ffi.PushCClosureK(state, registryTrampoline, nil, 1, nil)

	index := ffi.Ref(state, -1)
	c := &LuaChunk{vm: l, index: int(index)}
	runtime.SetFinalizer(c, func(c *LuaChunk) { ffi.Unref(state, index) })

	return c
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

	lua := &Lua{inner: state, compiler: compiler, fnRegistry: newFunctionRegistry()}
	runtime.SetFinalizer(lua, func(l *Lua) {
		if options.CollectGarbage {
			ffi.LuaGc(l.state(), ffi.LUA_GCCOLLECT, 0)
		}

		l.Close()
	})

	return lua
}
