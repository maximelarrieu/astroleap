package web

import "embed"

// FS embeds the WebAssembly game files.
//
//go:embed index.html wasm_exec.js game.wasm
var FS embed.FS
