//! [`lua.h`][h] - The Luau VM.
//!
//! [h]: https://github.com/luau-lang/luau/blob/master/VM/include/lua.h

use std::ffi::{c_char, c_double, c_float, c_int, c_uint, c_void};
use std::ptr;

use crate::luaconf;

/// Used to get all the results from a function call in [`lua_call`] and [`lua_pcall`].
pub const LUA_MULTRET: c_int = -1;

/// Pseudo-index for the registry.
pub const LUA_REGISTRYINDEX: c_int = -luaconf::LUAI_MAXCSTACK - 2000;
/// Pseudo-index for the environment of the running C function.
pub const LUA_ENVIRONINDEX: c_int = -luaconf::LUAI_MAXCSTACK - 2001;
/// Pseudo-index for the thread-environment.
pub const LUA_GLOBALSINDEX: c_int = -luaconf::LUAI_MAXCSTACK - 2002;

/// OK status.
pub const LUA_OK: c_int = 0;
/// The thread is suspended.
pub const LUA_YIELD: c_int = 1;
/// Runtime error.
pub const LUA_ERRRUN: c_int = 2;
/// Legacy error code, preserved for compatibility.
pub const LUA_ERRSYNTAX: c_int = 3;
/// Memory allocation error.
pub const LUA_ERRMEM: c_int = 4;
/// Error while running running the error handler function.
pub const LUA_ERRERR: c_int = 5;
/// Yielded for a debug breakpoint.
pub const LUA_BREAK: c_int = 6;

/// Coroutine is running.
pub const LUA_CORUN: c_int = 0;
/// Coroutine is suspended.
pub const LUA_COSUS: c_int = 1;
/// Coroutine is 'normal' (it resumed another coroutine)
pub const LUA_CONOR: c_int = 2;
/// Coroutine has finished.
pub const LUA_COFIN: c_int = 3;
/// Coorutine finished with an error.
pub const LUA_COERR: c_int = 4;

/// Type for an "empty" stack position.
pub const LUA_TNONE: c_int = -1;
/// Type `nil`
pub const LUA_TNIL: c_int = 0;
/// Type `boolean`
pub const LUA_TBOOLEAN: c_int = 1;
/// Type `lightuserdata`
pub const LUA_TLIGHTUSERDATA: c_int = 2;
/// Type `number`
pub const LUA_TNUMBER: c_int = 3;
/// Type `vector`
pub const LUA_TVECTOR: c_int = 4;
/// Type `string`
pub const LUA_TSTRING: c_int = 5;
/// Type `table`
pub const LUA_TTABLE: c_int = 6;
/// Type `function`
pub const LUA_TFUNCTION: c_int = 7;
/// Type `userdata`
pub const LUA_TUSERDATA: c_int = 8;
/// Type `thread`
pub const LUA_TTHREAD: c_int = 9;
/// Type `buffer`
pub const LUA_TBUFFER: c_int = 10;

/// Internal tag for GC objects.
pub const LUA_TPROTO: c_int = 11;
/// Internal tag for GC objects.
pub const LUA_TUPVAL: c_int = 12;
/// Internal tag for GC objects.
pub const LUA_TDEADKEY: c_int = 13;

/// The number of types that exist.
pub const LUA_T_COUNT: c_int = 11;

/// Stop garbage collection.
pub const LUA_GCSTOP: c_int = 0;
/// Resume garbage collection.
pub const LUA_GCRESTART: c_int = 1;
/// Run a full GC cycle. Not recommended for latency sensitive applications.
pub const LUA_GCCOLLECT: c_int = 2;
/// Returns the current amount of memory used in KB.
pub const LUA_GCCOUNT: c_int = 3;
/// Returns the remainder in bytes of the current amount of memory used.
///
/// This is the remainder after dividing the total amount of bytes by 1024.
pub const LUA_GCCOUNTB: c_int = 4;
/// Returns 1 if the GC is active (not stopped).
///
/// The GC may not be actively collecting even if it's running.
pub const LUA_GCISRUNNING: c_int = 5;
/// Performs an explicit GC step, with the step size specified in KB.
pub const LUA_GCSTEP: c_int = 6;
/// Set the goal GC parameter.
pub const LUA_GCSETGOAL: c_int = 7;
/// Set the step multiplier GC parameter
pub const LUA_GCSETSTEPMUL: c_int = 8;
/// Set the step size GC parameter.
pub const LUA_GCSETSTEPSIZE: c_int = 9;

/// Sentinel value indicating the absence of a registry reference.
pub const LUA_NOREF: c_int = -1;
/// Special reference indicating a `nil` value.
pub const LUA_REFNIL: c_int = 0;

/// Type of Luau numbers.
pub type lua_Number = c_double;
/// Type for integer functions.
pub type lua_Integer = c_int;
/// Unsigned integer type.
pub type lua_Unsigned = c_uint;

/// Type for C functions.
///
/// When called, the arguments passed to the function are available in the stack, with the first
/// argument at index 1 and the last argument at the top of the stack ([`lua_gettop`]). Return
/// values should be pushed onto the stack (the first result is pushed first), and the number of
/// return values should be returned.
pub type lua_CFunction = unsafe extern "C-unwind" fn(*mut lua_State) -> c_int;

/// Type for a continuation function.
///
/// See [`lua_pushcclosurek`].
pub type lua_Continuation = unsafe extern "C-unwind" fn(*mut lua_State, c_int) -> c_int;

/// Type for a memory allocation function.
///
/// This function is called with the `ud` passed to [`lua_newstate`], a pointer to the block being
/// allocated/reallocated/freed, the original size of the block, and the requested new size of the
/// block.
///
/// If the given new size is non-zero, a pointer to a block with the requested size should be
/// returned. If the request cannot be fulfilled or the requested size is zero, null should be
/// returned. If the given original size is not zero, the given block should be freed. It is
/// assumed this never fails if the requested size is smaller than the original size.
pub type lua_Alloc =
    unsafe extern "C-unwind" fn(*mut c_void, *mut c_void, usize, usize) -> *mut c_void;

