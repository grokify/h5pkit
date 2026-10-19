package h5p

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/grokify/h5pkit/schemas"
)

// libraryMultiChoice is the H5P library a Likert question is written as.
const libraryMultiChoice = "H5P.MultiChoice 1.16"

// Errors returned for Likert questions, wrapped with the offending question.
var (
	ErrLikertLibrary     = errors.New("likert question must be H5P.MultiChoice")
	ErrLikertTooFew      = errors.New("likert question needs at least two answers")
	ErrLikertNotCorrect  = errors.New("likert answers must all be marked correct")
	ErrLikertNotSingle   = errors.New("likert question must use behaviour type single")
	ErrLikertRandomized  = errors.New("likert question must set randomAnswers to false")
	ErrLikertCheckButton = errors.New("likert question must set enableCheckButton to false")
	ErrLikertResponse    = errors.New("likert response is out of range")
	ErrNotLikert         = errors.New("question is not a likert question")
)

// NewLikertParams returns MultiChoice params that present labels as an
// ordered scale. Every answer is marked correct, so H5P accepts the
// question and never grades it; answers are single choice, in order, with
// checking, solutions, retry and score points turned off.
func NewLikertParams(question string, labels []string) *schemas.MultiChoiceParams {
	answers := make([]schemas.AnswerOption, len(labels))
	for i, label := range labels {
		answers[i] = schemas.AnswerOption{Text: label, Correct: true}
	}
	return &schemas.MultiChoiceParams{
		Question: question,
		Answers:  answers,
		Behaviour: &schemas.Behaviour{
			Type:                  "single",
			RandomAnswers:         schemas.Bool(false),
			EnableCheckButton:     schemas.Bool(false),
			EnableSolutionsButton: schemas.Bool(false),
			EnableRetry:           schemas.Bool(false),
			ShowScorePoints:       schemas.Bool(false),
		},
	}
}

// ValidateLikertParams checks that params present an ordered, ungraded
// scale: at least two answers, all marked correct, single choice, not
// randomized, and no check button. These are the settings that, if wrong,
// silently change what a respondent sees or what a position means.
func ValidateLikertParams(params *schemas.MultiChoiceParams) error {
	if params == nil || len(params.Answers) < 2 {
		return ErrLikertTooFew
	}
	for i, a := range params.Answers {
		if !a.Correct {
			return fmt.Errorf("answers[%d]: %w", i, ErrLikertNotCorrect)
		}
	}
	b := params.Behaviour
	if b == nil || b.Type != "single" {
		return ErrLikertNotSingle
	}
	if b.RandomAnswers == nil || *b.RandomAnswers {
		return ErrLikertRandomized
	}
	if b.EnableCheckButton == nil || *b.EnableCheckButton {
		return ErrLikertCheckButton
	}
	return nil
}

// AddLikertQuestion adds a Likert question with the given scale labels.
// ext carries the question's metadata; its Scoring is set to Likert on
// scale, with reverse as given. QuestionSet.Validate checks the result.
func (b *QuestionSetBuilder) AddLikertQuestion(question string, scale string, labels []string, reverse bool, ext *Extension) *QuestionSetBuilder {
	if ext == nil {
		ext = &Extension{}
	}
	ext.Scoring = NewLikertScoring(scale, reverse)
	q := Question{
		Library:    libraryMultiChoice,
		Params:     NewLikertParams(question, labels),
		Extensions: &Extensions{H5PKit: ext},
	}
	b.questionSet.Questions = append(b.questionSet.Questions, q)
	return b
}

// likertParams returns the MultiChoice params of a Likert question, decoding
// them when the question was read from JSON.
func likertParams(q Question) (*schemas.MultiChoiceParams, error) {
	if !strings.HasPrefix(q.Library, "H5P.MultiChoice ") {
		return nil, fmt.Errorf("%w: got %q", ErrLikertLibrary, q.Library)
	}
	if p, ok := q.Params.(*schemas.MultiChoiceParams); ok {
		return p, nil
	}
	data, err := json.Marshal(q.Params)
	if err != nil {
		return nil, fmt.Errorf("encode params: %w", err)
	}
	var p schemas.MultiChoiceParams
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("decode params: %w", err)
	}
	return &p, nil
}

// validateLikert checks a question whose extension has Likert scoring.
func validateLikert(q Question) error {
	p, err := likertParams(q)
	if err != nil {
		return err
	}
	return ValidateLikertParams(p)
}

// LikertValue returns the scale value of the answer at position index (0
// based) on a scale of n answers: index+1, or n-index when reverse is true.
func LikertValue(index, n int, reverse bool) (int, error) {
	if n < 2 || index < 0 || index >= n {
		return 0, fmt.Errorf("%w: answer %d of %d", ErrLikertResponse, index, n)
	}
	if reverse {
		return n - index, nil
	}
	return index + 1, nil
}

// LikertScore is the result of scoring the Likert questions for one skill.
type LikertScore struct {
	// Count is the number of answered questions for the skill.
	Count int `json:"count"`

	// Mean is the average scale value, from 1 to the scale size.
	Mean float64 `json:"mean"`

	// Normalized is the average of each value mapped to 0 (lowest) to 1
	// (highest), so that skills on scales of different sizes compare.
	Normalized float64 `json:"normalized"`
}

// ScoreLikert scores responses to the Likert questions in qs by skill.
// responses maps a question's index in qs.Questions to the index of the
// chosen answer. Questions without a response are skipped. A response to a
// question that is not a Likert question, or an answer index out of range,
// is an error. A question with several skills counts toward each.
func (qs *QuestionSet) ScoreLikert(responses map[int]int) (map[SkillRef]LikertScore, error) {
	type sum struct {
		n          int
		value      float64
		normalized float64
	}
	sums := map[SkillRef]*sum{}
	for qi, ai := range responses {
		if qi < 0 || qi >= len(qs.Questions) {
			return nil, fmt.Errorf("question %d: %w", qi, ErrLikertResponse)
		}
		q := qs.Questions[qi]
		ext := q.Extensions.Extension()
		if ext == nil || !ext.Scoring.IsLikert() {
			return nil, fmt.Errorf("question %d: %w", qi, ErrNotLikert)
		}
		p, err := likertParams(q)
		if err != nil {
			return nil, fmt.Errorf("question %d: %w", qi, err)
		}
		n := len(p.Answers)
		v, err := LikertValue(ai, n, ext.Scoring.Reverse)
		if err != nil {
			return nil, fmt.Errorf("question %d: %w", qi, err)
		}
		for _, skill := range ext.Skills {
			s := sums[skill]
			if s == nil {
				s = &sum{}
				sums[skill] = s
			}
			s.n++
			s.value += float64(v)
			s.normalized += float64(v-1) / float64(n-1)
		}
	}
	scores := make(map[SkillRef]LikertScore, len(sums))
	for skill, s := range sums {
		scores[skill] = LikertScore{
			Count:      s.n,
			Mean:       s.value / float64(s.n),
			Normalized: s.normalized / float64(s.n),
		}
	}
	return scores, nil
}
