# Semantics API

The semantics package provides a Go implementation of the H5P semantics format, which is used to define form fields, validation rules, and UI behaviors for H5P content types.

## Overview

H5P semantics are JSON schemas that describe how content types should be configured and validated. The semantics package converts these JSON definitions into type-safe Go structs.

## Core Types

### Field

The `Field` struct represents a single field in an H5P semantics definition:

```go
type Field struct {
    Name        string      `json:"name,omitempty"`
    Type        string      `json:"type"`
    Label       string      `json:"label,omitempty"`
    Description string      `json:"description,omitempty"`
    Importance  string      `json:"importance,omitempty"`
    Optional    bool        `json:"optional,omitempty"`
    Default     interface{} `json:"default,omitempty"`
    
    // Polymorphic options field
    Options interface{} `json:"options,omitempty"`
    
    // Type-specific fields...
}
```

### SelectOption

Represents an option in a select field:

```go
type SelectOption struct {
    Value string `json:"value"`
    Label string `json:"label"`
}
```

### SemanticDefinition

The top-level semantics array:

```go
type SemanticDefinition []Field
```

## Polymorphic Options Field

The semantics implementation handles H5P's polymorphic use of the "options" field, which serves different purposes depending on the field type.

### Problem

In H5P semantics, the same JSON field name "options" is used for two completely different data structures:

- **Select fields**: Array of objects with "value" and "label" properties
- **Library fields**: Array of library names/versions as strings

### Solution

The Go implementation uses a single `Options interface{}` field with helper methods for type-safe access:

```go
// Field has a polymorphic Options field
type Field struct {
    Options interface{} `json:"options,omitempty"`
    // ...
}
```

## Helper Methods

### GetSelectOptions()

Retrieves options for select field types:

```go
func (f *Field) GetSelectOptions() []SelectOption
```

**Example:**
```go
field := Field{
    Type: "select",
    Options: []interface{}{
        map[string]interface{}{
            "value": "true",
            "label": "True",
        },
        map[string]interface{}{
            "value": "false", 
            "label": "False",
        },
    },
}

selectOptions := field.GetSelectOptions()
if selectOptions != nil {
    for _, opt := range selectOptions {
        fmt.Printf("Value: %s, Label: %s\n", opt.Value, opt.Label)
    }
}
```

### GetLibraryOptions()

Retrieves options for library field types:

```go
func (f *Field) GetLibraryOptions() []string
```

**Example:**
```go
field := Field{
    Type: "library",
    Options: []interface{}{
        "H5P.MultiChoice 1.16",
        "H5P.TrueFalse 1.8",
    },
}

libraryOptions := field.GetLibraryOptions()
if libraryOptions != nil {
    for _, lib := range libraryOptions {
        fmt.Printf("Available library: %s\n", lib)
    }
}
```

### SetSelectOptions()

Sets options for select fields:

```go
func (f *Field) SetSelectOptions(options []SelectOption)
```

**Example:**
```go
field := &Field{Type: "select"}
field.SetSelectOptions([]SelectOption{
    {Value: "small", Label: "Small"},
    {Value: "medium", Label: "Medium"}, 
    {Value: "large", Label: "Large"},
})
```

### SetLibraryOptions()

Sets options for library fields:

```go
func (f *Field) SetLibraryOptions(options []string)
```

**Example:**
```go
field := &Field{Type: "library"}
field.SetLibraryOptions([]string{
    "H5P.MultiChoice 1.16",
    "H5P.Essay 1.5",
    "H5P.TrueFalse 1.8",
})
```

## Field Types

The semantics format supports various field types:

### Basic Types
- `text` - Text input fields
- `number` - Numeric input fields  
- `boolean` - Checkbox fields
- `select` - Dropdown selection fields

### Complex Types
- `group` - Groups of related fields
- `list` - Dynamic arrays of items
- `library` - Content type selection fields

### Example Field Definitions

**Text Field:**
```go
field := Field{
    Name:        "title",
    Type:        "text",
    Label:       "Title",
    Description: "Enter a title for your content",
    Importance:  "high",
    MaxLength:   100,
}
```

**Select Field:**
```go
field := Field{
    Name:  "size",
    Type:  "select", 
    Label: "Size",
}
field.SetSelectOptions([]SelectOption{
    {Value: "small", Label: "Small"},
    {Value: "large", Label: "Large"},
})
```

**Library Field:**
```go
field := Field{
    Name:  "content",
    Type:  "library",
    Label: "Content Type",
}
field.SetLibraryOptions([]string{
    "H5P.MultiChoice 1.16",
    "H5P.TrueFalse 1.8",
})
```

## Type Validation

The helper methods include built-in validation:

- `GetSelectOptions()` returns `nil` if options are not in select format
- `GetLibraryOptions()` returns `nil` if options are not in library format
- Methods handle both JSON unmarshaling (`[]interface{}`) and direct assignment scenarios

```go
// Safe usage pattern
field := loadFieldFromJSON(data)

if selectOptions := field.GetSelectOptions(); selectOptions != nil {
    // Handle as select field
    processSelectField(selectOptions)
} else if libraryOptions := field.GetLibraryOptions(); libraryOptions != nil {
    // Handle as library field
    processLibraryField(libraryOptions)
} else {
    // Handle other field types or invalid options
    handleOtherFieldType(field)
}
```

## Advanced Usage

### Loading from JSON

```go
func loadSemantics(jsonData []byte) (SemanticDefinition, error) {
    var semantics SemanticDefinition
    err := json.Unmarshal(jsonData, &semantics)
    if err != nil {
        return nil, err
    }
    return semantics, nil
}
```

### Iterating Through Fields

```go
func processSemantics(semantics SemanticDefinition) {
    for _, field := range semantics {
        fmt.Printf("Processing field: %s (type: %s)\n", field.Name, field.Type)
        
        switch field.Type {
        case "select":
            if opts := field.GetSelectOptions(); opts != nil {
                fmt.Printf("  Select options: %d choices\n", len(opts))
            }
        case "library":
            if libs := field.GetLibraryOptions(); libs != nil {
                fmt.Printf("  Library options: %d libraries\n", len(libs))
            }
        case "group":
            if field.Fields != nil {
                fmt.Printf("  Group with %d sub-fields\n", len(field.Fields))
            }
        }
    }
}
```

## Standards Compliance

The semantics implementation follows the official H5P semantics specification:

- [H5P Semantics Documentation](https://h5p.org/semantics)
- [H5P Content Type Development](https://h5p.org/library-development)

The implementation handles all standard H5P semantic field types and attributes while providing type safety and validation in Go.