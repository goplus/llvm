// Part of the LLVM Project, under the Apache License v2.0 with LLVM Exceptions.
// See https://llvm.org/LICENSE.txt for license information.
// SPDX-License-Identifier: Apache-2.0 WITH LLVM-exception

package llvm

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFloatBits(t *testing.T) {
	// Parse independently specified IR, then rebuild in another context. Values
	// include precision below float64's significand and NaN payload/signaling bits.
	cases := []struct {
		name, ir string
		words    []uint64
	}{
		{"half", "half 0xH8000", []uint64{0x8000}},
		{"bfloat", "bfloat 0xR7FC1", []uint64{0x7fc1}},
		{"float_nan", "float 0x7FF82468A0000000", []uint64{0x7fc12345}},
		{"double_negative_zero", "double 0x8000000000000000", []uint64{0x8000000000000000}},
		{"double_snan", "double 0x7FF0000000000001", []uint64{0x7ff0000000000001}},
		{"double_infinity", "double 0x7FF0000000000000", []uint64{0x7ff0000000000000}},
		{"double_subnormal", "double 0x0000000000000001", []uint64{1}},
		{"fp80_precision", "x86_fp80 0xK3FFF8000000000000001", []uint64{0x8000000000000001, 0x3fff}},
		{"fp80_nan", "x86_fp80 0xK7FFFC000000000000123", []uint64{0xc000000000000123, 0x7fff}},
		{"fp128_precision", "fp128 0xL00000000000000013FFF000000000000", []uint64{1, 0x3fff000000000000}},
		{"fp128_nan", "fp128 0xL00000000000001237FFF800000000000", []uint64{0x123, 0x7fff800000000000}},
		{"ppc_fp128_low_only", "ppc_fp128 0xM00000000000000003FF0000000000000", []uint64{0, 0x3ff0000000000000}},
		{"ppc_fp128_negative_zero", "ppc_fp128 0xM80000000000000000000000000000000", []uint64{0x8000000000000000, 0}},
		{"ppc_fp128_noncanonical", "ppc_fp128 0xM3FF00000000000003FF0000000000000", []uint64{0x3ff0000000000000, 0x3ff0000000000000}},
		{"ppc_fp128_nan", "ppc_fp128 0xM7FF80000000001230000000000000000", []uint64{0x7ff8000000000123, 0}},
		{"ppc_fp128_precision", "ppc_fp128 0xM3FF00000000000003C90000000000000", []uint64{0x3ff0000000000000, 0x3c90000000000000}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srcCtx, dstCtx := NewContext(), NewContext()
			defer dstCtx.Dispose()
			path := filepath.Join(t.TempDir(), "float.ll")
			if err := os.WriteFile(path, []byte("@value = constant "+tc.ir+"\n"), 0600); err != nil {
				t.Fatal(err)
			}
			buf, err := NewMemoryBufferFromFile(path)
			if err != nil {
				t.Fatal(err)
			}
			src, err := srcCtx.ParseIR(buf)
			if err != nil {
				t.Fatal(err)
			}
			value := src.NamedGlobal("value").Initializer()
			words := value.FloatBits()
			if !reflect.DeepEqual(words, tc.words) {
				t.Fatalf("bits = %x, want %x", words, tc.words)
			}
			// Reparse just to obtain the identical destination-context type.
			buf, err = NewMemoryBufferFromFile(path)
			if err != nil {
				t.Fatal(err)
			}
			dst, err := dstCtx.ParseIR(buf)
			if err != nil {
				t.Fatal(err)
			}
			defer dst.Dispose()
			expected := dst.NamedGlobal("value").Initializer()
			cloned := ConstFloatFromBits(expected.Type(), words)
			if cloned != expected {
				t.Fatalf("rebuilt %s, want %s", cloned, expected)
			}
			src.Dispose()
			srcCtx.Dispose()
			if got := cloned.FloatBits(); !reflect.DeepEqual(got, tc.words) {
				t.Fatalf("after source disposal: %x", got)
			}
			words[0] ^= 1
			if got := cloned.FloatBits(); !reflect.DeepEqual(got, tc.words) {
				t.Fatalf("returned slice aliases constant: %x", got)
			}
			if err := VerifyModule(dst, ReturnStatusAction); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestFloatBitsInvalidInput(t *testing.T) {
	ctx := NewContext()
	defer ctx.Dispose()
	for name, call := range map[string]func(){
		"nil_type":              func() { ConstFloatFromBits(Type{}, []uint64{0}) },
		"integer_type":          func() { ConstFloatFromBits(ctx.Int64Type(), []uint64{0}) },
		"vector_type":           func() { ConstFloatFromBits(VectorType(ctx.DoubleType(), 2), []uint64{0, 0}) },
		"empty":                 func() { ConstFloatFromBits(ctx.DoubleType(), nil) },
		"too_few":               func() { ConstFloatFromBits(ctx.FP128Type(), []uint64{0}) },
		"too_many":              func() { ConstFloatFromBits(ctx.DoubleType(), []uint64{0, 0}) },
		"unused_high_bits":      func() { ConstFloatFromBits(ctx.FloatType(), []uint64{1 << 32}) },
		"fp80_unused_high_bits": func() { ConstFloatFromBits(ctx.X86FP80Type(), []uint64{0, 1 << 16}) },
		"nil_value":             func() { Value{}.FloatBits() },
		"integer_value":         func() { ConstInt(ctx.Int64Type(), 0, false).FloatBits() },
		"undef":                 func() { Undef(ctx.DoubleType()).FloatBits() },
	} {
		t.Run(name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("expected panic")
				}
			}()
			call()
		})
	}
}
