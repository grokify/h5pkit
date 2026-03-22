# @grokify/h5p-editor

H5P Quiz Editor - A TypeScript/React library for creating and editing H5P quiz content.

## Features

- **Multiple Choice** questions with multiple correct answers
- **True/False** questions
- **Fill in the Blanks** questions
- Undo/redo support
- Validation with detailed error messages
- Dark mode support
- React component or vanilla JS API

## Installation

```bash
npm install @grokify/h5p-editor
# or
pnpm add @grokify/h5p-editor
```

## Usage

### React

```tsx
import { QuizEditor, QuizEditorRef, H5PQuiz } from '@grokify/h5p-editor';
import '@grokify/h5p-editor/styles.css';
import { useRef } from 'react';

function App() {
  const editorRef = useRef<QuizEditorRef>(null);

  const handleSave = (quiz: H5PQuiz) => {
    console.log('Quiz saved:', quiz);
    // Send to your backend
  };

  return (
    <QuizEditor
      ref={editorRef}
      onSave={handleSave}
      allowedTypes={['multipleChoice', 'trueFalse', 'fillInBlanks']}
    />
  );
}
```

### CDN (Vanilla JS)

```html
<!DOCTYPE html>
<html>
<head>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.css">
</head>
<body>
  <div id="editor"></div>

  <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.umd.cjs"></script>
  <script>
    const editor = H5PEditor.createQuizEditor('#editor', {
      onSave: (quiz) => {
        console.log('Quiz saved:', quiz);
      }
    });

    // Get current content
    const quiz = editor.getContent();

    // Validate
    const { valid, errors } = editor.validate();

    // Clean up when done
    editor.destroy();
  </script>
</body>
</html>
```

## API

### QuizEditor Props

| Prop | Type | Description |
|------|------|-------------|
| `content` | `H5PQuiz` | Initial quiz content to edit |
| `onChange` | `(quiz: H5PQuiz) => void` | Called when content changes |
| `onSave` | `(quiz: H5PQuiz) => void` | Called when save button is clicked |
| `allowedTypes` | `Array<'multipleChoice' \| 'trueFalse' \| 'fillInBlanks'>` | Limit available question types |
| `theme` | `'light' \| 'dark'` | Color theme (default: 'light') |
| `readOnly` | `boolean` | Disable editing |
| `className` | `string` | Additional CSS class |

### QuizEditorRef Methods

| Method | Return | Description |
|--------|--------|-------------|
| `getContent()` | `H5PQuiz` | Get current quiz content |
| `setContent(quiz)` | `void` | Set quiz content |
| `validate()` | `ValidationResult` | Validate current content |
| `isDirty()` | `boolean` | Check if content has unsaved changes |

### Validation

```typescript
import { validateQuiz } from '@grokify/h5p-editor';

const result = validateQuiz(quizData);
if (!result.valid) {
  console.log('Errors:', result.errors);
  // [{ path: '/params/questions/0', message: 'Question text is required' }]
}
```

## H5P Quiz Format

The editor produces standard H5P QuestionSet JSON:

```json
{
  "library": "H5P.QuestionSet 1.20",
  "metadata": {
    "title": "My Quiz"
  },
  "params": {
    "progressType": "dots",
    "passPercentage": 70,
    "questions": [
      {
        "library": "H5P.MultiChoice 1.16",
        "params": {
          "question": "What is 2 + 2?",
          "answers": [
            { "text": "3", "correct": false },
            { "text": "4", "correct": true },
            { "text": "5", "correct": false }
          ]
        }
      }
    ]
  }
}
```

## Theming

The editor uses CSS custom properties for theming:

```css
.h5p-editor-root {
  --h5p-color-primary: #0066cc;
  --h5p-color-bg: #ffffff;
  --h5p-color-text: #1e293b;
  --h5p-color-border: #e2e8f0;
  /* ... */
}
```

Use `theme="dark"` for dark mode, or override CSS variables for custom themes.

## License

MIT
