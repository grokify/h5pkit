# Installation

## Requirements

- Go 1.19 or later
- Git (for fetching the module)

## Install via Go Modules

The easiest way to install the H5P Go SDK is using Go modules:

```bash
go get github.com/grokify/h5p-go
```

## Verify Installation

Create a simple test file to verify the installation:

```go
// test.go
package main

import (
    "fmt"
    "github.com/grokify/h5p-go"
)

func main() {
    builder := h5p.NewQuestionSetBuilder()
    fmt.Println("H5P Go SDK installed successfully!")
    fmt.Printf("Builder created: %T\n", builder)
}
```

Run the test:

```bash
go run test.go
```

You should see:
```
H5P Go SDK installed successfully!
Builder created: *h5p.QuestionSetBuilder
```

## Development Installation

If you want to contribute to the project or need the latest development version:

```bash
git clone https://github.com/grokify/h5p-go.git
cd h5p-go
go mod tidy
go test ./...
```

## Dependencies

The H5P Go SDK has minimal external dependencies:

- Standard library packages only for core functionality
- No runtime dependencies for basic usage
- Optional dependencies for advanced features (validation, etc.)

## Next Steps

- [Quick Start Tutorial](quick-start.md) - Your first H5P content
- [Basic Concepts](concepts.md) - Understanding H5P fundamentals