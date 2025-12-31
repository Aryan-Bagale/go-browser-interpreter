// main.go
//
// This file is compiled to WebAssembly (WASM) and runs directly
// inside the user's browser.
//
// It embeds the Yaegi Go interpreter so that we can execute
// arbitrary Go code written by the user at runtime — without
// any backend server.

package main

import (
	"bytes"      // Used to capture printed output (fmt.Println, etc.)

	// syscall/js allows Go code running in WebAssembly
	// to interact with JavaScript in the browser.
	"syscall/js"

	// Yaegi interpreter packages
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func main() {
	// ------------------------------------------------------------
	// IMPORTANT: Prevent the WASM program from exiting
	// ------------------------------------------------------------
	//
	// In a normal Go program, main() exits when it reaches the end.
	// In WebAssembly, exiting main() would immediately stop
	// our Go runtime and make exported functions unusable.
	//
	// We block forever using a channel so the program stays alive.
	keepAlive := make(chan struct{})

	// ------------------------------------------------------------
	// Expose a Go function to JavaScript
	// ------------------------------------------------------------
	//
	// js.Global() refers to the browser's global object (window).
	// js.FuncOf converts a Go function into a JavaScript function.
	//
	// After this, JavaScript can call:
	//     runGoCode("some go code")
	js.Global().Set("runGoCode", js.FuncOf(runGoCode))

	println("✅ Go WASM Interpreter initialized")

	// Block forever so WASM stays alive
	<-keepAlive
}

// runGoCode is called from JavaScript.
// It receives Go source code as a string, executes it using
// Yaegi, and returns all output back to JavaScript.
func runGoCode(this js.Value, args []js.Value) interface{} {

	// ------------------------------------------------------------
	// Read the code string passed from JavaScript
	// ------------------------------------------------------------
	code := args[0].String()

	// ------------------------------------------------------------
	// Create buffers to capture stdout and stderr
	// ------------------------------------------------------------
	//
	// Any fmt.Println(), panic(), or runtime error messages
	// produced by the interpreted code will be written here.
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	// ------------------------------------------------------------
	// Create a NEW Yaegi interpreter
	// ------------------------------------------------------------
	//
	// IMPORTANT:
	// In Yaegi v0.16.x, Stdout and Stderr can ONLY be set
	// at interpreter creation time via interp.Options.
	//
	// They CANNOT be modified later.
	i := interp.New(interp.Options{
		Stdout: &stdout,
		Stderr: &stderr,
	})

	// ------------------------------------------------------------
	// Load Go standard library symbols
	// ------------------------------------------------------------
	if err := i.Use(stdlib.Symbols); err != nil {
		return "Error loading standard library: " + err.Error()
	}

	// ------------------------------------------------------------
	// Execute the user's Go code
	// ------------------------------------------------------------
	_, err := i.Eval(code)

	// ------------------------------------------------------------
	// Combine all output into a single string
	// ------------------------------------------------------------
	output := stdout.String()

	if stderr.Len() > 0 {
		output += "\n--- STDERR ---\n" + stderr.String()
	}

	if err != nil {
		output += "\n--- ERROR ---\n" + err.Error()
	}

	return output
}
