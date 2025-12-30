# Go Client-Side Playground (WASM) 🚀

![Go](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)
![WASM](https://img.shields.io/badge/WebAssembly-Enabled-654FF0?style=flat&logo=webassembly)
![License](https://img.shields.io/badge/License-MIT-green)

A fully **serverless** Go playground that runs entirely in the browser.

Unlike standard playgrounds that send code to a backend server for execution, this project compiles the **Yaegi** interpreter into **WebAssembly**. This allows the browser to interpret and execute Go code locally on the user's machine with zero latency and offline capabilities.

## 🌟 Features

* **Serverless Architecture:** No backend API, Docker containers, or cloud functions required.
* **True Client-Side Execution:** Code runs 100% inside the user's browser tab.
* **Instant Feedback:** No network round-trips; execution happens immediately.
* **Safe Sandbox:** Runs inside the browser's WebAssembly sandbox.
* **Standard Library Support:** Supports `fmt`, `math`, `time`, and most standard Go libraries via Yaegi.

---

## 🏗 Architecture

The project consists of three main layers:

1. **The Frontend (`index.html`):** A simple UI that accepts user input and displays output.
2. **The Bridge (`wasm_exec.js`):** The official Go glue code that allows the browser to load the WASM binary.
3. **The Engine (`main.go` -> `main.wasm`):** A compiled Go binary embedding the **Yaegi** interpreter. It listens for input from JavaScript, evaluates it, and captures `Stdout`/`Stderr`.

---

## 🛠 Prerequisites

* **Go 1.18+** (Required for `GOOS=js` support)
* **Python 3** (Or any local static file server like `http-server`)
* **WSL / Linux / macOS** (Recommended for build commands)

---

## 🚀 Quick Start

### 1. Installation

Clone the repository (or create the folder) and initialize the dependencies:

```bash
# Initialize module
go mod init browser-go

# Install Yaegi (The Interpreter Engine)
go get [github.com/traefik/yaegi/interp](https://github.com/traefik/yaegi/interp)
go get [github.com/traefik/yaegi/stdlib](https://github.com/traefik/yaegi/stdlib)

### 2. Build the WASM Binary

You need to compile the Go code into a `.wasm` file and copy the JavaScript bridge provided by the Go toolchain.

**Standard Build (Recommended for Development):**
This keeps debug symbols intact, useful if you are modifying the interpreter code itself.

```bash
# 1. Copy the JS Bridge from your Go installation
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .

# 2. Compile the Interpreter to WASM
GOOS=js GOARCH=wasm go build -o main.wasm main.go

```

*Note: The resulting `main.wasm` file will be approximately 30-36MB. This is normal as it contains the Go runtime, Yaegi interpreter, and debug symbols.*

**Production Build (Optional - Smaller Size):**
If you want to reduce the file size significantly (often by ~50%), use the `-ldflags="-s -w"` flags to strip debug symbols.

```bash
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o main.wasm main.go

```

### 3. Run Locally

Due to browser **CORS policies**, you cannot open `index.html` directly from the file system. You must serve it via a local web server.

```bash
# Using Python 3 (Pre-installed on most systems)
python3 -m http.server 8080

```

Open your browser to: **[http://localhost:8080](https://www.google.com/search?q=http://localhost:8080)**

---

## 🧩 Technical Details

### How it works under the hood

1. **Initialization:** The `main()` function in Go creates a channel that blocks forever (`<-c`). This prevents the WASM process from exiting immediately, keeping the runtime alive in the browser memory.
2. **Exposing Functions:** We use `syscall/js` to bind our Go function `runGoCode` to the global JavaScript scope:

```go
js.Global().Set("runGoCode", js.FuncOf(runGoCode))

```

3. **Output Capture:** Since the browser has no "terminal," we replace `Stdout` and `Stderr` with `bytes.Buffer` inside the interpreter. When the code finishes, we convert these buffers to strings and return them to JavaScript.

### Optimization Tips (For Production)

The raw WASM file is large. If you deploy this to a static host (GitHub Pages, Netlify, Vercel), ensure you use compression:

* **Gzip:** Reduces size to ~4-5MB.
* **Brotli:** Reduces size to ~3MB.

Most modern CDNs handle this automatically.

---

## 📂 Project Structure

```text
.
├── go.mod           # Go module definition
├── go.sum           # Dependency checksums
├── index.html       # The Frontend UI
├── main.go          # The Go "Backend" (Interpreter logic)
├── main.wasm        # Compiled Binary (Generated after build)
└── wasm_exec.js     # Go-to-JS Bridge (Copied from GOROOT)

```

📖 Tutorial
Read the full step-by-step guide on Medium:

Building a Serverless Go Playground on Windows (using WSL + WebAssembly)

## 🤝 Contributing

Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## 📄 License

[MIT](https://choosealicense.com/licenses/mit/)

```
