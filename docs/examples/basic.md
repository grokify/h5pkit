# Basic Examples

This page provides practical examples for common H5P Go SDK usage patterns.

## Simple Quiz Creation

### Basic Multiple Choice Quiz

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Create a simple geography quiz
    quiz := createGeographyQuiz()
    
    // Export to JSON file
    jsonData, err := quiz.ToJSON()
    if err != nil {
        log.Fatal("JSON export failed:", err)
    }
    
    err = os.WriteFile("geography-quiz.json", jsonData, 0644)
    if err != nil {
        log.Fatal("File write failed:", err)
    }
    
    fmt.Println("Geography quiz created successfully!")
    fmt.Printf("Title: %s\n", quiz.Title)
    fmt.Printf("Questions: %d\n", len(quiz.Questions))
    fmt.Printf("Pass percentage: %d%%\n", quiz.PassPercentage)
}

func createGeographyQuiz() *h5p.QuestionSet {
    builder := h5p.NewQuestionSetBuilder()
    
    // Question 1: Capital cities
    capitalAnswers := []h5p.Answer{
        h5p.CreateAnswerWithFeedback("Paris", true, "Correct! Paris is the capital of France."),
        h5p.CreateAnswerWithFeedback("London", false, "London is the capital of the United Kingdom."),
        h5p.CreateAnswerWithFeedback("Berlin", false, "Berlin is the capital of Germany."),
        h5p.CreateAnswerWithFeedback("Madrid", false, "Madrid is the capital of Spain."),
    }
    
    // Question 2: Continents
    continentAnswers := []h5p.Answer{
        h5p.CreateAnswerWithFeedback("Asia", true, "Correct! Asia is the largest continent."),
        h5p.CreateAnswerWithFeedback("Africa", false, "Africa is the second largest continent."),
        h5p.CreateAnswerWithFeedback("North America", false, "North America is the third largest."),
        h5p.CreateAnswerWithFeedback("Europe", false, "Europe is much smaller than Asia."),
    }
    
    // Question 3: Oceans
    oceanAnswers := []h5p.Answer{
        h5p.CreateAnswerWithFeedback("Pacific Ocean", true, "Correct! The Pacific is the largest ocean."),
        h5p.CreateAnswerWithFeedback("Atlantic Ocean", false, "The Atlantic is the second largest."),
        h5p.CreateAnswerWithFeedback("Indian Ocean", false, "The Indian Ocean is the third largest."),
        h5p.CreateAnswerWithFeedback("Arctic Ocean", false, "The Arctic Ocean is the smallest."),
    }
    
    // Overall feedback based on performance
    feedback := []h5p.OverallFeedback{
        {
            From:     0,
            To:       40,
            Feedback: "Keep studying geography! There's so much to learn about our world.",
        },
        {
            From:     41,
            To:       70,
            Feedback: "Good job! You have solid basic geography knowledge.",
        },
        {
            From:     71,
            To:       100,
            Feedback: "Excellent! You're a geography expert!",
        },
    }
    
    // Build the complete quiz
    questionSet, err := builder.
        SetTitle("Basic Geography Quiz").
        SetProgressType("textual").
        SetPassPercentage(60).
        SetIntroduction("Test your knowledge of world geography with this fun quiz!").
        SetStartButtonText("Start Quiz").
        AddMultipleChoiceQuestion("What is the capital of France?", capitalAnswers).
        AddMultipleChoiceQuestion("Which is the largest continent by area?", continentAnswers).
        AddMultipleChoiceQuestion("What is the largest ocean on Earth?", oceanAnswers).
        AddOverallFeedback(feedback).
        Build()
    
    if err != nil {
        log.Fatal("Failed to build quiz:", err)
    }
    
    return questionSet
}
```

## Loading and Modifying Existing Content

### Load and Update Quiz

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/grokify/h5p-go"
)

func main() {
    // Load existing quiz from JSON file
    quiz, err := loadQuizFromFile("existing-quiz.json")
    if err != nil {
        log.Fatal("Failed to load quiz:", err)
    }
    
    fmt.Printf("Loaded quiz: %s\n", quiz.Title)
    fmt.Printf("Original questions: %d\n", len(quiz.Questions))
    
    // Add a new question
    newAnswers := []h5p.Answer{
        h5p.CreateAnswer("7", false),
        h5p.CreateAnswer("8", true),
        h5p.CreateAnswer("9", false),
        h5p.CreateAnswer("10", false),
    }
    
    // Create new question
    newQuestion := h5p.Question{
        Library: "H5P.MultiChoice 1.16",
        Params: map[string]interface{}{
            "question": "What is 4 + 4?",
            "answers":  convertAnswersToParams(newAnswers),
        },
    }
    
    // Add to existing quiz
    quiz.Questions = append(quiz.Questions, newQuestion)
    
    // Update title
    quiz.Title = "Updated " + quiz.Title
    
    // Validate modified quiz
    if err := quiz.Validate(); err != nil {
        log.Fatal("Modified quiz validation failed:", err)
    }
    
    // Save updated quiz
    jsonData, err := quiz.ToJSON()
    if err != nil {
        log.Fatal("JSON export failed:", err)
    }
    
    err = os.WriteFile("updated-quiz.json", jsonData, 0644)
    if err != nil {
        log.Fatal("File write failed:", err)
    }
    
    fmt.Printf("Updated quiz saved with %d questions\n", len(quiz.Questions))
}

func loadQuizFromFile(filename string) (*h5p.QuestionSet, error) {
    jsonData, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    quiz, err := h5p.FromJSON(jsonData)
    if err != nil {
        return nil, err
    }
    
    // Validate loaded content
    if err := quiz.Validate(); err != nil {
        return nil, fmt.Errorf("loaded quiz is invalid: %w", err)
    }
    
    return quiz, nil
}

func convertAnswersToParams(answers []h5p.Answer) []map[string]interface{} {
    params := make([]map[string]interface{}, len(answers))
    for i, answer := range answers {
        params[i] = map[string]interface{}{
            "text":    answer.Text,
            "correct": answer.Correct,
        }
        if answer.Feedback != "" {
            params[i]["feedback"] = answer.Feedback
        }
    }
    return params
}
```

