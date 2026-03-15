package h5p

import (
	"encoding/json"
	"os"
	"testing"
)

func TestExtensionsJSONParsing(t *testing.T) {
	// Test that we can marshal/unmarshal Extensions correctly
	ext := &Extensions{
		H5PGo: &H5PGoExtension{
			Section:           "1. Overview & Fundamentals",
			Topic:             "Introduction",
			Tags:              []string{"basics", "overview"},
			Difficulty:        "easy",
			QuestionNumber:    1,
			LearningObjective: "Understand basic concepts",
			Source:            "PRESENTATION.md",
		},
	}

	data, err := json.Marshal(ext)
	if err != nil {
		t.Fatalf("Failed to marshal Extensions: %v", err)
	}

	var parsed Extensions
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal Extensions: %v", err)
	}

	if parsed.H5PGo == nil {
		t.Fatal("H5PGo extension is nil after unmarshal")
	}

	if parsed.H5PGo.Section != ext.H5PGo.Section {
		t.Errorf("Section mismatch: got %q, want %q", parsed.H5PGo.Section, ext.H5PGo.Section)
	}

	if parsed.H5PGo.QuestionNumber != ext.H5PGo.QuestionNumber {
		t.Errorf("QuestionNumber mismatch: got %d, want %d", parsed.H5PGo.QuestionNumber, ext.H5PGo.QuestionNumber)
	}

	if len(parsed.H5PGo.Tags) != len(ext.H5PGo.Tags) {
		t.Errorf("Tags length mismatch: got %d, want %d", len(parsed.H5PGo.Tags), len(ext.H5PGo.Tags))
	}
}

func TestQuestionWithExtensions(t *testing.T) {
	// Test that Question struct correctly handles extensions
	q := &Question{
		Library: "H5P.MultiChoice 1.16",
		Params:  map[string]interface{}{"question": "Test question?"},
		Extensions: &Extensions{
			H5PGo: &H5PGoExtension{
				Section:        "Test Section",
				QuestionNumber: 42,
				Difficulty:     "medium",
			},
		},
	}

	data, err := json.Marshal(q)
	if err != nil {
		t.Fatalf("Failed to marshal Question: %v", err)
	}

	var parsed Question
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal Question: %v", err)
	}

	if parsed.Extensions == nil {
		t.Fatal("Extensions is nil after unmarshal")
	}

	if parsed.Extensions.H5PGo == nil {
		t.Fatal("H5PGo is nil after unmarshal")
	}

	if parsed.Extensions.H5PGo.QuestionNumber != 42 {
		t.Errorf("QuestionNumber mismatch: got %d, want 42", parsed.Extensions.H5PGo.QuestionNumber)
	}
}

func TestNewH5PGoExtension(t *testing.T) {
	ext := NewH5PGoExtension("Test Section", 5).
		WithTopic("Test Topic").
		WithDifficulty("hard").
		WithTags("tag1", "tag2").
		WithLearningObjective("Test objective").
		WithSource("test.md")

	if ext.Section != "Test Section" {
		t.Errorf("Section mismatch: got %q", ext.Section)
	}
	if ext.QuestionNumber != 5 {
		t.Errorf("QuestionNumber mismatch: got %d", ext.QuestionNumber)
	}
	if ext.Topic != "Test Topic" {
		t.Errorf("Topic mismatch: got %q", ext.Topic)
	}
	if ext.Difficulty != "hard" {
		t.Errorf("Difficulty mismatch: got %q", ext.Difficulty)
	}
	if len(ext.Tags) != 2 {
		t.Errorf("Tags length mismatch: got %d", len(ext.Tags))
	}
}

func TestQuestionSetWithExtensions(t *testing.T) {
	// Test parsing a QuestionSet with extensions
	jsonData := `{
		"title": "Test Quiz",
		"questions": [
			{
				"library": "H5P.MultiChoice 1.16",
				"params": {"question": "What is 2+2?"},
				"extensions": {
					"h5pGo": {
						"section": "Math Basics",
						"questionNumber": 1,
						"difficulty": "easy",
						"tags": ["math", "arithmetic"]
					}
				}
			}
		]
	}`

	var qs QuestionSet
	if err := json.Unmarshal([]byte(jsonData), &qs); err != nil {
		t.Fatalf("Failed to unmarshal QuestionSet: %v", err)
	}

	if len(qs.Questions) != 1 {
		t.Fatalf("Expected 1 question, got %d", len(qs.Questions))
	}

	if qs.Questions[0].Extensions == nil {
		t.Fatal("Extensions is nil")
	}

	if qs.Questions[0].Extensions.H5PGo == nil {
		t.Fatal("H5PGo is nil")
	}

	ext := qs.Questions[0].Extensions.H5PGo
	if ext.Section != "Math Basics" {
		t.Errorf("Section mismatch: got %q", ext.Section)
	}
	if ext.Difficulty != "easy" {
		t.Errorf("Difficulty mismatch: got %q", ext.Difficulty)
	}
	if len(ext.Tags) != 2 {
		t.Errorf("Tags mismatch: got %d tags", len(ext.Tags))
	}
}

// TestExternalQuizFile tests parsing an external quiz file if it exists
// This validates real-world usage of the extensions format
func TestExternalQuizFile(t *testing.T) {
	// Path to the LLM tuning quiz file
	quizPath := "../../../agentplexus/agentplexus-courses-internal/llm-tuning/quiz_llm_techniques.json"

	// Skip if file doesn't exist (allows test to pass in CI without the file)
	if _, err := os.Stat(quizPath); os.IsNotExist(err) {
		t.Skip("External quiz file not found, skipping validation test")
	}

	data, err := os.ReadFile(quizPath)
	if err != nil {
		t.Fatalf("Failed to read quiz file: %v", err)
	}

	var qs QuestionSet
	if err := json.Unmarshal(data, &qs); err != nil {
		t.Fatalf("Failed to unmarshal quiz file: %v", err)
	}

	// Validate basic structure
	if qs.Title == "" {
		t.Error("Quiz title is empty")
	}

	if len(qs.Questions) == 0 {
		t.Fatal("No questions found in quiz")
	}

	t.Logf("Parsed quiz: %q with %d questions", qs.Title, len(qs.Questions))

	// Validate each question has extensions
	sectionsFound := make(map[string]int)
	for i, q := range qs.Questions {
		if q.Library == "" {
			t.Errorf("Question %d: library is empty", i+1)
		}

		if q.Extensions == nil {
			t.Errorf("Question %d: extensions is nil", i+1)
			continue
		}

		if q.Extensions.H5PGo == nil {
			t.Errorf("Question %d: h5pGo extension is nil", i+1)
			continue
		}

		ext := q.Extensions.H5PGo
		if ext.Section == "" {
			t.Errorf("Question %d: section is empty", i+1)
		}
		if ext.QuestionNumber == 0 {
			t.Errorf("Question %d: questionNumber is 0", i+1)
		}

		sectionsFound[ext.Section]++
	}

	t.Logf("Sections found: %v", sectionsFound)

	// Expect 10 sections based on the quiz structure
	if len(sectionsFound) < 5 {
		t.Errorf("Expected at least 5 sections, found %d", len(sectionsFound))
	}
}
