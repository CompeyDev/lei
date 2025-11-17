package main

import (
	"fmt"

	"github.com/CompeyDev/lei/lua"
)

func main() {
	state := lua.New()
	memState := lua.GetMemoryState(state.RawState())
	memState.SetMemoryLimit(1) // FIXME: this no workie?

	table := state.CreateTable()
	key, value := state.CreateString("hello"), state.CreateString("world")
	table.Set(&key, &value)

	fmt.Printf("Used: %d, Limit: %d\n", memState.UsedMemory(), memState.MemoryLimit())

	fmt.Println(key.ToString(), table.Get(&key).(*lua.LuaString).ToString())
}
