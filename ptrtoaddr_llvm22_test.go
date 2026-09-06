//go:build llvm22 || (!llvm14 && !llvm15 && !llvm16 && !llvm17 && !llvm18 && !llvm19 && !llvm20 && !llvm21)

package llvm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPtrToAddrDataLayoutAndRoundTrip(t *testing.T) {
	ctx := NewContext()
	defer ctx.Dispose()
	mod := ctx.NewModule("ptrtoaddr")
	defer mod.Dispose()
	mod.SetDataLayout("e-p:64:64-p1:64:64:64:32")
	builder := ctx.NewBuilder()
	defer builder.Dispose()

	for _, tc := range []struct {
		name string
		ptr  Type
		addr Type
	}{
		{"scalar", PointerType(ctx.Int8Type(), 0), ctx.Int64Type()},
		{"narrow_index", PointerType(ctx.Int8Type(), 1), ctx.Int32Type()},
		{"vector", VectorType(PointerType(ctx.Int8Type(), 1), 2), VectorType(ctx.Int32Type(), 2)},
	} {
		fn := AddFunction(mod, tc.name, FunctionType(tc.addr, []Type{tc.ptr}, false))
		builder.SetInsertPointAtEnd(ctx.AddBasicBlock(fn, "entry"))
		addr := builder.CreatePtrToAddr(fn.Param(0), "addr")
		if addr.Type() != tc.addr || addr.InstructionOpcode() != PtrToAddr {
			t.Fatalf("%s: unexpected ptrtoaddr: %s", tc.name, addr.String())
		}
		builder.CreateRet(addr)
	}

	global := AddGlobal(mod, ctx.Int8Type(), "data")
	addr := ConstPtrToAddr(global, ctx.Int64Type())
	if addr.Opcode() != PtrToAddr {
		t.Fatalf("unexpected constant opcode: %s", addr.String())
	}
	alias := AddGlobal(mod, ctx.Int64Type(), "address")
	alias.SetInitializer(addr)
	if err := VerifyModule(mod, ReturnStatusAction); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "ptrtoaddr.ll")
	if err := os.WriteFile(path, []byte(mod.String()), 0600); err != nil {
		t.Fatal(err)
	}
	buf, err := NewMemoryBufferFromFile(path)
	if err != nil {
		t.Fatal(err)
	}
	roundTrip, err := ctx.ParseIR(buf) // ParseIR consumes buf.
	if err != nil {
		t.Fatal(err)
	}
	defer roundTrip.Dispose()
	if err := VerifyModule(roundTrip, ReturnStatusAction); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(roundTrip.String(), "ptrtoaddr (ptr @data to i64)") {
		t.Fatalf("constant expression lost during round trip:\n%s", roundTrip.String())
	}
}

func TestPtrToAddrDoesNotExposeProvenance(t *testing.T) {
	ctx := NewContext()
	defer ctx.Dispose()
	mod := ctx.NewModule("provenance")
	defer mod.Dispose()
	mod.SetDataLayout("e-p:64:64")
	b := ctx.NewBuilder()
	defer b.Dispose()
	exposeTy := FunctionType(ctx.VoidType(), []Type{ctx.Int64Type()}, false)
	expose := AddFunction(mod, "observe_address", exposeTy)
	for _, name := range []string{"address", "pointer"} {
		fn := AddFunction(mod, name, FunctionType(ctx.Int32Type(), nil, false))
		b.SetInsertPointAtEnd(ctx.AddBasicBlock(fn, "entry"))
		ptr := b.CreateAlloca(ctx.Int32Type(), "p")
		b.CreateStore(ConstInt(ctx.Int32Type(), 7, false), ptr)
		var addr Value
		if name == "address" {
			addr = b.CreatePtrToAddr(ptr, "addr")
		} else {
			addr = b.CreatePtrToInt(ptr, ctx.Int64Type(), "addr")
		}
		b.CreateCall(exposeTy, expose, []Value{addr}, "")
		b.CreateRet(b.CreateLoad(ctx.Int32Type(), ptr, "value"))
	}
	options := NewPassBuilderOptions()
	defer options.Dispose()
	if err := mod.RunPasses("default<O2>", TargetMachine{}, options); err != nil {
		t.Fatal(err)
	}
	if err := VerifyModule(mod, ReturnStatusAction); err != nil {
		t.Fatal(err)
	}
	if ir := mod.NamedFunction("address").String(); !strings.Contains(ir, "ret i32 7") || strings.Contains(ir, "load i32") {
		t.Fatalf("ptrtoaddr unexpectedly exposed provenance:\n%s", ir)
	}
	if ir := mod.NamedFunction("pointer").String(); !strings.Contains(ir, "load i32") {
		t.Fatalf("ptrtoint control did not expose provenance:\n%s", ir)
	}
}
