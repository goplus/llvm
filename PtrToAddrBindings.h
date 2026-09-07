#ifndef LLVM_BINDINGS_GO_PTRTOADDR_H
#define LLVM_BINDINGS_GO_PTRTOADDR_H

#include "llvm-c/Core.h"

#ifdef __cplusplus
extern "C" {
#endif

LLVMValueRef LLVMGoBuildPtrToAddr(LLVMBuilderRef B, LLVMValueRef V,
                                const char *Name);
LLVMValueRef LLVMGoConstPtrToAddr(LLVMValueRef V, LLVMTypeRef T);

#ifdef __cplusplus
}
#endif

#endif
