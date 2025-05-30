
#include "backports.h"
#include "llvm/IR/Instructions.h"
#if LLVM_VERSION_MAJOR >= 16
#include "llvm/IR/PassManager.h"
#include "llvm/Analysis/LoopAnalysisManager.h"
#include "llvm/Analysis/CGSCCPassManager.h"
#include "llvm/Passes/PassBuilder.h"
#include "llvm/Transforms/IPO/ThinLTOBitcodeWriter.h"
#else
#include "llvm/IR/LegacyPassManager.h"
#include "llvm/Transforms/IPO/PassManagerBuilder.h"
#endif
#include "llvm/IR/Module.h"
#include "llvm/Pass.h"
#include "llvm/Support/MemoryBuffer.h"
#include "llvm/Transforms/IPO.h"
#include "llvm-c/DebugInfo.h"

void LLVMGlobalObjectAddMetadata(LLVMValueRef Global, unsigned KindID, LLVMMetadataRef MD) {
  llvm::MDNode *N = MD ? llvm::unwrap<llvm::MDNode>(MD) : nullptr;
  llvm::GlobalObject *O = llvm::unwrap<llvm::GlobalObject>(Global);
  O->addMetadata(KindID, *N);
}

// See https://reviews.llvm.org/D119431
LLVMMemoryBufferRef LLVMGoWriteThinLTOBitcodeToMemoryBuffer(LLVMModuleRef M) {
  std::string Data;
  llvm::raw_string_ostream OS(Data);
#if LLVM_VERSION_MAJOR >= 16
  llvm::LoopAnalysisManager LAM;
  llvm::FunctionAnalysisManager FAM;
  llvm::CGSCCAnalysisManager CGAM;
  llvm::ModuleAnalysisManager MAM;
  llvm::PassBuilder PB;
  PB.registerModuleAnalyses(MAM);
  PB.registerCGSCCAnalyses(CGAM);
  PB.registerFunctionAnalyses(FAM);
  PB.registerLoopAnalyses(LAM);
  PB.crossRegisterProxies(LAM, FAM, CGAM, MAM);
  llvm::ModulePassManager MPM;
  MPM.addPass(llvm::ThinLTOBitcodeWriterPass(OS, nullptr));
  MPM.run(*llvm::unwrap(M), MAM);
#else
  llvm::legacy::PassManager PM;
  PM.add(createWriteThinLTOBitcodePass(OS));
  PM.run(*llvm::unwrap(M));
#endif
  return llvm::wrap(llvm::MemoryBuffer::getMemBufferCopy(OS.str()).release());
}

void LLVMGoDIBuilderInsertDbgValueRecordAtEnd(
    LLVMDIBuilderRef Builder, LLVMValueRef Val, LLVMMetadataRef VarInfo,
    LLVMMetadataRef Expr, LLVMMetadataRef DebugLoc, LLVMBasicBlockRef Block) {
#if LLVM_VERSION_MAJOR >= 19
  // Note: this returns a LLVMDbgRecordRef. Previously, InsertValueAtEnd would
  // return a Value. But since the type changed, and I'd like to keep the API
  // consistent across LLVM versions, I decided to drop the return value.
  LLVMDIBuilderInsertDbgValueRecordAtEnd(Builder, Val, VarInfo, Expr, DebugLoc, Block);
#else
  // Old llvm.dbg.* API.
  LLVMDIBuilderInsertDbgValueAtEnd(Builder, Val, VarInfo, Expr, DebugLoc, Block);
#endif
}

void LLVMGoDIBuilderInsertDbgDeclareRecordAtEnd(
    LLVMDIBuilderRef Builder, LLVMValueRef Val, LLVMMetadataRef VarInfo,
    LLVMMetadataRef Expr, LLVMMetadataRef DebugLoc, LLVMBasicBlockRef Block) {
#if LLVM_VERSION_MAJOR >= 19
  // Note: this returns a LLVMDbgRecordRef. Previously, InsertDeclareAtEnd would
  // return a Declare. But since the type changed, and I'd like to keep the API
  // consistent across LLVM versions, I decided to drop the return value.
  LLVMDIBuilderInsertDeclareRecordAtEnd(Builder, Val, VarInfo, Expr, DebugLoc, Block);
#else
  // Old llvm.dbg.* API.
  LLVMDIBuilderInsertDeclareAtEnd(Builder, Val, VarInfo, Expr, DebugLoc, Block);
#endif
}