# Quick Start

This tutorial will guide you through creating your first H5P content using the Go SDK.

## Your First Question Set

Let's create a simple geography quiz:

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Create a question set builder
    builder := h5p.NewQuestionSetBuilder()
    
    // Define answers for our question
    answers := []h5p.Answer{
        h5p.CreateAnswer("Paris", true),    // Correct answer
        h5p.CreateAnswer("London", false),  // Incorrect
        h5p.CreateAnswer("Berlin", false),  // Incorrect
        h5p.CreateAnswer("Madrid", false),  // Incorrect
    }
    
    // Build the question set
    questionSet, err := builder.
        SetTitle("Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Test your geography knowledge!").
        AddMultipleChoiceQuestion("What is the capital of France?", answers).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    // Convert to JSON
    jsonData, err := questionSet.ToJSON()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Created H5P Question Set:\n%s\n", string(jsonData))
}
```

## Adding Multiple Questions

Let's expand our quiz with more questions:

```go
func createGeographyQuiz() (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder()
    
    // Question 1: Capital of France
    answers1 := []h5p.Answer{
        h5p.CreateAnswer("Paris", true),
        h5p.CreateAnswer("London", false),
        h5p.CreateAnswer("Berlin", false),
        h5p.CreateAnswer("Madrid", false),
    }
    
    // Question 2: Largest country
    answers2 := []h5p.Answer{
        h5p.CreateAnswer("Russia", true),
        h5p.CreateAnswer("China", false),
        h5p.CreateAnswer("USA", false),
        h5p.CreateAnswer("Canada", false),
    }
    
    // Question 3: Longest river
    answers3 := []h5p.Answer{
        h5p.CreateAnswer("Nile", true),
        h5p.CreateAnswer("Amazon", false),
        h5p.CreateAnswer("Mississippi", false),
        h5p.CreateAnswer("Yangtze", false),
    }
    
    return builder.
        SetTitle("World Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(70).
        SetIntroduction("Test your knowledge of world geography!").
        AddMultipleChoiceQuestion("What is the capital of France?", answers1).
        AddMultipleChoiceQuestion("Which is the largest country by area?", answers2).
        AddMultipleChoiceQuestion("What is the longest river in the world?", answers3).
        Build()
}
```

## Adding Feedback

Enhance the user experience with feedback based on performance:

```go
func createQuizWithFeedback() (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder()
    
    // Create feedback ranges
    feedbackRanges := []h5p.OverallFeedback{
        {
            From: 0,
            To:   40,
            Feedback: "Keep studying! Geography takes practice.",
        },
        {
            From: 41,
            To:   70,
            Feedback: "Good job! You have solid geography knowledge.",
        },
        {
            From: 71,
            To:   100,
            Feedback: "Excellent! You're a geography expert!",
        },
    }
    
    // Build with feedback
    return builder.
        SetTitle("Geography Quiz with Feedback").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Test your geography knowledge!").
        // ... add questions ...
        AddOverallFeedback(feedbackRanges).
        Build()
}
```

## Validating Content

Always validate your content before deployment:

```go
func validateAndExport(questionSet *h5p.QuestionSet) error {
    // Validate the question set
    if err := questionSet.Validate(); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    
    // Export to JSON file
    jsonData, err := questionSet.ToJSON()
    if err != nil {
        return fmt.Errorf("JSON export failed: %w", err)
    }
    
    // Save to file
    err = os.WriteFile("quiz.json", jsonData, 0644)
    if err != nil {
        return fmt.Errorf("file write failed: %w", err)
    }
    
    fmt.Println("Quiz created and validated successfully!")
    return nil
}
```

## Complete Example

Here's a complete working example:

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Create quiz
    quiz, err := createCompleteQuiz()
    if err != nil {
        log.Fatal(err)
    }
    
    // Validate and export
    if err := validateAndExport(quiz); err != nil {
        log.Fatal(err)
    }
}

func createCompleteQuiz() (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder()
    
    // Geography questions with answers
    questions := []struct {
        question string
        answers  []h5p.Answer
    }{
        {
            "What is the capital of France?",
            []h5p.Answer{
                h5p.CreateAnswer("Paris", true),
                h5p.CreateAnswer("London", false),
                h5p.CreateAnswer("Berlin", false),
                h5p.CreateAnswer("Madrid", false),
            },
        },
        {
            "Which is the largest ocean?",
            []h5p.Answer{
                h5p.CreateAnswer("Pacific", true),
                h5p.CreateAnswer("Atlantic", false),
                h5p.CreateAnswer("Indian", false),
                h5p.CreateAnswer("Arctic", false),
            },
        },
    }
    
    // Start building
    builder = builder.
        SetTitle("Complete Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Welcome to our comprehensive geography quiz!")
    
    // Add all questions
    for _, q := range questions {
        builder = builder.AddMultipleChoiceQuestion(q.question, q.answers)
    }
    
    // Add feedback
    feedbackRanges := []h5p.OverallFeedback{
        {From: 0, To: 50, Feedback: "Keep studying geography!"},
        {From: 51, To: 80, Feedback: "Good knowledge of geography!"},
        {From: 81, To: 100, Feedback: "Excellent! You're a geography expert!"},
    }
    
    return builder.AddOverallFeedback(feedbackRanges).Build()
}

func validateAndExport(questionSet *h5p.QuestionSet) error {
    // Validate
    if err := questionSet.Validate(); err != nil {
        return err
    }
    
    // Export to JSON
    jsonData, err := questionSet.ToJSON()
    if err != nil {
        return err
    }
    
    // Save to file
    return os.WriteFile("complete-quiz.json", jsonData, 0644)
}
```

## What's Next?

Now that you've created your first H5P content, explore:

- [Basic Concepts](concepts.md) - Understand H5P fundamentals
- [Question Sets Guide](../guide/question-sets.md) - Advanced question set features
- [Typed Questions](../guide/typed-questions.md) - Using type-safe schemas
- [Validation](../guide/validation.md) - Ensuring content quality