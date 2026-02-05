package lua

import "github.com/CompeyDev/lei/ffi"

type LuaThread struct {
	vm    *Lua
	chunk *LuaChunk
	index int
}

func (t *LuaThread) State() *ffi.LuaState {
	state := t.vm.state()

	t.deref(t.vm)
	defer ffi.Pop(state, 1)

	return ffi.ToThread(state, -1)
}

func (t *LuaThread) Resume() ([]LuaValue, error) {
	threadState := t.State()
	t.pushMainFunction()

	status := int(ffi.Resume(threadState, nil, 0))
	return t.collectResults(threadState, status)
}

func (t *LuaThread) ResumeWith(args ...LuaValue) ([]LuaValue, error) {
	mainState := t.vm.state()
	threadState := t.State()

	// Push the function if required
	t.pushMainFunction()

	// Push args length and then the args
	argsCount := len(args)
	ffi.PushNumber(threadState, ffi.LuaNumber(argsCount))

	for _, arg := range args {
		arg.deref(t.vm)
		ffi.XMove(mainState, threadState, 1)
	}

	status := int(ffi.Resume(threadState, nil, int32(argsCount+1))) // +1 for count arg
	return t.collectResults(threadState, status)
}

func (t *LuaThread) collectResults(threadState *ffi.LuaState, status int) ([]LuaValue, error) {
	if status != ffi.LUA_OK && status != ffi.LUA_YIELD {
		// Return error if thread did not run successfully
		return nil, newLuaError(threadState, status)
	}

	nresults := int(ffi.GetTop(threadState))
	if nresults == 0 {
		return nil, nil
	}

	mainState := t.vm.state()
	results := make([]LuaValue, nresults)

	// Push arguments onto main thread and ref them into LuaValues
	for i := range nresults {
		ffi.PushValue(threadState, int32(i+1))
		ffi.XMove(threadState, mainState, 1)

		results[i] = intoLuaValue(t.vm, int32(ffi.GetTop(mainState)))
	}

	return results, nil
}

func (t *LuaThread) pushMainFunction() {
	if threadState := t.State(); t.Status() == ffi.LUA_OK && t.chunk != nil {
		// Reset the thread and push the coroutine function if the thread has
		// finished running and returned a non-resumable state
		ffi.ResetThread(threadState)
		t.chunk.pushToStack()
		ffi.XMove(t.vm.state(), threadState, 1)
	}
}

func (t *LuaThread) Status() int {
	threadState := t.State()
	return int(ffi.Status(threadState))
}

func (t *LuaThread) IsYielded() bool {
	return t.Status() == ffi.LUA_YIELD
}

func (t *LuaThread) IsFinished() bool {
	status := t.Status()
	threadState := t.State()
	return status == ffi.LUA_OK && ffi.GetTop(threadState) == 0
}

//
// LuaValue implementation
//

var _ LuaValue = (*LuaThread)(nil)

func (t *LuaThread) lua() *Lua { return t.vm }
func (t *LuaThread) ref() int  { return t.index }

func (t *LuaThread) deref(lua *Lua) int {
	return int(ffi.GetRef(lua.state(), int32(t.ref())))
}
