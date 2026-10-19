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

## Extension Builder

`Extension` setters return the extension for chaining. See the
[Extensions guide](../guide/extensions.md) for what each field means.

```go
ext := h5p.NewExtension("1. Overview & Fundamentals", 1).
    WithTopic("RAG Fundamentals").
    WithTags("rag", "retrieval", "llm").
    WithDifficulty("medium").
    WithDifficultyScore(0.4).
    WithSkills(h5p.SkillRef{Taxonomy: "example-skills", ID: "rag.basics"}).
    WithLearningObjective("Understand the difference between RAG and fine-tuning").
    WithSource("Chapter 1")
```

## MultiChoiceQuestion Builder

Create typed MultiChoice questions with or without an extension block:

```go
params := &schemas.MultiChoiceParams{
    Question: "What is 2 + 2?",
    Answers: []schemas.AnswerOption{
        {Text: "4", Correct: true},
        {Text: "5", Correct: false},
    },
}

// Without extensions
question := h5p.NewMultiChoiceQuestion(params)

// With extensions
question = h5p.NewMultiChoiceQuestionWithExtensions(params,
    h5p.NewExtension("Math", 1).WithDifficulty("easy"))

// Add extensions to an existing question
question.WithExtensions(h5p.NewExtension("Math", 1))
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
