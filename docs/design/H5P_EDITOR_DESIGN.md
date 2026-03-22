# H5P Editor Design Document

## Overview

Design for `@grokify/h5p-editor`, a TypeScript-based H5P quiz editor distributed as an npm package and CDN bundle.

## Scope

### In Scope (v1.0)

- **Quiz content types only:**
  - Multiple Choice
  - True/False
  - Fill in the Blanks
  - Question Set (collection of questions)

### Out of Scope (Future)

- Interactive Video
- Drag and Drop
- Course Presentation
- Other H5P content types

---

## Package Structure

```
@grokify/h5p-editor/
├── package.json
├── tsconfig.json
├── vite.config.ts
├── src/
│   ├── index.ts              # Main entry point
│   ├── editor/
│   │   ├── H5PEditor.tsx     # Main editor component
│   │   ├── EditorContext.tsx # State management
│   │   └── EditorToolbar.tsx
│   ├── content-types/
│   │   ├── index.ts
│   │   ├── MultipleChoice.tsx
│   │   ├── TrueFalse.tsx
│   │   ├── FillInBlanks.tsx
│   │   └── QuestionSet.tsx
│   ├── schema/
│   │   ├── h5p-quiz.schema.json
│   │   └── types.ts          # TypeScript types from schema
│   ├── utils/
│   │   ├── validation.ts
│   │   └── export.ts
│   └── styles/
│       └── editor.css
├── dist/
│   ├── h5p-editor.js         # UMD bundle (CDN)
│   ├── h5p-editor.esm.js     # ES module
│   ├── h5p-editor.css
│   └── index.d.ts            # Type definitions
└── examples/
    ├── react-app/
    ├── vanilla-js/
    └── cdn-embed/
```

---

## API Design

### CDN Usage (Vanilla JS)

```html
<!-- Load from CDN -->
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor@1/dist/h5p-editor.css">
<script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor@1/dist/h5p-editor.js"></script>

<!-- Container -->
<div id="quiz-editor"></div>

<script>
  const editor = new H5PEditor.QuizEditor('#quiz-editor', {
    // Optional: Load existing quiz
    content: existingQuizJson,

    // Callbacks
    onChange: (quiz) => console.log('Quiz changed:', quiz),
    onSave: (quiz) => saveToBackend(quiz),

    // Options
    allowedTypes: ['multipleChoice', 'trueFalse', 'fillInBlanks'],
    theme: 'light', // or 'dark'
    locale: 'en',
  });

  // Programmatic access
  const quizJson = editor.getContent();
  editor.setContent(newQuizJson);
  editor.validate(); // Returns { valid: boolean, errors: [] }
  editor.destroy();
</script>
```

### React Usage

```tsx
import { QuizEditor, QuizEditorRef } from '@grokify/h5p-editor';
import '@grokify/h5p-editor/dist/h5p-editor.css';

function MyQuizEditor() {
  const editorRef = useRef<QuizEditorRef>(null);
  const [quiz, setQuiz] = useState<H5PQuiz | null>(null);

  const handleSave = async () => {
    const content = editorRef.current?.getContent();
    if (content) {
      await api.saveQuiz(content);
    }
  };

  return (
    <div>
      <QuizEditor
        ref={editorRef}
        content={quiz}
        onChange={setQuiz}
        allowedTypes={['multipleChoice', 'trueFalse']}
      />
      <button onClick={handleSave}>Save Quiz</button>
    </div>
  );
}
```

---

## H5P Quiz JSON Schema

### Question Set Structure

```json
{
  "$schema": "https://h5p.org/schemas/quiz.json",
  "library": "H5P.QuestionSet 1.20",
  "metadata": {
    "title": "Chapter 1 Quiz",
    "license": "U",
    "authors": [],
    "changes": []
  },
  "params": {
    "introPage": {
      "showIntroPage": true,
      "title": "Chapter 1 Quiz",
      "introduction": "Test your knowledge of Chapter 1."
    },
    "progressType": "dots",
    "passPercentage": 70,
    "showResultPage": true,
    "questions": [
      {
        "library": "H5P.MultiChoice 1.16",
        "params": {
          "question": "What is the capital of France?",
          "answers": [
            { "text": "London", "correct": false },
            { "text": "Paris", "correct": true },
            { "text": "Berlin", "correct": false }
          ],
          "behaviour": {
            "singleAnswer": true,
            "enableRetry": true,
            "showSolutionsRequiresInput": true
          }
        }
      },
      {
        "library": "H5P.TrueFalse 1.8",
        "params": {
          "question": "The Earth is flat.",
          "correct": false,
          "behaviour": {
            "enableRetry": true
          }
        }
      }
    ],
    "endGame": {
      "showResultPage": true,
      "showSolutionButton": true,
      "showRetryButton": true,
      "message": "Your score: @score / @total"
    }
  }
}
```

