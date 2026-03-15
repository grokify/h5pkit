# Builder API

The Builder API provides a fluent interface for creating H5P content.

## QuestionSetBuilder

The `QuestionSetBuilder` provides a fluent API for constructing question sets.

### Creating a Builder

```go
builder := h5p.NewQuestionSetBuilder()
```

### Basic Methods

```go
// SetTitle sets the question set title
func (b *QuestionSetBuilder) SetTitle(title string) *QuestionSetBuilder

// SetProgressType sets progress display type ("textual" or "dots")
func (b *QuestionSetBuilder) SetProgressType(progressType string) *QuestionSetBuilder

// SetPassPercentage sets the passing threshold (0-100)
func (b *QuestionSetBuilder) SetPassPercentage(percentage int) *QuestionSetBuilder

// SetIntroduction sets the introduction text
func (b *QuestionSetBuilder) SetIntroduction(text string) *QuestionSetBuilder

// AddMultipleChoiceQuestion adds a multiple choice question
func (b *QuestionSetBuilder) AddMultipleChoiceQuestion(question string, answers []Answer) *QuestionSetBuilder

// AddOverallFeedback adds performance-based feedback
func (b *QuestionSetBuilder) AddOverallFeedback(feedback []OverallFeedback) *QuestionSetBuilder

// Build creates the final question set
func (b *QuestionSetBuilder) Build() (*QuestionSet, error)
```

### Example

```go
builder := h5p.NewQuestionSetBuilder()

answers := []h5p.Answer{
    h5p.CreateAnswer("Paris", true),
    h5p.CreateAnswer("London", false),
    h5p.CreateAnswer("Berlin", false),
}

questionSet, err := builder.
    SetTitle("Geography Quiz").
    SetProgressType("textual").
    SetPassPercentage(70).
    AddMultipleChoiceQuestion("What is the capital of France?", answers).
    Build()

if err != nil {
    log.Fatal(err)
}
```

## H5PGoExtension Builder

The `H5PGoExtension` type provides builder methods for constructing question metadata extensions.

### Creating an Extension

```go
// Create with section and question number
ext := h5p.NewH5PGoExtension("1. Overview", 1)

// Or create empty extensions container
extensions := h5p.NewExtensions()
```

### Builder Methods

All builder methods return the extension for method chaining:

```go
// WithTopic sets the topic within a section
func (e *H5PGoExtension) WithTopic(topic string) *H5PGoExtension

// WithTags sets categorization tags (variadic)
func (e *H5PGoExtension) WithTags(tags ...string) *H5PGoExtension

// WithDifficulty sets the difficulty level
func (e *H5PGoExtension) WithDifficulty(difficulty string) *H5PGoExtension

// WithLearningObjective sets what the question tests
func (e *H5PGoExtension) WithLearningObjective(objective string) *H5PGoExtension

// WithSource sets the origin of the question content
func (e *H5PGoExtension) WithSource(source string) *H5PGoExtension
```

### Complete Example

```go
// Build extension with all metadata
ext := h5p.NewH5PGoExtension("1. Overview & Fundamentals", 1).
    WithTopic("RAG Fundamentals").
    WithTags("rag", "retrieval", "llm").
    WithDifficulty("medium").
    WithLearningObjective("Understand the difference between RAG and fine-tuning").
    WithSource("PRESENTATION.md")

// Create question with extensions
question := h5p.NewMultiChoiceQuestionWithExtensions(
    "What does RAG stand for?",
    []h5p.Answer{
        {Text: "Retrieval-Augmented Generation", Correct: true},
        {Text: "Random Answer Generator", Correct: false},
        {Text: "Real-time AI Gateway", Correct: false},
    },
    &h5p.Extensions{H5PGo: ext},
)
```

## MultiChoiceQuestion Builder

Create MultiChoice questions with or without extensions:

```go
// Without extensions
question := h5p.NewMultiChoiceQuestion(
    "What is 2 + 2?",
    []h5p.Answer{
        {Text: "4", Correct: true},
        {Text: "5", Correct: false},
    },
)

// With extensions
question := h5p.NewMultiChoiceQuestionWithExtensions(
    "What is 2 + 2?",
    []h5p.Answer{
        {Text: "4", Correct: true},
        {Text: "5", Correct: false},
    },
    &h5p.Extensions{
        H5PGo: h5p.NewH5PGoExtension("Math", 1).
            WithDifficulty("easy"),
    },
)

// Add extensions to existing question
question.WithExtensions(&h5p.Extensions{
    H5PGo: h5p.NewH5PGoExtension("Math", 1),
})
```

## Error Handling

The builder pattern collects errors during construction:

```go
questionSet, err := builder.
    SetTitle("").  // Invalid: empty title
    SetPassPercentage(150).  // Invalid: > 100
    Build()

if err != nil {
    // Handle validation errors
    log.Printf("Build failed: %v", err)
}
```

## Best Practices

1. **Chain methods for readability**: Use method chaining to create readable construction code.

2. **Validate before use**: Always check the error returned by `Build()`.

3. **Use extensions for metadata**: Store categorization, difficulty, and learning objectives in extensions rather than custom params.

4. **Consistent difficulty values**: Use "easy", "medium", "hard" for interoperability.

5. **Descriptive learning objectives**: Write clear, testable learning objectives.
