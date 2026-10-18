package h5p

import (
	"encoding/json"
	"testing"
)

func TestExtensionsJSONParsing(t *testing.T) {
	// Test that we can marshal/unmarshal Extensions correctly
	ext := &Extensions{
		H5PKit: &Extension{
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

	if parsed.H5PKit == nil {
		t.Fatal("H5PKit extension is nil after unmarshal")
	}

	if parsed.H5PKit.Section != ext.H5PKit.Section {
		t.Errorf("Section mismatch: got %q, want %q", parsed.H5PKit.Section, ext.H5PKit.Section)
	}

	if parsed.H5PKit.QuestionNumber != ext.H5PKit.QuestionNumber {
		t.Errorf("QuestionNumber mismatch: got %d, want %d", parsed.H5PKit.QuestionNumber, ext.H5PKit.QuestionNumber)
	}

	if len(parsed.H5PKit.Tags) != len(ext.H5PKit.Tags) {
		t.Errorf("Tags length mismatch: got %d, want %d", len(parsed.H5PKit.Tags), len(ext.H5PKit.Tags))
	}
}

func TestQuestionWithExtensions(t *testing.T) {
	// Test that Question struct correctly handles extensions
	q := &Question{
		Library: "H5P.MultiChoice 1.16",
		Params:  map[string]interface{}{"question": "Test question?"},
		Extensions: &Extensions{
			H5PKit: &Extension{
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

	if parsed.Extensions.H5PKit == nil {
		t.Fatal("H5PKit is nil after unmarshal")
	}

	if parsed.Extensions.H5PKit.QuestionNumber != 42 {
		t.Errorf("QuestionNumber mismatch: got %d, want 42", parsed.Extensions.H5PKit.QuestionNumber)
	}
}

func TestNewExtension(t *testing.T) {
	ext := NewExtension("Test Section", 5).
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
					"h5pkit": {
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

	if qs.Questions[0].Extensions.H5PKit == nil {
		t.Fatal("H5PKit is nil")
	}

	ext := qs.Questions[0].Extensions.H5PKit
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
