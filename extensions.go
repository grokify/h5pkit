package h5p

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// Extensions holds metadata that is not part of the H5P specification. It is
// serialized as a sibling of a question's "library" and "params", so H5P
// players, which ignore unknown fields, never see it inside library
// parameters.
//
// The h5pkit block is written under the "h5pkit" key.
type Extensions struct {
	// H5PKit is the h5pkit metadata block, serialized under "h5pkit".
	H5PKit *Extension `json:"h5pkit,omitempty"`
}

// Extension returns the h5pkit block, or nil when none is set, including on
// a nil receiver.
func (e *Extensions) Extension() *Extension {
	if e == nil {
		return nil
	}
	return e.H5PKit
}

// Extension is the h5pkit metadata block for a question. It is not part of
// the H5P specification; it adds organization and categorization that
// question banks and authoring tools can share.
type Extension struct {
	// Section identifies the section or category this question belongs to.
	// Example: "1. Overview & Fundamentals", "3. Vector RAG"
	Section string `json:"section,omitempty"`

	// Topic provides a more specific topic within a section.
	// Example: "RAG Fundamentals", "Chunking Strategies"
	Topic string `json:"topic,omitempty"`

	// Tags provide flexible categorization with multiple labels.
	// Example: ["rag", "retrieval", "inference-time"]
	Tags []string `json:"tags,omitempty"`

	// Difficulty is a display label for the question's difficulty.
	// Suggested values: "easy", "medium", "hard"
	Difficulty string `json:"difficulty,omitempty"`

	// DifficultyScore is the question's difficulty relative to the other
	// questions in the same bank, from 0 (easiest) to 1 (hardest). It is not
	// comparable across banks. Nil means unscored.
	DifficultyScore *float64 `json:"difficultyScore,omitempty"`

	// Skills are the skills the question exercises.
	Skills []SkillRef `json:"skills,omitempty"`

	// Standards are the curriculum or competency standards the question is
	// aligned to.
	Standards []StandardRef `json:"standards,omitempty"`

	// QuestionNumber provides explicit ordering within a question set.
	QuestionNumber int `json:"questionNumber,omitempty"`

	// LearningObjective describes what the question tests.
	// Example: "Understand the difference between RAG and fine-tuning"
	LearningObjective string `json:"learningObjective,omitempty"`

	// Source indicates where the question content originated.
	// Example: "PRESENTATION.md", "Chapter 3"
	Source string `json:"source,omitempty"`

	// Scoring describes how a question with no correct answer is scored,
	// such as a Likert item. Nil means the question is scored by H5P.
	Scoring *Scoring `json:"scoring,omitempty"`

	// Custom holds consumer-private metadata. h5pkit never interprets it.
	Custom map[string]any `json:"custom,omitempty"`
}

// SkillRef names a skill within a taxonomy. Taxonomy identifies the skill
// framework, and ID is the skill's identifier within it.
type SkillRef struct {
	Taxonomy string `json:"taxonomy"`
	ID       string `json:"id"`
}

// StandardRef names a standard within a framework. Framework identifies the
// standards framework, and Code is the standard's code within it.
type StandardRef struct {
	Framework string `json:"framework"`
	Code      string `json:"code"`
}

// ScoringMethodLikert scores a question by the position of the chosen
// answer on an ordered scale. See Scoring.
const ScoringMethodLikert = "likert"

// Scoring describes how a consumer scores a question that H5P cannot score.
//
// For ScoringMethodLikert, the question is an H5P.MultiChoice whose answers
// are an ordered scale, all marked correct so that H5P accepts them. The
// value of the chosen answer is its position, 1 to N. When Reverse is true,
// the value is N+1 minus the position.
type Scoring struct {
	// Method is the scoring method. The only method is "likert".
	Method string `json:"method"`

	// Scale identifies the answer scale, so that questions sharing a scale
	// can be compared. Example: "agreement-5"
	Scale string `json:"scale"`

	// Reverse marks a reverse-keyed item, whose scale runs the other way.
	Reverse bool `json:"reverse,omitempty"`
}

// NewLikertScoring returns Likert scoring on the given scale.
func NewLikertScoring(scale string, reverse bool) *Scoring {
	return &Scoring{Method: ScoringMethodLikert, Scale: scale, Reverse: reverse}
}

