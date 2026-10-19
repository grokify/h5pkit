# Installation

## Package Manager

Install via npm or pnpm:

=== "pnpm"

    ```bash
    pnpm add @grokify/h5p-editor
    ```

=== "npm"

    ```bash
    npm install @grokify/h5p-editor
    ```

=== "yarn"

    ```bash
    yarn add @grokify/h5p-editor
    ```

### Peer Dependencies

The editor requires React 18+:

```bash
pnpm add react react-dom
```

## CDN

For vanilla JavaScript usage without a build step:

```html
<!-- Styles (required) -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.css">

<!-- React (peer dependency) -->
<script src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
<script src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>

<!-- Editor -->
<script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.umd.cjs"></script>
```

The editor is available as `window.H5PEditor` when loaded via CDN.

## TypeScript

The package includes TypeScript definitions. No additional `@types/` package needed.

```typescript
import type {
  H5PQuiz,
  QuizEditorRef,
  QuizEditorProps,
  ValidationResult
} from '@grokify/h5p-editor';
```

## Local Development

For local development without publishing to npm, use pnpm link:

```bash
# In the h5p-go/ts directory
pnpm install
pnpm build

# In your consuming project
pnpm link /path/to/h5p-go/ts/packages/editor
```

## Verify Installation

```typescript
import { QuizEditor, validateQuiz } from '@grokify/h5p-editor';

// Check exports are available
console.log('QuizEditor:', typeof QuizEditor);  // 'function'
console.log('validateQuiz:', typeof validateQuiz);  // 'function'
```
