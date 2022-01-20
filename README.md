# builder
builder for gopher 

[![GoDoc](https://godoc.org/github.com/things-go/builder?status.svg)](https://godoc.org/github.com/things-go/builder)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/builder?tab=doc)
[![codecov](https://codecov.io/gh/things-go/builder/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/builder)
[![Tests](https://github.com/things-go/builder/actions/workflows/ci.yml/badge.svg)](https://github.com/things-go/builder/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/builder)](https://goreportcard.com/report/github.com/things-go/builder)
[![Licence](https://img.shields.io/github/license/things-go/builder)](https://raw.githubusercontent.com/things-go/builder/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/builder)](https://github.com/things-go/builder/tags)

## Usage

### Installation

Use go get.
```bash
    go get github.com/things-go/builder
```

Then import the package into your own code.
```bash
    import "github.com/things-go/builder"
```

### Example

see [Makefile](_example/Makefile)

[embedmd]:# (_example/main.go go)
```go
package main

import (
	"github.com/things-go/builder"
)

func main() {
	builder.Println()
}
```

## References

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
