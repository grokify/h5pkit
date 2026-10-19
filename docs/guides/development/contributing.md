# Contributing

We welcome contributions to the H5P Go SDK! This guide will help you get started with contributing to the project.

## Getting Started

### Prerequisites

- Go 1.19 or later
- Git
- Basic understanding of H5P concepts
- Familiarity with Go development practices

### Development Setup

1. **Fork and clone the repository:**
   ```bash
   git clone https://github.com/yourusername/h5p-go.git
   cd h5p-go
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Run tests to ensure everything works:**
   ```bash
   go test ./...
   ```

4. **Build the project:**
   ```bash
   go build ./...
   ```

## Development Workflow

### Branch Strategy

- `main` - Stable release branch
- `develop` - Development branch (if used)
- `feature/feature-name` - Feature branches
- `fix/issue-description` - Bug fix branches
- `docs/description` - Documentation updates

### Making Changes

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Write code following Go best practices
   - Add tests for new functionality
   - Update documentation as needed
   - Follow existing code style and conventions

3. **Test your changes:**
   ```bash
   go test ./...
   go vet ./...
   gofmt -s -w .
   ```

4. **Commit your changes:**
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

### Commit Message Format

Follow conventional commit format:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation updates
- `style:` - Code style changes
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:
```
feat: add support for Essay content type
fix: resolve validation issue with empty questions
docs: update API documentation for semantics
test: add comprehensive validation tests
```

## Code Guidelines

### Go Style

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow effective Go guidelines
- Use meaningful variable and function names
- Write clear, concise comments
- Handle errors appropriately

### Package Structure

```
h5p-go/
├── *.go              # Core functionality
├── schemas/          # Typed schema implementations
├── semantics/        # H5P semantics format
├── testdata/         # Test data files
├── docsrc/          # Documentation source
└── docs/            # Generated documentation
```

### Testing

- Write unit tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Include both positive and negative test cases

Example test:
```go
func TestQuestionSetBuilder(t *testing.T) {
    tests := []struct {
        name      string
        builder   func() *QuestionSetBuilder
        wantError bool
    }{
        {
            name: "valid question set",
            builder: func() *QuestionSetBuilder {
                return NewQuestionSetBuilder().
                    SetTitle("Test Quiz").
                    SetPassPercentage(70).
                    AddMultipleChoiceQuestion("Question?", validAnswers)
            },
            wantError: false,
        },
        // Add more test cases...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := tt.builder().Build()
            if (err != nil) != tt.wantError {
                t.Errorf("Build() error = %v, wantError %v", err, tt.wantError)
            }
        })
    }
}
```

## Types of Contributions

### Code Contributions

**New Features:**
- H5P content type support
- Builder pattern improvements
- Validation enhancements
- Package management features

**Bug Fixes:**
- Validation issues
- JSON marshaling/unmarshaling problems
- Build failures
- Documentation errors

**Performance Improvements:**
- Optimization of validation logic
- Memory usage improvements
- Build time optimizations

### Documentation

**API Documentation:**
- Function and method documentation
- Code examples
- Usage patterns

**User Guides:**
- Tutorials and how-tos
- Best practices
- Integration examples

**Reference Material:**
- H5P specification compliance
- Schema documentation
- Troubleshooting guides

### Testing

**Test Coverage:**
- Unit tests for new features
- Integration tests
- Error condition testing
- Performance benchmarks

**Test Data:**
- Sample H5P content
- Validation test cases
- Edge case scenarios

## Specific Contribution Areas

### H5P Content Type Support

To add support for a new H5P content type:

1. **Add schema definition:**
   ```go
   // schemas/newtype_types.go
   type NewTypeParams struct {
       // Define parameters based on H5P semantics
   }
   
   func (p *NewTypeParams) Validate() error {
       // Implement validation logic
   }
   ```

2. **Add semantics file:**
   ```
   schemas/newtype_semantics.json
   ```

3. **Add tests:**
   ```go
   // schemas/newtype_test.go
   func TestNewTypeValidation(t *testing.T) {
       // Test validation logic
   }
   ```

4. **Update documentation:**
   - Add to API reference
   - Create usage examples
   - Update README if major feature

### Validation Improvements

When improving validation:

1. **Identify validation gaps:**
   - Missing required field checks
   - Business logic validation
   - Cross-field validation

2. **Implement validation logic:**
   ```go
   func (qs *QuestionSet) Validate() error {
       var errors ValidationErrors
       
       if qs.Title == "" {
           errors = append(errors, ValidationError{
               Field:   "title",
               Message: "title is required",
           })
       }
       
       // Add more validation...
       
       if len(errors) > 0 {
           return errors
       }
       return nil
   }
   ```

3. **Add comprehensive tests:**
   ```go
   func TestValidation(t *testing.T) {
       // Test all validation scenarios
   }
   ```

### Semantics Engine Improvements

For semantics-related contributions:

1. **Study H5P semantics specification**
2. **Identify missing field types or attributes**
3. **Implement polymorphic handling as needed**
4. **Add validation for new semantic elements**
5. **Update documentation with examples**

## Pull Request Process

### Before Submitting

1. **Ensure all tests pass:**
   ```bash
   go test ./... -v
   ```

2. **Check code formatting:**
   ```bash
   gofmt -s -w .
   ```

3. **Run static analysis:**
   ```bash
   go vet ./...
   ```

4. **Update documentation:**
   - API documentation
   - User guides if needed
   - CHANGELOG.md entry

### Submitting the Pull Request

1. **Push your branch:**
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create pull request:**
   - Use descriptive title
   - Explain what changes were made
   - Reference related issues
   - Include testing information

3. **PR Description Template:**
   ```markdown
   ## Description
   Brief description of changes made.
   
   ## Type of Change
   - [ ] Bug fix
   - [ ] New feature
   - [ ] Breaking change
   - [ ] Documentation update
   
   ## Testing
   - [ ] Unit tests added/updated
   - [ ] All tests passing
   - [ ] Manual testing completed
   
   ## Checklist
   - [ ] Code follows project style guidelines
   - [ ] Self-review completed
   - [ ] Documentation updated
   - [ ] No breaking changes (or properly documented)
   ```

### Review Process

1. **Automated checks** will run (tests, linting)
2. **Maintainer review** will be requested
3. **Address feedback** by updating your branch
4. **Approval and merge** once ready

## Development Environment

### Recommended Tools

- **IDE:** VS Code with Go extension, GoLand
- **Testing:** Go test tools, testify library
- **Linting:** golangci-lint
- **Documentation:** godoc, mkdocs

### Local Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test ./schemas/ -v

# Run benchmarks
go test -bench=. ./...
```

## Community Guidelines

### Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers get started
- Follow project maintainer guidance

### Communication

- **GitHub Issues:** Bug reports, feature requests
- **Pull Requests:** Code contributions
- **Discussions:** Questions, ideas, general discussion

### Getting Help

If you need help:

1. Check existing documentation
2. Search closed issues for similar problems
3. Create a new issue with detailed information
4. Be patient and provide additional context when requested

## Recognition

Contributors will be:
- Listed in project contributors
- Credited in release notes for significant contributions
- Invited to participate in project decisions for regular contributors

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to the H5P Go SDK! Your contributions help make H5P content creation more accessible to Go developers.