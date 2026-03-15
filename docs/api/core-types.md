# Core Types

This reference covers the fundamental types and interfaces in the H5P Go SDK.

## Question Set Types

### QuestionSet

The main structure for H5P question sets:

```go
type QuestionSet struct {
    Title            string             `json:"title"`
    ProgressType     string             `json:"progressType"`
    PassPercentage   int                `json:"passPercentage"`
    Questions        []Question         `json:"questions"`
    IntroPage        *IntroPage         `json:"introPage,omitempty"`
    TextualProgress  string             `json:"textualProgress,omitempty"`
    EndGame          *EndGame           `json:"endGame,omitempty"`
    OverallFeedback  []OverallFeedback  `json:"overallFeedback,omitempty"`
    Behaviour        *Behaviour         `json:"behaviour,omitempty"`
    L10n             *L10n              `json:"l10n,omitempty"`
}
```

#### Methods

```go
// ToJSON serializes the question set to JSON
func (qs *QuestionSet) ToJSON() ([]byte, error)

// Validate checks if the question set is valid
func (qs *QuestionSet) Validate() error

// AddQuestion adds a question to the set
func (qs *QuestionSet) AddQuestion(question Question) error
```

### Question

Represents a single question in a question set:

```go
type Question struct {
    Library    string      `json:"library"`
    Params     interface{} `json:"params"`
    SubContentId string    `json:"subContentId,omitempty"`
    Metadata   *Metadata   `json:"metadata,omitempty"`
}
```

### Answer

Represents an answer option in multiple choice questions:

```go
type Answer struct {
    Text     string  `json:"text"`
    Correct  bool    `json:"correct"`
    Feedback string  `json:"feedback,omitempty"`
    Weight   float64 `json:"weight,omitempty"`
}
```

#### Helper Functions

```go
// CreateAnswer creates a new answer with text and correctness
func CreateAnswer(text string, correct bool) Answer

// CreateAnswerWithFeedback creates an answer with feedback
func CreateAnswerWithFeedback(text string, correct bool, feedback string) Answer
```

## Builder Types

### QuestionSetBuilder

Fluent interface for building question sets:

```go
type QuestionSetBuilder struct {
    questionSet *QuestionSet
    errors      []error
}
```

#### Methods

```go
// SetTitle sets the question set title
func (b *QuestionSetBuilder) SetTitle(title string) *QuestionSetBuilder

// SetProgressType sets progress display type ("textual" or "dots")
func (b *QuestionSetBuilder) SetProgressType(progressType string) *QuestionSetBuilder

// SetPassPercentage sets the passing threshold (0-100)
func (b *QuestionSetBuilder) SetPassPercentage(percentage int) *QuestionSetBuilder

// SetIntroduction sets the introduction text
func (b *QuestionSetBuilder) SetIntroduction(text string) *QuestionSetBuilder

// AddMultipleChoiceQuestion adds a multiple choice question
func (b *QuestionSetBuilder) AddMultipleChoiceQuestion(question string, answers []Answer) *QuestionSetBuilder

// AddOverallFeedback adds performance-based feedback
func (b *QuestionSetBuilder) AddOverallFeedback(feedback []OverallFeedback) *QuestionSetBuilder

// Build creates the final question set
func (b *QuestionSetBuilder) Build() (*QuestionSet, error)
```

## Package Types

### H5PPackage

Represents a complete H5P package:

```go
type H5PPackage struct {
    packageDefinition *PackageDefinition
    content          *Content
    libraries        map[string]*Library
}
```

#### Methods

```go
// SetPackageDefinition sets the package metadata
func (pkg *H5PPackage) SetPackageDefinition(def *PackageDefinition)

// SetContent sets the main content
func (pkg *H5PPackage) SetContent(content *Content)

// AddLibrary adds a library to the package
func (pkg *H5PPackage) AddLibrary(lib *Library)

// CreateZipFile exports the package as a .h5p file
func (pkg *H5PPackage) CreateZipFile(filename string) error

// Validate validates the complete package
func (pkg *H5PPackage) Validate() error
```

### PackageDefinition

Defines package metadata (h5p.json):

```go
type PackageDefinition struct {
    Title       string              `json:"title"`
    Language    string              `json:"language"`
    MainLibrary string              `json:"mainLibrary"`
    EmbedTypes  []string           `json:"embedTypes"`
    
    PreloadedDependencies []LibraryDependency `json:"preloadedDependencies"`
    DynamicDependencies   []LibraryDependency `json:"dynamicDependencies,omitempty"`
    EditorDependencies    []LibraryDependency `json:"editorDependencies,omitempty"`
    
    License       string   `json:"license,omitempty"`
    Authors       []Author `json:"authors,omitempty"`
    Changes       []Change `json:"changes,omitempty"`
}
```

### Library

Represents an H5P library:

```go
type Library struct {
    MachineName string                `json:"-"`
    Definition  *LibraryDefinition    `json:"-"`
    Files       map[string][]byte     `json:"-"`
}
```

