//go:build llvm22 || (!llvm14 && !llvm15 && !llvm16 && !llvm17 && !llvm18 && !llvm19 && !llvm20 && !llvm21)

package llvm

/*
#include "PtrToAddrBindings.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

const PtrToAddr Opcode = C.LLVMPtrToAddr

// CreatePtrToAddr extracts the address of val without capturing its provenance.
// The result uses the pointer address space's index width from the containing
// module's DataLayout (and preserves the shape of pointer vectors). The builder
// must have an insertion point in a module with the intended DataLayout.
// Available with LLVM 22 and later.
func (b Builder) CreatePtrToAddr(val Value, name string) (v Value) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	v.C = C.LLVMGoBuildPtrToAddr(b.C, val.C, cname)
	return
}

// ConstPtrToAddr extracts a constant pointer's address without capturing its
// provenance. t must have the pointer address space's DataLayout index width
// and the same scalar or vector shape as val. Available with LLVM 22 and later.
func ConstPtrToAddr(val Value, t Type) (v Value) {
	v.C = C.LLVMGoConstPtrToAddr(val.C, t.C)
	return
}
