package h5p

import (
	"encoding/json"
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
// The h5pkit block is written under the "h5pkit" key. Content written before
// the rename used the "h5pGo" key; it is still read, and promoted into
// H5PKit, but it is never written.
type Extensions struct {
	// H5PKit is the h5pkit metadata block, serialized under "h5pkit".
	H5PKit *Extension `json:"h5pkit,omitempty"`

	// H5PGo holds a block read from the legacy "h5pGo" key.
	//
	// Deprecated: read-only alias for content written before the rename to
	// "h5pkit". Unmarshaling promotes it into H5PKit when "h5pkit" is absent,
	// and marshaling never writes it. Use H5PKit.
	H5PGo *Extension `json:"h5pGo,omitempty"`
}

// extensionsJSON mirrors Extensions without its methods, so the custom
// marshalers can use the default encoding.
type extensionsJSON struct {
	H5PKit *Extension `json:"h5pkit,omitempty"`
	H5PGo  *Extension `json:"h5pGo,omitempty"`
}

// UnmarshalJSON reads both the "h5pkit" and the legacy "h5pGo" keys. When
// only "h5pGo" is present, its block is promoted into H5PKit.
func (e *Extensions) UnmarshalJSON(data []byte) error {
	var raw extensionsJSON
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	e.H5PKit = raw.H5PKit
	e.H5PGo = raw.H5PGo
	if e.H5PKit == nil && e.H5PGo != nil {
		e.H5PKit = e.H5PGo
	}
	return nil
}

// MarshalJSON writes the block under "h5pkit" only. If H5PKit is nil but
// the legacy H5PGo field is set, that block is written under "h5pkit".
func (e Extensions) MarshalJSON() ([]byte, error) {
	return json.Marshal(extensionsJSON{H5PKit: e.Extension()})
}

// Extension returns the h5pkit block: H5PKit, or the legacy H5PGo block
// when H5PKit is nil. It returns nil when neither is set, including on a nil
// receiver.
func (e *Extensions) Extension() *Extension {
	if e == nil {
		return nil
	}
	if e.H5PKit != nil {
		return e.H5PKit
	}
	return e.H5PGo
}

// MigrateLegacyExtensions rewrites every question in qs that carries a
// legacy "h5pGo" block so the block is held in H5PKit, and clears the legacy
// field. A question that already has an H5PKit block keeps it. It returns
// the number of questions changed.
func MigrateLegacyExtensions(qs *QuestionSet) int {
	if qs == nil {
		return 0
	}
	n := 0
	for i := range qs.Questions {
		ext := qs.Questions[i].Extensions
		if ext == nil || ext.H5PGo == nil {
			continue
		}
		if ext.H5PKit == nil {
			ext.H5PKit = ext.H5PGo
		}
		ext.H5PGo = nil
		n++
	}
	return n
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

// Errors returned by Extension.Validate, wrapped with the offending field.
var (
	ErrDifficultyScoreRange  = errors.New("difficultyScore must be between 0 and 1")
	ErrSkillRefIncomplete    = errors.New("skill reference needs a taxonomy and an id")
	ErrStandardRefIncomplete = errors.New("standard reference needs a framework and a code")
)

// Validate checks the block: DifficultyScore, when set, is within [0, 1];
// every SkillRef has a Taxonomy and an ID; and every StandardRef has a
// Framework and a Code. A nil Extension is valid.
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
	return nil
}

// H5PGoExtension is the former name of Extension.
//
// Deprecated: use Extension.
type H5PGoExtension = Extension

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

// NewH5PGoExtension creates an Extension with the given section and question
// number.
//
// Deprecated: use NewExtension.
func NewH5PGoExtension(section string, questionNumber int) *Extension {
	return NewExtension(section, questionNumber)
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

// WithStandards sets the standards and returns the extension for chaining.
func (e *Extension) WithStandards(standards ...StandardRef) *Extension {
	e.Standards = standards
	return e
}