### Content

Main content structure (content/content.json):

```go
type Content struct {
    Params   interface{} `json:"params"`
    Metadata *Metadata   `json:"metadata,omitempty"`
}
```

## Feedback Types

### OverallFeedback

Performance-based feedback for question sets:

```go
type OverallFeedback struct {
    From     int    `json:"from"`
    To       int    `json:"to"`
    Feedback string `json:"feedback"`
}
```

### IntroPage

Introduction page for question sets:

```go
type IntroPage struct {
    ShowIntroPage      bool   `json:"showIntroPage"`
    Title              string `json:"title,omitempty"`
    Introduction       string `json:"introduction,omitempty"`
    StartButtonText    string `json:"startButtonText,omitempty"`
}
```

### EndGame

End screen configuration:

```go
type EndGame struct {
    ShowResultPage       bool   `json:"showResultPage"`
    ShowSolutionButton   bool   `json:"showSolutionButton"`
    ShowRetryButton      bool   `json:"showRetryButton"`
    NoResultMessage      string `json:"noResultMessage,omitempty"`
    Message              string `json:"message,omitempty"`
    OverallFeedback      []OverallFeedback `json:"overallFeedback,omitempty"`
    SolutionButtonText   string `json:"solutionButtonText,omitempty"`
    RetryButtonText      string `json:"retryButtonText,omitempty"`
}
```

## Metadata Types

### Author

Author information:

```go
type Author struct {
    Name string `json:"name"`
    Role string `json:"role"`
}
```

### Change

Change log entry:

```go
type Change struct {
    Date   string `json:"date"`
    Author string `json:"author"`
    Log    string `json:"log"`
}
```

### LibraryDependency

Library dependency specification:

```go
type LibraryDependency struct {
    MachineName  string `json:"machineName"`
    MajorVersion int    `json:"majorVersion"`
    MinorVersion int    `json:"minorVersion"`
}
```

## Utility Types

### ValidationError

Structured validation error:

```go
type ValidationError struct {
    Field   string      `json:"field"`
    Message string      `json:"message"`
    Value   interface{} `json:"value,omitempty"`
}
```

### ValidationErrors

Collection of validation errors:

```go
type ValidationErrors []ValidationError
```

## Factory Functions

### Package Creation

```go
// NewH5PPackage creates a new empty H5P package
func NewH5PPackage() *H5PPackage

// LoadH5PPackage loads an existing .h5p file
func LoadH5PPackage(filename string) (*H5PPackage, error)
```

### Builder Creation

```go
// NewQuestionSetBuilder creates a new question set builder
func NewQuestionSetBuilder() *QuestionSetBuilder
```

### Content Parsing

```go
// FromJSON parses a question set from JSON data
func FromJSON(jsonData []byte) (*QuestionSet, error)
```

## Constants

### Progress Types

```go
const (
    ProgressTypeTextual = "textual"
    ProgressTypeDots    = "dots"
)
```

### Library Names

```go
const (
    LibraryMultiChoice  = "H5P.MultiChoice"
    LibraryQuestionSet  = "H5P.QuestionSet"
    LibraryTrueFalse   = "H5P.TrueFalse"
    LibraryEssay       = "H5P.Essay"
)
```

## Interface Specifications

### Validator

Types that can be validated implement this interface:

```go
type Validator interface {
    Validate() error
}
```

### JSONMarshaler

Types that can be serialized to JSON:

```go
type JSONMarshaler interface {
    ToJSON() ([]byte, error)
}
```

## Usage Examples

### Basic Question Set Creation

```go
builder := h5p.NewQuestionSetBuilder()

answers := []h5p.Answer{
    h5p.CreateAnswer("Correct", true),
    h5p.CreateAnswer("Wrong", false),
}

questionSet, err := builder.
    SetTitle("Basic Quiz").
    SetPassPercentage(70).
    AddMultipleChoiceQuestion("Question?", answers).
    Build()

if err != nil {
    log.Fatal(err)
}
```

### Package Creation

```go
pkg := h5p.NewH5PPackage()

packageDef := &h5p.PackageDefinition{
    Title:       "My Package",
    Language:    "en",
    MainLibrary: "H5P.QuestionSet",
    EmbedTypes:  []string{"div"},
}

pkg.SetPackageDefinition(packageDef)
pkg.SetContent(&h5p.Content{Params: questionSet})

err := pkg.CreateZipFile("package.h5p")
if err != nil {
    log.Fatal(err)
}
```

## Type Relationships

```
H5PPackage
├── PackageDefinition (h5p.json)
├── Content (content/content.json)
│   ├── Params (QuestionSet or other content)
│   └── Metadata
└── Libraries (H5P.*/library.json + files)

QuestionSet
├── Questions[]
│   ├── Library (string)
│   └── Params (interface{})
├── OverallFeedback[]
├── IntroPage
└── EndGame
```

This type system provides a complete foundation for creating, manipulating, and validating H5P content in Go applications.