## Working with Different Question Types

### True/False Questions

```go
func createTrueFalseQuiz() *h5p.QuestionSet {
    builder := h5p.NewQuestionSetBuilder()
    
    // True/False questions have exactly 2 answers
    question1Answers := []h5p.Answer{
        h5p.CreateAnswerWithFeedback("True", true, "Correct! The Earth is indeed round."),
        h5p.CreateAnswerWithFeedback("False", false, "Incorrect. The Earth is approximately spherical."),
    }
    
    question2Answers := []h5p.Answer{
        h5p.CreateAnswerWithFeedback("True", false, "Incorrect. Water boils at 100°C, not freezes."),
        h5p.CreateAnswerWithFeedback("False", true, "Correct! Water freezes at 0°C."),
    }
    
    questionSet, err := builder.
        SetTitle("True or False Science Quiz").
        SetProgressType("dots").
        SetPassPercentage(50).
        SetIntroduction("Test your science knowledge with these true/false questions.").
        AddMultipleChoiceQuestion("The Earth is round.", question1Answers).
        AddMultipleChoiceQuestion("Water freezes at 100 degrees Celsius.", question2Answers).
        Build()
    
    if err != nil {
        log.Fatal("Failed to build true/false quiz:", err)
    }
    
    return questionSet
}
```

### Multi-Select Questions

```go
func createMultiSelectQuiz() *h5p.QuestionSet {
    builder := h5p.NewQuestionSetBuilder()
    
    // Multiple correct answers
    programmingAnswers := []h5p.Answer{
        h5p.CreateAnswer("Go", true),        // Correct
        h5p.CreateAnswer("Python", true),    // Correct  
        h5p.CreateAnswer("Java", true),      // Correct
        h5p.CreateAnswer("HTML", false),     // Not a programming language
        h5p.CreateAnswer("CSS", false),      // Not a programming language
    }
    
    // Single correct answer for comparison
    capitalAnswers := []h5p.Answer{
        h5p.CreateAnswer("Tokyo", true),
        h5p.CreateAnswer("Osaka", false),
        h5p.CreateAnswer("Kyoto", false),
    }
    
    questionSet, err := builder.
        SetTitle("Mixed Question Types").
        SetProgressType("textual").
        SetPassPercentage(70).
        SetIntroduction("This quiz has both single-answer and multiple-answer questions.").
        AddMultipleChoiceQuestion("Which of these are programming languages? (Select all that apply)", programmingAnswers).
        AddMultipleChoiceQuestion("What is the capital of Japan?", capitalAnswers).
        Build()
    
    if err != nil {
        log.Fatal("Failed to build multi-select quiz:", err)
    }
    
    return questionSet
}
```

## Error Handling Patterns

### Comprehensive Error Handling

```go
func createQuizWithErrorHandling() (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder()
    
    // Validate inputs before building
    title := "Programming Quiz"
    if title == "" {
        return nil, fmt.Errorf("quiz title cannot be empty")
    }
    
    passPercentage := 75
    if passPercentage < 0 || passPercentage > 100 {
        return nil, fmt.Errorf("pass percentage must be between 0 and 100, got %d", passPercentage)
    }
    
    // Create answers with validation
    answers, err := createValidatedAnswers()
    if err != nil {
        return nil, fmt.Errorf("failed to create answers: %w", err)
    }
    
    // Build with error handling
    questionSet, err := builder.
        SetTitle(title).
        SetProgressType("textual").
        SetPassPercentage(passPercentage).
        SetIntroduction("Welcome to the programming quiz!").
        AddMultipleChoiceQuestion("What is Go?", answers).
        Build()
    
    if err != nil {
        return nil, fmt.Errorf("failed to build question set: %w", err)
    }
    
    // Additional validation
    if len(questionSet.Questions) == 0 {
        return nil, fmt.Errorf("question set must have at least one question")
    }
    
    // Final validation
    if err := questionSet.Validate(); err != nil {
        return nil, fmt.Errorf("question set validation failed: %w", err)
    }
    
    return questionSet, nil
}

func createValidatedAnswers() ([]h5p.Answer, error) {
    answers := []h5p.Answer{
        {Text: "A programming language", Correct: true},
        {Text: "A database", Correct: false},
        {Text: "An operating system", Correct: false},
        {Text: "A web browser", Correct: false},
    }
    
    // Validate answers
    if len(answers) < 2 {
        return nil, fmt.Errorf("need at least 2 answers")
    }
    
    hasCorrect := false
    for i, answer := range answers {
        if answer.Text == "" {
            return nil, fmt.Errorf("answer %d has empty text", i+1)
        }
        if answer.Correct {
            hasCorrect = true
        }
    }
    
    if !hasCorrect {
        return nil, fmt.Errorf("at least one answer must be correct")
    }
    
    return answers, nil
}
```

