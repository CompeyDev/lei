package ffi

/*
#cgo CFLAGS: -Iluau/Common/include
#include "Luau/Bytecode.h"

enum LuauBytecodeTag LuauBytecodeTag;
*/
import "C"

//
// Version Constants
//

const (
	LBC_VERSION_MIN    = C.LBC_VERSION_MIN
	LBC_VERSION_MAX    = C.LBC_VERSION_MAX
	LBC_VERSION_TARGET = C.LBC_VERSION_TARGET

	LBC_TYPE_VERSION_MIN    = C.LBC_TYPE_VERSION_MIN
	LBC_TYPE_VERSION_MAX    = C.LBC_TYPE_VERSION_MAX
	LBC_TYPE_VERSION_TARGET = C.LBC_TYPE_VERSION_TARGET
)
