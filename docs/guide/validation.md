# Validation

The H5P Go SDK includes comprehensive validation to ensure your content meets H5P standards and functions correctly across platforms.

## Overview

Validation occurs at multiple levels:
- **Parameter validation** - Field types, required values, constraints
- **Structure validation** - Proper relationships between components
- **Schema compliance** - Conformance with official H5P specifications
- **Business logic validation** - Sensible values and configurations

## Validation Levels

### Content Validation

```go
// Validate question set parameters
questionSet, err := builder.Build()
if err != nil {
    // Build-time validation failed
    log.Fatal("Build failed:", err)
}

// Runtime validation
if err := questionSet.Validate(); err != nil {
    // Content validation failed
    log.Fatal("Validation failed:", err)
}
```

### Schema Validation

```go
import "github.com/grokify/h5p-go/schemas"

// Validate typed parameters
params := &schemas.MultiChoiceParams{
    Question: "What is 2+2?",
    Answers: []schemas.AnswerOption{
        {Text: "4", Correct: true},
    },
}

if err := params.Validate(); err != nil {
    log.Fatal("Schema validation failed:", err)
}
```

### Package Validation

```go
// Validate complete H5P package
pkg := h5p.NewH5PPackage()
// ... configure package ...

if err := pkg.Validate(); err != nil {
    log.Fatal("Package validation failed:", err)
}
```

## Common Validation Rules

### Question Sets

**Required Fields:**
- Title must not be empty
- At least one question required
- Pass percentage between 0-100

**Business Logic:**
- Pass percentage should be reasonable (typically 50-90)
- Questions should have meaningful text
- Feedback ranges should not overlap

```go
// Valid question set
questionSet, err := h5p.NewQuestionSetBuilder().
    SetTitle("My Quiz").              // Required
    SetPassPercentage(70).           // 0-100
    AddMultipleChoiceQuestion(       // At least 1 question
        "Question text?", 
        answers,
    ).
    Build()
```

### Multiple Choice Questions

**Required Fields:**
- Question text must not be empty
- At least 2 answers required
- At least 1 correct answer required

**Constraints:**
- Maximum answer length limits
- Feedback text length limits
- Pass percentage 0-100

```go
// Valid multiple choice
params := &schemas.MultiChoiceParams{
    Question: "What is the capital?",  // Required, non-empty
    Answers: []schemas.AnswerOption{   // Min 2 answers
        {Text: "Paris", Correct: true},     // At least 1 correct
        {Text: "London", Correct: false},   // Text required
    },
}

if err := params.Validate(); err != nil {
    // Validation failed
}
```

### Answer Options

**Rules:**
- Text field is required and non-empty
- At least one answer must be marked correct
- Feedback text has length limits

```go
// Valid answer
answer := schemas.AnswerOption{
    Text:    "Correct Answer",        // Required
    Correct: true,                    // At least 1 must be true
    TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
        ChosenFeedback: "Good job!",  // Optional, but limited length
    },
}
```

## Error Types

### Validation Errors

The SDK provides structured validation errors:

```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}
```

### Multiple Errors

Validation can return multiple errors at once:

```go
type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
    messages := make([]string, len(e))
    for i, err := range e {
        messages[i] = err.Error()
    }
    return strings.Join(messages, "; ")
}
```

### Error Handling

```go
questionSet, err := builder.Build()
if err != nil {
    if validationErrors, ok := err.(h5p.ValidationErrors); ok {
        fmt.Println("Multiple validation errors:")
        for _, validationError := range validationErrors {
            fmt.Printf("- %s: %s\n", validationError.Field, validationError.Message)
        }
    } else {
        fmt.Printf("Build error: %s\n", err)
    }
    return
}
```

## Custom Validation

### Adding Custom Rules

You can implement custom validation for specific use cases:

```go
func validateQuizComplexity(questionSet *h5p.QuestionSet) error {
    if len(questionSet.Questions) < 3 {
        return fmt.Errorf("quiz should have at least 3 questions for proper assessment")
    }
    
    // Check for variety in answer counts
    answerCounts := make(map[int]int)
    for _, question := range questionSet.Questions {
        if params, ok := question.Params.(map[string]interface{}); ok {
            if answers, exists := params["answers"]; exists {
                if answerList, ok := answers.([]interface{}); ok {
                    answerCounts[len(answerList)]++
                }
            }
        }
    }
    
    if len(answerCounts) == 1 {
        return fmt.Errorf("consider varying the number of answers for better engagement")
    }
    
    return nil
}

// Use custom validation
if err := validateQuizComplexity(questionSet); err != nil {
    log.Printf("Warning: %s", err)
}
```

## Validation Best Practices

### 1. Validate Early and Often

