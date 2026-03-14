# idgen

Tiny utility for generating stable, human-friendly document identifiers for Couchbase and other document stores. The package exposes a single entry point, `DocID`, that normalizes prefixes and appends a short ID in one of several formats (UUID, timestamp, cryptographic random strings, etc.).

## Features
- Turns any prefix (or the default `doc`) into lowercase, underscore-safe key segments.
- Supports multiple suffix formats: short crypto-safe strings, UUIDs, timestamps, and longer crypto tokens.
- Deterministic prefix handling keeps keys predictable across services.
- Includes a reusable `cryptoRandomString` helper for other secure ID needs.

## Installation

```bash
go get github.com/m-theory-io/idgen
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/m-theory-io/idgen"
)

func main() {
	docID := idgen.DocID("invoice", "short")
	fmt.Println(docID) // invoice:pa7D9sL2wX (example output)
}
```

### Format Options

`DocID(prefix, format)` accepts the following `format` values:

| Format     | Description                                            | Example                     |
|------------|--------------------------------------------------------|-----------------------------|
| `short`    | 10-char crypto-safe string (default)                   | `invoice:Kx8n2AD4Zj`        |
| `crypto`   | 16-char crypto-safe string for higher entropy          | `invoice:dP0sQz81LmN4sQhV`  |
| `uuid`     | RFC 4122 UUID                                          | `invoice:7f8d...`           |
| `timestamp`| Unix nanosecond timestamp                              | `invoice:1736680530123456`  |

Pass `""` for `prefix` to fall back to `doc`.

## `cryptoRandomString`

`cryptoRandomString(n int)` returns a cryptographically secure string drawn from `[a-zA-Z0-9]`. It validates input length and removes modulo bias by only accepting random bytes below 248 before mapping them onto the 62-character charset.

## Developing
- Run tests: `go test ./...`
- Tidy dependencies: `go mod tidy`

## Python Bindings (pybind11)

This repository now includes a `pybind11` binding that exposes `DocID` to Python via a Go `c-shared` library.

### Prerequisites
- Go toolchain with CGO enabled
- CMake (>= 3.18)
- A C++ compiler
- Python with `pybind11` installed (`python3 -m pip install pybind11`)

### Build

From the repository root:

```bash
./python/build.sh
```

The script performs:
1. `go build -buildmode=c-shared` to generate `python/.cgo/libidgen.dylib` and header.
2. CMake build of the Python extension module `_idgen` in `python/build`.

### Use From Python

```python
import sys
sys.path.append("python/build")

import _idgen

print(_idgen.doc_id("invoice", "short"))
print(_idgen.doc_id("invoice", "uuid"))
```

`doc_id(prefix, format)` maps directly to the Go `DocID(prefix, format)` function.

### Smoke Test

After building, run:

```bash
python3 python/smoke_test.py
```

This validates import + basic `DocID` behavior for `short`, `uuid`, and `timestamp` formats.

## License

MIT License © 2026 William J House
