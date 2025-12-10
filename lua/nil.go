package lua

import "github.com/CompeyDev/lei/ffi"

type LuaNil struct{ vm *Lua }

//
// LuaValue Implementation
//

var _ LuaValue = (*LuaNil)(nil)

func (n *LuaNil) lua() *Lua  { return n.vm }
func (n *LuaNil) ref() int   { return ffi.LUA_REFNIL }
func (n *LuaNil) deref() int { return 0 }
