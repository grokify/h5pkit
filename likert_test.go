package h5p

import (
	"errors"
	"math"
	"strings"
	"testing"

	"github.com/grokify/h5pkit/schemas"
)

var agreement5 = []string{"Strongly disagree", "Disagree", "Neutral", "Agree", "Strongly agree"}

func likertSet(t *testing.T) *QuestionSet {
	t.Helper()
	innovation := SkillRef{Taxonomy: "example", ID: "innovation"}
	people := SkillRef{Taxonomy: "example", ID: "people"}
	qs, err := NewQuestionSetBuilder().
		SetTitle("Reflection").
		AddLikertQuestion("I make room for experiments.", "agreement-5", agreement5, false,
			NewExtension("Innovation", 1).WithSkills(innovation)).
		AddLikertQuestion("I prefer proven approaches.", "agreement-5", agreement5, true,
			NewExtension("Innovation", 2).WithSkills(innovation)).
		AddLikertQuestion("I coach others.", "agreement-5", agreement5, false,
			NewExtension("People", 3).WithSkills(people)).
		Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return qs
}

func TestNewLikertParamsIsValidMultiChoice(t *testing.T) {
	p := NewLikertParams("I coach others.", agreement5)
	if err := p.Validate(); err != nil {
		t.Fatalf("MultiChoice Validate: %v", err)
	}
	if err := ValidateLikertParams(p); err != nil {
		t.Fatalf("ValidateLikertParams: %v", err)
	}
}

func TestValidateLikertParams(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*schemas.MultiChoiceParams)
		want   error
	}{
		{"ok", func(*schemas.MultiChoiceParams) {}, nil},
		{"one answer", func(p *schemas.MultiChoiceParams) { p.Answers = p.Answers[:1] }, ErrLikertTooFew},
		{"not correct", func(p *schemas.MultiChoiceParams) { p.Answers[2].Correct = false }, ErrLikertNotCorrect},
		{"no behaviour", func(p *schemas.MultiChoiceParams) { p.Behaviour = nil }, ErrLikertNotSingle},
		{"auto type", func(p *schemas.MultiChoiceParams) { p.Behaviour.Type = "auto" }, ErrLikertNotSingle},
		{"random default", func(p *schemas.MultiChoiceParams) { p.Behaviour.RandomAnswers = nil }, ErrLikertRandomized},
		{"random on", func(p *schemas.MultiChoiceParams) { p.Behaviour.RandomAnswers = schemas.Bool(true) }, ErrLikertRandomized},
		{"check default", func(p *schemas.MultiChoiceParams) { p.Behaviour.EnableCheckButton = nil }, ErrLikertCheckButton},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewLikertParams("Q", agreement5)
			tt.mutate(p)
			err := ValidateLikertParams(p)
			if !errors.Is(err, tt.want) || (tt.want == nil && err != nil) {
				t.Errorf("ValidateLikertParams() = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestQuestionSetValidateChecksLikert(t *testing.T) {
	qs := likertSet(t)
	if err := qs.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
	qs.Questions[1].Params.(*schemas.MultiChoiceParams).Behaviour.RandomAnswers = nil
	if err := qs.Validate(); !errors.Is(err, ErrLikertRandomized) {
		t.Fatalf("Validate() = %v, want ErrLikertRandomized", err)
	}
}

func TestQuestionSetValidateLikertLibrary(t *testing.T) {
	qs := &QuestionSet{Questions: []Question{{
		Library:    "H5P.TrueFalse 1.8",
		Extensions: &Extensions{H5PKit: &Extension{Scoring: NewLikertScoring("s", false)}},
	}}}
	if err := qs.Validate(); !errors.Is(err, ErrLikertLibrary) {
		t.Fatalf("Validate() = %v, want ErrLikertLibrary", err)
	}
}

func TestLikertJSONRoundTrip(t *testing.T) {
	data, err := likertSet(t).ToJSON()
	if err != nil {
		t.Fatalf("ToJSON: %v", err)
	}
	for _, want := range []string{`"randomAnswers": false`, `"enableCheckButton": false`, `"method": "likert"`, `"reverse": true`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("JSON missing %s", want)
		}
	}
	qs, err := FromJSON(data)
	if err != nil {
		t.Fatalf("FromJSON: %v", err)
	}
	if err := qs.Validate(); err != nil {
		t.Fatalf("Validate after round trip: %v", err)
	}
	scores, err := qs.ScoreLikert(map[int]int{0: 4})
	if err != nil {
		t.Fatalf("ScoreLikert after round trip: %v", err)
	}
	if got := scores[SkillRef{Taxonomy: "example", ID: "innovation"}].Mean; got != 5 {
		t.Errorf("Mean = %v, want 5", got)
	}
}

func TestLikertValue(t *testing.T) {
	tests := []struct {
		index, n int
		reverse  bool
		want     int
		err      bool
	}{
		{0, 5, false, 1, false},
		{4, 5, false, 5, false},
		{0, 5, true, 5, false},
		{4, 5, true, 1, false},
		{5, 5, false, 0, true},
		{-1, 5, false, 0, true},
		{0, 1, false, 0, true},
	}
	for _, tt := range tests {
		got, err := LikertValue(tt.index, tt.n, tt.reverse)
		if (err != nil) != tt.err || got != tt.want {
			t.Errorf("LikertValue(%d, %d, %v) = %d, %v", tt.index, tt.n, tt.reverse, got, err)
		}
	}
}

func TestScoreLikert(t *testing.T) {
	qs := likertSet(t)
	// Q1 "Agree" = 4; Q2 reverse, "Disagree" (index 1) = 4; Q3 "Neutral" = 3.
	scores, err := qs.ScoreLikert(map[int]int{0: 3, 1: 1, 2: 2})
	if err != nil {
		t.Fatalf("ScoreLikert: %v", err)
	}
	inn := scores[SkillRef{Taxonomy: "example", ID: "innovation"}]
	if inn.Count != 2 || inn.Mean != 4 || math.Abs(inn.Normalized-0.75) > 1e-9 {
		t.Errorf("innovation = %+v, want count 2, mean 4, normalized 0.75", inn)
	}
	ppl := scores[SkillRef{Taxonomy: "example", ID: "people"}]
	if ppl.Count != 1 || ppl.Mean != 3 || ppl.Normalized != 0.5 {
		t.Errorf("people = %+v, want count 1, mean 3, normalized 0.5", ppl)
	}
}

func TestScoreLikertErrors(t *testing.T) {
	qs := likertSet(t)
	qs.Questions = append(qs.Questions, Question{Library: "H5P.MultiChoice 1.16"})
	tests := []struct {
		name      string
		responses map[int]int
		want      error
	}{
		{"question out of range", map[int]int{9: 0}, ErrLikertResponse},
		{"answer out of range", map[int]int{0: 5}, ErrLikertResponse},
		{"not likert", map[int]int{3: 0}, ErrNotLikert},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := qs.ScoreLikert(tt.responses); !errors.Is(err, tt.want) {
				t.Errorf("ScoreLikert() = %v, want %v", err, tt.want)
			}
		})
	}
}
