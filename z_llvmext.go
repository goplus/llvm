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
*/
import "C"
import "runtime"

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

func CreateBinOp(b Builder, op Opcode, lhs, rhs Value) (v Value) {
	v.C = C.LLVMBuildBinOp(b.C, C.LLVMOpcode(op), lhs.C, rhs.C, nil)
	return
}

func CreateICmp(b Builder, op IntPredicate, lhs, rhs Value) (v Value) {
	v.C = C.LLVMBuildICmp(b.C, C.LLVMIntPredicate(op), lhs.C, rhs.C, nil)
	return
}

func CreateFCmp(b Builder, op FloatPredicate, lhs, rhs Value) (v Value) {
	v.C = C.LLVMBuildFCmp(b.C, C.LLVMRealPredicate(op), lhs.C, rhs.C, nil)
	return
}

func CreateAlloca(b Builder, t Type) (v Value) {
	v.C = C.LLVMBuildAlloca(b.C, t.C, nil)
	return
}

func CreateStructGEP(b Builder, t Type, p Value, i int) (v Value) {
	v.C = C.LLVMBuildStructGEP2(b.C, t.C, p.C, C.unsigned(i), nil)
	return
}

func CreateLoad(b Builder, t Type, p Value) (v Value) {
	v.C = C.LLVMBuildLoad2(b.C, t.C, p.C, nil)
	return
}
