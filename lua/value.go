package lua

import (
	"fmt"
	"reflect"

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

func As[T any](v LuaValue) (T, error) {
	var zero T

	targetType := reflect.TypeOf(zero)
	reflectValue, err := asReflectValue(v, targetType)

	return reflectValue.Interface().(T), err
}

func asReflectValue(v LuaValue, t reflect.Type) (reflect.Value, error) {
	// TODO: allow annotations to override field names

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

				field := res.FieldByName(keyStr.ToString())
				if !field.IsValid() || !field.CanSet() {
					continue
				}

				vVal, err := asReflectValue(value, field.Type())
				if err != nil {
					return zero, err
				}

				field.Set(vVal)
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
		return &LuaNil{}
	default:
		panic("unsupported Lua type")
	}
}
