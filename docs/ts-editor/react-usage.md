# React Usage

## Basic Usage

```tsx
import { QuizEditor } from '@grokify/h5p-editor';
import '@grokify/h5p-editor/styles.css';

function App() {
  return <QuizEditor />;
}
```

## With Save Handler

```tsx
import { QuizEditor, H5PQuiz } from '@grokify/h5p-editor';
import '@grokify/h5p-editor/styles.css';

function App() {
  const handleSave = (quiz: H5PQuiz) => {
    // Send to your backend
    fetch('/api/quizzes', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(quiz),
    });
  };

  return <QuizEditor onSave={handleSave} />;
}
```

## Ref API for Programmatic Control

Use a ref to access editor methods:

```tsx
import { useRef } from 'react';
import { QuizEditor, QuizEditorRef } from '@grokify/h5p-editor';
import '@grokify/h5p-editor/styles.css';

function App() {
  const editorRef = useRef<QuizEditorRef>(null);

  const handleExport = () => {
    const quiz = editorRef.current?.getContent();
    if (quiz) {
      const blob = new Blob([JSON.stringify(quiz, null, 2)], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      // Download logic...
    }
  };

  const handleValidate = () => {
    const result = editorRef.current?.validate();
    if (result && !result.valid) {
      console.log('Errors:', result.errors);
    }
  };

  return (
    <div>
      <button onClick={handleExport}>Export JSON</button>
      <button onClick={handleValidate}>Validate</button>
      <QuizEditor ref={editorRef} />
    </div>
  );
}
```

## Loading Existing Content

```tsx
import { QuizEditor, H5PQuiz } from '@grokify/h5p-editor';

const existingQuiz: H5PQuiz = {
  library: 'H5P.QuestionSet 1.20',
  metadata: { title: 'My Quiz' },
  params: {
    progressType: 'dots',
    passPercentage: 70,
    questions: [
      {
        library: 'H5P.MultiChoice 1.16',
        params: {
          question: 'What is 2 + 2?',
          answers: [
            { text: '3', correct: false },
            { text: '4', correct: true },
          ],
          behaviour: {},
        },
      },
    ],
  },
};

function App() {
  return <QuizEditor content={existingQuiz} />;
}
```

## Restricting Question Types

Limit which question types users can add:

```tsx
<QuizEditor
  allowedTypes={['multipleChoice', 'trueFalse']}  // No fill-in-blanks
/>
```

## Theming

### Built-in Themes

```tsx
<QuizEditor theme="light" />  // Default
<QuizEditor theme="dark" />
```

### Custom Themes via CSS

Override CSS custom properties:

```css
.h5p-editor-root {
  --h5p-color-primary: #0066cc;
  --h5p-color-bg: #ffffff;
  --h5p-color-text: #1e293b;
  --h5p-color-border: #e2e8f0;
  --h5p-color-error: #dc2626;
  --h5p-color-success: #16a34a;
}
```

## Read-Only Mode

Display quiz content without editing:

```tsx
<QuizEditor content={quiz} readOnly />
```

## Change Tracking

Track changes in real-time:

```tsx
function App() {
  const editorRef = useRef<QuizEditorRef>(null);
  const [hasChanges, setHasChanges] = useState(false);

  const handleChange = (quiz: H5PQuiz) => {
    setHasChanges(editorRef.current?.isDirty() ?? false);
  };

  return (
    <div>
      {hasChanges && <span>Unsaved changes</span>}
      <QuizEditor ref={editorRef} onChange={handleChange} />
    </div>
  );
}
```

## Props Reference

| Prop | Type | Default | Description |
|------|------|---------|-------------|
| `content` | `H5PQuiz` | Empty quiz | Initial content to edit |
| `onChange` | `(quiz: H5PQuiz) => void` | - | Called on every change |
| `onSave` | `(quiz: H5PQuiz) => void` | - | Called when save button clicked |
| `allowedTypes` | `('multipleChoice' \| 'trueFalse' \| 'fillInBlanks')[]` | All | Limit question types |
| `theme` | `'light' \| 'dark'` | `'light'` | Color theme |
| `readOnly` | `boolean` | `false` | Disable editing |
| `className` | `string` | - | Additional CSS class |

## Ref Methods

| Method | Returns | Description |
|--------|---------|-------------|
| `getContent()` | `H5PQuiz` | Get current quiz content |
| `setContent(quiz)` | `void` | Replace quiz content |
| `validate()` | `ValidationResult` | Validate and return errors |
| `isDirty()` | `boolean` | Check for unsaved changes |