### TypeScript Types

```typescript
// Generated from JSON Schema

export interface H5PQuiz {
  library: 'H5P.QuestionSet 1.20';
  metadata: H5PMetadata;
  params: QuestionSetParams;
}

export interface H5PMetadata {
  title: string;
  license?: string;
  authors?: Author[];
  changes?: Change[];
}

export interface QuestionSetParams {
  introPage?: IntroPage;
  progressType: 'dots' | 'text' | 'none';
  passPercentage: number;
  showResultPage: boolean;
  questions: Question[];
  endGame: EndGame;
}

export type Question =
  | MultiChoiceQuestion
  | TrueFalseQuestion
  | FillInBlanksQuestion;

export interface MultiChoiceQuestion {
  library: 'H5P.MultiChoice 1.16';
  params: {
    question: string;
    answers: Answer[];
    behaviour: MultiChoiceBehaviour;
  };
}

export interface Answer {
  text: string;
  correct: boolean;
  tipsAndFeedback?: {
    tip?: string;
    chosenFeedback?: string;
    notChosenFeedback?: string;
  };
}

// ... more types
```

---

## Editor Components

### QuizEditor (Main Component)

```tsx
interface QuizEditorProps {
  content?: H5PQuiz;
  onChange?: (quiz: H5PQuiz) => void;
  onSave?: (quiz: H5PQuiz) => void;
  allowedTypes?: QuestionType[];
  theme?: 'light' | 'dark';
  locale?: string;
  readOnly?: boolean;
}

// Internal state
interface EditorState {
  quiz: H5PQuiz;
  selectedQuestion: number | null;
  isDirty: boolean;
  validationErrors: ValidationError[];
}
```

### Question Type Components

Each question type has:

1. **Editor** - For creating/editing questions
2. **Preview** - For previewing how it will look
3. **Validator** - For validating the question

```tsx
// MultipleChoice/Editor.tsx
export function MultiChoiceEditor({
  question,
  onChange
}: {
  question: MultiChoiceQuestion;
  onChange: (q: MultiChoiceQuestion) => void;
}) {
  return (
    <div className="h5p-multichoice-editor">
      <TextEditor
        label="Question"
        value={question.params.question}
        onChange={(text) => onChange({
          ...question,
          params: { ...question.params, question: text }
        })}
      />
      <AnswerList
        answers={question.params.answers}
        onChange={(answers) => onChange({
          ...question,
          params: { ...question.params, answers }
        })}
      />
      <BehaviourSettings
        behaviour={question.params.behaviour}
        onChange={(behaviour) => onChange({
          ...question,
          params: { ...question.params, behaviour }
        })}
      />
    </div>
  );
}
```

---

## Validation

### Client-Side Validation

```typescript
import Ajv from 'ajv';
import schema from './schema/h5p-quiz.schema.json';

const ajv = new Ajv({ allErrors: true });
const validate = ajv.compile(schema);

export function validateQuiz(quiz: unknown): ValidationResult {
  const valid = validate(quiz);

  if (!valid) {
    return {
      valid: false,
      errors: validate.errors?.map(err => ({
        path: err.instancePath,
        message: err.message || 'Invalid',
        keyword: err.keyword,
      })) || [],
    };
  }

  // Additional semantic validation
  const semanticErrors = validateSemantics(quiz as H5PQuiz);

  return {
    valid: semanticErrors.length === 0,
    errors: semanticErrors,
  };
}

function validateSemantics(quiz: H5PQuiz): ValidationError[] {
  const errors: ValidationError[] = [];

  // Check at least one question
  if (quiz.params.questions.length === 0) {
    errors.push({
      path: '/params/questions',
      message: 'Quiz must have at least one question',
    });
  }

  // Check each question has a correct answer
  quiz.params.questions.forEach((q, i) => {
    if (q.library === 'H5P.MultiChoice 1.16') {
      const hasCorrect = q.params.answers.some(a => a.correct);
      if (!hasCorrect) {
        errors.push({
          path: `/params/questions/${i}/params/answers`,
          message: 'Multiple choice question must have at least one correct answer',
        });
      }
    }
  });

  return errors;
}
```

