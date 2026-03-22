/**
 * @grokify/h5p-editor
 *
 * H5P Quiz Editor - Create and edit H5P quiz content.
 *
 * @example React usage
 * ```tsx
 * import { QuizEditor } from '@grokify/h5p-editor';
 * import '@grokify/h5p-editor/styles.css';
 *
 * function App() {
 *   const editorRef = useRef<QuizEditorRef>(null);
 *
 *   const handleSave = (quiz: H5PQuiz) => {
 *     console.log('Saving:', quiz);
 *   };
 *
 *   return (
 *     <QuizEditor
 *       ref={editorRef}
 *       onSave={handleSave}
 *     />
 *   );
 * }
 * ```
 *
 * @example Vanilla JS (CDN)
 * ```html
 * <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/styles.css">
 * <script src="https://cdn.jsdelivr.net/npm/@grokify/h5p-editor/dist/h5p-editor.umd.cjs"></script>
 * <script>
 *   const editor = H5PEditor.createQuizEditor('#container', {
 *     onSave: (quiz) => console.log('Saving:', quiz)
 *   });
 * </script>
 * ```
 */

// React components
export { QuizEditor } from './editor/QuizEditor';
export type { QuizEditorRef, QuizEditorProps } from './editor/QuizEditor';
export { EditorProvider, useEditor } from './editor/EditorContext';

// Content type editors (for custom implementations)
export { MultiChoiceEditor } from './content-types/MultiChoiceEditor';
export { TrueFalseEditor } from './content-types/TrueFalseEditor';
export { FillInBlanksEditor } from './content-types/FillInBlanksEditor';

// Schema and types
export * from './schema/types';
export { validateQuiz, isSupportedQuestionType } from './schema/validation';

// Vanilla JS API for CDN usage
import { QuizEditor } from './editor/QuizEditor';
import type { QuizEditorRef, QuizEditorProps } from './editor/QuizEditor';
import type { H5PQuiz } from './schema/types';
import { createRoot, type Root } from 'react-dom/client';
import React, { createElement, createRef } from 'react';

interface VanillaEditorOptions extends Omit<QuizEditorProps, 'ref'> {
  onSave?: (quiz: H5PQuiz) => void;
}

interface VanillaEditorInstance {
  getContent: () => H5PQuiz;
  setContent: (quiz: H5PQuiz) => void;
  validate: () => { valid: boolean; errors: { path: string; message: string }[] };
  isDirty: () => boolean;
  destroy: () => void;
}

/**
 * Creates a quiz editor instance for vanilla JS usage.
 * Use this when not using React.
 */
export function createQuizEditor(
  container: string | HTMLElement,
  options: VanillaEditorOptions = {}
): VanillaEditorInstance {
  const element =
    typeof container === 'string'
      ? document.querySelector(container)
      : container;

  if (!element) {
    throw new Error(`Container not found: ${container}`);
  }

  const ref = createRef<QuizEditorRef>();
  let root: Root | null = null;

  // Render React component
  root = createRoot(element);
  root.render(
    createElement(QuizEditor, {
      ...options,
      ref: ref as React.RefObject<QuizEditorRef>,
    })
  );

  return {
    getContent: () => {
      if (!ref.current) throw new Error('Editor not initialized');
      return ref.current.getContent();
    },
    setContent: (quiz: H5PQuiz) => {
      if (!ref.current) throw new Error('Editor not initialized');
      ref.current.setContent(quiz);
    },
    validate: () => {
      if (!ref.current) throw new Error('Editor not initialized');
      return ref.current.validate();
    },
    isDirty: () => {
      if (!ref.current) throw new Error('Editor not initialized');
      return ref.current.isDirty();
    },
    destroy: () => {
      if (root) {
        root.unmount();
        root = null;
      }
    },
  };
}
