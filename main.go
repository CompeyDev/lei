package main

import (
	"fmt"

	"github.com/CompeyDev/lei/lua"
)

func main() {
	mem := lua.NewMemoryState()
	// mem.SetLimit(250 * 1024) // 250KB max
	state := lua.NewWith(lua.StdLibALLSAFE, lua.LuaOptions{InitMemoryState: mem, CatchPanics: true, EnableCodegen: true})

	table := state.CreateTable()
	key, value := state.CreateString("hello"), state.CreateString("lei")
	table.Set(key, value)

	fmt.Printf("Used: %d, Limit: %d\n", mem.Used(), mem.Limit())

	fmt.Println(key.ToString(), table.Get(key).(*lua.LuaString).ToString())
	chunk, err := state.Load("main", []byte("print('hello, lei!!!!'); return {['mrrp'] = 'foo', ['meow'] = 'bar'}, 'baz'"))
	if err != nil {
		fmt.Println(err)
		return
	}

	values, returnErr := chunk.Call()
	if returnErr != nil {
		fmt.Println(returnErr)
		return
	}

	for i, value := range values {
		fmt.Print(i, ": ")

		switch v := value.(type) {
		case *lua.LuaString:
			fmt.Println(v.ToString())
		case *lua.LuaTable:
			fmt.Println()

			for key, val := range v.Iterable() {
				k, kErr := lua.As[string](key)
				v, vErr := lua.As[string](val)

				if kErr != nil || vErr != nil {
					fmt.Println("  (non-string key or value)")
				}

				fmt.Printf("  %v: %v\n", k, v)
			}
		}
	}

	iterable, iterErr := lua.As[map[string]string](table)
	if iterErr != nil {
		fmt.Println(iterErr)
		return
	}

	for k, v := range iterable { // or, we can use `.Iterable`
		fmt.Printf("%s %s\n", k, v)
	}

	cFnChunk := state.CreateFunction(func(luaState *lua.Lua, args ...lua.LuaValue) ([]lua.LuaValue, error) {
		someNumber := lua.LuaNumber(22713)
		return []lua.LuaValue{
			luaState.CreateString("Hello"),
			luaState.CreateString("from"),
			luaState.CreateString(fmt.Sprintf("Go, %s!", args[0].(*lua.LuaString).ToString())),
			&someNumber,
		}, nil
	})

	returns, callErr := cFnChunk.Call(state.CreateString("Lua"))
	if callErr != nil {
		fmt.Println(callErr)
		return
	}

	for i, ret := range returns {
		str, err := lua.As[string](ret)
		if err == nil {
			fmt.Printf("Return %d: %s\n", i+1, str)
		} else {
			num, _ := lua.As[float64](ret)
			fmt.Printf("Return %d: %f\n", i+1, num)
		}
	}

	class := &Class{value: 420.0}
	classUd := state.CreateUserData(class)
	state.SetGlobal("classUd", classUd)

	got := state.GetGlobal("classUd").(*lua.LuaUserData).Downcast()
	fmt.Println(got.(*Class).value)

	udChunk, udErr := state.Load("udChunk", []byte("print(tostring(classUd), classUd.toggle); classUd.flip(); print(classUd.toggle, classUd.fakeToggle)"))
	if udErr != nil {
		fmt.Println(udErr)
		return
	}

	_, udCallErr := udChunk.Call()
	if udCallErr != nil {
		fmt.Println(udCallErr)
		return
	}
}

type Class struct{ value float64 }

var _ lua.IntoUserData = (*Class)(nil)

func (c *Class) Fields(fields *lua.FieldMap) {
	// NOTE: this references takes a copy of the value and mutations hence do
	// not persist here. Instead we need a getter which captures the class
	// itself
	funnyNumber := lua.LuaNumber(c.value)
	fields.Insert("fakeToggle", &funnyNumber)

	fields.Insert("toggle", func(*lua.Lua) lua.LuaValue {
		value := lua.LuaNumber(c.value)
		return &value
	})
}

func (c *Class) MetaMethods(metaMethods *lua.MethodMap) {
	metaMethods.Insert("__tostring", func(vm *lua.Lua, _ ...lua.LuaValue) ([]lua.LuaValue, error) {
		return []lua.LuaValue{vm.CreateString("Class")}, nil
	})
}

func (c *Class) Methods(methods *lua.MethodMap) {
	methods.Insert("flip", func(_G *lua.Lua, args ...lua.LuaValue) ([]lua.LuaValue, error) {
		c.toggle()
		return []lua.LuaValue{}, nil
	})
}

func (c *Class) toggle() { c.value = 69.0 }
