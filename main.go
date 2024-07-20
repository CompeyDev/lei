package main

import lualib "github.com/CompeyDev/lei/internal"

func main() {
	lua := lualib.LNewState()
	println("Lua VM Address: ", lua)

	lualib.PushCFunction(lua, func(L *lualib.LuaState) int32 {
		println("hi from closure?")
		return 0
	})

	lualib.PushString(lua, "123")
	lualib.PushNumber(lua, lualib.ToNumber(lua, 2))

	if !lualib.IsCFunction(lua, 1) {
		panic("CFunction was not correctly pushed onto stack")
	}

	if !lualib.IsNumber(lua, 3) {
		panic("Number was not correctly pushed onto stack")
	}
}
