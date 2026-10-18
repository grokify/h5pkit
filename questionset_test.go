package h5p

import (
	"os"
	"strings"
	"testing"
)

func TestQuestionSetBuilder(t *testing.T) {
	builder := NewQuestionSetBuilder()

	answers := []Answer{
		CreateAnswer("Correct answer", true),
		CreateAnswer("Wrong answer", false),
	}

	qs, err := builder.
		SetTitle("Test Quiz").
		SetPassPercentage(70).
		AddMultipleChoiceQuestion("Test question?", answers).
		Build()

	if err != nil {
		t.Fatalf("Failed to build question set: %v", err)
	}

	if qs.Title != "Test Quiz" {
		t.Errorf("Expected title 'Test Quiz', got '%s'", qs.Title)
	}

	if qs.PassPercentage != 70 {
		t.Errorf("Expected pass percentage 70, got %d", qs.PassPercentage)
	}

	if len(qs.Questions) != 1 {
		t.Errorf("Expected 1 question, got %d", len(qs.Questions))
	}
}

func TestQuestionSetValidation(t *testing.T) {
	qs := &QuestionSet{
		PassPercentage: 150,
		Questions:      []Question{},
	}

	err := qs.Validate()
	if err == nil {
		t.Error("Expected validation error for invalid pass percentage and no questions")
	}
}

func TestJSONMarshalUnmarshal(t *testing.T) {
	builder := NewQuestionSetBuilder()

	answers := []Answer{
		CreateAnswer("Answer 1", true),
		CreateAnswer("Answer 2", false),
	}

	original, err := builder.
		SetTitle("JSON Test").
		AddMultipleChoiceQuestion("Test question?", answers).
		Build()

	if err != nil {
		t.Fatalf("Failed to build question set: %v", err)
	}

	jsonData, err := original.ToJSON()
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	restored, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	if restored.Title != original.Title {
		t.Errorf("Title mismatch after JSON round-trip: expected '%s', got '%s'", original.Title, restored.Title)
	}

	if len(restored.Questions) != len(original.Questions) {
		t.Errorf("Question count mismatch after JSON round-trip: expected %d, got %d", len(original.Questions), len(restored.Questions))
	}
}

func TestCreateAnswerFunctions(t *testing.T) {
	answer := CreateAnswer("Test answer", true)
	if answer.Text != "Test answer" || !answer.Correct {
		t.Error("CreateAnswer did not create answer correctly")
	}

	answerWithFeedback := CreateAnswerWithFeedback("Test answer", false, "Feedback text")
	if answerWithFeedback.Feedback != "Feedback text" {
		t.Error("CreateAnswerWithFeedback did not set feedback correctly")
	}
}

func TestFeedbackRange(t *testing.T) {
	fr := CreateFeedbackRange(80, 100, "Excellent!")
	if fr.From != 80 || fr.To != 100 || fr.Text != "Excellent!" {
		t.Error("CreateFeedbackRange did not create range correctly")
	}
}

func TestReadSampleQuestionSetJSON(t *testing.T) {
	jsonData, err := os.ReadFile("testdata/sample_questionset.json")
	if err != nil {
		t.Fatalf("Failed to read sample JSON file: %v", err)
	}

	questionSet, err := FromJSON(jsonData)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if questionSet.Title != "Mixed Quiz Types" {
		t.Errorf("Expected title 'Mixed Quiz Types', got '%s'", questionSet.Title)
	}

	if len(questionSet.Questions) != 4 {
		t.Errorf("Expected 4 questions, got %d", len(questionSet.Questions))
	}

	singleAnswerCount := 0
	multiAnswerCount := 0

	for i, question := range questionSet.Questions {
		// Since params are now interface{}, we need to handle them differently
		if paramsMap, ok := question.Params.(map[string]interface{}); ok {
			if behaviourInterface, exists := paramsMap["behaviour"]; exists {
				if behaviour, ok := behaviourInterface.(map[string]interface{}); ok {
					if questionType, exists := behaviour["type"]; exists {
						if questionType == "single" {
							singleAnswerCount++
						} else if questionType == "multi" {
							multiAnswerCount++
						}
					}
				}
			}
		}

		if question.Library != "H5P.MultiChoice 1.16" {
			t.Errorf("Question %d: expected library 'H5P.MultiChoice 1.16', got '%s'", i, question.Library)
		}

		if paramsMap, ok := question.Params.(map[string]interface{}); ok {
			if answersInterface, exists := paramsMap["answers"]; exists {
				if answers, ok := answersInterface.([]interface{}); ok {
					if len(answers) == 0 {
						t.Errorf("Question %d: no answers found", i)
					}
				}
			}
		}
	}

	if singleAnswerCount != 2 {
		t.Errorf("Expected 2 single-answer questions, got %d", singleAnswerCount)
	}

	if multiAnswerCount != 2 {
		t.Errorf("Expected 2 multi-answer questions, got %d", multiAnswerCount)
	}

	err = questionSet.Validate()
	if err != nil {
		t.Errorf("Question set validation failed: %v", err)
	}

	t.Logf("Successfully parsed question set with %d questions (%d single-answer, %d multi-answer)",
		len(questionSet.Questions), singleAnswerCount, multiAnswerCount)
}

func TestValidateReportsFeedbackRangeIndex(t *testing.T) {
	qs := &QuestionSet{
		Questions:       []Question{{Library: "H5P.MultiChoice 1.16"}},
		OverallFeedback: []FeedbackRange{{From: 0, To: 50}, {From: 90, To: 10}},
	}
	err := qs.Validate()
	if err == nil || !strings.Contains(err.Error(), "index 1") {
		t.Fatalf("got %v, want an error naming index 1", err)
	}
}
