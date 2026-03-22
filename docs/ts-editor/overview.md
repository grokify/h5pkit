# TypeScript Quiz Editor

The `@grokify/h5p-editor` package provides a TypeScript/React library for creating and editing H5P quiz content in the browser.

## Features

- **Multiple Question Types**: Multiple Choice, True/False, Fill in the Blanks
- **React Component**: Use as a React component with ref-based API
- **Vanilla JS API**: Use via CDN without React knowledge
- **Undo/Redo**: Full history support for editing operations
- **Validation**: Real-time validation with detailed error messages
- **Theming**: Light and dark mode support with CSS custom properties
- **H5P Compatible**: Outputs standard H5P QuestionSet 1.20 JSON format

## Quick Example

=== "React"

    ```tsx
    import { QuizEditor, QuizEditorRef, H5PQuiz } from '@grokify/h5p-editor';
    import '@grokify/h5p-editor/styles.css';
    import { useRef } from 'react';

    function App() {
      const editorRef = useRef<QuizEditorRef>(null);

      const handleSave = (quiz: H5PQuiz) => {
        console.log('Quiz saved:', quiz);
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

=== "Vanilla JS (CDN)"

    ```html
    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.css">
    <script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
    <script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.umd.cjs"></script>

    <div id="editor"></div>

    <script>
      const editor = H5PEditor.createQuizEditor('#editor', {
        onSave: (quiz) => console.log('Quiz saved:', quiz)
      });
    </script>
    ```

## Architecture

The editor is built with:

- **React 18+** for the component layer
- **TypeScript** for type safety
- **Vite** for building ES modules and UMD bundles
- **CSS Custom Properties** for theming

The package outputs:

| File | Format | Use Case |
|------|--------|----------|
| `h5p-editor.js` | ES Module | Modern bundlers (Vite, webpack, etc.) |
| `h5p-editor.umd.cjs` | UMD | CDN/script tag usage |
| `h5p-editor.css` | CSS | Styles (required) |
| `index.d.ts` | TypeScript | Type definitions |

## Relationship to Go SDK

This TypeScript editor complements the Go SDK:

| Component | Language | Purpose |
|-----------|----------|---------|
| `github.com/grokify/h5p-go` | Go | Server-side creation, validation, package management |
| `@grokify/h5p-editor` | TypeScript | Browser-based visual editing |

Both produce the same H5P QuestionSet JSON format, ensuring full interoperability.
