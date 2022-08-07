//go:build !byollvm && darwin
// +build !byollvm,darwin

package llvm

// Automatically generated by `make config BUILDDIR=`, do not edit.

// #cgo amd64 CPPFLAGS: -I/usr/local/opt/llvm@14/include   -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
// #cgo amd64 CXXFLAGS: -std=c++14
// #cgo amd64 LDFLAGS: -L/usr/local/opt/llvm@14/lib -Wl,-search_paths_first -Wl,-headerpad_max_install_names -lLLVM -lz -lcurses -lm -lxml2 -L/usr/local/opt/libffi/lib -lffi
// #cgo arm64 CPPFLAGS: -I/opt/homebrew/opt/llvm@14/include   -D__STDC_CONSTANT_MACROS -D__STDC_FORMAT_MACROS -D__STDC_LIMIT_MACROS
// #cgo arm64 CXXFLAGS: -std=c++14
// #cgo arm64 LDFLAGS: -L/opt/homebrew/opt/llvm@14/lib -Wl,-search_paths_first -Wl,-headerpad_max_install_names -lLLVM -lz -lcurses -lm -lxml2 -L/opt/homebrew/opt/libffi/lib -lffi
import "C"

type run_build_sh int
