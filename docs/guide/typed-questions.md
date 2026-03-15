# Typed Questions

Typed questions provide strongly-typed, schema-validated interfaces for creating H5P content. They ensure type safety and compliance with official H5P content type specifications.

## Overview

The SDK provides typed implementations for official H5P content types:
- **MultiChoice** - Single and multiple answer questions
- **Essay** - Open-ended text responses (future)
- **TrueFalse** - Simple true/false questions (future)

## MultiChoice Questions

MultiChoice is the most commonly used H5P content type for quizzes and assessments.

### Basic Usage

```go
import "github.com/grokify/h5p-go/schemas"

// Create parameters using typed structs
params := &schemas.MultiChoiceParams{
    Question: "What is the capital of France?",
    Answers: []schemas.AnswerOption{
        {
            Text:    "Paris",
            Correct: true,
            TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                ChosenFeedback: "Correct! Paris is the capital of France.",
            },
        },
        {
            Text:    "London", 
            Correct: false,
            TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                ChosenFeedback: "Incorrect. London is the capital of the UK.",
            },
        },
    },
}

// Create the typed question
question := h5p.NewMultiChoiceQuestion(params)
```

### Answer Configuration

#### Single Answer Questions

```go
params := &schemas.MultiChoiceParams{
    Question: "Which planet is closest to the Sun?",
    Answers: []schemas.AnswerOption{
        {Text: "Mercury", Correct: true},
        {Text: "Venus", Correct: false},
        {Text: "Earth", Correct: false},
        {Text: "Mars", Correct: false},
    },
    Behaviour: &schemas.Behaviour{
        Type: "single", // Single answer mode
    },
}
```

#### Multiple Answer Questions

```go
params := &schemas.MultiChoiceParams{
    Question: "Which of these are programming languages?",
    Answers: []schemas.AnswerOption{
        {Text: "Go", Correct: true},
        {Text: "Python", Correct: true},
        {Text: "HTML", Correct: false},
        {Text: "CSS", Correct: false},
    },
    Behaviour: &schemas.Behaviour{
        Type: "multi", // Multiple answer mode
    },
}
```

### Answer Feedback

Provide detailed feedback for each answer option:

```go
answers := []schemas.AnswerOption{
    {
        Text:    "Correct Answer",
        Correct: true,
        TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
            Tip:            "Think about the basic principles.",
            ChosenFeedback: "Excellent! You understood the concept perfectly.",
            NotChosenFeedback: "This was actually the correct answer.",
        },
    },
    {
        Text:    "Incorrect Answer",
        Correct: false,
        TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
            Tip:            "Consider the context carefully.",
            ChosenFeedback: "Not quite. Review the material on this topic.",
            NotChosenFeedback: "Good job avoiding this incorrect option.",
        },
    },
}
```

### Behavior Configuration

Control how the question behaves:

```go
behaviour := &schemas.Behaviour{
    Type:                  "single",        // "single" or "multi"
    EnableRetry:           true,            // Allow retrying
    EnableSolutionsButton: true,            // Show solutions
    EnableCheckButton:     true,            // Show check button
    RandomAnswers:         true,            // Randomize answer order
    ShowScorePoints:       true,            // Display points
    AutoCheck:             false,           // Auto-check answers
    PassPercentage:        75,              // Pass threshold
    ShowResultsOnQuizEnd:  true,            // Show results at end
}

params.Behaviour = behaviour
```

### Overall Feedback

Provide feedback based on score ranges:

```go
overallFeedback := []schemas.OverallFeedback{
    {
        From: 0,
        To:   25,
        Feedback: "You might want to review the material and try again.",
    },
    {
        From: 26,
        To:   75,
        Feedback: "Good effort! You're making progress.",
    },
    {
        From: 76,
        To:   100,
        Feedback: "Excellent work! You've mastered this topic.",
    },
}

params.OverallFeedback = overallFeedback
```

### Media Support

Add images, videos, or audio to questions:

```go
params := &schemas.MultiChoiceParams{
    Question: "What landmark is shown in the image?",
    Media: &schemas.MediaField{
        Type: "image",
        Copyright: &schemas.Copyright{
            License: "CC BY-SA",
            Author:  "Example Author",
        },
        // File would be included in H5P package
    },
    Answers: answers,
}
```

### UI Configuration

Customize the user interface:

```go
ui := &schemas.UI{
    CheckAnswerButton: "Check Answer",
    SubmitAnswerButton: "Submit",
    ShowSolutionButton: "Show Solution", 
    TryAgainButton: "Try Again",
    ScoreBarLabel: "Progress",
    A11yCheck: "Check the answers",
    A11yShowSolution: "Show the solution",
    A11yRetry: "Retry the task",
}

params.UI = ui
```

## Validation

Typed questions include comprehensive validation:

