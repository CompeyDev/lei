/* Code generated by cmd/cgo; DO NOT EDIT. */

#include <stdlib.h>
#include "_cgo_export.h"

#pragma GCC diagnostic ignored "-Wunknown-pragmas"
#pragma GCC diagnostic ignored "-Wpragmas"
#pragma GCC diagnostic ignored "-Waddress-of-packed-member"
#pragma GCC diagnostic ignored "-Wunknown-warning-option"
#pragma GCC diagnostic ignored "-Wunaligned-access"
extern void crosscall2(void (*fn)(void *), void *, int, size_t);
extern size_t _cgo_wait_runtime_init_done(void);
extern void _cgo_release_context(size_t);

extern char* _cgo_topofstack(void);
#define CGO_NO_SANITIZE_THREAD
#define _cgo_tsan_acquire()
#define _cgo_tsan_release()


#define _cgo_msan_write(addr, sz)

extern void _cgoexp_4543809e40d5_go_allocf(void *);

CGO_NO_SANITIZE_THREAD
GoUintptr go_allocf(GoUintptr fp, GoUintptr ptr, GoUint64 osize, GoUint64 nsize)
{
	size_t _cgo_ctxt = _cgo_wait_runtime_init_done();
	typedef struct {
		GoUintptr p0;
		GoUintptr p1;
		GoUint64 p2;
		GoUint64 p3;
		GoUintptr r0;
	} __attribute__((__packed__, __gcc_struct__)) _cgo_argtype;
	static _cgo_argtype _cgo_zero;
	_cgo_argtype _cgo_a = _cgo_zero;
	_cgo_a.p0 = fp;
	_cgo_a.p1 = ptr;
	_cgo_a.p2 = osize;
	_cgo_a.p3 = nsize;
	_cgo_tsan_release();
	crosscall2(_cgoexp_4543809e40d5_go_allocf, &_cgo_a, 40, _cgo_ctxt);
	_cgo_tsan_acquire();
	_cgo_release_context(_cgo_ctxt);
	return _cgo_a.r0;
}

CGO_NO_SANITIZE_THREAD
void _cgo_4543809e40d5_Cfunc__Cmalloc(void *v) {
	struct {
		unsigned long long p0;
		void *r1;
	} __attribute__((__packed__, __gcc_struct__)) *a = v;
	void *ret;
	_cgo_tsan_acquire();
	ret = malloc(a->p0);
	if (ret == 0 && a->p0 == 0) {
		ret = malloc(1);
	}
	a->r1 = ret;
	_cgo_tsan_release();
}