/// Destructor function for a userdata. Called before the userdata is garbage collected.
pub type lua_Destructor = unsafe extern "C-unwind" fn(L: *mut lua_State, userdata: *mut c_void);

/// Functions to be called by the debugger in specific events.
pub type lua_Hook = unsafe extern "C-unwind" fn(L: *mut lua_State, ar: *mut lua_Debug);

/// Callback function for [`lua_getcoverage`]. Receives the coverage information.
pub type lua_Coverage = unsafe extern "C-unwind" fn(
    context: *mut c_void,
    function: *const c_char,
    linedefined: c_int,
    depth: c_int,
    hits: *const c_int,
    size: usize,
);

/// A Luau thread.
///
/// This is an opaque type and always exists behind a pointer (like `*mut lua_State`).
#[repr(C)]
pub struct lua_State {
    _data: (),
    _marker: core::marker::PhantomData<core::marker::PhantomPinned>,
}

/// Activation record. Contains debug information.
#[repr(C)]
#[derive(Debug, Clone, Copy)]
pub struct lua_Debug {
    /// Name of the function.
    pub name: *const c_char,
    /// One of `Lua`, `C`, `main`, or `tail`.
    pub what: *const c_char,
    /// The source. Usually the filename.
    pub source: *const c_char,
    /// Short chunk identifier.
    pub short_src: *const c_char,
    /// The line where the function is defined.
    pub linedefined: c_int,
    /// The current line.
    pub currentline: c_int,
    /// The number of upvalues.
    pub nupvals: c_uint,
    /// The number of parameters.
    pub nparams: c_uint,
    /// Whether or not the function is variadic.
    pub isvararg: c_char,
    /// Userdata.
    pub userdata: *mut c_void,
    /// Buffer for `short_src`.
    pub ssbuf: [c_char; luaconf::LUA_IDSIZE],
}

impl lua_Debug {
    pub fn new() -> lua_Debug {
        lua_Debug {
            name: ptr::null(),
            what: ptr::null(),
            source: ptr::null(),
            short_src: ptr::null(),
            linedefined: 0,
            currentline: 0,
            nupvals: 0,
            nparams: 0,
            isvararg: 0,
            userdata: ptr::null_mut(),
            ssbuf: [0; luaconf::LUA_IDSIZE],
        }
    }
}

impl Default for lua_Debug {
    fn default() -> lua_Debug {
        lua_Debug::new()
    }
}

/// Callbacks that can be used to reconfigure behavior of the VM dynamically.
///
/// This can be retrieved using [`lua_callbacks`].
///
/// **Note:** `interrupt` is safe to set from an arbitrary thread. All other callbacks should only
/// be changed when the VM is not running any code.
#[repr(C)]
#[derive(Debug, Default, Clone, Copy)]
pub struct lua_Callbacks {
    /// Arbitrary userdata pointer. Never overwritten by Luau.
    pub userdata: *mut c_void,

    /// Called at safepoints (e.g. loops, calls/returns, garbage collection).
    pub interrupt: Option<unsafe extern "C-unwind" fn(L: *mut lua_State, gc: c_int)>,
    /// Called when an unprotected error is raised (if longjmp is used).
    pub panic: Option<unsafe extern "C-unwind" fn(L: *mut lua_State, errcode: c_int)>,

    /// Called when `L` is created (`LP` is the parent), or destroyed (`LP` is null).
    pub userthread: Option<unsafe extern "C-unwind" fn(LP: *mut lua_State, L: *mut lua_State)>,
    /// Called when a string is created. Returned atom can be retrieved via [`lua_tostringatom`].
    pub useratom: Option<unsafe extern "C-unwind" fn(s: *const c_char, l: usize) -> i16>,

    /// Called when a `BREAK` instruction is encountered.
    pub debugbreak: Option<unsafe extern "C-unwind" fn(L: *mut lua_State, ar: *mut lua_Debug)>,
    /// Called after each instruction in single step mode.
    pub debugstep: Option<unsafe extern "C-unwind" fn(L: *mut lua_State, ar: *mut lua_Debug)>,
    /// Called when thread execution is interrupted by break in another thread.
    pub debuginterrupt: Option<unsafe extern "C-unwind" fn(L: *mut lua_State, ar: *mut lua_Debug)>,
    /// Called when a protected call results in an error.
    pub debugprotectederror: Option<unsafe extern "C-unwind" fn(L: *mut lua_State)>,

    /// Called when memory is allocated.
    pub onallocate:
        Option<unsafe extern "C-unwind" fn(L: *mut lua_State, osize: usize, nsize: usize)>,
}

