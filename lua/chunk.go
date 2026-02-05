package lua

import "github.com/CompeyDev/lei/ffi"

// NOTE: `bytecode` and `index` are expected to be mutually exclusive

type LuaChunk struct {
	vm *Lua

	name     *string
	bytecode []byte

	index int
}

func (c *LuaChunk) Call(args ...LuaValue) ([]LuaValue, error) {
	state := c.vm.state()

	initialStack := ffi.GetTop(state) // Track initial stack size
	c.pushToStack()

	argsCount := len(args)
	if c.bytecode == nil {
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

	if c.bytecode == nil {
		// Chunk is of a C function, need to deref
		ffi.GetRef(state, int32(c.index))
	} else {
		// Chunk is bytecode, load it into the VM
		hasLoaded := ffi.LuauLoad(state, *c.name, c.bytecode, uint64(len(c.bytecode)), 0)
		if !hasLoaded {
			// Miscellaneous error is denoted with a -1 code
			return &LuaError{Code: -1, Message: ffi.ToLString(state, -1, nil)}
		}

		// Apply native code generation if requested
		if ffi.LuauCodegenSupported() && c.vm.codegenEnabled {
			ffi.LuauCodegenCompile(state, -1)
		}
	}

	return nil
}
