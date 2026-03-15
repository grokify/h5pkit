# Basic Concepts

Understanding the key concepts behind H5P and this Go SDK will help you create better interactive content.

## What is H5P?

H5P (HTML5 Package) is an open-source framework for creating interactive content such as quizzes, presentations, and games. H5P content is:

- **Reusable** - Content can be embedded in multiple platforms
- **Responsive** - Works on desktop, tablet, and mobile devices  
- **Standards-based** - Built on web standards (HTML5, CSS, JavaScript)
- **Accessible** - Designed with accessibility in mind

## Core Components

### Content Types

H5P organizes interactive elements into **content types**:

- **MultiChoice** - Single or multiple answer questions
- **TrueFalse** - Simple true/false questions
- **Essay** - Open-ended text responses
- **Question Set** - Collections of questions with scoring

### H5P Package Structure

An H5P package (`.h5p` file) is a ZIP archive containing:

```
package.h5p
├── h5p.json              # Package metadata
├── content/
│   └── content.json     # Content parameters
└── H5P.ContentType-1.0/ # Library files
    ├── library.json     # Library definition
    ├── semantics.json   # Content schema
    ├── js/             # JavaScript files
    └── css/            # Stylesheet files
```

### Semantics

**Semantics** define the structure and validation rules for content types. They specify:

- What fields are available for editing
- Field types (text, select, number, etc.)
- Validation rules and constraints
- UI behavior and layout

## SDK Architecture

### Builder Pattern

The SDK uses the builder pattern for creating content:

```go
questionSet, err := h5p.NewQuestionSetBuilder().
    SetTitle("My Quiz").
    SetPassPercentage(80).
    AddMultipleChoiceQuestion("Question?", answers).
    Build()
```

Benefits:
- **Fluent API** - Easy to read and write
- **Validation** - Catches errors during construction
- **Flexibility** - Optional parameters are simple to handle

### Type Safety

The SDK provides two levels of type safety:

1. **Generic Types** - Work with any content type
2. **Typed Schemas** - Strongly-typed for specific content types

```go
// Generic approach
question := h5p.Question{
    Library: "H5P.MultiChoice 1.16",
    Params:  map[string]interface{}{"question": "What is 2+2?"},
}

// Typed approach  
params := &schemas.MultiChoiceParams{
    Question: "What is 2+2?",
    Answers: []schemas.AnswerOption{
        {Text: "4", Correct: true},
    },
}
```

### Validation

The SDK includes comprehensive validation:

- **Structure validation** - Required fields, correct types
- **Business logic validation** - Sensible values, relationships
- **H5P compliance** - Matches official specifications

## Key Concepts

### Questions vs Question Sets

- **Question** - A single interactive element (e.g., one multiple choice question)
- **Question Set** - A collection of questions with overall scoring and feedback

### Parameters vs Content

- **Parameters** - The configuration data for content (JSON)
- **Content** - The complete H5P content including metadata and parameters

### Libraries vs Content Types

- **Library** - The code (JS/CSS) that renders content
- **Content Type** - A specific type of interactive content (MultiChoice, Essay, etc.)

## Working with Content

### Creating Content

1. **Define your content structure** (questions, answers, feedback)
2. **Use the builder** to construct your content
3. **Validate** the result
4. **Export** to JSON or H5P package

### Loading Content

1. **Parse** JSON or extract H5P package
2. **Validate** the loaded content
3. **Access** parameters and metadata
4. **Modify** as needed

### Best Practices

- Always validate content before deployment
- Use typed schemas when available for better safety
- Follow H5P naming conventions for compatibility
- Test content in H5P-compatible environments

## Next Steps

- [Question Sets Guide](../guide/question-sets.md) - Working with question collections
- [Typed Questions](../guide/typed-questions.md) - Using type-safe schemas
- [H5P Packages](../guide/h5p-packages.md) - Creating complete H5P files
- [API Reference](../api/core-types.md) - Detailed API documentation