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
	CatchPanics     bool
	EnableCodegen   bool
	EnableSandbox   bool
	Compiler        *Compiler
}

type Lua struct {
	inner          *StateWithMemory
	compiler       *Compiler
	fnRegistry     *functionRegistry
	codegenEnabled bool
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

func (l *Lua) SetCodegen(enabled bool) bool {
	// NOTE: disabling codegen if it was enabled before still has the codegen
	// backend enabled for the state since we already called LuauCodegenCreate
	// during state initialization

	supported := ffi.LuauCodegenSupported()
	if supported {
		l.codegenEnabled = enabled
	}

	return supported
}

func (l *Lua) GetGlobal(global string) LuaValue {
	state := l.state()

	ffi.GetGlobal(state, global)
	value := intoLuaValue(l, -1)

	ffi.Pop(state, 1)
	return value
}

func (l *Lua) SetGlobal(global string, value LuaValue) {
	value.deref(l)
	ffi.SetGlobal(l.state(), global)
}

func (l *Lua) CreateTable() *LuaTable {
	state := l.inner.luaState

	ffi.NewTable(state)
	index := ffi.Ref(state, -1)

	t := &LuaTable{vm: l, index: int(index)}
	runtime.SetFinalizer(t, valueUnrefer[*LuaTable](l))

	return t
}

func (l *Lua) CreateString(str string) *LuaString {
	state := l.inner.luaState

	ffi.PushString(state, str)
	index := ffi.Ref(state, -1)

	s := &LuaString{vm: l, index: int(index)}
	runtime.SetFinalizer(s, valueUnrefer[*LuaString](l))

	return s
}

func (l *Lua) CreateFunction(fn GoFunction) *LuaChunk {
	state := l.state()

	entry := l.fnRegistry.register(fn)
	pushUpvalue(state, entry, registryTrampolineDtor)

	ffi.PushCClosureK(state, registryTrampoline, nil, 1, nil)

	index := ffi.Ref(state, -1)
	c := &LuaChunk{vm: l, index: int(index)}
	runtime.SetFinalizer(c, func(c *LuaChunk) { ffi.Unref(state, index) })

	return c
}

func (l *Lua) CreateUserData(value IntoUserData) *LuaUserData {
	state := l.state()
	userdata := &LuaUserData{vm: l, inner: value}

	// TOOD: custom destructor support
	ud := ffi.NewUserdata(state, uint64(unsafe.Sizeof(uintptr(0))))
	*(*IntoUserData)(unsafe.Pointer(ud)) = value

	if ffi.LNewMetatable(state, "") {
		fieldsMap := newFieldMap()
		methodsMap := newMethodMap(l.fnRegistry)
		metaMethodsMap := newMethodMap(l.fnRegistry)

		value.Fields(fieldsMap)
		value.Methods(methodsMap)
		value.MetaMethods(metaMethodsMap)

		pushUpvalue(state, fieldsMap, fieldMapDtor)
		pushUpvalue(state, methodsMap, methodMapDtor)

		ffi.PushCClosureK(state, indexMt, nil, 2, nil)
		ffi.SetField(state, -2, "__index")

		for method, impl := range metaMethodsMap.inner {
			if method == "__index" {
				panic("Cannot have a manual __index implementation")
			}

			pushUpvalue(state, impl, registryTrampolineDtor)
			ffi.PushCClosureK(state, registryTrampoline, nil, 1, nil)
			ffi.SetField(state, -2, method)
		}
	}

	ffi.SetMetatable(state, -2)

	userdata.index = int(ffi.Ref(state, -1))
	runtime.SetFinalizer(userdata, valueUnrefer[*LuaUserData](l))

	return userdata
}

func (l *Lua) CreateBuffer(size uint64) *LuaBuffer {
	state := l.state()

	ffi.NewBuffer(state, size)
	index := ffi.Ref(state, -1)

	b := &LuaBuffer{vm: l, index: int(index), size: size}
	runtime.SetFinalizer(b, valueUnrefer[*LuaBuffer](l))

	return b
}

func (l *Lua) CreateThread(chunk *LuaChunk) (*LuaThread, error) {
	mainState := l.state()
	threadState := ffi.NewThread(mainState)

	chunk.pushToStack()
	ffi.XMove(mainState, threadState, 1)

	index := ffi.Ref(mainState, -1)
	t := &LuaThread{vm: l, chunk: chunk, index: int(index)}

	runtime.SetFinalizer(t, func(t *LuaThread) {
		ffi.LuaClose(t.state())
		ffi.Unref(l.state(), int32(t.ref()))
	})

	return t, nil
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
		CatchPanics:    true,
		EnableCodegen:  true,
		Compiler:       DefaultCompiler(),
	})
}

func NewWith(libs StdLib, options LuaOptions) *Lua {
	state := newStateWithAllocator(options.InitMemoryState)
	if state == nil {
		panic("Failed to create Lua state")
	}

	ffi.OpenBase(state.luaState)
	luaLibs := map[StdLib]func(*ffi.LuaState){
		StdLibCOROUTINE: ffi.OpenCoroutine,
		StdLibTABLE:     ffi.OpenTable,
		StdLibOS:        ffi.OpenOs,
		StdLibSTRING:    ffi.OpenString,
		StdLibUTF8:      ffi.OpenUtf8,
		StdLibBIT:       ffi.OpenBit32,
		StdLibBUFFER:    ffi.OpenBuffer,
		StdLibMATH:      ffi.OpenMath,
		StdLibDEBUG:     ffi.OpenDebug,
		StdLibVECTOR:    ffi.OpenVector,
	}

	for library, opener := range luaLibs {
		if (!options.IsSafe || StdLibALLSAFE.Contains(library)) && libs.Contains(library) {
			opener(state.luaState)
		}
	}

	if options.EnableSandbox {
		ffi.LSandbox(state.luaState)
	}

	compiler := options.Compiler
	if compiler == nil {
		compiler = DefaultCompiler()
	}

	fnReg := newFunctionRegistry()
	fnReg.recoverPanics = options.CatchPanics

	lua := &Lua{inner: state, compiler: compiler, fnRegistry: fnReg, codegenEnabled: false}
	if options.EnableCodegen && ffi.LuauCodegenSupported() {
		ffi.LuauCodegenCreate(state.luaState)
		lua.codegenEnabled = true
	}

	runtime.SetFinalizer(lua, func(l *Lua) {
		if options.CollectGarbage {
			ffi.LuaGc(l.state(), ffi.LUA_GCCOLLECT, 0)
		}

		l.Close()
	})

	return lua
}

func pushUpvalue[T any](state *ffi.LuaState, ptr *T, dtor unsafe.Pointer) *uintptr {
	var up *uintptr

	sz := uint64(unsafe.Sizeof(uintptr(0)))
	if dtor != nil {
		up = (*uintptr)(ffi.NewUserdataDtor(state, sz, dtor))
	} else {
		up = (*uintptr)(ffi.NewUserdata(state, sz))
	}

	*up = uintptr(cgo.NewHandle(ptr))

	return up
}