---

## Build Configuration

### Vite Config

```typescript
// vite.config.ts
import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import dts from 'vite-plugin-dts';

export default defineConfig({
  plugins: [
    react(),
    dts({ include: ['src'] }),
  ],
  build: {
    lib: {
      entry: 'src/index.ts',
      name: 'H5PEditor',
      formats: ['es', 'umd'],
      fileName: (format) => `h5p-editor.${format === 'es' ? 'esm' : format}.js`,
    },
    rollupOptions: {
      external: ['react', 'react-dom'],
      output: {
        globals: {
          react: 'React',
          'react-dom': 'ReactDOM',
        },
      },
    },
  },
});
```

### Package.json

```json
{
  "name": "@grokify/h5p-editor",
  "version": "1.0.0",
  "description": "H5P Quiz Editor - Create and edit H5P quiz content",
  "main": "dist/h5p-editor.umd.js",
  "module": "dist/h5p-editor.esm.js",
  "types": "dist/index.d.ts",
  "exports": {
    ".": {
      "import": "./dist/h5p-editor.esm.js",
      "require": "./dist/h5p-editor.umd.js",
      "types": "./dist/index.d.ts"
    },
    "./dist/h5p-editor.css": "./dist/h5p-editor.css"
  },
  "files": ["dist"],
  "scripts": {
    "dev": "vite",
    "build": "tsc && vite build",
    "test": "vitest",
    "lint": "eslint src"
  },
  "peerDependencies": {
    "react": ">=18.0.0",
    "react-dom": ">=18.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "typescript": "^5.0.0",
    "vite": "^5.0.0",
    "vitest": "^1.0.0"
  },
  "dependencies": {
    "ajv": "^8.12.0"
  }
}
```

---

## Integration with AcademyOS

### QuizEditor Page Update

```tsx
// academyos-ui/src/pages/QuizEditor.tsx

import { useEffect, useRef } from 'react';

// Dynamic import for CDN or npm
const loadH5PEditor = async () => {
  if (window.H5PEditor) return window.H5PEditor;

  // Load from CDN
  await Promise.all([
    loadScript('https://cdn.jsdelivr.net/npm/@grokify/h5p-editor@1/dist/h5p-editor.js'),
    loadCSS('https://cdn.jsdelivr.net/npm/@grokify/h5p-editor@1/dist/h5p-editor.css'),
  ]);

  return window.H5PEditor;
};

export function QuizEditor() {
  const containerRef = useRef<HTMLDivElement>(null);
  const editorRef = useRef<H5PEditorInstance | null>(null);

  useEffect(() => {
    let mounted = true;

    loadH5PEditor().then((H5PEditor) => {
      if (!mounted || !containerRef.current) return;

      editorRef.current = new H5PEditor.QuizEditor(containerRef.current, {
        content: existingQuiz?.h5p_content ? JSON.parse(existingQuiz.h5p_content) : undefined,
        onChange: (quiz) => setH5pContent(JSON.stringify(quiz)),
        allowedTypes: ['multipleChoice', 'trueFalse', 'fillInBlanks'],
      });
    });

    return () => {
      mounted = false;
      editorRef.current?.destroy();
    };
  }, []);

  return (
    <div>
      <div ref={containerRef} className="h5p-editor-container" />
      <button onClick={handleSave}>Save Quiz</button>
    </div>
  );
}
```

---

## Roadmap

### v1.0.0

- [ ] Multiple Choice questions
- [ ] True/False questions
- [ ] Fill in the Blanks questions
- [ ] Question Set container
- [ ] JSON Schema validation
- [ ] CDN bundle
- [ ] React component
- [ ] Basic theming (light/dark)

### v1.1.0

- [ ] Question bank / library
- [ ] Import from existing H5P files
- [ ] Export to H5P package (.h5p)
- [ ] Localization (i18n)

### v2.0.0

- [ ] Interactive Video support
- [ ] Drag and Drop support
- [ ] Course Presentation support
- [ ] Rich text editor (images, formatting)