## Batch Processing

### Create Multiple Quizzes

```go
func createMultipleQuizzes() error {
    quizData := []struct {
        title     string
        questions []QuestionData
    }{
        {
            title: "Math Quiz",
            questions: []QuestionData{
                {
                    text: "What is 2 + 2?",
                    answers: []AnswerData{
                        {text: "4", correct: true},
                        {text: "3", correct: false},
                        {text: "5", correct: false},
                    },
                },
            },
        },
        {
            title: "Science Quiz",
            questions: []QuestionData{
                {
                    text: "What is H2O?",
                    answers: []AnswerData{
                        {text: "Water", correct: true},
                        {text: "Hydrogen", correct: false},
                        {text: "Oxygen", correct: false},
                    },
                },
            },
        },
    }
    
    for i, data := range quizData {
        quiz, err := buildQuizFromData(data)
        if err != nil {
            return fmt.Errorf("failed to create quiz %d: %w", i+1, err)
        }
        
        filename := fmt.Sprintf("quiz-%d.json", i+1)
        if err := exportQuiz(quiz, filename); err != nil {
            return fmt.Errorf("failed to export quiz %d: %w", i+1, err)
        }
        
        fmt.Printf("Created %s\n", filename)
    }
    
    return nil
}

type QuestionData struct {
    text    string
    answers []AnswerData
}

type AnswerData struct {
    text     string
    correct  bool
    feedback string
}

func buildQuizFromData(data struct {
    title     string
    questions []QuestionData
}) (*h5p.QuestionSet, error) {
    builder := h5p.NewQuestionSetBuilder().
        SetTitle(data.title).
        SetProgressType("textual").
        SetPassPercentage(60)
    
    for _, q := range data.questions {
        answers := make([]h5p.Answer, len(q.answers))
        for i, a := range q.answers {
            answers[i] = h5p.Answer{
                Text:     a.text,
                Correct:  a.correct,
                Feedback: a.feedback,
            }
        }
        
        builder = builder.AddMultipleChoiceQuestion(q.text, answers)
    }
    
    return builder.Build()
}

func exportQuiz(quiz *h5p.QuestionSet, filename string) error {
    jsonData, err := quiz.ToJSON()
    if err != nil {
        return err
    }
    
    return os.WriteFile(filename, jsonData, 0644)
}
```

## Integration Examples

### Web API Integration

```go
// HTTP handler for creating quizzes
func createQuizHandler(w http.ResponseWriter, r *http.Request) {
    var request struct {
        Title      string `json:"title"`
        Questions  []struct {
            Text    string `json:"text"`
            Answers []struct {
                Text     string `json:"text"`
                Correct  bool   `json:"correct"`
                Feedback string `json:"feedback,omitempty"`
            } `json:"answers"`
        } `json:"questions"`
        PassPercentage int `json:"passPercentage"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    // Build quiz from request
    builder := h5p.NewQuestionSetBuilder().
        SetTitle(request.Title).
        SetProgressType("textual").
        SetPassPercentage(request.PassPercentage)
    
    for _, q := range request.Questions {
        answers := make([]h5p.Answer, len(q.Answers))
        for i, a := range q.Answers {
            answers[i] = h5p.Answer{
                Text:     a.Text,
                Correct:  a.Correct,
                Feedback: a.Feedback,
            }
        }
        
        builder = builder.AddMultipleChoiceQuestion(q.Text, answers)
    }
    
    quiz, err := builder.Build()
    if err != nil {
        http.Error(w, fmt.Sprintf("Build failed: %v", err), http.StatusBadRequest)
        return
    }
    
    // Validate
    if err := quiz.Validate(); err != nil {
        http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusBadRequest)
        return
    }
    
    // Return JSON
    w.Header().Set("Content-Type", "application/json")
    jsonData, _ := quiz.ToJSON()
    w.Write(jsonData)
}
```

These examples demonstrate the most common usage patterns for the H5P Go SDK. They show how to create content, handle errors properly, and integrate with larger applications.

## Next Steps

- [Advanced Examples](advanced.md) - More complex usage patterns
- [Package Examples](packages.md) - Creating complete H5P packages  
- [User Guide](../guide/question-sets.md) - Comprehensive documentation