```go
params := &schemas.MultiChoiceParams{
    Question: "What is 2 + 2?",
    Answers: []schemas.AnswerOption{
        {Text: "4", Correct: true},
        {Text: "5", Correct: false},
    },
}

// Validate the parameters
if err := params.Validate(); err != nil {
    log.Fatal("Validation failed:", err)
}

// Create the question (also validates)
question := h5p.NewMultiChoiceQuestion(params)
```

Common validation rules:
- Question text is required
- At least 2 answers required
- At least 1 correct answer required
- Pass percentage must be 0-100
- Feedback ranges must not overlap

## Converting to Generic Questions

Convert typed questions to generic format for question sets:

```go
// Create typed question
typedQuestion := h5p.NewMultiChoiceQuestion(params)

// Convert to generic question
genericQuestion := typedQuestion.ToQuestion()

// Add to question set
builder := h5p.NewQuestionSetBuilder()
builder.AddQuestion(genericQuestion)
```

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/grokify/h5p-go"
    "github.com/grokify/h5p-go/schemas"
)

func main() {
    // Create a comprehensive multiple choice question
    question := createAdvancedQuestion()
    
    // Convert to generic format
    genericQuestion := question.ToQuestion()
    
    // Add to a question set
    questionSet, err := h5p.NewQuestionSetBuilder().
        SetTitle("Advanced Programming Quiz").
        SetProgressType("textual").
        SetPassPercentage(70).
        AddQuestion(genericQuestion).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Export and validate
    if err := questionSet.Validate(); err != nil {
        log.Fatal("Validation failed:", err)
    }
    
    jsonData, _ := questionSet.ToJSON()
    fmt.Printf("Created advanced question:\n%s\n", string(jsonData))
}

func createAdvancedQuestion() *h5p.MultiChoiceQuestion {
    params := &schemas.MultiChoiceParams{
        Question: "Which of the following are characteristics of Go programming language?",
        
        Answers: []schemas.AnswerOption{
            {
                Text:    "Statically typed",
                Correct: true,
                TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                    Tip:            "Think about compile-time vs runtime type checking.",
                    ChosenFeedback: "Correct! Go uses static typing with type inference.",
                    NotChosenFeedback: "This is actually a key feature of Go.",
                },
            },
            {
                Text:    "Compiled language",
                Correct: true,
                TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                    Tip:            "Consider how Go code is executed.",
                    ChosenFeedback: "Right! Go compiles to native machine code.",
                    NotChosenFeedback: "Go is indeed compiled, not interpreted.",
                },
            },
            {
                Text:    "Object-oriented with classes",
                Correct: false,
                TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                    Tip:            "Think about Go's approach to OOP.",
                    ChosenFeedback: "Not quite. Go doesn't have classes, but uses structs and interfaces.",
                    NotChosenFeedback: "Good! Go doesn't use traditional class-based OOP.",
                },
            },
            {
                Text:    "Has built-in concurrency support",
                Correct: true,
                TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
                    Tip:            "Consider Go's approach to concurrent programming.",
                    ChosenFeedback: "Excellent! Goroutines and channels are core features.",
                    NotChosenFeedback: "This is actually one of Go's strongest features.",
                },
            },
        },
        
        Behaviour: &schemas.Behaviour{
            Type:                  "multi",
            EnableRetry:           true,
            EnableSolutionsButton: true,
            RandomAnswers:         true,
            ShowScorePoints:       true,
            PassPercentage:        75,
        },
        
        OverallFeedback: []schemas.OverallFeedback{
            {
                From:     0,
                To:       40,
                Feedback: "Review the fundamentals of Go programming.",
            },
            {
                From:     41,
                To:       80,
                Feedback: "Good understanding! Keep learning about Go's unique features.",
            },
            {
                From:     81,
                To:       100,
                Feedback: "Excellent! You have a solid grasp of Go's characteristics.",
            },
        },
        
        UI: &schemas.UI{
            CheckAnswerButton:  "Check My Answers",
            SubmitAnswerButton: "Submit Answers",
            TryAgainButton:     "Try Again",
            ShowSolutionButton: "Show Solutions",
            ScoreBarLabel:      "Your Progress",
        },
    }
    
    return h5p.NewMultiChoiceQuestion(params)
}
```

## Best Practices

1. **Use Validation** - Always validate parameters before creating questions
2. **Meaningful Feedback** - Provide helpful tips and explanations
3. **Clear Questions** - Write unambiguous question text
4. **Balanced Options** - Include plausible distractors
5. **Accessibility** - Use proper UI labels for screen readers
6. **Test Thoroughly** - Validate both individual questions and complete question sets

## Schema Compliance

Typed questions ensure compliance with official H5P schemas:
- Field validation matches H5P semantics
- Required fields are enforced
- Value ranges are checked
- Relationships between fields are validated

## Next Steps

- [H5P Packages](h5p-packages.md) - Creating complete H5P files
- [Validation](validation.md) - Understanding validation rules
- [API Reference](../api/schemas.md) - Detailed schema documentation