package main

import (
	"bytes"
	"syscall/js"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func main() {
	// 1. Create a channel to keep the WASM running indefinitely
	c := make(chan struct{}, 0)

	// 2. Expose the "runGoCode" function to JavaScript
	js.Global().Set("runGoCode", js.FuncOf(runGoCode))

	println("✅ WASM Interpreter Initialized")
	
	<-c // Block forever
}

// This function is called from JavaScript
func runGoCode(this js.Value, args []js.Value) interface{} {
	// Get the code string passed from JS
	code := args[0].String()

	// Create buffers to capture Standard Output and Error
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer

	// Initialize Yaegi
	i := interp.New(interp.Options{
		Stdout: &outBuf,
		Stderr: &errBuf,
	})

	// Load standard libraries (fmt, math, etc.)
	if err := i.Use(stdlib.Symbols); err != nil {
		return "Error loading stdlib: " + err.Error()
	}

	// Execute the code
	_, err := i.Eval(code)

	// Combine output
	output := outBuf.String()
	if errBuf.Len() > 0 {
		output += "\n--- STDERR ---\n" + errBuf.String()
	}
	if err != nil {
		output += "\n--- ERROR ---\n" + err.Error()
	}

	return output
}