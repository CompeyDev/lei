package lua

import "github.com/CompeyDev/lei/ffi"

type ChunkMode int

const (
	// Raw text source code that must be compiled before executing
	ChunkModeSOURCE = iota

	// Compiled bytecode that can be directly executed
	ChunkModeBYTECODE

	// A C function pointer loaded onto the stack
	ChunkModeFUNCTION
)

type LuaChunk struct {
	vm   *Lua
	env  *LuaTable
	mode ChunkMode

	// Values only applicable for source or bytecode types
	name     *string
	data     []byte
	compiler *Compiler

	// An index is held for chunks of the function type
	index int
}

func (c *LuaChunk) Environment() *LuaTable       { return c.env }
func (c *LuaChunk) SetEnvironment(env *LuaTable) { c.env = env }

func (c *LuaChunk) Mode() ChunkMode        { return c.mode }
func (c *LuaChunk) SetMode(mode ChunkMode) { c.mode = mode }

func (c *LuaChunk) Compiler() *Compiler            { return c.compiler }
func (c *LuaChunk) SetCompiler(compiler *Compiler) { c.compiler = compiler }

func (c *LuaChunk) Call(args ...LuaValue) ([]LuaValue, error) {
	state := c.vm.state()

	initialStack := ffi.GetTop(state) // Track initial stack size
	c.pushToStack()

	argsCount := len(args)
	if c.mode == ChunkModeFUNCTION {
		// Chunk is a C function, push length and args
		ffi.PushNumber(state, ffi.LuaNumber(argsCount))
		argsCount++
		for _, arg := range args {
			arg.deref(c.vm)
		}
	}

	status := ffi.Pcall(state, int32(argsCount), -1, 0)
	if status != ffi.LUA_OK {
		return nil, newLuaError(state, int(status))
	}

	stackNow := ffi.GetTop(state)
	resultsCount := stackNow - initialStack

	if resultsCount == 0 {
		return nil, nil
	}

	results := make([]LuaValue, resultsCount)
	for i := range resultsCount {
		// The stack has grown by the number of returns of the chunk from the
		// initial value tracked at the beginning. We add one to that due to
		// Lua's 1-based indexing system
		stackIndex := int32(initialStack + i + 1)
		results[i] = intoLuaValue(c.vm, stackIndex)
	}

	return results, nil
}

func (c *LuaChunk) pushToStack() error {
	state := c.vm.state()

	if c.data == nil {
		// Chunk is of a C function, need to deref
		ffi.GetRef(state, int32(c.index))
	} else {
		// Chunk is bytecode, load it into the VM
		var bytecode []byte

		if c.mode == ChunkModeSOURCE {
			// Need to compile
			var err error
			if bytecode, err = c.compiler.Compile(string(c.data)); err != nil {
				return err
			}
		} else {
			// Already compiled
			bytecode = c.data
		}

		hasLoaded := ffi.LuauLoad(state, *c.name, bytecode, uint64(len(bytecode)), 0)
		if !hasLoaded {
			// Miscellaneous error is denoted with a -1 code
			return &LuaError{Code: -1, Message: ffi.ToLString(state, -1, nil)}
		}

		// Apply native code generation if requested
		if ffi.LuauCodegenSupported() && c.vm.codegenEnabled {
			ffi.LuauCodegenCompile(state, -1)
		}
	}

	if c.env != nil {
		// If a custom environment was provided, set it for the loaded value
		c.env.deref(c.vm)
		if ok := ffi.Setfenv(c.vm.state(), -2); !ok {
			return &LuaError{Code: -1, Message: "Failed to set environment for chunk"}
		}
	}

	return nil
}
