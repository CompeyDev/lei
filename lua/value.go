package lua

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/CompeyDev/lei/ffi"
)

type LuaValue interface {
	lua() *Lua
	ref() int
	deref() int
}

func TypeName(val LuaValue) string {
	lua := val.lua().state()
	return ffi.TypeName(lua, ffi.Type(lua, int32(val.ref())))
}

//
// Lua<->Go Type Conversion
//

func As[T any](v LuaValue) (T, error) {
	var zero T

	targetType := reflect.TypeOf(zero)
	reflectValue, err := asReflectValue(v, targetType)

	return reflectValue.Interface().(T), err
}

func asReflectValue(v LuaValue, t reflect.Type) (reflect.Value, error) {
	zero := reflect.Zero(t)

	switch val := v.(type) {
	case *LuaString:
		if t.Kind() == reflect.String {
			str := reflect.ValueOf(val.ToString()).Convert(t)
			return str, nil
		}

	case *LuaTable:
		switch t.Kind() {
		case reflect.Map:
			res := reflect.MakeMap(t)
			for key, value := range val.Iterable() {
				var kVal reflect.Value
				var vVal reflect.Value
				var err error

				// Key conversion
				if t.Key() == reflect.TypeOf((*LuaValue)(nil)).Elem() {
					kVal = reflect.ValueOf(key)
				} else {
					kVal, err = asReflectValue(key, t.Key())
					if err != nil {
						return zero, err
					}
				}

				// Value conversion
				if t.Elem() == reflect.TypeOf((*LuaValue)(nil)).Elem() {
					vVal = reflect.ValueOf(value)
				} else {
					vVal, err = asReflectValue(value, t.Elem())
					if err != nil {
						return zero, err
					}
				}

				res.SetMapIndex(kVal, vVal)
			}

			return res, nil

		case reflect.Struct:
			res := reflect.New(t).Elem()
			for key, value := range val.Iterable() {
				keyStr, ok := key.(*LuaString)
				if !ok {
					continue
				}

				luaKey := keyStr.ToString()
				var field reflect.Value
				var found bool

				for i := 0; i < t.NumField(); i++ {
					// Annotation-based field name overrides (eg: `lua:"field_name"`)
					structField := t.Field(i)
					tagVal, ok := structField.Tag.Lookup("lua")
					if ok && tagVal == luaKey {
						field = res.Field(i)
						found = true
						break
					}

					// Exact matches
					if structField.Name == luaKey {
						field = res.Field(i)
						found = true
						break
					}

					// If field is exported, try also using lowercase first character
					if name := structField.Name; structField.IsExported() {
						lower := strings.ToLower(name[:1]) + name[1:]
						if lower == luaKey {
							field = res.Field(i)
							found = true
							break
						}
					}
				}

				if found && field.IsValid() && field.CanSet() {
					// Recursively convert value to a reflect value
					vVal, err := asReflectValue(value, field.Type())
					if err != nil {
						return zero, err
					}

					field.Set(vVal)
				}
			}

			return res, nil

		}

	case *LuaNil:
		return zero, nil

	}

	return zero, fmt.Errorf("cannot convert LuaValue(%T) into %T", v, zero)
}

func intoLuaValue(lua *Lua, index int32) LuaValue {
	state := lua.state()

	switch ffi.Type(state, index) {
	case ffi.LUA_TSTRING:
		ref := ffi.Ref(state, index)
		return &LuaString{vm: lua, index: int(ref)}
	case ffi.LUA_TTABLE:
		ref := ffi.Ref(state, index)
		return &LuaTable{vm: lua, index: int(ref)}
	case ffi.LUA_TNIL:
		return &LuaNil{vm: lua}
	default:
		panic("unsupported Lua type")
	}
}

func valueUnrefer[T LuaValue](lua *Lua) func(T) {
	return func(value T) {
		ffi.Unref(lua.state(), int32(value.ref()))
	}
}
