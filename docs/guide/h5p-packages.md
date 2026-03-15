# H5P Packages

H5P packages are complete, deployable `.h5p` files that contain everything needed to run interactive content. This guide covers creating, managing, and deploying H5P packages.

## Package Structure

An H5P package is a ZIP file with a specific structure:

```
example.h5p
├── h5p.json                    # Package definition
├── content/
│   └── content.json           # Content parameters  
└── H5P.LibraryName-1.0/       # Library directories
    ├── library.json          # Library definition
    ├── semantics.json        # Content schema
    ├── js/                   # JavaScript files
    ├── css/                  # Stylesheet files
    └── assets/               # Images, fonts, etc.
```

## Creating Packages

### Basic Package Creation

```go
import "github.com/grokify/h5p-go"

// Create a new H5P package
pkg := h5p.NewH5PPackage()

// Set package definition
packageDef := &h5p.PackageDefinition{
    Title:       "My Interactive Quiz",
    Language:    "en", 
    MainLibrary: "H5P.MultiChoice",
    EmbedTypes:  []string{"div"},
    PreloadedDependencies: []h5p.LibraryDependency{
        {
            MachineName:  "H5P.MultiChoice",
            MajorVersion: 1,
            MinorVersion: 16,
        },
    },
}

pkg.SetPackageDefinition(packageDef)
```

### Adding Content

```go
// Create content using typed schemas
params := &schemas.MultiChoiceParams{
    Question: "What is the answer?",
    Answers: []schemas.AnswerOption{
        {Text: "Correct", Correct: true},
        {Text: "Wrong", Correct: false},
    },
}

content := &h5p.Content{
    Params: params,
}

pkg.SetContent(content)
```

### Adding Libraries

```go
// Add the MultiChoice library
lib := &h5p.Library{
    MachineName: "H5P.MultiChoice-1.16",
    Definition: &h5p.LibraryDefinition{
        Title:        "Multiple Choice",
        MachineName:  "H5P.MultiChoice", 
        MajorVersion: 1,
        MinorVersion: 16,
        PatchVersion: 4,
        Runnable:     true,
        CoreAPI: &h5p.CoreAPI{
            MajorVersion: 1,
            MinorVersion: 24,
        },
    },
    Files: map[string][]byte{
        "js/multichoice.js":  loadFile("multichoice.js"),
        "css/multichoice.css": loadFile("multichoice.css"),
    },
}

pkg.AddLibrary(lib)
```

## Package Definition

The `h5p.json` file contains package metadata:

```go
type PackageDefinition struct {
    Title       string              `json:"title"`
    Language    string              `json:"language"`
    MainLibrary string              `json:"mainLibrary"`
    EmbedTypes  []string           `json:"embedTypes"`
    
    PreloadedDependencies []LibraryDependency `json:"preloadedDependencies"`
    DynamicDependencies   []LibraryDependency `json:"dynamicDependencies,omitempty"`
    EditorDependencies    []LibraryDependency `json:"editorDependencies,omitempty"`
    
    License       string `json:"license,omitempty"`
    DefaultLanguage string `json:"defaultLanguage,omitempty"`
    Authors       []Author `json:"authors,omitempty"`
    Changes       []Change `json:"changes,omitempty"`
}
```

### Example Package Definition

```go
packageDef := &h5p.PackageDefinition{
    Title:       "Advanced Programming Quiz",
    Language:    "en",
    MainLibrary: "H5P.QuestionSet",
    EmbedTypes:  []string{"div"},
    
    PreloadedDependencies: []h5p.LibraryDependency{
        {
            MachineName:  "H5P.QuestionSet",
            MajorVersion: 1,
            MinorVersion: 20,
        },
        {
            MachineName:  "H5P.MultiChoice", 
            MajorVersion: 1,
            MinorVersion: 16,
        },
    },
    
    License: "CC BY-SA",
    Authors: []h5p.Author{
        {
            Name: "Your Name",
            Role: "Author",
        },
    },
    
    Changes: []h5p.Change{
        {
            Date:        "2024-01-15",
            Author:      "Your Name", 
            Log:         "Initial version",
        },
    },
}
```

## Library Management

### Library Definition

Each library needs a `library.json` file:

```go
type LibraryDefinition struct {
    Title        string   `json:"title"`
    Description  string   `json:"description,omitempty"`
    MajorVersion int      `json:"majorVersion"`
    MinorVersion int      `json:"minorVersion"`
    PatchVersion int      `json:"patchVersion"`
    MachineName  string   `json:"machineName"`
    Runnable     bool     `json:"runnable"`
    
    CoreAPI      *CoreAPI `json:"coreApi,omitempty"`
    
    PreloadedJS  []string `json:"preloadedJs,omitempty"`
    PreloadedCSS []string `json:"preloadedCss,omitempty"`
    
    PreloadedDependencies []LibraryDependency `json:"preloadedDependencies,omitempty"`
    DynamicDependencies   []LibraryDependency `json:"dynamicDependencies,omitempty"`
    EditorDependencies    []LibraryDependency `json:"editorDependencies,omitempty"`
}
```

### Adding Library Files

```go
lib := &h5p.Library{
    MachineName: "H5P.MultiChoice-1.16",
    Definition:  libDefinition,
    Files: map[string][]byte{
        // JavaScript files
        "js/multichoice.js": loadJSFile(),
        
        // CSS files
        "css/multichoice.css": loadCSSFile(),
        
        // Assets
        "images/icon.png": loadImageFile(),
        
        // Semantics definition
        "semantics.json": loadSemanticsFile(),
    },
}
```

## Content Management

### Content Structure

The `content/content.json` file contains the actual content parameters:

```go
type Content struct {
    Params   interface{} `json:"params"`
    Metadata *Metadata   `json:"metadata,omitempty"`
}

type Metadata struct {
    Title           string `json:"title,omitempty"`
    Authors         []Author `json:"authors,omitempty"`  
    Source          string `json:"source,omitempty"`
    License         string `json:"license,omitempty"`
    LicenseVersion  string `json:"licenseVersion,omitempty"`
    YearFrom        int    `json:"yearFrom,omitempty"`
    YearTo          int    `json:"yearTo,omitempty"`
    DefaultLanguage string `json:"defaultLanguage,omitempty"`
}
```

### Setting Content

```go
content := &h5p.Content{
    Params: questionSetParams,  // Your content parameters
    Metadata: &h5p.Metadata{
        Title:   "My Quiz",
        License: "CC BY-SA",
        Authors: []h5p.Author{
            {Name: "Your Name", Role: "Author"},
        },
    },
}

pkg.SetContent(content)
```

## Package Export

### Create ZIP File

```go
// Export as .h5p file
err := pkg.CreateZipFile("my-content.h5p")
if err != nil {
    log.Fatal("Failed to create H5P file:", err)
}

fmt.Println("Successfully created my-content.h5p")
```

### Export to Directory

```go
// Export to directory structure (for development)
err := pkg.ExportToDirectory("./h5p-export")
if err != nil {
    log.Fatal("Failed to export:", err)
}
```

## Loading Packages

### Load from ZIP File

```go
// Load existing H5P package
pkg, err := h5p.LoadH5PPackage("existing-content.h5p")
if err != nil {
    log.Fatal("Failed to load package:", err)
}

// Access package components
packageDef := pkg.GetPackageDefinition()
content := pkg.GetContent()
libraries := pkg.GetLibraries()
```

### Modify Existing Package

```go
// Load and modify
pkg, err := h5p.LoadH5PPackage("quiz.h5p")
if err != nil {
    log.Fatal(err)
}

// Update title
packageDef := pkg.GetPackageDefinition()
packageDef.Title = "Updated Quiz Title"
pkg.SetPackageDefinition(packageDef)

// Save changes
err = pkg.CreateZipFile("updated-quiz.h5p")
if err != nil {
    log.Fatal(err)
}
```

## Validation

### Package Validation

```go
// Validate complete package
if err := pkg.Validate(); err != nil {
    log.Fatal("Package validation failed:", err)
}
```

### Individual Component Validation

```go
// Validate package definition
if err := packageDef.Validate(); err != nil {
    log.Fatal("Package definition invalid:", err)
}

// Validate content
if err := content.Validate(); err != nil {
    log.Fatal("Content invalid:", err)
}

// Validate libraries
for _, lib := range libraries {
    if err := lib.Validate(); err != nil {
        log.Fatal("Library invalid:", err)
    }
}
```

## Complete Example

