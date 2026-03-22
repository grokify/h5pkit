package h5p

import (
	"encoding/json"
	"errors"

	"github.com/grokify/h5pkit/schemas"
)

type QuestionSetBuilder struct {
	questionSet *QuestionSet
}

func NewQuestionSetBuilder() *QuestionSetBuilder {
	return &QuestionSetBuilder{
		questionSet: &QuestionSet{
			Questions: make([]Question, 0),
		},
	}
}

func (b *QuestionSetBuilder) SetProgressType(progressType string) *QuestionSetBuilder {
	b.questionSet.ProgressType = progressType
	return b
}

func (b *QuestionSetBuilder) SetPassPercentage(percentage int) *QuestionSetBuilder {
	b.questionSet.PassPercentage = percentage
	return b
}

func (b *QuestionSetBuilder) SetTitle(title string) *QuestionSetBuilder {
	b.questionSet.Title = title
	return b
}

func (b *QuestionSetBuilder) SetIntroduction(introduction string) *QuestionSetBuilder {
	b.questionSet.Introduction = introduction
	b.questionSet.ShowIntroPage = true
	return b
}

func (b *QuestionSetBuilder) SetStartButtonText(text string) *QuestionSetBuilder {
	b.questionSet.StartButtonText = text
	return b
}

func (b *QuestionSetBuilder) SetBackgroundImage(path, mime string) *QuestionSetBuilder {
	b.questionSet.BackgroundImage = &BackgroundImage{
		Path: path,
		Mime: mime,
	}
	return b
}

func (b *QuestionSetBuilder) AddMultipleChoiceQuestion(question string, answers []Answer) *QuestionSetBuilder {
	// Convert legacy Answer to schemas.AnswerOption
	schemaAnswers := make([]schemas.AnswerOption, len(answers))
	for i, answer := range answers {
		schemaAnswers[i] = schemas.AnswerOption{
			Text:    answer.Text,
			Correct: answer.Correct,
			TipsAndFeedback: &schemas.AnswerTipsAndFeedback{
				ChosenFeedback: answer.Feedback,
			},
		}
	}

	params := &schemas.MultiChoiceParams{
		Question: question,
		Answers:  schemaAnswers,
	}

	q := Question{
		Library: "H5P.MultiChoice 1.16",
		Params:  params,
	}

	b.questionSet.Questions = append(b.questionSet.Questions, q)
	return b
}

func (b *QuestionSetBuilder) AddOverallFeedback(ranges []FeedbackRange) *QuestionSetBuilder {
	b.questionSet.OverallFeedback = ranges
	return b
}

func (b *QuestionSetBuilder) Build() (*QuestionSet, error) {
	if len(b.questionSet.Questions) == 0 {
		return nil, errors.New("question set must have at least one question")
	}
	return b.questionSet, nil
}

func CreateAnswer(text string, correct bool) Answer {
	return Answer{
		Text:    text,
		Correct: correct,
	}
}

func CreateAnswerWithFeedback(text string, correct bool, feedback string) Answer {
	return Answer{
		Text:     text,
		Correct:  correct,
		Feedback: feedback,
	}
}

func CreateFeedbackRange(from, to int, text string) FeedbackRange {
	return FeedbackRange{
		From: from,
		To:   to,
		Text: text,
	}
}

func (qs *QuestionSet) ToJSON() ([]byte, error) {
	return json.MarshalIndent(qs, "", "  ")
}

func FromJSON(data []byte) (*QuestionSet, error) {
	var qs QuestionSet
	err := json.Unmarshal(data, &qs)
	if err != nil {
		return nil, err
	}
	return &qs, nil
}

func (qs *QuestionSet) Validate() error {
	if len(qs.Questions) == 0 {
		return errors.New("question set must have at least one question")
	}

	if qs.PassPercentage < 0 || qs.PassPercentage > 100 {
		return errors.New("pass percentage must be between 0 and 100")
	}

	for i, feedback := range qs.OverallFeedback {
		if feedback.From > feedback.To {
			return errors.New("feedback range 'from' cannot be greater than 'to' at index " + string(rune(i)))
		}
	}

	return nil
}
