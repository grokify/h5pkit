package h5p

// Extensions provides a namespace for vendor-specific extensions to H5P content.
// Standard H5P parsers will ignore unknown fields, making this safe for interoperability.
// Each vendor should use their own key within the Extensions map.
type Extensions struct {
	H5PGo *H5PGoExtension `json:"h5pGo,omitempty"`
}

// H5PGoExtension contains h5p-go specific metadata for questions.
// These fields are not part of the official H5P specification but provide
// useful organization and categorization capabilities.
type H5PGoExtension struct {
	// Section identifies the section/category this question belongs to
	// Example: "1. Overview & Fundamentals", "3. Vector RAG"
	Section string `json:"section,omitempty"`

	// Topic provides a more specific topic within a section
	// Example: "RAG Fundamentals", "Chunking Strategies"
	Topic string `json:"topic,omitempty"`

	// Tags provide flexible categorization with multiple labels
	// Example: ["rag", "retrieval", "inference-time"]
	Tags []string `json:"tags,omitempty"`

	// Difficulty indicates the question difficulty level
	// Suggested values: "easy", "medium", "hard"
	Difficulty string `json:"difficulty,omitempty"`

	// QuestionNumber provides explicit ordering within a question set
	QuestionNumber int `json:"questionNumber,omitempty"`

	// LearningObjective describes what the question tests
	// Example: "Understand the difference between RAG and fine-tuning"
	LearningObjective string `json:"learningObjective,omitempty"`

	// Source indicates where the question content originated
	// Example: "PRESENTATION.md", "Chapter 3"
	Source string `json:"source,omitempty"`

	// Custom allows arbitrary additional metadata as key-value pairs
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// NewExtensions creates a new Extensions struct with an initialized H5PGoExtension
func NewExtensions() *Extensions {
	return &Extensions{
		H5PGo: &H5PGoExtension{},
	}
}

// NewH5PGoExtension creates a new H5PGoExtension with the given section and question number
func NewH5PGoExtension(section string, questionNumber int) *H5PGoExtension {
	return &H5PGoExtension{
		Section:        section,
		QuestionNumber: questionNumber,
	}
}

// WithTopic sets the topic and returns the extension for chaining
func (e *H5PGoExtension) WithTopic(topic string) *H5PGoExtension {
	e.Topic = topic
	return e
}

// WithTags sets the tags and returns the extension for chaining
func (e *H5PGoExtension) WithTags(tags ...string) *H5PGoExtension {
	e.Tags = tags
	return e
}

// WithDifficulty sets the difficulty and returns the extension for chaining
func (e *H5PGoExtension) WithDifficulty(difficulty string) *H5PGoExtension {
	e.Difficulty = difficulty
	return e
}

// WithLearningObjective sets the learning objective and returns the extension for chaining
func (e *H5PGoExtension) WithLearningObjective(objective string) *H5PGoExtension {
	e.LearningObjective = objective
	return e
}

// WithSource sets the source and returns the extension for chaining
func (e *H5PGoExtension) WithSource(source string) *H5PGoExtension {
	e.Source = source
	return e
}
