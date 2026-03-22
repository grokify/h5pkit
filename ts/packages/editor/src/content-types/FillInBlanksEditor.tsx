import { useCallback } from 'react';
import { useEditor } from '../editor/EditorContext';
import type { FillInBlanksQuestion } from '../schema/types';

interface FillInBlanksEditorProps {
  question: FillInBlanksQuestion;
  index: number;
  readOnly?: boolean;
}

export function FillInBlanksEditor({
  question,
  index,
  readOnly,
}: FillInBlanksEditorProps) {
  const { updateQuestion } = useEditor();

  const update = useCallback(
    (updates: Partial<FillInBlanksQuestion['params']>) => {
      updateQuestion(index, {
        ...question,
        params: { ...question.params, ...updates },
      });
    },
    [updateQuestion, index, question]
  );

  const handleTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    update({ text: e.target.value });
  };

  const handleBehaviourChange = (
    key: keyof NonNullable<FillInBlanksQuestion['params']['behaviour']>,
    value: boolean
  ) => {
    update({
      behaviour: { ...question.params.behaviour, [key]: value },
    });
  };

  // Count blanks in text (words surrounded by asterisks)
  const blankCount = (question.params.text.match(/\*[^*]+\*/g) || []).length;

  return (
    <div className="h5p-question-editor h5p-fillblanks-editor">
      <h3 className="h5p-question-editor-header">
        <span className="h5p-question-editor-icon">___</span>
        Fill in the Blanks
      </h3>

      {/* Instructions */}
      <div className="h5p-fillblanks-instructions">
        <p>
          Mark blank words with <code>*asterisks*</code>. For multiple accepted
          answers, use <code>*answer1/answer2*</code>.
        </p>
        <p>
          Example: <code>The capital of France is *Paris*.</code>
        </p>
      </div>

      {/* Text with blanks */}
      <div className="h5p-question-editor-field">
        <label htmlFor={`text-${index}`}>Text with Blanks</label>
        <textarea
          id={`text-${index}`}
          value={question.params.text}
          onChange={handleTextChange}
          disabled={readOnly}
          className="h5p-editor-textarea h5p-editor-textarea-large"
          placeholder="Enter text with *blanks* marked by asterisks..."
          rows={6}
        />
        <div className="h5p-fillblanks-blank-count">
          {blankCount} blank{blankCount !== 1 ? 's' : ''} detected
        </div>
      </div>

      {/* Preview */}
      {question.params.text && (
        <div className="h5p-question-editor-field">
          <label>Preview</label>
          <div className="h5p-fillblanks-preview">
            {renderPreview(question.params.text)}
          </div>
        </div>
      )}

      {/* Behaviour Settings */}
      <div className="h5p-question-editor-field">
        <label>Behaviour</label>
        <div className="h5p-question-editor-behaviour">
          <label className="h5p-editor-checkbox">
            <input
              type="checkbox"
              checked={question.params.behaviour?.caseSensitive ?? false}
              onChange={(e) =>
                handleBehaviourChange('caseSensitive', e.target.checked)
              }
              disabled={readOnly}
            />
            Case sensitive
          </label>
          <label className="h5p-editor-checkbox">
            <input
              type="checkbox"
              checked={question.params.behaviour?.enableRetry ?? true}
              onChange={(e) =>
                handleBehaviourChange('enableRetry', e.target.checked)
              }
              disabled={readOnly}
            />
            Allow retry
          </label>
          <label className="h5p-editor-checkbox">
            <input
              type="checkbox"
              checked={question.params.behaviour?.enableSolutionsButton ?? true}
              onChange={(e) =>
                handleBehaviourChange('enableSolutionsButton', e.target.checked)
              }
              disabled={readOnly}
            />
            Show solution button
          </label>
          <label className="h5p-editor-checkbox">
            <input
              type="checkbox"
              checked={question.params.behaviour?.acceptSpellingErrors ?? false}
              onChange={(e) =>
                handleBehaviourChange('acceptSpellingErrors', e.target.checked)
              }
              disabled={readOnly}
            />
            Accept minor spelling errors
          </label>
        </div>
      </div>
    </div>
  );
}

/**
 * Renders a preview of the fill-in-blanks text with blanks shown as input fields.
 */
function renderPreview(text: string): React.ReactNode {
  const parts: React.ReactNode[] = [];
  let lastIndex = 0;
  const regex = /\*([^*]+)\*/g;
  let match;

  while ((match = regex.exec(text)) !== null) {
    // Add text before the blank
    if (match.index > lastIndex) {
      parts.push(
        <span key={`text-${lastIndex}`}>
          {text.slice(lastIndex, match.index)}
        </span>
      );
    }

    // Add the blank (show first answer option as placeholder)
    const answers = match[1].split('/');
    parts.push(
      <span key={`blank-${match.index}`} className="h5p-fillblanks-preview-blank">
        {answers[0]}
      </span>
    );

    lastIndex = match.index + match[0].length;
  }

  // Add remaining text
  if (lastIndex < text.length) {
    parts.push(<span key={`text-${lastIndex}`}>{text.slice(lastIndex)}</span>);
  }

  return parts;
}