```go
// Validate at build time
questionSet, err := builder.Build()
if err != nil {
    return err  // Fail fast
}

// Validate before export
if err := questionSet.Validate(); err != nil {
    return err
}

// Validate before package creation
if err := pkg.Validate(); err != nil {
    return err
}
```

### 2. Handle Validation Gracefully

```go
func buildQuestionSet(title string, questions []QuestionData) (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder().SetTitle(title)
    
    for i, q := range questions {
        if q.Text == "" {
            return nil, fmt.Errorf("question %d has empty text", i+1)
        }
        
        if len(q.Answers) < 2 {
            return nil, fmt.Errorf("question %d needs at least 2 answers", i+1)
        }
        
        // Convert and add question
        answers := convertAnswers(q.Answers)
        builder = builder.AddMultipleChoiceQuestion(q.Text, answers)
    }
    
    return builder.Build()
}
```

### 3. Provide Meaningful Error Messages

```go
func validatePassPercentage(percentage int) error {
    if percentage < 0 || percentage > 100 {
        return fmt.Errorf("pass percentage must be between 0 and 100, got %d", percentage)
    }
    
    if percentage < 30 {
        return fmt.Errorf("pass percentage %d seems too low - consider 50-80 range", percentage)
    }
    
    if percentage > 95 {
        return fmt.Errorf("pass percentage %d seems too high - consider 60-90 range", percentage)
    }
    
    return nil
}
```

## Testing Validation

### Unit Testing

```go
func TestQuestionSetValidation(t *testing.T) {
    tests := []struct {
        name      string
        builder   func() *h5p.QuestionSetBuilder
        wantError bool
        errorText string
    }{
        {
            name: "valid question set",
            builder: func() *h5p.QuestionSetBuilder {
                return h5p.NewQuestionSetBuilder().
                    SetTitle("Test Quiz").
                    SetPassPercentage(70).
                    AddMultipleChoiceQuestion("Question?", validAnswers)
            },
            wantError: false,
        },
        {
            name: "missing title",
            builder: func() *h5p.QuestionSetBuilder {
                return h5p.NewQuestionSetBuilder().
                    SetPassPercentage(70).
                    AddMultipleChoiceQuestion("Question?", validAnswers)
            },
            wantError: true,
            errorText: "title is required",
        },
        {
            name: "invalid pass percentage",
            builder: func() *h5p.QuestionSetBuilder {
                return h5p.NewQuestionSetBuilder().
                    SetTitle("Test Quiz").
                    SetPassPercentage(150).  // Invalid
                    AddMultipleChoiceQuestion("Question?", validAnswers)
            },
            wantError: true,
            errorText: "pass percentage must be between 0 and 100",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := tt.builder().Build()
            
            if tt.wantError {
                if err == nil {
                    t.Errorf("expected error but got none")
                    return
                }
                if !strings.Contains(err.Error(), tt.errorText) {
                    t.Errorf("error %q should contain %q", err.Error(), tt.errorText)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

### Integration Testing

```go
func TestCompleteWorkflow(t *testing.T) {
    // Test the complete workflow with validation
    questionSet, err := createTestQuestionSet()
    if err != nil {
        t.Fatalf("Failed to create question set: %v", err)
    }
    
    // Validate content
    if err := questionSet.Validate(); err != nil {
        t.Fatalf("Question set validation failed: %v", err)
    }
    
    // Create package
    pkg := h5p.NewH5PPackage()
    pkg.SetContent(&h5p.Content{Params: questionSet})
    
    // Validate package
    if err := pkg.Validate(); err != nil {
        t.Fatalf("Package validation failed: %v", err)
    }
    
    // Export should succeed
    tempFile := filepath.Join(t.TempDir(), "test.h5p")
    if err := pkg.CreateZipFile(tempFile); err != nil {
        t.Fatalf("Failed to export package: %v", err)
    }
    
    // Re-load and validate
    loadedPkg, err := h5p.LoadH5PPackage(tempFile)
    if err != nil {
        t.Fatalf("Failed to load package: %v", err)
    }
    
    if err := loadedPkg.Validate(); err != nil {
        t.Fatalf("Loaded package validation failed: %v", err)
    }
}
```

## Validation Checklist

Before deploying H5P content, ensure:

- [ ] All required fields are populated
- [ ] Field values are within acceptable ranges  
- [ ] Questions have meaningful text
- [ ] At least one correct answer per question
- [ ] Pass percentage is reasonable (50-90)
- [ ] Feedback messages are helpful and appropriate
- [ ] Package includes all required libraries
- [ ] Library versions are compatible
- [ ] No validation errors or warnings

## Next Steps

- [Examples](../examples/basic.md) - See validation in practice
- [Testing](../development/testing.md) - Learn about testing strategies
- [Standards Compliance](../reference/standards.md) - Official H5P specifications