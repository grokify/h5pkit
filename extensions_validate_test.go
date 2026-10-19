package h5p

import (
	"encoding/json"
	"errors"
	"math"
	"testing"
)

func score(f float64) *float64 { return &f }

func TestExtensionValidate(t *testing.T) {
	tests := []struct {
		name string
		ext  *Extension
		want error
	}{
		{"nil", nil, nil},
		{"empty", &Extension{}, nil},
		{"score zero", &Extension{DifficultyScore: score(0)}, nil},
		{"score one", &Extension{DifficultyScore: score(1)}, nil},
		{"score below", &Extension{DifficultyScore: score(-0.01)}, ErrDifficultyScoreRange},
		{"score above", &Extension{DifficultyScore: score(1.01)}, ErrDifficultyScoreRange},
		{"score NaN", &Extension{DifficultyScore: score(math.NaN())}, ErrDifficultyScoreRange},
		{"skill ok", &Extension{Skills: []SkillRef{{Taxonomy: "tax", ID: "fractions.add"}}}, nil},
		{"skill no taxonomy", &Extension{Skills: []SkillRef{{ID: "x"}}}, ErrSkillRefIncomplete},
		{"skill no id", &Extension{Skills: []SkillRef{{Taxonomy: "tax", ID: " "}}}, ErrSkillRefIncomplete},
		{"standard ok", &Extension{Standards: []StandardRef{{Framework: "fw", Code: "1.A"}}}, nil},
		{"standard no framework", &Extension{Standards: []StandardRef{{Code: "1.A"}}}, ErrStandardRefIncomplete},
		{"standard no code", &Extension{Standards: []StandardRef{{Framework: "fw"}}}, ErrStandardRefIncomplete},
		{"scoring likert", &Extension{Scoring: NewLikertScoring("agreement-5", true)}, nil},
		{"scoring no method", &Extension{Scoring: &Scoring{Scale: "agreement-5"}}, ErrScoringMethod},
		{"scoring unknown method", &Extension{Scoring: &Scoring{Method: "rank", Scale: "s"}}, ErrScoringMethod},
		{"scoring no scale", &Extension{Scoring: &Scoring{Method: ScoringMethodLikert, Scale: " "}}, ErrScoringScale},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.ext.Validate()
			if !errors.Is(err, tt.want) || (tt.want == nil && err != nil) {
				t.Errorf("Validate() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestQuestionSetValidateChecksExtensions(t *testing.T) {
	qs := &QuestionSet{Questions: []Question{
		{Library: "H5P.MultiChoice 1.16", Extensions: &Extensions{H5PKit: NewExtension("S", 1).WithDifficultyScore(0.5)}},
		{Library: "H5P.MultiChoice 1.16", Extensions: &Extensions{H5PKit: &Extension{DifficultyScore: score(2)}}},
	}}
	if err := qs.Validate(); !errors.Is(err, ErrDifficultyScoreRange) {
		t.Fatalf("Validate() = %v, want ErrDifficultyScoreRange", err)
	}
}

func TestExtensionNewFieldsRoundTrip(t *testing.T) {
	ext := NewExtension("S", 1).
		WithDifficultyScore(0.25).
		WithSkills(SkillRef{Taxonomy: "tax", ID: "a"}).
		WithStandards(StandardRef{Framework: "fw", Code: "1.A"})
	data, err := json.Marshal(Extensions{H5PKit: ext})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	want := `{"h5pkit":{"section":"S","difficultyScore":0.25,"skills":[{"taxonomy":"tax","id":"a"}],"standards":[{"framework":"fw","code":"1.A"}],"questionNumber":1}}`
	if string(data) != want {
		t.Fatalf("got  %s\nwant %s", data, want)
	}
	var back Extensions
	if err := json.Unmarshal(data, &back); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if back.H5PKit.DifficultyScore == nil || *back.H5PKit.DifficultyScore != 0.25 ||
		back.H5PKit.Skills[0] != ext.Skills[0] || back.H5PKit.Standards[0] != ext.Standards[0] {
		t.Errorf("round trip mismatch: %+v", back.H5PKit)
	}
}
