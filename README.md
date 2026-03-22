# H5P Kit

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/h5pkit/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/h5pkit/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/h5pkit/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/h5pkit/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/h5pkit/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/h5pkit/actions/workflows/go-sast-codeql.yaml
 [goreport-svg]: https://goreportcard.com/badge/github.com/grokify/h5pkit
 [goreport-url]: https://goreportcard.com/report/github.com/grokify/h5pkit
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/h5pkit
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/h5pkit
 [viz-svg]: https://img.shields.io/badge/visualizaton-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fh5pkit
 [loc-svg]: https://tokei.rs/b1/github/grokify/h5pkit
 [repo-url]: https://github.com/grokify/h5pkit
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/h5pkit/blob/master/LICENSE

A toolkit for creating, editing, and validating H5P quiz content. Includes a **Go SDK** for server-side operations and a **TypeScript Editor** for browser-based quiz editing.

| Component | Language | Package |
|-----------|----------|---------|
| Go SDK | Go | `github.com/grokify/h5pkit` |
| Quiz Editor | TypeScript | `@grokify/h5p-editor` |

## Go SDK Features

- **Full H5P Package Support** - Create and extract `.h5p` ZIP files
- **Type-Safe Schema Implementation** - Official H5P content type schemas
- **Question Set Builder** - Fluent API for building interactive content
- **Validation** - Built-in H5P compliance validation
- **Multiple Question Types** - Support for various H5P content types
- **JSON Serialization** - Complete marshaling/unmarshaling support

## TypeScript Editor Features

- **Visual Quiz Editor** - React component for browser-based quiz creation
- **Multiple Question Types** - Multiple Choice, True/False, Fill in the Blanks
- **Undo/Redo** - Full history support for editing operations
- **Theming** - Light and dark mode with CSS custom properties

## Quick Start

### Go SDK

```bash
go get github.com/grokify/h5pkit
```

```go
package main

import (
    "fmt"
    "log"
    "github.com/grokify/h5pkit"
)

func main() {
    // Create answers
    answers := []h5p.Answer{
        h5p.CreateAnswer("Paris", true),
        h5p.CreateAnswer("London", false),
        h5p.CreateAnswer("Berlin", false),
    }

    // Build question set
    questionSet, err := h5p.NewQuestionSetBuilder().
        SetTitle("Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        AddMultipleChoiceQuestion("What is the capital of France?", answers).
        Build()

    if err != nil {
        log.Fatal(err)
    }

    // Export to JSON
    jsonData, _ := questionSet.ToJSON()
    fmt.Printf("Generated H5P content:\n%s\n", string(jsonData))
}
```

### TypeScript Editor

```bash
pnpm add @grokify/h5p-editor react react-dom
```

```tsx
import { QuizEditor, H5PQuiz } from '@grokify/h5p-editor';
import '@grokify/h5p-editor/styles.css';

function App() {
  const handleSave = (quiz: H5PQuiz) => {
    console.log('Quiz:', quiz);
  };

  return <QuizEditor onSave={handleSave} />;
}
```

## Documentation

**[Full Documentation](https://grokify.github.io/h5pkit/)**

### Quick Links

- **[Go SDK Installation](https://grokify.github.io/h5pkit/getting-started/installation/)** - Setup instructions
- **[TypeScript Editor](https://grokify.github.io/h5pkit/ts-editor/overview/)** - Browser-based editor
- **[API Reference](https://grokify.github.io/h5pkit/api/core-types/)** - Complete API documentation
- **[Examples](https://grokify.github.io/h5pkit/examples/basic/)** - Real-world usage examples

## Architecture

```
h5pkit/
├── schemas/          # Official H5P content type schemas
├── semantics/        # Universal H5P semantics format
├── ts/               # TypeScript editor (@grokify/h5p-editor)
├── builder.go        # Fluent API for content creation
├── questionset.go    # Core question set functionality
└── h5p_package.go    # Complete H5P package management
```

## Testing

```bash
go test ./...
```

## Standards Compliance

This library implements the official H5P specifications:

- [H5P File Format](https://h5p.org/documentation/developers/h5p-specification)
- [H5P Semantics](https://h5p.org/semantics)
- [H5P Content Types](https://github.com/h5p/)

## Contributing

See our [Contributing Guide](https://grokify.github.io/h5pkit/development/contributing/) for detailed information.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
