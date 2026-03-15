# Question Sets

Question sets are collections of interactive questions with overall scoring and feedback. They're perfect for quizzes, assessments, and educational content.

## Overview

A question set consists of:
- Multiple questions of various types
- Overall scoring and pass percentage
- Introduction and conclusion text
- Feedback based on performance ranges
- Navigation and progress indicators

## Creating Question Sets

### Using the Builder

The recommended way to create question sets is using the fluent builder API:

```go
builder := h5p.NewQuestionSetBuilder()

questionSet, err := builder.
    SetTitle("My Quiz").
    SetProgressType("textual").
    SetPassPercentage(70).
    SetIntroduction("Welcome to the quiz!").
    AddMultipleChoiceQuestion("Question 1?", answers1).
    AddMultipleChoiceQuestion("Question 2?", answers2).
    Build()
```

### Configuration Options

#### Title and Metadata
```go
builder.
    SetTitle("Advanced Geography Quiz").
    SetProgressType("textual").          // "textual" or "dots"
    SetPassPercentage(80).               // 0-100
    SetIntroduction("Test your knowledge!").
    SetStartButtonText("Begin Quiz")
```

#### Question Management
```go
// Add multiple choice questions
builder.AddMultipleChoiceQuestion("What is the capital?", answers)

// Questions are added in order and will be presented sequentially
```

#### Feedback Configuration
```go
feedbackRanges := []h5p.OverallFeedback{
    {From: 0, To: 40, Feedback: "Keep studying!"},
    {From: 41, To: 70, Feedback: "Good job!"},
    {From: 71, To: 100, Feedback: "Excellent work!"},
}

builder.AddOverallFeedback(feedbackRanges)
```

## Question Types

### Multiple Choice Questions

Support both single-answer and multi-answer questions:

```go
// Single correct answer
answers := []h5p.Answer{
    h5p.CreateAnswer("Correct answer", true),
    h5p.CreateAnswer("Wrong answer 1", false),
    h5p.CreateAnswer("Wrong answer 2", false),
}

// Multiple correct answers  
answers := []h5p.Answer{
    h5p.CreateAnswer("Correct 1", true),
    h5p.CreateAnswer("Correct 2", true),
    h5p.CreateAnswer("Wrong", false),
}

builder.AddMultipleChoiceQuestion("Question text", answers)
```

### Answer Options

Each answer can include additional metadata:

```go
type Answer struct {
    Text     string  `json:"text"`
    Correct  bool    `json:"correct"`
    Feedback string  `json:"feedback,omitempty"`
    Weight   float64 `json:"weight,omitempty"`
}

answer := h5p.Answer{
    Text:     "Paris",
    Correct:  true,
    Feedback: "Correct! Paris is the capital of France.",
    Weight:   1.0,
}
```

## Scoring and Feedback

### Pass Percentage

Set the minimum score required to pass:

```go
builder.SetPassPercentage(75) // 75% required to pass
```

### Overall Feedback

Provide different feedback messages based on score ranges:

```go
feedback := []h5p.OverallFeedback{
    {
        From:     0,
        To:       30,
        Feedback: "You might want to review the material and try again.",
    },
    {
        From:     31,
        To:       60,
        Feedback: "Good effort! You're getting there.",
    },
    {
        From:     61,
        To:       100,
        Feedback: "Excellent work! You've mastered this topic.",
    },
}

builder.AddOverallFeedback(feedback)
```

## Progress Indicators

Choose how progress is displayed to users:

```go
// Textual progress (e.g., "Question 1 of 5")
builder.SetProgressType("textual")

// Dots progress indicator
builder.SetProgressType("dots")
```

## Advanced Features

### Custom Start Button

```go
builder.SetStartButtonText("Start the Challenge!")
```

### Randomization

Currently, question order is determined by the order they're added to the builder. Future versions may support randomization.

## Validation

Always validate your question sets before deployment:

```go
questionSet, err := builder.Build()
if err != nil {
    log.Fatal("Build failed:", err)
}

// Validate the complete question set
if err := questionSet.Validate(); err != nil {
    log.Fatal("Validation failed:", err)
}
```

Common validation errors:
- Missing required fields (title, questions)
- Invalid pass percentage (not 0-100)
- Empty feedback ranges
- Questions without answers

## Export and Usage

### JSON Export
```go
jsonData, err := questionSet.ToJSON()
if err != nil {
    log.Fatal(err)
}

// Save to file
os.WriteFile("quiz.json", jsonData, 0644)
```

### Loading from JSON
```go
jsonData, err := os.ReadFile("quiz.json") 
if err != nil {
    log.Fatal(err)
}

questionSet, err := h5p.FromJSON(jsonData)
if err != nil {
    log.Fatal(err)
}
```

## Best Practices

1. **Clear Questions** - Write unambiguous questions
2. **Balanced Difficulty** - Mix easy and challenging questions
3. **Meaningful Feedback** - Provide helpful feedback for both correct and incorrect answers
4. **Appropriate Length** - Keep quizzes focused (5-15 questions typically)
5. **Test Thoroughly** - Validate and test your content before deployment
6. **Progressive Difficulty** - Start with easier questions to build confidence

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Create a comprehensive geography quiz
    quiz, err := createGeographyQuiz()
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate before export
    if err := quiz.Validate(); err != nil {
        log.Fatal("Validation failed:", err)
    }
    
    // Export to JSON
    jsonData, err := quiz.ToJSON()
    if err != nil {
        log.Fatal(err)
    }
    
    // Save to file
    err = os.WriteFile("geography-quiz.json", jsonData, 0644)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Geography quiz created successfully!")
}

func createGeographyQuiz() (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder()
    
    // Question 1: Capitals
    capitals := []h5p.Answer{
        {Text: "Paris", Correct: true, Feedback: "Correct!"},
        {Text: "London", Correct: false, Feedback: "London is the capital of the UK."},
        {Text: "Berlin", Correct: false, Feedback: "Berlin is the capital of Germany."},
        {Text: "Madrid", Correct: false, Feedback: "Madrid is the capital of Spain."},
    }
    
    // Question 2: Geography
    continents := []h5p.Answer{
        {Text: "Asia", Correct: true, Feedback: "Correct! Asia is the largest continent."},
        {Text: "Africa", Correct: false, Feedback: "Africa is the second largest."},
        {Text: "North America", Correct: false, Feedback: "North America is the third largest."},
        {Text: "Europe", Correct: false, Feedback: "Europe is smaller than Asia."},
    }
    
    // Feedback ranges
    feedback := []h5p.OverallFeedback{
        {From: 0, To: 50, Feedback: "Keep studying geography!"},
        {From: 51, To: 80, Feedback: "Good knowledge of world geography!"},
        {From: 81, To: 100, Feedback: "Excellent! You're a geography expert!"},
    }
    
    return builder.
        SetTitle("World Geography Challenge").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Test your knowledge of world geography with this challenging quiz!").
        SetStartButtonText("Start Challenge").
        AddMultipleChoiceQuestion("What is the capital of France?", capitals).
        AddMultipleChoiceQuestion("Which is the largest continent by area?", continents).
        AddOverallFeedback(feedback).
        Build()
}
```

## Next Steps

- [Typed Questions](typed-questions.md) - Using strongly-typed content schemas
- [H5P Packages](h5p-packages.md) - Creating complete H5P files  
- [Validation](validation.md) - Ensuring content quality