unsafe extern "C-unwind" {
    /// Creates a new independent state using the given allocation function and `ud`.
    pub unsafe fn lua_newstate(f: lua_Alloc, ud: *mut c_void) -> *mut lua_State;
    /// Closes the given state.
    ///
    /// Destroys all objects in the state and frees all dynamic memory used by the state.
    pub unsafe fn lua_close(L: *mut lua_State);
    /// Pushes a new thread to the stack.
    ///
    /// Returns a pointer to a [`lua_State`] that represents the created thread. The new state
    /// shared all global objects (such as tables) with the original state, but has an
    /// independent execution stack.
    pub unsafe fn lua_newthread(L: *mut lua_State) -> *mut lua_State;
    /// Returns a pointer to a [`lua_State`], which represents the main thread of the given thread.
    pub unsafe fn lua_mainthread(L: *mut lua_State) -> *mut lua_State;
    /// Resets the given thread to its original state.
    ///
    /// Clears all the thread state, call frames, and the stack.
    pub unsafe fn lua_resetthread(L: *mut lua_State);
    /// Returns true if the given thread is in its original state.
    pub unsafe fn lua_isthreadreset(L: *mut lua_State) -> c_int;

    /// Converts the given index into an absolute index.
    pub unsafe fn lua_absindex(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the index of the top element in the stack.
    ///
    /// This is equal to the number of elements in the stack.
    pub unsafe fn lua_gettop(L: *mut lua_State) -> c_int;
    /// Sets the top of the stack to the given index.
    ///
    /// If the new top is greater than the old one, new elements are filled with `nil`. If the new
    /// top is smaller, elements will be removed from the old top.
    pub unsafe fn lua_settop(L: *mut lua_State, idx: c_int);
    /// Pushes a copy of the element at the given index onto the stack.
    pub unsafe fn lua_pushvalue(L: *mut lua_State, idx: c_int);
    /// Removes the element at the given index from the stack.
    ///
    /// Elements above the given index are shifted down to fill the gap. This function cannot be
    /// called with a pseudo-index.
    pub unsafe fn lua_remove(L: *mut lua_State, idx: c_int);
    /// Moves the top element into the given index.
    ///
    /// Elements above the given index are shifted up to make space. This function cannot be called
    /// with a pseudo-index.
    pub unsafe fn lua_insert(L: *mut lua_State, idx: c_int);
    /// Replaces the element at the given index with the top element.
    ///
    /// The top element will be popped.
    pub unsafe fn lua_replace(L: *mut lua_State, idx: c_int);
    /// Ensures the stack has space for the given number of additional elements.
    ///
    /// Returns `0` if the request cannot be fulfilled, either because it would be bigger than the
    /// maximum stack size, or because it cannot allocate memory for the extra space.
    pub unsafe fn lua_checkstack(L: *mut lua_State, sz: c_int) -> c_int;
    /// Like [`lua_checkstack`], but allows for an infinite amount of stack frames.
    pub unsafe fn lua_rawcheckstack(L: *mut lua_State, sz: c_int);

    /// Move values between different threads of the same state.
    ///
    /// Pops `n` values from the stack of `from`, and pushes them onto the stack of `to`.
    pub unsafe fn lua_xmove(from: *mut lua_State, to: *mut lua_State, n: c_int);
    /// Copy a value to a different thread of the same state.
    ///
    /// The element at the given index is copied from `from`, and pushed onto the stack of `to`.
    pub unsafe fn lua_xpush(from: *mut lua_State, to: *mut lua_State, idx: c_int);

    /// Checks if a value can be converted to a number.
    ///
    /// Returns `1` if the value at the given index is a number, or a string that can be converted
    /// to a number, and `0` otherwise.
    pub unsafe fn lua_isnumber(L: *mut lua_State, idx: c_int) -> c_int;
    /// Checks if a value can be converted to a string.
    ///
    /// Returns `1` if the value at the given index is a string or a number, and `0` otherwise.
    pub unsafe fn lua_isstring(L: *mut lua_State, idx: c_int) -> c_int;
    /// Checks if a value is a C function.
    ///
    /// Returns `1` if the value at the given index is a C function, and `0` otherwise.
    pub unsafe fn lua_iscfunction(L: *mut lua_State, idx: c_int) -> c_int;
    /// Checks if a value is a Luau function.
    ///
    /// Returns `1` if the value at the given index is a Luau function, and `0` otherwise.
    pub unsafe fn lua_isLfunction(L: *mut lua_State, idx: c_int) -> c_int;
    /// Checks if a value is a userdata.
    ///
    /// Returns `1` if the value at the given index is a userdata, and `0` otherwise.
    pub unsafe fn lua_isuserdata(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the type of the value at the given index.
    pub unsafe fn lua_type(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the name of the given type.
    pub unsafe fn lua_typename(L: *mut lua_State, tp: c_int) -> *const c_char;

    /// Checks if two values are equal.
    ///
    /// Returns `1` if the two values in the given indices are equal, and `0` otherwise. This
    /// follows the behavior of the `==` operator, and may call metamethods.
    pub unsafe fn lua_equal(L: *mut lua_State, idx1: c_int, idx2: c_int) -> c_int;
    /// Like [`lua_equal`], but will not invoke metamethods.
    pub unsafe fn lua_rawequal(L: *mut lua_State, idx1: c_int, idx2: c_int) -> c_int;
    /// Checks if one value is less than another value.
    ///
    /// Returns `1` if the value at `idx1` is smaller than the value at `idx2`, and `0` otherwise.
    /// This follows the behavior of the `<` operator, and may call metamethods.
    pub unsafe fn lua_lessthan(L: *mut lua_State, idx1: c_int, idx2: c_int) -> c_int;

    /// Checks if the given value can be converted to a number, and returns the number.
    ///
    /// If the value at the given index is a number, or a string that can be converted to one, the
    /// number will be returned and `isnum` will be `1`. Otherwise, `0` will be returned and `isnum`
    /// will be `0`.
    pub unsafe fn lua_tonumberx(L: *mut lua_State, idx: c_int, isnum: *mut c_int) -> lua_Number;
    /// Checks if the given value can be converted to an integer, and returns the integer.
    ///
    /// If the value at the given index is a number, or a string that can be converted to one, the
    /// number will be truncated and returned and `isnum` will be `1`. Otherwise, `0` will be
    /// returned and `isnum` will be `0`.
    pub unsafe fn lua_tointegerx(L: *mut lua_State, idx: c_int, isnum: *mut c_int) -> lua_Integer;
    /// Checks if the given value can be converted to an unsigned integer, and returns it.
    ///
    /// If the value at the given index is a number, or a string that can be converted to one, the
    /// number will be truncated and returned and `isnum` will be `1`. Otherwise, `0` will be
    /// returned and `isnum` will be `0`.
    ///
    /// If the integer is negative, it will be forced into an unsigned integer in an unspecified
    /// way.
    pub unsafe fn lua_tounsignedx(L: *mut lua_State, idx: c_int, isnum: *mut c_int)
    -> lua_Unsigned;
    /// Get the vector in the given value.
    ///
    /// If the value at the given index is a vetor, returns a pointer to the first component of
    /// that vector (`x`), otherwise returns null.
    pub unsafe fn lua_tovector(L: *mut lua_State, idx: c_int) -> *const c_float;
    /// Get the boolean in the given value.
    ///
    /// Returns `1` if the given value is truthy, and `0` if it is falsy. A value is falsy if it is
    /// `false` or `nil`.
    pub unsafe fn lua_toboolean(L: *mut lua_State, idx: c_int) -> c_int;
    /// Converts the given value to a string.
    ///
    /// The value at the given index is converted into a string and `len` will be the length of the
    /// string. The value must be a string or number, otherwise null is returned. If it is a
    /// number, the actual value in the stack will be changed to a string.
    ///
    /// An aligned pointer to the string within the Luau state is returned. This string is null
    /// terminated, but may also contain nulls within its body. If the string has been removed from
    /// the stack, it may get garbage collected causing the pointer to dangle.
    pub unsafe fn lua_tolstring(L: *mut lua_State, idx: c_int, len: *mut usize) -> *const c_char;
    /// Like [`lua_tolstringatom`] but doesn't get the length.
    pub unsafe fn lua_tostringatom(
        L: *mut lua_State,
        idx: c_int,
        atom: *mut c_int,
    ) -> *const c_char;
    /// Get a string and its atom.
    ///
    /// Returns an aligned pointer to the string within the Luau state. If the value is not a
    /// string, null is returned, no conversion is performed for non-strings. The string is null
    /// terminated, but may also contain nulls within its body. If the string has been removed from
    /// the stack, it may get garbage collected causing the pointer to dangle.
    ///
    /// The `atom` will be the result of [`lua_Callbacks::useratom`] or -1 if not set. The same
    /// string will have the same atom.
    pub unsafe fn lua_tolstringatom(
        L: *mut lua_State,
        idx: c_int,
        len: *mut usize,
        atom: *mut c_int,
    ) -> *const c_char;
    /// Gets the current namecall string and its atom.
    ///
    /// When called during a namecall, a pointer to the string containing the name of the method
    /// that was called is returned, otherwise null is returned. See [`lua_tolstringatom`] for more
    /// information.
    pub unsafe fn lua_namecallatom(L: *mut lua_State, atom: *mut c_int) -> *const c_char;
    /// Get the length of the given value.
    ///
    /// Returns the "length" of the value at the given index. For strings, it is the length of the
    /// string. For tables, it is the result of the `#` operator. For userdata, this is the size of
    /// the block of memory allocated for the userdata. For other values, it is `0`.
    pub unsafe fn lua_objlen(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the C function in the given value.
    ///
    /// Returns the C function in the value at the given index, otherwise returns null.
    pub unsafe fn lua_tocfunction(L: *mut lua_State, idx: c_int) -> lua_CFunction;
    /// Gets the light userdata in the given value.
    ///
    /// Returns the light userdata in the value at the given index, otherwise returns null.
    pub unsafe fn lua_tolightuserdata(L: *mut lua_State, idx: c_int) -> *mut c_void;
    /// Gets the light userdata with the given tag in the given value.
    ///
    /// Returns the light userdata in the value at the given index if its tag matches the given
    /// tag, otherwise returns null.
    pub unsafe fn lua_tolightuserdatatagged(
        L: *mut lua_State,
        idx: c_int,
        tag: c_int,
    ) -> *mut c_void;
    /// Gets the userdata in the given value.
    ///
    /// Returns the userdata in the value at the given index, otherwise returns null.
    pub unsafe fn lua_touserdata(L: *mut lua_State, idx: c_int) -> *mut c_void;
    /// Gets the userdata with the given tag in the given value.
    ///
    /// Returns the userdata in the value at the given index if its tag matches the given tag,
    /// otherwise returns null.
    pub unsafe fn lua_touserdatatagged(L: *mut lua_State, idx: c_int, tag: c_int) -> *mut c_void;
    /// Gets the tag of the userdata in the given value.
    ///
    /// Returns the tag of the userdata in the value at the given index, or `-1` if the value is
    /// not a userdata.
    pub unsafe fn lua_userdatatag(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the tag of the light userdata in the given value.
    ///
    /// Returns the tag of the light userdata in the value at the given index, or `-1` if the value
    /// is not a userdata.
    pub unsafe fn lua_lightuserdatatag(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the thread in the given value.
    ///
    /// Returns a pointer to a [`lua_State`] that represents the thread in the value at the given
    /// index, or null if the value is not a thread.
    pub unsafe fn lua_tothread(L: *mut lua_State, idx: c_int) -> *mut lua_State;
    /// Gets the buffer in the given value.
    ///
    /// Returns a pointer to the first byte of the buffer in the value at the given index, or null
    /// if the value is not a buffer. The `len` will be the length of the buffer.
    pub unsafe fn lua_tobuffer(L: *mut lua_State, idx: c_int, len: *mut usize) -> *mut c_void;
    /// Converts the given value into a generic pointer.
    ///
    /// Returns a pointer that represents the given value. The value can be a userdata, a table, a
    /// thread or a function; otherwise, null is returned. Different objects give different
    /// pointers. There is no way to convert a pointer back into its original value.
    ///
    /// This function is typically used for debugging purposes.
    pub unsafe fn lua_topointer(L: *mut lua_State, idx: c_int) -> *const c_void;

    /// Pushes `nil` value onto the stack.
    pub unsafe fn lua_pushnil(L: *mut lua_State);
    /// Pushes the given number onto the stack.
    pub unsafe fn lua_pushnumber(L: *mut lua_State, n: c_double);
    /// Pushes the given integer onto the stack.
    pub unsafe fn lua_pushinteger(L: *mut lua_State, n: c_int);
    /// Pushes the given unsigned integer onto the stack.
    pub unsafe fn lua_pushunsigned(L: *mut lua_State, n: c_uint);

    /// Pushes the given vector onto the stack.
    #[cfg(not(feature = "vector4"))]
    pub unsafe fn lua_pushvector(L: *mut lua_State, x: c_float, y: c_float, z: c_float);
    /// Pushes the given vector onto the stack.
    #[cfg(feature = "vector4")]
    pub unsafe fn lua_pushvector(L: *mut lua_State, x: c_float, y: c_float, z: c_float, w: c_float);

    /// Pushes the pointed-to string with the given length onto the stack.
    pub unsafe fn lua_pushlstring(L: *mut lua_State, s: *const c_char, l: usize);
    /// Pushes the given null-terminated string onto the stack.
    pub unsafe fn lua_pushstring(L: *mut lua_State, s: *const c_char);
    /// Pushes a formatted string onto the stack.
    ///
    /// The `fmt` string can contain the following specifiers:
    ///   - `%%`: inserts a literal `%`.
    ///   - `%s`: inserts a zero terminated string.
    ///   - `%f`: inserts a [`lua_Number`].
    ///   - `%p`: inserts a pointer as a hexadecimal numeral.
    ///   - `%d`: inserts an int.
    ///   - `%c`: inserts an int as a character.
    pub unsafe fn lua_pushfstringL(L: *mut lua_State, fmt: *const c_char, ...) -> *const c_char;
    /// Pushes the given C closure onto the stack, with a continuation function.
    ///
    /// Takes the C function, the function name for debugging, the number of upvalues and
    /// an optional continuation function.
    ///
    /// The upvalues will be popped from the stack. The upvalues can then be retrieved within the
    /// function by subtracting the index of the upvalue (starting from 1) from [`LUA_GLOBALSINDEX`].
    /// For example, to get the first upvalue you would do `LUA_GLOBALSINDEX - 1`.
    ///
    /// See `examples/continuations.rs` for more information about the continutation function.
    pub unsafe fn lua_pushcclosurek(
        L: *mut lua_State,
        fn_: lua_CFunction,
        debugname: *const c_char,
        nup: c_int,
        cont: Option<lua_Continuation>,
    );
    /// Pushes the given boolean onto the stack.
    pub unsafe fn lua_pushboolean(L: *mut lua_State, b: c_int);
    /// Pushes the current thread onto the stack.
    ///
    /// Returns `1` if this thread is the main thread of its state.
    pub unsafe fn lua_pushthread(L: *mut lua_State) -> c_int;

    /// Pushes a tagged light userdata onto the stack.
    pub unsafe fn lua_pushlightuserdatatagged(L: *mut lua_State, p: *mut c_void, tag: c_int);
    /// Allocates a new tagged userdata and pushes it onto the stack.
    ///
    /// Returns a pointer to the allocated block of memory with the given size.
    pub unsafe fn lua_newuserdatatagged(L: *mut lua_State, sz: usize, tag: c_int) -> *mut c_void;
    /// Allocates a new tagged userdata with the metatable and pushes it onto the stack.
    ///
    /// The userdata will have the given tag. The metatable will be set to the metatable that has
    /// been set using [`lua_setuserdatametatable`].
    pub unsafe fn lua_newuserdatataggedwithmetatable(
        L: *mut lua_State,
        sz: usize,
        tag: c_int,
    ) -> *mut c_void;
    /// Allocates a new userdata with the given destructor and pushes it onto the stack.
    ///
    /// The specified destructor function will be called before the userdata is garbage collected.
    pub unsafe fn lua_newuserdatadtor(
        L: *mut lua_State,
        sz: usize,
        dtor: unsafe extern "C-unwind" fn(*mut c_void),
    ) -> *mut c_void;

    /// Allocates a new Luau buffer and pushes it onto the stack.
    ///
    /// Returns a pointer to the allocated block of memory with the given size.
    pub unsafe fn lua_newbuffer(L: *mut lua_State, sz: usize) -> *mut c_void;

    /// Gets an element from the given table and pushes it onto the stack.
    ///
    /// Pushes the value `t[k]` onto the stack, where `t` is the table at the given index and `k`
    /// is the value at the top of the stack. The key is popped from the stack. This may call the
    /// `__index` metamethod.
    ///
    /// Returns the type of the pushed value.
    pub unsafe fn lua_gettable(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets a field from the given table and pushes it onto the stack.
    ///
    /// Pushes the value `t[k]` onto the stack, where `t` is the table at the given index and `k`
    /// is the given string. This may call the `__index` metamethod.
    ///
    /// Returns the type of the pushed value.
    pub unsafe fn lua_getfield(L: *mut lua_State, idx: c_int, k: *const c_char) -> c_int;
    /// Gets a field from the given table without invoking metamethods.
    ///
    /// Like [`lua_getfield`], but does not invoke metamethods.
    pub unsafe fn lua_rawgetfield(L: *mut lua_State, idx: c_int, k: *const c_char) -> c_int;
    /// Gets an element from the given table without invoking metamethods.
    ///
    /// Like [`lua_gettable`], but does not invoke metamethods.
    pub unsafe fn lua_rawget(L: *mut lua_State, idx: c_int) -> c_int;
    /// Gets the element at the given index in the given table without invoking metamethods.
    ///
    /// Pushes the value `t[n]` onto the stack, where `t` is the table at the given index and `n`
    /// is the given index.
    ///
    /// Returns the type of the pushed value.
    pub unsafe fn lua_rawgeti(L: *mut lua_State, idx: c_int, n: c_int) -> c_int;
    /// Pushes a new empty table with the given size hints onto the stack.
    ///
    /// `narr` is a hint for how many elements the array portion the table will contain, `nrec` is
    /// a hint for how many other elements it will contain. The hints may be used to preallocate
    /// the memory for the table.
    ///
    /// If you do not know the number of elements the table will contain in advance, you may use
    /// [`lua_newtable`] instead.
    pub unsafe fn lua_createtable(L: *mut lua_State, narr: c_int, nrec: c_int);

    /// Sets the `readonly`` flag of the given table.
    pub unsafe fn lua_setreadonly(L: *mut lua_State, idx: c_int, enabled: c_int);
    /// Gets the `readonly` flag of the given table.
    pub unsafe fn lua_getreadonly(L: *mut lua_State, idx: c_int) -> c_int;
    /// Sets the `safeenv` flag of the given table.
    pub unsafe fn lua_setsafeenv(L: *mut lua_State, idx: c_int, enabled: c_int);

    /// Pushes the metatable of the value at the given index onto the stack.
    ///
    /// If the index is not valid or if it doesn't have a metatable, zero is returned.
    pub unsafe fn lua_getmetatable(L: *mut lua_State, objindex: c_int) -> c_int;
    /// Pushes the environment table of the value at the given index onto the stack.
    pub unsafe fn lua_getfenv(L: *mut lua_State, idx: c_int);

    /// Sets an element in the given table.
    ///
    /// Does the equivalent of `t[k] = v`, where `t` is the value at the given index, `v` is the
    /// value at the top of the stack, and `k` is the value just below the top. The key and value
    /// will be popped from the stack. This may call the `__newindex` metamethod.
    pub unsafe fn lua_settable(L: *mut lua_State, idx: c_int);
    /// Sets a field in the given table.
    ///
    /// Does the equivalent of `t[k] = v`, where `t` is the value at the given index, `v` is the
    /// value at the top of the stack, and `k` is the given string. The value will be popped from
    /// the stack. This may call the `__newindex` metamethod.
    pub unsafe fn lua_setfield(L: *mut lua_State, idx: c_int, k: *const c_char);
    /// Sets a field in the given table without invoking metamethods.
    ///
    /// Like [`lua_setfield`], but does not invoke metamethods.
    pub unsafe fn lua_rawsetfield(L: *mut lua_State, idx: c_int, k: *const c_char);
    /// Sets an element in the given table without invoking metamethods.
    ///
    /// Like [`lua_settable`], but does not invoke metamethods.
    pub unsafe fn lua_rawset(L: *mut lua_State, idx: c_int);
    /// Sets the element at the given index in the given table without invoking metamethods.
    ///
    /// Does the equivalent of `t[n] = v` onto the stack, where `t` is the table at the given index
    /// and `v` is the value at the top of the stack.
    pub unsafe fn lua_rawseti(L: *mut lua_State, idx: c_int, n: c_int);
    /// Sets the metatable of the given value.
    ///
    /// Pops a table from the stack and sets it as the metatable for the value at the given index.
    ///
    /// Always returns 1 for historical reasons.
    pub unsafe fn lua_setmetatable(L: *mut lua_State, objindex: c_int) -> c_int;
    /// Sets the environment for the given value.
    ///
    /// Pops a table from the stack and sets it as the new environment for the value at the given
    /// index. Returns 1 if the value is a function, thread, or userdata. Otherwise, returns 0.
    pub unsafe fn lua_setfenv(L: *mut lua_State, idx: c_int) -> c_int;

    /// Loads the given Luau chunk and pushes it as a function onto the stack.
    ///
    /// Takes the name of the chunk for debugging, a pointer to the bytecode, the size of the
    /// bytecode, and a value representing the environment of the chunk. If `env` is 0, the current
    /// environment will be used, otherwise the table on the stack at index `env` will be used.
    ///
    /// Returns 1 if there was an error, 0 otherwise. The error message will be pushed onto the
    /// stack if there was an error.
    pub unsafe fn luau_load(
        L: *mut lua_State,
        chunkname: *const c_char,
        data: *const c_char,
        size: usize,
        env: c_int,
    ) -> c_int;
    /// Calls the given function.
    ///
    /// This will pop `nargs` values from the stack, which will be passed into the function as
    /// arguments. The first argument is the first value that got pushed. The results of the
    /// function are pushed onto the stack. The number of results will be adjusted to `nresults`,
    /// unless `nresults` is [`LUA_MULTRET`], in which case all results are pushed. The first
    /// result is pushed first, the last result will be at the top of the stack.
    pub unsafe fn lua_call(L: *mut lua_State, nargs: c_int, nresults: c_int);
    /// Calls the given function in protected mode.
    ///
    /// Like [`lua_call`]. However, if there is an error, it will get catched, the error message
    /// will be pushed onto the stack, and the error code will be returned.
    ///
    /// If `errfunc` is zero, the error message pushed on the stack will not be modified. Otherwise,
    /// `errfunc` is the stack index of an *error handler function*. In case of an error, this
    /// function will be called with the error message and its return value will be the error
    /// messaged pushed onto the stack.
    ///
    /// Returns zero if successful, otherwise returns one of [`LUA_ERRRUN`], [`LUA_ERRMEM`], or
    /// [`LUA_ERRERR`].
    pub unsafe fn lua_pcall(
        L: *mut lua_State,
        nargs: c_int,
        nresults: c_int,
        errfunc: c_int,
    ) -> c_int;

    /// Yields the current coroutine.
    ///
    /// This function should be called as the return of a C function.
    pub unsafe fn lua_yield(L: *mut lua_State, nresults: c_int) -> c_int;
    /// Breaks execution, as if a debug breakpoint has been reached.
    ///
    /// This function should be called as the return of a C function.
    pub unsafe fn lua_break(L: *mut lua_State) -> c_int;
    /// Starts and resumes a coroutine in the given thread `L`.
    ///
    /// To start a coroutine, push the main function plus any arguments onto the stack, then call
    /// this function with `nargs` being the number of arguments. Once the coroutine suspends or
    /// finishes, the results of the resumption will be pushed onto the stack. The first result
    /// will be pushed first, the last result will be at the top of the stack. If the coroutine
    /// yielded, [`LUA_YIELD`] is returned. If there was an error an error code is returned, and
    /// the error value will be pushed onto the stack. Otherwise, [`LUA_OK`] is returned.
    ///
    /// To resume a coroutine, push values to be returned from the `yield` call onto the stack, and
    /// call this function. Make sure to remove the results pushed from your previous resume call
    /// before resuming again.
    ///
    /// The `from` parameter represents the coroutine that is resuming `L`. This may be null if
    /// there is no such coroutine.
    pub unsafe fn lua_resume(L: *mut lua_State, from: *mut lua_State, narg: c_int) -> c_int;
    /// Error a coroutine in the given thread `L`.
    ///
    /// Like [`lua_resume`], but this will pop a single value from the stack which will be used as
    /// the error value. If the coroutine is currently inside of a `pcall`, the error will be
    /// catched by that `pcall`. Otherwise, this call will return an error code and the error value
    /// will be pushed onto the stack.
    pub unsafe fn lua_resumeerror(L: *mut lua_State, from: *mut lua_State) -> c_int;
    /// Gets the status of the given thread `L`.
    ///
    /// The status can be [`LUA_OK`] for a normal thread, an error code if the thread finished
    /// execution of a `lua_resume` with an error, or [`LUA_YIELD`] if the thread is suspended.
    ///
    /// You can only call functions in threads that are [`LUA_OK`]. You can resume threads that are
    /// [`LUA_OK`] or [`LUA_YIELD`].
    pub unsafe fn lua_status(L: *mut lua_State) -> c_int;
    /// Returns 1 if the given coroutine can yield, and 0 otherwise.
    pub unsafe fn lua_isyieldable(L: *mut lua_State) -> c_int;
    /// Gets the thread data of the given thread.
    ///
    /// The thread data can be set with [`lua_setthreaddata`].
    pub unsafe fn lua_getthreaddata(L: *mut lua_State) -> *mut c_void;
    /// Sets the thread data of the given thread.
    ///
    /// The thread data can be retrieved afterwards with [`lua_getthreaddata`].
    pub unsafe fn lua_setthreaddata(L: *mut lua_State, data: *mut c_void);
    /// Gets the status of the given coroutine `co`.
    pub unsafe fn lua_costatus(L: *mut lua_State, co: *mut lua_State) -> c_int;

    /// Controls the garbage collector.
    ///
    /// Takes a GC operation `what` and data to use for the operation. The different operations are
    /// documented inside [lua.h].
    ///
    /// [lua.h]: https://github.com/luau-lang/luau/blob/master/VM/include/lua.h#L249
    pub unsafe fn lua_gc(L: *mut lua_State, what: c_int, data: c_int) -> c_int;

    /// Set the memory category used for memory statistics.
    pub unsafe fn lua_setmemcat(L: *mut lua_State, category: c_int);
    /// Get the total amount of memory used, in bytes.
    ///
    /// If `category < 0`, then the total amount of memory is returned.
    pub unsafe fn lua_totalbytes(L: *mut lua_State, category: c_int) -> usize;

    /// Throws a Luau error.
    ///
    /// The error value must be at the top of the stack.
    pub unsafe fn lua_error(L: *mut lua_State) -> !;

    /// Gets the next key-value pair in the given table.
    ///
    /// This will pop a key from the stack, and pushes the next key-value pair onto the stack of
    /// the table at the given index. If the end of the table has been reached, nothing is pushed
    /// and zero is returned. If a non-zero value is returned, then the key is at stack index `-2`,
    /// and the value is at `-1`.
    ///
    /// While traversing a table, avoid calling [`lua_tolstring`] directly on a key, unless you
    /// know the key is actually a string. [`lua_tolstring`] will change the value to a string if
    /// it's a number, which will confuse [`lua_next`].
    pub unsafe fn lua_next(L: *mut lua_State, idx: c_int) -> c_int;
    /// Iterate over the given table.
    ///
    /// This function should be repeatedly called, where `iter` is the previous return value, or `0`
    /// on the first iteration. If `-1` is returned, the end of the table has been reached.
    /// Otherwise, the key-value pair is pushed onto the stack. The key will be at index `-2` and
    /// the value at `-1`.
    ///
    /// This function is similar to [`lua_next`], however the order of items from this function
    /// matches the order from a normal `for i, v in t do` loop (array portion first, then the
    /// remaining keys).
    pub unsafe fn lua_rawiter(L: *mut lua_State, idx: c_int, iter: c_int) -> c_int;

    /// Pops `n` values from the stack and concatenates them.
    ///
    /// The concatenated string is pushed onto the stack. If `n` is zero, the result is an empty
    /// string. Follows the same semantics as the `..` operator in Luau.
    pub unsafe fn lua_concat(L: *mut lua_State, n: c_int);

    /// Obfuscates the given pointer.
    ///
    /// The same `p` will return the same value on the same state.
    pub unsafe fn lua_encodepointer(L: *mut lua_State, p: usize) -> usize;

    /// Returns a high-precision timestamp (in seconds).
    ///
    /// This is the same as `os.clock()` from the Luau standard library. This can be used to
    /// measure a duration with sub-microsecond precision.
    pub unsafe fn lua_clock() -> c_double;

    /// Sets the tag of the userdata at the given index.
    pub unsafe fn lua_setuserdatatag(L: *mut lua_State, idx: c_int, tag: c_int);

    /// Set the destructor function for userdata objects with the given tag.
    pub unsafe fn lua_setuserdatadtor(L: *mut lua_State, tag: c_int, dtor: Option<lua_Destructor>);
    /// Gets the destructor function for userdata objects with the given tag.
    pub unsafe fn lua_getuserdatadtor(L: *mut lua_State, tag: c_int) -> lua_Destructor;

    /// Sets the metatable for userdata objects with the given tag.
    ///
    /// This will pop a table from the top of the stack, and use it as the metatable for any
    /// userdata objects created with the given tag. This cannot be called on a tag that already
    /// has a metatable set.
    pub unsafe fn lua_setuserdatametatable(L: *mut lua_State, tag: c_int);
    /// Gets the metatable for userdata objects with the given tag.
    pub unsafe fn lua_getuserdatametatable(L: *mut lua_State, tag: c_int);

    /// Sets the name for light userdata objects with the given tag.
    ///
    /// This cannot be called on a tag that already has a name set.
    pub unsafe fn lua_setlightuserdataname(L: *mut lua_State, tag: c_int, name: *const c_char);
    /// Gets the name for light userdata objects wiwth the given tag.
    pub unsafe fn lua_getlightuserdataname(L: *mut lua_State, tag: c_int) -> *const c_char;

    /// Pushes a clone of the function at the given index onto the stack.
    pub unsafe fn lua_clonefunction(L: *mut lua_State, idx: c_int);

    /// Clears the table at the given index.
    pub unsafe fn lua_cleartable(L: *mut lua_State, idx: c_int);
    /// Pushes a clone of the table at the given index onto the stack.
    pub unsafe fn lua_clonetable(L: *mut lua_State, idx: c_int);

    /// Gets the allocator function of the given state.
    ///
    /// Will update `ud` to be the value passed to [`lua_newstate`].
    pub unsafe fn lua_getallocf(L: *mut lua_State, ud: *mut *mut c_void) -> lua_Alloc;

    /// Stores the given value in the registry and gets a reference to it.
    ///
    /// This will store a copy of the value at the given index into the registry. The reference,
    /// which is an integer that identifies this value, is returned. A reference can then be
    /// removed from the registry using [`lua_unref`], allowing it to be garbage collected.
    /// References can be reused once freed.
    ///
    /// If the given value is `nil`, the special reference [`LUA_REFNIL`] is returned. The sentinel
    /// reference [`LUA_NOREF`] will never be returned by this function.
    pub unsafe fn lua_ref(L: *mut lua_State, idx: c_int) -> c_int;
    /// Frees a reference created by [`lua_ref`] from the registry.
    pub unsafe fn lua_unref(L: *mut lua_State, ref_: c_int);

    /// Gets the depth of the call stack.
    pub unsafe fn lua_stackdepth(L: *mut lua_State) -> c_int;
    /// Gets information about a specific function or function invocation.
    ///
    /// If `level < 0`, it is assumed to be a stack index, and information about the function at
    /// that index is retrieved. Otherwise, it is assumed to be a level in the call stack, and
    /// information about the invocation at that level is retrieved.
    ///
    /// The `ar` argument will be filled with the values requested in `what`. Each character in
    /// `what` selects some fields in `ar` to be filled.
    ///
    ///   - `s`: fills `source`, `what`, `linedefined`, and `short_src`
    ///   - `l`: fills `currentline`
    ///   - `u`: fills `nupvals`
    ///   - `a`: fills `isvararg` and `nparams`
    ///   - `n`: fills `name`
    ///   - `f`: pushes the function onto the stack
    ///
    /// Returns `1` if `ar` was updated, `0` otherwise.
    pub unsafe fn lua_getinfo(
        L: *mut lua_State,
        level: c_int,
        what: *const c_char,
        ar: *mut lua_Debug,
    ) -> c_int;
    /// Pushes a copy of the `n`th argument at the given level onto the stack.
    ///
    /// Returns `1` if a value was pushed, `0` otherwise. Always returns `0` for invocations to
    /// native functions.
    pub unsafe fn lua_getargument(L: *mut lua_State, level: c_int, n: c_int) -> c_int;
    /// Pushes a copy of the `n`th local variable at the given level onto the stack.
    ///
    /// Returns a pointer to the name of the variable. If no value was pushed, null is returned.
    /// Always returns null for invocations to native functions.
    pub unsafe fn lua_getlocal(L: *mut lua_State, level: c_int, n: c_int) -> *const c_char;
    /// Pops a value from the stack and sets it as the `n`th local variable at the given level.
    ///
    /// Returns a pointer to the name of the variable. If no variable was set, null is returned.
    /// Always returns null for invocations to native functions.
    pub unsafe fn lua_setlocal(L: *mut lua_State, level: c_int, n: c_int) -> *const c_char;
    /// Pushes a copy of the `n`th upvalue of the function at the given index onto the stack.
    pub unsafe fn lua_getupvalue(L: *mut lua_State, funcindex: c_int, n: c_int) -> *const c_char;
    /// Pops a value from the stack and sets it as the `n`th upvalue of the given function.
    pub unsafe fn lua_setupvalue(L: *mut lua_State, funcindex: c_int, n: c_int) -> *const c_char;

    /// Enables single stepping for the given thread.
    pub unsafe fn lua_singlestep(L: *mut lua_State, enabled: c_int);
    /// Sets a breakpoint for the given function at the given line.
    pub unsafe fn lua_breakpoint(
        L: *mut lua_State,
        funcindex: c_int,
        line: c_int,
        enabled: c_int,
    ) -> c_int;

    /// Collects coverage information for the given function.
    pub unsafe fn lua_getcoverage(
        L: *mut lua_State,
        funcindex: c_int,
        context: *mut c_void,
        callback: lua_Coverage,
    );

    /// Returns a string representation of the stack trace for debugging.
    ///
    /// This is **NOT thread-safe**, the result is stored in a shared global buffer.
    pub unsafe fn lua_debugtrace(L: *mut lua_State) -> *const c_char;

    /// Gets the callbacks used by the given state.
    ///
    /// These are shared between all coroutines.
    pub unsafe fn lua_callbacks(L: *mut lua_State) -> *mut lua_Callbacks;
}
