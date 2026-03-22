import { useCallback, useImperativeHandle, forwardRef } from 'react';
import { EditorProvider, useEditor } from './EditorContext';
import { QuestionList } from './QuestionList';
import { QuestionEditor } from './QuestionEditor';
import { EditorToolbar } from './EditorToolbar';
import { QuizSettings } from './QuizSettings';
import type { H5PQuiz, ValidationResult } from '../schema/types';
import '../styles/editor.css';

// =============================================================================
// Public API (Ref)
// =============================================================================

export interface QuizEditorRef {
  getContent: () => H5PQuiz;
  setContent: (quiz: H5PQuiz) => void;
  validate: () => ValidationResult;
  isDirty: () => boolean;
}

// =============================================================================
// Props
// =============================================================================

export interface QuizEditorProps {
  /** Initial quiz content to edit */
  content?: H5PQuiz;
  /** Called when quiz content changes */
  onChange?: (quiz: H5PQuiz) => void;
  /** Called when save is requested (toolbar button) */
  onSave?: (quiz: H5PQuiz) => void;
  /** Question types to allow */
  allowedTypes?: Array<'multipleChoice' | 'trueFalse' | 'fillInBlanks'>;
  /** Theme */
  theme?: 'light' | 'dark';
  /** Read-only mode */
  readOnly?: boolean;
  /** Custom class name */
  className?: string;
}

// =============================================================================
// Inner Component (uses context)
// =============================================================================

interface InnerEditorProps {
  onSave?: (quiz: H5PQuiz) => void;
  allowedTypes?: Array<'multipleChoice' | 'trueFalse' | 'fillInBlanks'>;
  readOnly?: boolean;
}

function InnerEditor({ onSave, allowedTypes, readOnly }: InnerEditorProps) {
  const { state, getQuiz, validate } = useEditor();
  const { selectedQuestionIndex } = state;

  const handleSave = useCallback(() => {
    const quiz = getQuiz();
    const result = validate();
    if (result.valid) {
      onSave?.(quiz);
    }
  }, [getQuiz, validate, onSave]);

  const selectedQuestion =
    selectedQuestionIndex !== null
      ? state.quiz.params.questions[selectedQuestionIndex]
      : null;

  return (
    <div className="h5p-editor">
      <EditorToolbar onSave={onSave ? handleSave : undefined} readOnly={readOnly} />

      <div className="h5p-editor-layout">
        {/* Left Panel: Question List */}
        <div className="h5p-editor-sidebar">
          <QuizSettings readOnly={readOnly} />
          <QuestionList
            allowedTypes={allowedTypes}
            readOnly={readOnly}
          />
        </div>

        {/* Right Panel: Question Editor */}
        <div className="h5p-editor-main">
          {selectedQuestion ? (
            <QuestionEditor
              question={selectedQuestion}
              index={selectedQuestionIndex!}
              readOnly={readOnly}
            />
          ) : (
            <div className="h5p-editor-empty">
              <p>Select a question to edit, or add a new question.</p>
            </div>
          )}
        </div>
      </div>

      {/* Validation Errors */}
      {!state.validationResult.valid && (
        <div className="h5p-editor-errors">
          <h4>Validation Errors</h4>
          <ul>
            {state.validationResult.errors.map((error, i) => (
              <li key={i}>
                <code>{error.path}</code>: {error.message}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}

// =============================================================================
// Main Component
// =============================================================================

export const QuizEditor = forwardRef<QuizEditorRef, QuizEditorProps>(
  function QuizEditor(
    { content, onChange, onSave, allowedTypes, theme = 'light', readOnly, className },
    ref
  ) {
    return (
      <div className={`h5p-editor-root h5p-theme-${theme} ${className || ''}`}>
        <EditorProvider initialQuiz={content} onChange={onChange}>
          <QuizEditorInner ref={ref} onSave={onSave} allowedTypes={allowedTypes} readOnly={readOnly} />
        </EditorProvider>
      </div>
    );
  }
);

// Inner component that can use useEditor and forward ref
const QuizEditorInner = forwardRef<QuizEditorRef, InnerEditorProps>(
  function QuizEditorInner({ onSave, allowedTypes, readOnly }, ref) {
    const { getQuiz, setQuiz, validate, state } = useEditor();

    useImperativeHandle(ref, () => ({
      getContent: getQuiz,
      setContent: setQuiz,
      validate,
      isDirty: () => state.isDirty,
    }), [getQuiz, setQuiz, validate, state.isDirty]);

    return (
      <InnerEditor
        onSave={onSave}
        allowedTypes={allowedTypes}
        readOnly={readOnly}
      />
    );
  }
);
