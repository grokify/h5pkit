package schemas

import (
	"fmt"
)

// MultiChoiceParams represents the parameters for H5P.MultiChoice content type
// This struct is generated from the official H5P MultiChoice semantics.json schema
type MultiChoiceParams struct {
	Media           *MediaGroup      `json:"media,omitempty"`
	Question        string           `json:"question"`
	Answers         []AnswerOption   `json:"answers"`
	OverallFeedback *OverallFeedback `json:"overallFeedback,omitempty"`
	Behaviour       *Behaviour       `json:"behaviour,omitempty"`
	UI              *UITranslations  `json:"UI,omitempty"`
}

// MediaGroup represents optional media content (images, videos)
type MediaGroup struct {
	Type                string `json:"type,omitempty"`
	DisableImageZooming bool   `json:"disableImageZooming,omitempty"`
}

// AnswerOption represents a single answer choice
type AnswerOption struct {
	Text            string                 `json:"text"`
	Correct         bool                   `json:"correct"`
	TipsAndFeedback *AnswerTipsAndFeedback `json:"tipsAndFeedback,omitempty"`
}

// AnswerTipsAndFeedback provides hints and feedback for individual answers
type AnswerTipsAndFeedback struct {
	Tip               string `json:"tip,omitempty"`
	ChosenFeedback    string `json:"chosenFeedback,omitempty"`
	NotChosenFeedback string `json:"notChosenFeedback,omitempty"`
}

// OverallFeedback provides score-based feedback ranges
type OverallFeedback struct {
	OverallFeedback []FeedbackRange `json:"overallFeedback,omitempty"`
}

// FeedbackRange defines feedback for a score range
type FeedbackRange struct {
	From     int    `json:"from"`
	To       int    `json:"to"`
	Feedback string `json:"feedback,omitempty"`
}

// Behaviour controls how the MultiChoice question behaves.
//
// Fields whose H5P default is true are pointers, so that an explicit false
// is serialized. A nil pointer is omitted and the player uses its default.
type Behaviour struct {
	EnableRetry           *bool  `json:"enableRetry,omitempty"`           // H5P default: true
	EnableSolutionsButton *bool  `json:"enableSolutionsButton,omitempty"` // H5P default: true
	EnableCheckButton     *bool  `json:"enableCheckButton,omitempty"`     // H5P default: true
	Type                  string `json:"type,omitempty"`                  // "auto", "multi", "single"
	SinglePoint           bool   `json:"singlePoint,omitempty"`
	RandomAnswers         *bool  `json:"randomAnswers,omitempty"` // H5P default: true
	PassPercentage        int    `json:"passPercentage,omitempty"`
	ShowScorePoints       *bool  `json:"showScorePoints,omitempty"` // H5P default: true
}

// Bool returns a pointer to b, for the pointer fields of Behaviour.
func Bool(b bool) *bool { return &b }

// UITranslations contains user interface text labels
type UITranslations struct {
	CheckAnswerButton  string `json:"checkAnswerButton,omitempty"`
	ShowSolutionButton string `json:"showSolutionButton,omitempty"`
	TryAgainButton     string `json:"tryAgainButton,omitempty"`
	TipsLabel          string `json:"tipsLabel,omitempty"`
	ScoreBarLabel      string `json:"scoreBarLabel,omitempty"`
	TipAvailable       string `json:"tipAvailable,omitempty"`
	FeedbackAvailable  string `json:"feedbackAvailable,omitempty"`
	ReadFeedback       string `json:"readFeedback,omitempty"`
	WrongAnswer        string `json:"wrongAnswer,omitempty"`
	CorrectAnswer      string `json:"correctAnswer,omitempty"`
}

// Validate checks if the MultiChoiceParams are valid according to H5P semantics
func (p *MultiChoiceParams) Validate() error {
	if p.Question == "" {
		return fmt.Errorf("question text is required")
	}

	if len(p.Answers) < 1 {
		return fmt.Errorf("at least one answer is required")
	}

	correctCount := 0
	for _, answer := range p.Answers {
		if answer.Text == "" {
			return fmt.Errorf("answer text cannot be empty")
		}
		if answer.Correct {
			correctCount++
		}
	}

	if correctCount == 0 {
		return fmt.Errorf("at least one answer must be marked as correct")
	}

	// Validate behavior settings
	if p.Behaviour != nil {
		if p.Behaviour.Type != "" {
			validTypes := map[string]bool{"auto": true, "multi": true, "single": true}
			if !validTypes[p.Behaviour.Type] {
				return fmt.Errorf("invalid question type: %s", p.Behaviour.Type)
			}
		}

		if p.Behaviour.PassPercentage < 0 || p.Behaviour.PassPercentage > 100 {
			return fmt.Errorf("pass percentage must be between 0 and 100")
		}
	}

	return nil
}