// IsLikert reports whether s is Likert scoring. It is false for nil.
func (s *Scoring) IsLikert() bool {
	return s != nil && s.Method == ScoringMethodLikert
}

// Errors returned by Extension.Validate, wrapped with the offending field.
var (
	ErrDifficultyScoreRange  = errors.New("difficultyScore must be between 0 and 1")
	ErrSkillRefIncomplete    = errors.New("skill reference needs a taxonomy and an id")
	ErrStandardRefIncomplete = errors.New("standard reference needs a framework and a code")
	ErrScoringMethod         = errors.New("scoring method must be likert")
	ErrScoringScale          = errors.New("scoring needs a scale")
)

// Validate checks the block: DifficultyScore, when set, is within [0, 1];
// every SkillRef has a Taxonomy and an ID; every StandardRef has a
// Framework and a Code; and Scoring, when set, has the likert method and a
// scale. A nil Extension is valid.
func (e *Extension) Validate() error {
	if e == nil {
		return nil
	}
	if e.DifficultyScore != nil {
		d := *e.DifficultyScore
		if math.IsNaN(d) || d < 0 || d > 1 {
			return fmt.Errorf("%w: got %v", ErrDifficultyScoreRange, d)
		}
	}
	for i, s := range e.Skills {
		if strings.TrimSpace(s.Taxonomy) == "" || strings.TrimSpace(s.ID) == "" {
			return fmt.Errorf("skills[%d]: %w", i, ErrSkillRefIncomplete)
		}
	}
	for i, s := range e.Standards {
		if strings.TrimSpace(s.Framework) == "" || strings.TrimSpace(s.Code) == "" {
			return fmt.Errorf("standards[%d]: %w", i, ErrStandardRefIncomplete)
		}
	}
	if e.Scoring != nil {
		if e.Scoring.Method != ScoringMethodLikert {
			return fmt.Errorf("scoring: %w: got %q", ErrScoringMethod, e.Scoring.Method)
		}
		if strings.TrimSpace(e.Scoring.Scale) == "" {
			return fmt.Errorf("scoring: %w", ErrScoringScale)
		}
	}
	return nil
}

// NewExtensions creates an Extensions with an empty h5pkit block.
func NewExtensions() *Extensions {
	return &Extensions{
		H5PKit: &Extension{},
	}
}

// NewExtension creates an Extension with the given section and question
// number.
func NewExtension(section string, questionNumber int) *Extension {
	return &Extension{
		Section:        section,
		QuestionNumber: questionNumber,
	}
}

// WithTopic sets the topic and returns the extension for chaining.
func (e *Extension) WithTopic(topic string) *Extension {
	e.Topic = topic
	return e
}

// WithTags sets the tags and returns the extension for chaining.
func (e *Extension) WithTags(tags ...string) *Extension {
	e.Tags = tags
	return e
}

// WithDifficulty sets the difficulty label and returns the extension for
// chaining.
func (e *Extension) WithDifficulty(difficulty string) *Extension {
	e.Difficulty = difficulty
	return e
}

// WithLearningObjective sets the learning objective and returns the
// extension for chaining.
func (e *Extension) WithLearningObjective(objective string) *Extension {
	e.LearningObjective = objective
	return e
}

// WithSource sets the source and returns the extension for chaining.
func (e *Extension) WithSource(source string) *Extension {
	e.Source = source
	return e
}

// WithDifficultyScore sets the relative difficulty score and returns the
// extension for chaining. See DifficultyScore for the scale.
func (e *Extension) WithDifficultyScore(score float64) *Extension {
	e.DifficultyScore = &score
	return e
}

// WithSkills sets the skills and returns the extension for chaining.
func (e *Extension) WithSkills(skills ...SkillRef) *Extension {
	e.Skills = skills
	return e
}

// WithScoring sets the scoring and returns the extension for chaining.
func (e *Extension) WithScoring(scoring *Scoring) *Extension {
	e.Scoring = scoring
	return e
}

// WithStandards sets the standards and returns the extension for chaining.
func (e *Extension) WithStandards(standards ...StandardRef) *Extension {
	e.Standards = standards
	return e
}
