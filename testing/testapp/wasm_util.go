package testapp

import (
	wasmtime "github.com/bytecodealliance/wasmtime-go/v20"

	"github.com/ODIN-PROTOCOL/odin-core/x/oracle/types"
)

func compile(code []byte) []byte {
	compiled, err := OwasmVM.Compile(code, types.MaxCompiledWasmCodeSize)
	if err != nil {
		panic(err)
	}
	return compiled
}

func wat2wasm(wat string) []byte {
	wasm, err := wasmtime.Wat2Wasm(wat)
	if err != nil {
		panic(err)
	}

	return wasm
}
