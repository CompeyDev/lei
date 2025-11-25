package lua

type LuaNil struct{ vm *Lua }

//
// LuaValue Implementation
//

func (n *LuaNil) lua() *Lua  { return n.vm }
func (n *LuaNil) ref() int   { return 0 }
func (n *LuaNil) deref() int { return 0 }
