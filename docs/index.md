# H5P Go SDK

[![Build Status](https://github.com/grokify/h5p-go/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/grokify/h5p-go/actions/workflows/ci.yaml)
[![Lint Status](https://github.com/grokify/h5p-go/actions/workflows/lint.yaml/badge.svg?branch=main)](https://github.com/grokify/h5p-go/actions/workflows/lint.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/grokify/h5p-go)](https://goreportcard.com/report/github.com/grokify/h5p-go)
[![Docs](https://pkg.go.dev/badge/github.com/grokify/h5p-go)](https://pkg.go.dev/github.com/grokify/h5p-go)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/grokify/h5p-go/blob/master/LICENSE)

A Go library for creating, manipulating, and validating H5P (HTML5 Package) content, with support for the official H5P file format and schemas.

## ✨ Features

- 📦 **Full H5P Package Support**: Create and extract `.h5p` ZIP files with proper structure
- 🔒 **Type-Safe Schema Implementation**: Official H5P MultiChoice schema with Go structs
- 🏗️ **Question Set Builder**: Fluent API for building interactive question sets
- ✅ **Validation**: Built-in validation for H5P compliance
- 🎯 **Multiple Question Types**: Support for single-answer and multi-answer questions
- 📋 **Official Schema Compliance**: Uses actual H5P semantics definitions
- 🔄 **JSON Serialization**: Full marshaling/unmarshaling support

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

## 📚 What's Next?

- **[Installation Guide](getting-started/installation.md)** - Detailed installation instructions
- **[Quick Start Tutorial](getting-started/quick-start.md)** - Step-by-step tutorial 
- **[Basic Concepts](getting-started/concepts.md)** - Understanding H5P and this library
- **[User Guide](guide/question-sets.md)** - Comprehensive usage documentation
- **[API Reference](api/core-types.md)** - Complete API documentation
- **[Examples](examples/basic.md)** - Real-world usage examples

## 🏛️ Architecture

The H5P Go SDK is organized into several key components:

- **Core Types** (`questionset.go`, `h5p_package.go`) - Fundamental data structures
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