import { useCallback } from 'react';
import { useEditor } from '../editor/EditorContext';
import type { TrueFalseQuestion } from '../schema/types';

interface TrueFalseEditorProps {
  question: TrueFalseQuestion;
  index: number;
  readOnly?: boolean;
}

export function TrueFalseEditor({
  question,
  index,
  readOnly,
}: TrueFalseEditorProps) {
  const { updateQuestion } = useEditor();

  const update = useCallback(
    (updates: Partial<TrueFalseQuestion['params']>) => {
      updateQuestion(index, {
        ...question,
        params: { ...question.params, ...updates },
      });
    },
    [updateQuestion, index, question]
  );

  const handleQuestionTextChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    update({ question: e.target.value });
  };

  const handleCorrectChange = (value: 'true' | 'false') => {
    update({ correct: value });
  };

  const handleFeedbackChange = (
    field: 'feedbackOnCorrect' | 'feedbackOnWrong',
    value: string
  ) => {
    update({ [field]: value });
  };

  return (
    <div className="h5p-question-editor h5p-truefalse-editor">
      <h3 className="h5p-question-editor-header">
        <span className="h5p-question-editor-icon">⚖</span>
        True/False
      </h3>

      {/* Question Text */}
      <div className="h5p-question-editor-field">
        <label htmlFor={`question-text-${index}`}>Statement</label>
        <textarea
          id={`question-text-${index}`}
          value={question.params.question}
          onChange={handleQuestionTextChange}
          disabled={readOnly}
          className="h5p-editor-textarea"
          placeholder="Enter a statement that is either true or false..."
          rows={3}
        />
      </div>

      {/* Correct Answer */}
      <div className="h5p-question-editor-field">
        <label>Correct Answer</label>
        <div className="h5p-truefalse-options">
          <label className="h5p-truefalse-option">
            <input
              type="radio"
              name={`correct-${index}`}
              checked={question.params.correct === 'true'}
              onChange={() => handleCorrectChange('true')}
              disabled={readOnly}
            />
            <span className="h5p-truefalse-option-label h5p-truefalse-true">
              True
            </span>
          </label>
          <label className="h5p-truefalse-option">
            <input
              type="radio"
              name={`correct-${index}`}
              checked={question.params.correct === 'false'}
              onChange={() => handleCorrectChange('false')}
              disabled={readOnly}
            />
            <span className="h5p-truefalse-option-label h5p-truefalse-false">
              False
            </span>
          </label>
        </div>
      </div>

      {/* Feedback */}
      <div className="h5p-question-editor-field">
        <label htmlFor={`feedback-correct-${index}`}>Feedback (Correct)</label>
        <input
          id={`feedback-correct-${index}`}
          type="text"
          value={question.params.feedbackOnCorrect || ''}
          onChange={(e) => handleFeedbackChange('feedbackOnCorrect', e.target.value)}
          disabled={readOnly}
          className="h5p-editor-input"
          placeholder="Optional feedback when answered correctly"
        />
      </div>

      <div className="h5p-question-editor-field">
        <label htmlFor={`feedback-wrong-${index}`}>Feedback (Wrong)</label>
        <input
          id={`feedback-wrong-${index}`}
          type="text"
          value={question.params.feedbackOnWrong || ''}
          onChange={(e) => handleFeedbackChange('feedbackOnWrong', e.target.value)}
          disabled={readOnly}
          className="h5p-editor-input"
          placeholder="Optional feedback when answered incorrectly"
        />
      </div>
    </div>
  );
}
