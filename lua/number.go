package lua

import "github.com/CompeyDev/lei/ffi"

type LuaNumber struct {
	vm    *Lua
	inner float64
}

//
// LuaValue implementation
//

var _ LuaValue = (*LuaNumber)(nil)

// Numbers are cheap to copy, so we don't store the reference index

func (n *LuaNumber) lua() *Lua { return n.vm }
func (n *LuaNumber) ref() int  { return ffi.LUA_NOREF }
func (n *LuaNumber) deref() int {
	state := n.vm.state()

	ffi.PushNumber(state, ffi.LuaNumber(n.inner))
	return int(ffi.GetTop(state))
}
