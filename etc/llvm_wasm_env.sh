#!/usr/bin/env bash

# Set this to whatever your LLVM build dir is.
# --> See llvm-mirror/llvm/bindings/go/README.txt.
# If you use llvm-mirror/llvm/bindings/go/build.sh to build, your build path will look
# something like this:
llvm_build=/Users/cwolff/workspace/github.com/llvm-mirror/llvm/bindings/go/llvm/workdir/llvm_build

export CGO_CPPFLAGS="`${llvm_build}/bin/llvm-config --cppflags`"
export CGO_CXXFLAGS=-std=c++11
export CGO_LDFLAGS="`${llvm_build}/bin/llvm-config --ldflags --libs --system-libs all`"

# Set this to whatever your emsdk repo is:
# --> See https://webassembly.org/getting-started/developers-guide/
path_to_emsdk=/Users/cwolff/workspace/github.com/emscripten-core/emsdk
source ${path_to_emsdk}/emsdk_env.sh

# NOTE: emsdk and Go LLVM bindings use different versions of LLVM :( so it's required to first set the
#   CGO environemt and then set up the emsdk environment afterwards.  Otherwise emsdk_env.sh will
#   pollute the environment with its own value of LLVM_ROOT.
#
# To serve WASM from Flux:
#   go run -tags byollvm ./cmd/flux serve
