package h5p

import (
	"encoding/json"
	"testing"

	"github.com/grokify/h5pkit/schemas"
)

func TestTypedMultiChoiceQuestion(t *testing.T) {
	// Create a typed MultiChoice question using the schemas
	params := &schemas.MultiChoiceParams{
		Question: "What is the capital of France?",
		Answers: []schemas.AnswerOption{
			{
				Text:    "Paris",
				Correct: true,
				TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
					ChosenFeedback:    "Correct! Paris is the capital of France.",
					NotChosenFeedback: "Paris is the correct answer.",
				},
			},
			{
				Text:    "London",
				Correct: false,
				TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
					ChosenFeedback:    "Incorrect. London is the capital of the UK.",
					NotChosenFeedback: "London is not the capital of France.",
				},
			},
			{
				Text:    "Berlin",
				Correct: false,
				TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
					ChosenFeedback:    "Incorrect. Berlin is the capital of Germany.",
					NotChosenFeedback: "Berlin is not the capital of France.",
				},
			},
		},
		Behaviour: &schemas.Behaviour{
			EnableRetry:           true,
			EnableSolutionsButton: true,
			Type:                  "single",
			SinglePoint:           false,
			RandomAnswers:         true,
			PassPercentage:        100,
			ShowScorePoints:       true,
		},
		UI: &schemas.UITranslations{
			CheckAnswerButton:  "Check",
			ShowSolutionButton: "Show solution",
			TryAgainButton:     "Retry",
		},
	}

	// Test validation
	err := params.Validate()
	if err != nil {
		t.Errorf("Valid params failed validation: %v", err)
	}

	// Create typed question
	mcQuestion := NewMultiChoiceQuestion(params)
	if mcQuestion.Library != "H5P.MultiChoice 1.16" {
		t.Errorf("Expected library 'H5P.MultiChoice 1.16', got '%s'", mcQuestion.Library)
	}

	// Test conversion to generic Question
	genericQuestion := mcQuestion.ToQuestion()
	if genericQuestion.Library != mcQuestion.Library {
		t.Errorf("Library mismatch after conversion")
	}

	// Test JSON marshaling
	jsonData, err := json.MarshalIndent(mcQuestion, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal typed question: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaled MultiChoiceQuestion
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal typed question: %v", err)
	}

	if unmarshaled.Params.Question != params.Question {
		t.Errorf("Question text mismatch after round-trip")
	}

	if len(unmarshaled.Params.Answers) != len(params.Answers) {
		t.Errorf("Answer count mismatch after round-trip")
	}

	t.Logf("Successfully created and validated typed MultiChoice question")
}

func TestMultiChoiceParamsValidation(t *testing.T) {
	// Test empty question
	params := &schemas.MultiChoiceParams{}
	err := params.Validate()
	if err == nil {
		t.Error("Expected validation error for empty question")
	}

	// Test no answers
	params.Question = "Test question?"
	err = params.Validate()
	if err == nil {
		t.Error("Expected validation error for no answers")
	}

	// Test no correct answers
	params.Answers = []schemas.AnswerOption{
		{Text: "Answer 1", Correct: false},
		{Text: "Answer 2", Correct: false},
	}
	err = params.Validate()
	if err == nil {
		t.Error("Expected validation error for no correct answers")
	}

	// Test empty answer text
	params.Answers = []schemas.AnswerOption{
		{Text: "", Correct: true},
	}
	err = params.Validate()
	if err == nil {
		t.Error("Expected validation error for empty answer text")
	}

	// Test invalid question type
	params.Answers = []schemas.AnswerOption{
		{Text: "Valid answer", Correct: true},
	}
	params.Behaviour = &schemas.Behaviour{
		Type: "invalid",
	}
	err = params.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid question type")
	}

	// Test invalid pass percentage
	params.Behaviour.Type = "single"
	params.Behaviour.PassPercentage = 150
	err = params.Validate()
	if err == nil {
		t.Error("Expected validation error for pass percentage > 100")
	}

	// Test valid params
	params.Behaviour.PassPercentage = 80
	err = params.Validate()
	if err != nil {
		t.Errorf("Valid params failed validation: %v", err)
	}

	t.Log("MultiChoice params validation tests completed")
}

func TestQuestionSetWithTypedQuestions(t *testing.T) {
	// Create typed MultiChoice question
	params := &schemas.MultiChoiceParams{
		Question: "What is 2 + 2?",
		Answers: []schemas.AnswerOption{
			{Text: "3", Correct: false},
			{Text: "4", Correct: true},
			{Text: "5", Correct: false},
		},
		Behaviour: &schemas.Behaviour{
			Type:        "single",
			SinglePoint: false,
		},
	}

	mcQuestion := NewMultiChoiceQuestion(params)

	// Add to question set
	qs := &QuestionSet{
		Title:     "Math Quiz",
		Questions: []Question{*mcQuestion.ToQuestion()},
	}

	// Test JSON serialization
	jsonData, err := qs.ToJSON()
	if err != nil {
		t.Fatalf("Failed to serialize question set: %v", err)
	}

	// Test deserialization
	loadedQS, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to deserialize question set: %v", err)
	}

	if len(loadedQS.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(loadedQS.Questions))
	}

	// Validate the question params can be cast back to MultiChoiceParams
	if questionParams, ok := loadedQS.Questions[0].Params.(map[string]interface{}); ok {
		if question, exists := questionParams["question"]; exists {
			if question != "What is 2 + 2?" {
				t.Errorf("Question text mismatch after round-trip")
			}
		}
	}

	t.Log("Question set with typed questions test completed")
}
