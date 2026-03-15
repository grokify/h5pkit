# H5P Go SDK

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/h5p-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/h5p-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/h5p-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/h5p-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/h5p-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/h5p-go/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/h5p-go
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/h5p-go
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/h5p-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/h5p-go
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fh5p-go
 [loc-svg]: https://tokei.rs/b1/github/grokify/h5p-go
 [repo-url]: https://github.com/grokify/h5p-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/h5p-go/blob/master/LICENSE

A Go library for creating, manipulating, and validating H5P (HTML5 Package) content, with support for the official H5P file format and schemas.

## ✨ Features

- 📦 **Full H5P Package Support**: Create and extract `.h5p` ZIP files with proper structure
- 🔒 **Type-Safe Schema Implementation**: Official H5P MultiChoice schema with Go structs
- 🏗️ **Question Set Builder**: Fluent API for building interactive question sets
- ✅ **Validation**: Built-in validation for H5P compliance
- 🎯 **Multiple Question Types**: Support for single-answer and multi-answer questions
- 📋 **Official Schema Compliance**: Uses actual H5P semantics definitions
- 🔄 **JSON Serialization**: Full marshaling/unmarshaling support
- 🏷️ **Extensions Support**: Vendor-specific metadata with h5pGo namespace (v0.4.0+)

## 🚀 Quick Start

```bash
go get github.com/grokify/h5p-go
```

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Create a question set using the builder pattern
    builder := h5p.NewQuestionSetBuilder()
    
    answers := []h5p.Answer{
        h5p.CreateAnswer("Paris", true),
        h5p.CreateAnswer("London", false),
        h5p.CreateAnswer("Berlin", false),
        h5p.CreateAnswer("Madrid", false),
    }
    
    questionSet, err := builder.
        SetTitle("Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Welcome to our geography quiz!").
        AddMultipleChoiceQuestion("What is the capital of France?", answers).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Export to JSON
    jsonData, _ := questionSet.ToJSON()
    fmt.Printf("Generated H5P Question Set:\n%s\n", string(jsonData))
}
```

## 🆕 What's New in v0.4.0

Extensions support for vendor-specific question metadata! Add categorization, difficulty levels, learning objectives, and custom metadata to your questions while maintaining full H5P compatibility.

```go
ext := h5p.NewH5PGoExtension("1. Overview", 1).
    WithTopic("RAG Fundamentals").
    WithTags("rag", "retrieval").
    WithDifficulty("medium").
    WithLearningObjective("Understand RAG basics")

question := h5p.NewMultiChoiceQuestionWithExtensions(text, answers, &h5p.Extensions{H5PGo: ext})
```

See the [v0.4.0 release notes](releases/v0.4.0.md) for details.

## 📚 What's Next?

- **[Installation Guide](getting-started/installation.md)** - Detailed installation instructions
- **[Quick Start Tutorial](getting-started/quick-start.md)** - Step-by-step tutorial
- **[Basic Concepts](getting-started/concepts.md)** - Understanding H5P and this library
- **[User Guide](guide/question-sets.md)** - Comprehensive usage documentation
- **[API Reference](api/core-types.md)** - Complete API documentation
- **[Examples](examples/basic.md)** - Real-world usage examples
- **[Release Notes](releases/v0.4.0.md)** - Version history and changes

## 🏛️ Architecture

The H5P Go SDK is organized into several key components:

- **Core Types** (`questionset.go`, `h5p_package.go`) - Fundamental data structures
- **Extensions** (`extensions.go`) - Vendor-specific metadata support
- **Builder Pattern** (`builder.go`) - Fluent API for creating content
- **Official Schemas** (`schemas/`) - Type-safe H5P content type definitions
- **Semantics Engine** (`semantics/`) - Universal H5P semantics format support
- **Package Management** - Full `.h5p` file creation and extraction

## 🔧 Key Capabilities

### Question Sets
Build interactive question sets with multiple question types, scoring, and rich feedback systems.

### Type-Safe Schemas
Work with official H5P content type schemas using strongly-typed Go structs.

### H5P Package Management
Create, modify, and validate complete H5P packages ready for deployment.

### Standards Compliance
Full compliance with official H5P specifications and content type definitions.

## 🤝 Contributing

We welcome contributions! See our [Contributing Guide](development/contributing.md) for details.

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](https://github.com/grokify/h5p-go/blob/master/LICENSE) file for details.