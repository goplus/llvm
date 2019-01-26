
all: help

.PHONY: all help config update

COMPONENTS = all-targets analysis asmparser asmprinter bitreader bitwriter codegen core coroutines debuginfodwarf executionengine instrumentation interpreter ipo irreader linker mc mcjit objcarcopts option profiledata scalaropts support target

VERSION=7.0.0
VERSION_MAJOR=$(firstword $(subst ., ,$(VERSION)))

SRCDIR=llvm-$(VERSION)
UNAME_S:=$(shell uname -s)

GOOS := $(shell go env GOOS)

LDFLAGS = $(shell $(CONFIG) --ldflags --libs --system-libs $(COMPONENTS))
ifeq ($(BUILDDIR),)
	ifeq ($(UNAME_S),Darwin)
	CONFIG = /usr/local/Cellar/llvm/$(VERSION)/bin/llvm-config
	LDFLAGS += -L/usr/local/opt/libffi/lib -lffi
	else
	CONFIG = llvm-config-$(VERSION_MAJOR)
	endif
else
CONFIG=$(BUILDDIR)/bin/llvm-config
endif

help:
	@echo "Use one of the following commands:"
	@echo "  config:    update llvm_config_*.go to use a locally installed LLVM version"
	@echo "  update:    update local .go/.cpp/.h files (see SRCDIR)"
	@echo "  clean:     remove local .go/.cpp/.h files"
	@echo "  clean-all: remove all LLVM files, including downloads"
	@echo ""
	@echo "Provided flags:"
	@echo "  VERSION=$(VERSION)"
	@echo "    LLVM version to use"
	@echo "  SRCDIR=$(SRCDIR)"
	@echo "    LLVM sources root"
	@echo "  BUILDDIR=$(BUILDDIR)"
	@echo "    LLVM sources root"

config:
	@echo "// +build !byollvm" > llvm_config_$(GOOS).go
	@echo "" >> llvm_config_$(GOOS).go
	@echo "package llvm" >> llvm_config_$(GOOS).go
	@echo "" >> llvm_config_$(GOOS).go
	@echo "// Automatically generated by \`make config BUILDDIR=$(BUILDDIR)\`, do not edit." >> llvm_config_$(GOOS).go
	@echo "" >> llvm_config_$(GOOS).go
	@echo "// #cgo CPPFLAGS: $(shell $(CONFIG) --cppflags)" >> llvm_config_$(GOOS).go
	@echo "// #cgo CXXFLAGS: -std=c++11" >> llvm_config_$(GOOS).go
	@echo "// #cgo LDFLAGS: $(LDFLAGS)" >> llvm_config_$(GOOS).go
	@echo "import \"C\"" >> llvm_config_$(GOOS).go
	@echo "" >> llvm_config_$(GOOS).go
	@echo "type (run_build_sh int)" >> llvm_config_$(GOOS).go

update: $(SRCDIR) clean
	@cp -rp "$(SRCDIR)"/bindings/go/llvm/*.go .
	@cp -rp "$(SRCDIR)"/bindings/go/llvm/*.cpp .
	@cp -rp "$(SRCDIR)"/bindings/go/llvm/*.h .

clean:
	@rm -f llvm_config_*.go
	@rm -f $(sort $(filter-out backport%,$(wildcard *.go)))
	@rm -f $(sort $(filter-out backport%,$(wildcard *.cpp)))
	@rm -f $(sort $(filter-out backport%,$(wildcard *.h)))

clean-all: clean
	@rm -rf llvm-$(VERSION)
	@rm -f llvm-$(VERSION).src.tar.xz

# Downlaod the given LLVM release.
llvm-$(VERSION).src.tar.xz:
	wget -O llvm-$(VERSION).src.tar.xz http://releases.llvm.org/$(VERSION)/llvm-$(VERSION).src.tar.xz

# Extract the given LLVM release to a temporary directory.
llvm-$(VERSION): llvm-$(VERSION).src.tar.xz
	mkdir -p $@
	tar -C $@ --strip-components=1 -xf $< $@.src
