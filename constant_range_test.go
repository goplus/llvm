package llvm

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestConstantRangeAttribute(t *testing.T) {
	ctx := NewContext()
	defer ctx.Dispose()
	kind := AttributeKindID("range")
	major, _ := strconv.Atoi(strings.SplitN(Version, ".", 2)[0])
	if major < 19 {
		if !ctx.CreateConstantRangeAttribute(kind, 32, []uint64{0}, []uint64{1}).IsNil() {
			t.Fatal("constant-range attributes must be unavailable before LLVM 19")
		}
		return
	}
	if kind == 0 {
		t.Fatal("range attribute kind not found")
	}
	for _, test := range []struct {
		bits         int
		lower, upper []uint64
		want         string
	}{
		{1, []uint64{0}, []uint64{1}, "range(i1 0, -1)"},
		{32, []uint64{0}, []uint64{1 << 31}, "range(i32 0, -2147483648)"},
		{64, []uint64{0}, []uint64{1 << 63}, "range(i64 0, -9223372036854775808)"},
		{65, []uint64{3, 0}, []uint64{9, 1}, "range(i65 3, -18446744073709551607)"},
		{128, []uint64{3, 2}, []uint64{9, 4}, "range(i128 36893488147419103235, 73786976294838206473)"},
	} {
		t.Run(fmt.Sprint(test.bits), func(t *testing.T) {
			mod := ctx.NewModule("range")
			defer mod.Dispose()
			fn := AddFunction(mod, "length", FunctionType(ctx.IntType(test.bits), nil, false))
			attr := ctx.CreateConstantRangeAttribute(kind, test.bits, test.lower, test.upper)
			if attr.IsNil() {
				t.Fatal("constant-range attribute is nil")
			}
			fn.AddAttributeAtIndex(0, attr)
			if attrs := fn.GetAttributesAtIndex(0); len(attrs) != 1 || attrs[0] != attr {
				t.Fatal("return attribute did not round-trip")
			}
			// The context owns the APInts; it must not retain Go slice storage.
			for i := range test.lower {
				test.lower[i] = 0
				test.upper[i] = 0
			}
			if ir := mod.String(); !strings.Contains(ir, test.want) {
				t.Fatalf("missing %q:\n%s", test.want, ir)
			}
			if err := VerifyModule(mod, ReturnStatusAction); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestConstantRangeAttributeInvalidBounds(t *testing.T) {
	ctx := NewContext()
	defer ctx.Dispose()
	for _, test := range []struct {
		bits         int
		lower, upper []uint64
	}{
		{0, nil, nil}, {-1, nil, nil},
		{32, nil, []uint64{1}}, {32, []uint64{0}, nil},
		{65, []uint64{0}, []uint64{1}},
		{64, []uint64{0, 0}, []uint64{1}},
	} {
		t.Run(fmt.Sprint(test.bits, "/", len(test.lower), "/", len(test.upper)), func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("invalid bounds did not panic")
				}
			}()
			ctx.CreateConstantRangeAttribute(AttributeKindID("range"), test.bits, test.lower, test.upper)
		})
	}
}