```go
package main

import (
    "log"
    "os"
    
    "github.com/grokify/h5p-go"
    "github.com/grokify/h5p-go/schemas"
)

func main() {
    // Create a complete H5P package
    pkg := createCompletePackage()
    
    // Validate before export
    if err := pkg.Validate(); err != nil {
        log.Fatal("Validation failed:", err)
    }
    
    // Export as H5P file
    err := pkg.CreateZipFile("complete-quiz.h5p")
    if err != nil {
        log.Fatal("Export failed:", err)
    }
    
    fmt.Println("Successfully created complete-quiz.h5p")
}

func createCompletePackage() *h5p.H5PPackage {
    pkg := h5p.NewH5PPackage()
    
    // Package definition
    packageDef := &h5p.PackageDefinition{
        Title:       "Complete Programming Quiz",
        Language:    "en",
        MainLibrary: "H5P.QuestionSet",
        EmbedTypes:  []string{"div"},
        PreloadedDependencies: []h5p.LibraryDependency{
            {
                MachineName:  "H5P.QuestionSet",
                MajorVersion: 1,
                MinorVersion: 20,
            },
            {
                MachineName:  "H5P.MultiChoice",
                MajorVersion: 1, 
                MinorVersion: 16,
            },
        },
        License: "CC BY-SA",
        Authors: []h5p.Author{
            {Name: "Quiz Creator", Role: "Author"},
        },
    }
    pkg.SetPackageDefinition(packageDef)
    
    // Create question set content
    questionSet := createQuestionSet()
    content := &h5p.Content{
        Params: questionSet,
        Metadata: &h5p.Metadata{
            Title:   "Programming Quiz",
            License: "CC BY-SA",
        },
    }
    pkg.SetContent(content)
    
    // Add required libraries
    addLibraries(pkg)
    
    return pkg
}

func createQuestionSet() *h5p.QuestionSet {
    // Create questions using the builder
    answers := []h5p.Answer{
        h5p.CreateAnswer("Go", true),
        h5p.CreateAnswer("JavaScript", false),
    }
    
    questionSet, err := h5p.NewQuestionSetBuilder().
        SetTitle("Programming Knowledge Test").
        SetProgressType("textual").
        SetPassPercentage(70).
        AddMultipleChoiceQuestion("Which language is statically typed?", answers).
        Build()
    
    if err != nil {
        log.Fatal(err)
    }
    
    return questionSet
}

func addLibraries(pkg *h5p.H5PPackage) {
    // Add QuestionSet library
    questionSetLib := &h5p.Library{
        MachineName: "H5P.QuestionSet-1.20",
        Definition: &h5p.LibraryDefinition{
            Title:        "Question Set",
            MachineName:  "H5P.QuestionSet",
            MajorVersion: 1,
            MinorVersion: 20,
            PatchVersion: 0,
            Runnable:     true,
        },
        Files: map[string][]byte{
            "js/questionset.js":  loadLibraryFile("questionset.js"),
            "css/questionset.css": loadLibraryFile("questionset.css"),
            "semantics.json":     loadLibraryFile("questionset-semantics.json"),
        },
    }
    pkg.AddLibrary(questionSetLib)
    
    // Add MultiChoice library  
    multiChoiceLib := &h5p.Library{
        MachineName: "H5P.MultiChoice-1.16",
        Definition: &h5p.LibraryDefinition{
            Title:        "Multiple Choice",
            MachineName:  "H5P.MultiChoice",
            MajorVersion: 1,
            MinorVersion: 16,
            PatchVersion: 4,
            Runnable:     true,
        },
        Files: map[string][]byte{
            "js/multichoice.js":  loadLibraryFile("multichoice.js"),
            "css/multichoice.css": loadLibraryFile("multichoice.css"),
            "semantics.json":     loadLibraryFile("multichoice-semantics.json"),
        },
    }
    pkg.AddLibrary(multiChoiceLib)
}

func loadLibraryFile(filename string) []byte {
    // In practice, load actual H5P library files
    // For example: return os.ReadFile("libraries/" + filename)
    return []byte("// Placeholder library content")
}
```

## Best Practices

1. **Version Management** - Use proper semantic versioning for libraries
2. **Dependency Management** - Specify exact versions for dependencies
3. **Asset Organization** - Keep assets organized in appropriate directories
4. **Validation** - Always validate packages before deployment
5. **Licensing** - Include proper license information
6. **Documentation** - Document any custom libraries or modifications

## Deployment

H5P packages can be deployed to:
- H5P-compatible LMS platforms
- WordPress with H5P plugin
- Drupal with H5P module
- Custom web applications with H5P Core

## Next Steps

- [Validation](validation.md) - Understanding validation rules
- [Examples](../examples/packages.md) - Package creation examples
- [File Structure](../reference/file-structure.md) - Detailed file format reference