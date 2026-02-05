package lua

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/CompeyDev/lei/ffi"
)

type LuaValue interface {
	// Optionally returns the Lua VM this value belongs to
	lua() *Lua
	// Returns the reference index of this value in the Lua registry
	ref() int
	// Dereferences this value onto the Lua stack, returning the stack index
	deref(*Lua) int
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
	case *LuaNumber:
		// Map of all numeric types for O(1) lookup
		var numericKinds = map[reflect.Kind]bool{
			reflect.Int: true, reflect.Int8: true, reflect.Int16: true, reflect.Int32: true, reflect.Int64: true,
			reflect.Uint: true, reflect.Uint8: true, reflect.Uint16: true, reflect.Uint32: true, reflect.Uint64: true,
			reflect.Uintptr: true,
			reflect.Float32: true, reflect.Float64: true,
			reflect.Complex64: true, reflect.Complex128: true,
		}

		if kind := t.Kind(); numericKinds[kind] {
			num := reflect.ValueOf(*val).Convert(t)
			return num, nil
		}

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
			fieldSet := make(map[int]bool)

			for key, value := range val.Iterable() {
				keyStr, ok := key.(*LuaString)
				if !ok {
					continue
				}

				luaKey := keyStr.ToString()
				var field reflect.Value
				var found bool
				var priority int // 0 = explicit annotation, 1 = direct match, 2 = lowercase fallback
				var fieldIndex int

				for i := 0; i < t.NumField(); i++ {
					// Annotation-based field name overrides (eg: `lua:"field_name"`)
					structField := t.Field(i)
					tagVal, ok := structField.Tag.Lookup("lua")
					if ok && tagVal == luaKey {
						field = res.Field(i)
						found, priority = true, 0
						fieldIndex = i
						break
					}

					// Exact matches
					if structField.Name == luaKey {
						if !found || priority > 1 {
							field = res.Field(i)
							found, priority = true, 1
							fieldIndex = i
						}
						continue
					}

					// If field is exported, try also using lowercase first character
					if name := structField.Name; structField.IsExported() {
						lower := strings.ToLower(name[:1]) + name[1:]
						if lower == luaKey {
							if !found || priority > 2 {
								field = res.Field(i)
								found, priority = true, 2
								fieldIndex = i
							}
						}
					}
				}

				if found && field.IsValid() && field.CanSet() {
					// We keep track of whether the field has been found, its priority, and the
					// index at which it was found within the struct. If there is an explicit
					// annotation, we set the field value directly, otherwise we check that
					// the field hasn't already been set in another match, and only set it then
					if !fieldSet[fieldIndex] || priority == 0 {
						// Recursively convert value to a reflect value
						vVal, err := asReflectValue(value, field.Type())
						if err != nil {
							return zero, err
						}

						field.Set(vVal)
						fieldSet[fieldIndex] = true
					}
				}
			}

			return res, nil

		}

	case *LuaNil:
		return zero, nil

	case *LuaUserData:
		if downcasted := val.Downcast(); downcasted != nil {
			return reflect.ValueOf(downcasted).Convert(t), nil
		}

		return zero, fmt.Errorf("value isn't userdata")

	case *LuaVector:
		return reflect.ValueOf(v), nil

	case *LuaBuffer:
		kind := t.Kind()
		if kind == reflect.Array {
			return reflect.ValueOf(val.Read(0, uint64(t.Len()))).Convert(t), nil
		}

		if kind == reflect.Slice {
			return reflect.ValueOf(val.Read(0, val.size)).Convert(t), nil
		}
	}

	return zero, fmt.Errorf("cannot convert LuaValue(%T) into %T", v, zero.Type().Name())
}

func intoLuaValue(lua *Lua, index int32) LuaValue {
	state := lua.state()

	switch ffi.Type(state, index) {
	case ffi.LUA_TNUMBER:
		num := ffi.ToNumber(state, index)
		li := LuaNumber(float64(num))
		return &li
	case ffi.LUA_TSTRING:
		ref := ffi.Ref(state, index)
		return &LuaString{vm: lua, index: int(ref)}
	case ffi.LUA_TTABLE:
		ref := ffi.Ref(state, index)
		return &LuaTable{vm: lua, index: int(ref)}
	case ffi.LUA_TNIL:
		return &LuaNil{}
	case ffi.LUA_TUSERDATA:
		ref := ffi.Ref(state, index)
		return &LuaUserData{vm: lua, index: int(ref)}
	case ffi.LUA_TVECTOR:
		x, y, z := ffi.ToVector(state, index)
		return &LuaVector{*x, *y, *z}
	case ffi.LUA_TBUFFER:
		ref := ffi.Ref(state, index)
		return &LuaBuffer{vm: lua, index: int(ref), size: ffi.ObjLen(state, ref)}
	case ffi.LUA_TTHREAD:
		ref := ffi.Ref(state, index)
		return &LuaThread{vm: lua, index: int(ref)} // NOTE: no chunk, can only be executed once
	default:
		panic("unsupported Lua type")
	}
}

func valueUnrefer[T LuaValue](lua *Lua) func(T) {
	return func(value T) {
		ffi.Unref(lua.state(), int32(value.ref()))
	}
}
