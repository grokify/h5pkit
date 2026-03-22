# Vanilla JS Usage

Use the editor without a React build setup via CDN.

## Setup

Include the required scripts in your HTML:

```html
<!DOCTYPE html>
<html>
<head>
  <!-- Editor styles -->
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.css">
</head>
<body>
  <!-- Container for the editor -->
  <div id="editor"></div>

  <!-- React (peer dependency) -->
  <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>

  <!-- Editor library -->
  <script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.umd.cjs"></script>

  <script>
    // Your code here
  </script>
</body>
</html>
```

## Creating an Editor

```javascript
const editor = H5PEditor.createQuizEditor('#editor', {
  onSave: (quiz) => {
    console.log('Quiz saved:', quiz);
  }
});
```

### With a DOM Element

```javascript
const container = document.getElementById('editor');
const editor = H5PEditor.createQuizEditor(container, {
  onSave: (quiz) => console.log(quiz)
});
```

## Editor Methods

### Get Content

```javascript
const quiz = editor.getContent();
console.log(JSON.stringify(quiz, null, 2));
```

### Set Content

```javascript
editor.setContent({
  library: 'H5P.QuestionSet 1.20',
  metadata: { title: 'My Quiz' },
  params: {
    progressType: 'dots',
    passPercentage: 70,
    questions: []
  }
});
```

### Validate

```javascript
const result = editor.validate();
if (!result.valid) {
  result.errors.forEach(error => {
    console.error(`${error.path}: ${error.message}`);
  });
}
```

### Check for Changes

```javascript
if (editor.isDirty()) {
  if (confirm('You have unsaved changes. Discard?')) {
    // proceed
  }
}
```

### Destroy

Clean up when removing the editor:

```javascript
editor.destroy();
```

## Complete Example

```html
<!DOCTYPE html>
<html>
<head>
  <title>Quiz Editor</title>
  <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.css">
  <style>
    body { font-family: system-ui, sans-serif; margin: 0; padding: 20px; }
    .toolbar { margin-bottom: 16px; }
    .toolbar button { margin-right: 8px; padding: 8px 16px; }
    #editor { border: 1px solid #e2e8f0; border-radius: 8px; }
  </style>
</head>
<body>
  <div class="toolbar">
    <button id="save-btn">Save Quiz</button>
    <button id="export-btn">Export JSON</button>
    <button id="validate-btn">Validate</button>
  </div>

  <div id="editor"></div>

  <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
  <script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.umd.cjs"></script>

  <script>
    // Create editor
    const editor = H5PEditor.createQuizEditor('#editor', {
      theme: 'light',
      onSave: (quiz) => {
        // Send to backend
        fetch('/api/quizzes', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(quiz)
        })
        .then(res => res.json())
        .then(data => alert('Quiz saved! ID: ' + data.id))
        .catch(err => alert('Save failed: ' + err.message));
      }
    });

    // Save button
    document.getElementById('save-btn').addEventListener('click', () => {
      const result = editor.validate();
      if (result.valid) {
        // This triggers onSave callback
        const quiz = editor.getContent();
        fetch('/api/quizzes', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(quiz)
        });
      } else {
        alert('Please fix validation errors first');
      }
    });

    // Export button
    document.getElementById('export-btn').addEventListener('click', () => {
      const quiz = editor.getContent();
      const json = JSON.stringify(quiz, null, 2);
      const blob = new Blob([json], { type: 'application/json' });
      const url = URL.createObjectURL(blob);

      const a = document.createElement('a');
      a.href = url;
      a.download = 'quiz.json';
      a.click();

      URL.revokeObjectURL(url);
    });

    // Validate button
    document.getElementById('validate-btn').addEventListener('click', () => {
      const result = editor.validate();
      if (result.valid) {
        alert('Quiz is valid!');
      } else {
        const messages = result.errors.map(e => `${e.path}: ${e.message}`).join('\n');
        alert('Validation errors:\n' + messages);
      }
    });

    // Warn before leaving with unsaved changes
    window.addEventListener('beforeunload', (e) => {
      if (editor.isDirty()) {
        e.preventDefault();
        e.returnValue = '';
      }
    });
  </script>
</body>
</html>
```

## Options Reference

```typescript
interface VanillaEditorOptions {
  // Initial quiz content
  content?: H5PQuiz;

  // Called when content changes
  onChange?: (quiz: H5PQuiz) => void;

  // Called when save button is clicked
  onSave?: (quiz: H5PQuiz) => void;

  // Allowed question types
  allowedTypes?: Array<'multipleChoice' | 'trueFalse' | 'fillInBlanks'>;

  // Theme: 'light' or 'dark'
  theme?: 'light' | 'dark';

  // Disable editing
  readOnly?: boolean;
}
```

## Instance Methods

```typescript
interface VanillaEditorInstance {
  // Get current quiz content
  getContent(): H5PQuiz;

  // Set quiz content
  setContent(quiz: H5PQuiz): void;

  // Validate current content
  validate(): { valid: boolean; errors: { path: string; message: string }[] };

  // Check for unsaved changes
  isDirty(): boolean;

  // Clean up and unmount
  destroy(): void;
}
```
