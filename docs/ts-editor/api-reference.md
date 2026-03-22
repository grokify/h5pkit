# API Reference

## Types

### H5PQuiz

The root quiz structure following H5P QuestionSet 1.20 format:

```typescript
interface H5PQuiz {
  library: 'H5P.QuestionSet 1.20';
  metadata: H5PMetadata;
  params: QuestionSetParams;
}

interface H5PMetadata {
  title: string;
  license?: string;
  authors?: Array<{ name: string; role?: string }>;
}

interface QuestionSetParams {
  progressType: 'dots' | 'textual' | 'none';
  passPercentage: number;
  questions: Question[];
  introductionText?: string;
  overallFeedback?: OverallFeedback[];
}
```

### Question Types

```typescript
type Question = MultiChoiceQuestion | TrueFalseQuestion | FillInBlanksQuestion;

// Multiple Choice
interface MultiChoiceQuestion {
  library: 'H5P.MultiChoice 1.16';
  params: {
    question: string;
    answers: Array<{ text: string; correct: boolean; tips?: string }>;
    behaviour: MultiChoiceBehaviour;
  };
}

interface MultiChoiceBehaviour {
  enableRetry?: boolean;
  enableSolutionsButton?: boolean;
  randomAnswers?: boolean;
  singleAnswer?: boolean;
}

// True/False
interface TrueFalseQuestion {
  library: 'H5P.TrueFalse 1.8';
  params: {
    question: string;
    correct: boolean;
    behaviour: TrueFalseBehaviour;
  };
}

// Fill in the Blanks
interface FillInBlanksQuestion {
  library: 'H5P.Blanks 1.14';
  params: {
    text: string;  // Use *answer* syntax for blanks
    behaviour: BlanksBehaviour;
  };
}
```

## Components

### QuizEditor

Main editor component.

```tsx
import { QuizEditor, QuizEditorRef, QuizEditorProps } from '@grokify/h5p-editor';

<QuizEditor
  ref={editorRef}
  content={initialQuiz}
  onChange={handleChange}
  onSave={handleSave}
  allowedTypes={['multipleChoice', 'trueFalse']}
  theme="light"
  readOnly={false}
  className="my-editor"
/>
```

### Content Type Editors

Individual question type editors (for custom implementations):

```tsx
import {
  MultiChoiceEditor,
  TrueFalseEditor,
  FillInBlanksEditor
} from '@grokify/h5p-editor';

<MultiChoiceEditor
  question={multiChoiceQuestion}
  index={0}
  readOnly={false}
/>
```

### EditorProvider & useEditor

Context for building custom editor UIs:

```tsx
import { EditorProvider, useEditor } from '@grokify/h5p-editor';

function CustomEditor() {
  const {
    state,           // Current editor state
    getQuiz,         // Get quiz content
    setQuiz,         // Set quiz content
    addQuestion,     // Add a question
    updateQuestion,  // Update a question
    removeQuestion,  // Remove a question
    moveQuestion,    // Reorder questions
    selectQuestion,  // Select for editing
    undo,            // Undo last action
    redo,            // Redo undone action
    validate,        // Validate quiz
  } = useEditor();

  // ...
}

<EditorProvider initialQuiz={quiz} onChange={handleChange}>
  <CustomEditor />
</EditorProvider>
```

## Functions

### validateQuiz

Standalone validation function:

```typescript
import { validateQuiz } from '@grokify/h5p-editor';

const result = validateQuiz(quizData);

interface ValidationResult {
  valid: boolean;
  errors: Array<{
    path: string;    // JSON path to error location
    message: string; // Human-readable error
  }>;
}
```

Validation checks:

- Quiz has at least one question
- Each question has required text content
- Multiple choice questions have at least 2 answers
- Multiple choice questions have at least 1 correct answer
- Fill-in-blanks text contains at least one blank (`*answer*`)

### isSupportedQuestionType

Check if a question type is supported:

```typescript
import { isSupportedQuestionType } from '@grokify/h5p-editor';

isSupportedQuestionType('H5P.MultiChoice 1.16');  // true
isSupportedQuestionType('H5P.DragDrop 1.0');      // false
```

### createQuizEditor

Create a vanilla JS editor instance:

```typescript
import { createQuizEditor } from '@grokify/h5p-editor';

const editor = createQuizEditor('#container', {
  onSave: (quiz) => console.log(quiz),
  theme: 'dark',
});

// Methods
editor.getContent();
editor.setContent(quiz);
editor.validate();
editor.isDirty();
editor.destroy();
```

## CSS Custom Properties

Available theme variables:

```css
.h5p-editor-root {
  /* Colors */
  --h5p-color-primary: #0066cc;
  --h5p-color-primary-hover: #0052a3;
  --h5p-color-bg: #ffffff;
  --h5p-color-bg-secondary: #f8fafc;
  --h5p-color-text: #1e293b;
  --h5p-color-text-secondary: #64748b;
  --h5p-color-border: #e2e8f0;
  --h5p-color-error: #dc2626;
  --h5p-color-success: #16a34a;

  /* Spacing */
  --h5p-spacing-xs: 4px;
  --h5p-spacing-sm: 8px;
  --h5p-spacing-md: 16px;
  --h5p-spacing-lg: 24px;

  /* Typography */
  --h5p-font-family: system-ui, -apple-system, sans-serif;
  --h5p-font-size-sm: 14px;
  --h5p-font-size-base: 16px;

  /* Border radius */
  --h5p-radius-sm: 4px;
  --h5p-radius-md: 8px;
}
```

## Supported H5P Libraries

| Library | Version | Description |
|---------|---------|-------------|
| `H5P.QuestionSet` | 1.20 | Quiz container |
| `H5P.MultiChoice` | 1.16 | Multiple choice questions |
| `H5P.TrueFalse` | 1.8 | True/false questions |
| `H5P.Blanks` | 1.14 | Fill in the blanks |
