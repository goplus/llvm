#include "PtrToAddrBindings.h"
#include "llvm/Config/llvm-config.h"

#if LLVM_VERSION_MAJOR >= 22
#include "llvm/IR/Constants.h"
#include "llvm/IR/IRBuilder.h"

using namespace llvm;

LLVMValueRef LLVMGoBuildPtrToAddr(LLVMBuilderRef B, LLVMValueRef V,
                                const char *Name) {
  return wrap(unwrap(B)->CreatePtrToAddr(unwrap(V), Name));
}

LLVMValueRef LLVMGoConstPtrToAddr(LLVMValueRef V, LLVMTypeRef T) {
  return wrap(ConstantExpr::getPtrToAddr(unwrap<Constant>(V), unwrap(T)));
}
#endif
