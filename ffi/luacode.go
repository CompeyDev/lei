package ffi

//go:generate go run ../build buildProject Luau.VM Luau.Compiler Luau.Ast

/*
#cgo CFLAGS: -Iluau/Compiler/include
#cgo LDFLAGS: -L_obj -lLuau.Compiler -lLuau.Ast -lm -lstdc++
#include <stdlib.h>
#include <luacode.h>
*/
import "C"
import "unsafe"

type CompileConstant *C.void

type CompileOptions struct {
	OptimizationLevel int
	DebugLevel        int
	TypeInfoLevel     int
	CoverageLevel     int

	VectorLib  string
	VectorCtor string
	VectorType string

	MutableGlobals []string
	UserdataTypes  []string

	LibrariesWithKnownMembers []string
	LibraryMemberTypeCb       unsafe.Pointer
	LibraryMemberConstantCb   unsafe.Pointer

	DisabledBuiltins []string
}

func LuauCompile(source string, size int, options *CompileOptions, outsize *int) []byte {
	var goArrToC = func(goArr []string) **C.char {
		if len(goArr) == 0 {
			return nil
		}

		// Allocate space for N+1 pointers (extra for NULL terminator)
		arr := C.malloc(C.size_t(len(goArr)+1) * C.size_t(unsafe.Sizeof(uintptr(0))))
		slice := (*[1 << 30]*C.char)(arr)[: len(goArr)+1 : len(goArr)+1]

		for i, s := range goArr {
			slice[i] = C.CString(s)
		}
		slice[len(goArr)] = nil // NULL terminator
		return (**C.char)(arr)
	}

	var freeCArr = func(arr **C.char) {
		if arr == nil {
			return
		}
		// Free strings until we hit NULL
		for i := 0; ; i++ {
			ptr := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(arr)) + uintptr(i)*unsafe.Sizeof(uintptr(0))))
			if ptr == nil {
				break
			}
			C.free(unsafe.Pointer(ptr))
		}
		C.free(unsafe.Pointer(arr))
	}

	csource := C.CString(source)
	coutsize := C.size_t(*outsize)
	coptions := (*C.lua_CompileOptions)(C.malloc(C.sizeof_lua_CompileOptions))

	coptions.optimizationLevel = C.int(options.OptimizationLevel)
	coptions.debugLevel = C.int(options.DebugLevel)
	coptions.typeInfoLevel = C.int(options.TypeInfoLevel)
	coptions.coverageLevel = C.int(options.CoverageLevel)

	coptions.vectorLib = C.CString(options.VectorLib)
	coptions.vectorCtor = C.CString(options.VectorCtor)
	coptions.vectorType = C.CString(options.VectorType)

	coptions.mutableGlobals = goArrToC(options.MutableGlobals)
	coptions.userdataTypes = goArrToC(options.UserdataTypes)
	coptions.librariesWithKnownMembers = goArrToC(options.LibrariesWithKnownMembers)

	coptions.libraryMemberTypeCb = C.lua_LibraryMemberTypeCallback(options.LibraryMemberTypeCb)
	coptions.libraryMemberConstantCb = C.lua_LibraryMemberConstantCallback(options.LibraryMemberConstantCb)

	coptions.disabledBuiltins = goArrToC(options.DisabledBuiltins)

	defer C.free(unsafe.Pointer(csource))
	defer C.free(unsafe.Pointer(coptions.vectorLib))
	defer C.free(unsafe.Pointer(coptions.vectorCtor))
	defer C.free(unsafe.Pointer(coptions.vectorType))
	defer C.free(unsafe.Pointer(coptions))

	defer freeCArr(coptions.mutableGlobals)
	defer freeCArr(coptions.userdataTypes)
	defer freeCArr(coptions.librariesWithKnownMembers)
	defer freeCArr(coptions.disabledBuiltins)

	bytecode := C.luau_compile(csource, C.size_t(size), coptions, &coutsize)
	defer C.free(unsafe.Pointer(bytecode))

	*outsize = int(coutsize)
	result := make([]byte, coutsize)

	copy(result, (*[1 << 30]byte)(unsafe.Pointer(bytecode))[:coutsize:coutsize])

	return result
}

func LuauSetCompileConstantNil(constant unsafe.Pointer) {
	C.luau_set_compile_constant_nil((*C.lua_CompileConstant)(constant))
}

func LuauSetCompileConstantBoolean(constant unsafe.Pointer, b bool) {
	var cBool C.int
	if b {
		cBool = 1
	} else {
		cBool = 0
	}
	C.luau_set_compile_constant_boolean((*C.lua_CompileConstant)(constant), cBool)
}

func LuauSetCompileConstantNumber(constant unsafe.Pointer, n float64) {
	C.luau_set_compile_constant_number((*C.lua_CompileConstant)(constant), C.double(n))
}

func LuauSetCompileConstantVector(constant unsafe.Pointer, x, y, z, w float32) {
	C.luau_set_compile_constant_vector(
		(*C.lua_CompileConstant)(constant),
		C.float(x),
		C.float(y),
		C.float(z),
		C.float(w),
	)
}

func LuauSetCompileConstantString(constant unsafe.Pointer, s string) {
	if len(s) == 0 {
		C.luau_set_compile_constant_string((*C.lua_CompileConstant)(constant), nil, 0)
		return
	}

	bytes := []byte(s)
	ptr := (*C.char)(unsafe.Pointer(&bytes[0]))
	size := C.size_t(len(s))

	C.luau_set_compile_constant_string((*C.lua_CompileConstant)(constant), ptr, size)
}

func LuauSetCompileConstantStringBytes(constant unsafe.Pointer, data []byte) {
	if len(data) == 0 {
		C.luau_set_compile_constant_string((*C.lua_CompileConstant)(constant), nil, 0)
		return
	}

	ptr := (*C.char)(unsafe.Pointer(&data[0]))
	size := C.size_t(len(data))

	C.luau_set_compile_constant_string((*C.lua_CompileConstant)(constant), ptr, size)
}
