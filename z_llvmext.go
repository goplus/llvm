/*
 * Copyright (c) 2024 The GoPlus Authors (goplus.org). All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llvm

/*
#include "llvm-c/Core.h"
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func (c Context) Finalize() {
	runtime.SetFinalizer(c.C, func(c C.LLVMContextRef) {
		Context{c}.Dispose()
	})
}

func (m Module) Finalize() {
	runtime.SetFinalizer(m.C, func(m C.LLVMModuleRef) {
		Module{m}.Dispose()
	})
}

func (b Builder) Finalize() {
	runtime.SetFinalizer(b.C, func(b C.LLVMBuilderRef) {
		Builder{b}.Dispose()
	})
}

var (
	emptyCStr [1]C.char
)

func CreateBinOp(b Builder, op Opcode, lhs, rhs Value) (v Value) {
	v.C = C.LLVMBuildBinOp(b.C, C.LLVMOpcode(op), lhs.C, rhs.C, &emptyCStr[0])
	return
}

func CreatePHI(b Builder, t Type) (v Value) {
	v.C = C.LLVMBuildPhi(b.C, t.C, &emptyCStr[0])
	return
}

func CreateICmp(b Builder, op IntPredicate, lhs, rhs Value) (v Value) {
	v.C = C.LLVMBuildICmp(b.C, C.LLVMIntPredicate(op), lhs.C, rhs.C, &emptyCStr[0])
	return
}

func CreateFCmp(b Builder, op FloatPredicate, lhs, rhs Value) (v Value) {
	v.C = C.LLVMBuildFCmp(b.C, C.LLVMRealPredicate(op), lhs.C, rhs.C, &emptyCStr[0])
	return
}

func CreateIntCast(b Builder, v Value, t Type) (r Value) {
	r.C = C.LLVMBuildIntCast(b.C, v.C, t.C, &emptyCStr[0])
	return
}

func CreateTrunc(b Builder, v Value, t Type) (r Value) {
	r.C = C.LLVMBuildTrunc(b.C, v.C, t.C, &emptyCStr[0])
	return
}

func CreatePtrToInt(b Builder, v Value, t Type) (r Value) {
	r.C = C.LLVMBuildPtrToInt(b.C, v.C, t.C, &emptyCStr[0])
	return
}

func CreateIntToPtr(b Builder, v Value, t Type) (r Value) {
	r.C = C.LLVMBuildIntToPtr(b.C, v.C, t.C, &emptyCStr[0])
	return
}

func CreatePointerCast(b Builder, v Value, t Type) (r Value) {
	r.C = C.LLVMBuildPointerCast(b.C, v.C, t.C, &emptyCStr[0])
	return
}

func CreateAlloca(b Builder, t Type) (v Value) {
	v.C = C.LLVMBuildAlloca(b.C, t.C, &emptyCStr[0])
	return
}

func CreateArrayAlloca(b Builder, t Type, n Value) (v Value) {
	v.C = C.LLVMBuildArrayAlloca(b.C, t.C, n.C, &emptyCStr[0])
	return
}

func CreateGEP(b Builder, t Type, p Value, indices []Value) (v Value) {
	ptr, nvals := llvmValueRefs(indices)
	v.C = C.LLVMBuildGEP2(b.C, t.C, p.C, ptr, nvals, &emptyCStr[0])
	return
}

func CreateInBoundsGEP(b Builder, t Type, p Value, indices []Value) (v Value) {
	ptr, nvals := llvmValueRefs(indices)
	v.C = C.LLVMBuildInBoundsGEP2(b.C, t.C, p.C, ptr, nvals, &emptyCStr[0])
	return
}

func CreateStructGEP(b Builder, t Type, p Value, i int) (v Value) {
	v.C = C.LLVMBuildStructGEP2(b.C, t.C, p.C, C.unsigned(i), &emptyCStr[0])
	return
}

func CreateExtractValue(b Builder, agg Value, i int) (v Value) {
	v.C = C.LLVMBuildExtractValue(b.C, agg.C, C.unsigned(i), &emptyCStr[0])
	return
}

func CreateLoad(b Builder, t Type, p Value) (v Value) {
	v.C = C.LLVMBuildLoad2(b.C, t.C, p.C, &emptyCStr[0])
	return
}

func CreateCall(b Builder, t Type, fn Value, args []Value) (v Value) {
	ptr, nvals := llvmValueRefs(args)
	v.C = C.LLVMBuildCall2(b.C, t.C, fn.C, ptr, nvals, &emptyCStr[0])
	return
}

func CreateGlobalStringPtr(b Builder, str string) (v Value) {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))
	v.C = C.LLVMBuildGlobalStringPtr(b.C, cstr, &emptyCStr[0])
	return
}
