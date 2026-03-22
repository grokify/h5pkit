import { useCallback } from 'react';
import { useEditor } from '../editor/EditorContext';
import type { MultiChoiceQuestion } from '../schema/types';

interface MultiChoiceEditorProps {
  question: MultiChoiceQuestion;
  index: number;
  readOnly?: boolean;
}

export function MultiChoiceEditor({
  question,
  index,
  readOnly,
}: MultiChoiceEditorProps) {
  const { updateQuestion } = useEditor();

  const update = useCallback(
    (updates: Partial<MultiChoiceQuestion['params']>) => {
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

  const handleAnswerTextChange = (answerIndex: number, text: string) => {
    const answers = [...question.params.answers];
    answers[answerIndex] = { ...answers[answerIndex], text };
    update({ answers });
  };

  const handleAnswerCorrectChange = (answerIndex: number, correct: boolean) => {
    const answers = [...question.params.answers];
    answers[answerIndex] = { ...answers[answerIndex], correct };
    update({ answers });
  };

  const handleAddAnswer = () => {
    const answers = [...question.params.answers, { text: '', correct: false }];
    update({ answers });
  };

  const handleRemoveAnswer = (answerIndex: number) => {
    const answers = question.params.answers.filter((_, i) => i !== answerIndex);
    update({ answers });
  };

  const handleMoveAnswer = (from: number, to: number) => {
    const answers = [...question.params.answers];
    const [moved] = answers.splice(from, 1);
    answers.splice(to, 0, moved);
    update({ answers });
  };

  const handleBehaviourChange = (key: keyof MultiChoiceQuestion['params']['behaviour'], value: boolean) => {
    update({
      behaviour: { ...question.params.behaviour, [key]: value },
    });
  };

  return (
    <div className="h5p-question-editor h5p-multichoice-editor">
      <h3 className="h5p-question-editor-header">
        <span className="h5p-question-editor-icon">☑</span>
        Multiple Choice
      </h3>

      {/* Question Text */}
      <div className="h5p-question-editor-field">
        <label htmlFor={`question-text-${index}`}>Question</label>
        <textarea
          id={`question-text-${index}`}
          value={question.params.question}
          onChange={handleQuestionTextChange}
          disabled={readOnly}
          className="h5p-editor-textarea"
          placeholder="Enter your question here..."
          rows={3}
        />
      </div>

      {/* Answers */}
      <div className="h5p-question-editor-field">
        <label>Answers</label>
        <ul className="h5p-multichoice-answers">
          {question.params.answers.map((answer, answerIndex) => (
            <li key={answerIndex} className="h5p-multichoice-answer">
              <div className="h5p-multichoice-answer-correct">
                <input
                  type="checkbox"
                  checked={answer.correct}
                  onChange={(e) =>
                    handleAnswerCorrectChange(answerIndex, e.target.checked)
                  }
                  disabled={readOnly}
                  title="Mark as correct answer"
                />
              </div>
              <input
                type="text"
                value={answer.text}
                onChange={(e) =>
                  handleAnswerTextChange(answerIndex, e.target.value)
                }
                disabled={readOnly}
                className="h5p-editor-input h5p-multichoice-answer-text"
                placeholder={`Answer ${answerIndex + 1}`}
              />
              {!readOnly && (
                <div className="h5p-multichoice-answer-actions">
                  <button
                    type="button"
                    onClick={() => handleMoveAnswer(answerIndex, answerIndex - 1)}
                    disabled={answerIndex === 0}
                    className="h5p-editor-btn-icon"
                    title="Move up"
                  >
                    ↑
                  </button>
                  <button
                    type="button"
                    onClick={() => handleMoveAnswer(answerIndex, answerIndex + 1)}
                    disabled={answerIndex === question.params.answers.length - 1}
                    className="h5p-editor-btn-icon"
                    title="Move down"
                  >
                    ↓
                  </button>
                  <button
                    type="button"
                    onClick={() => handleRemoveAnswer(answerIndex)}
                    disabled={question.params.answers.length <= 2}
                    className="h5p-editor-btn-icon h5p-editor-btn-danger"
                    title="Remove"
                  >
                    ×
                  </button>
                </div>
              )}
            </li>
          ))}
        </ul>

        {!readOnly && (
          <button
            type="button"
            onClick={handleAddAnswer}
            className="h5p-editor-btn h5p-editor-btn-secondary"
          >
            + Add Answer
          </button>
        )}
      </div>

      {/* Behaviour Settings */}
      <div className="h5p-question-editor-field">
        <label>Behaviour</label>
        <div className="h5p-question-editor-behaviour">
          <label className="h5p-editor-checkbox">
            <input
              type="checkbox"
              checked={question.params.behaviour.enableRetry ?? true}
              onChange={(e) => handleBehaviourChange('enableRetry', e.target.checked)}
              disabled={readOnly}
            />
            Allow retry
          </label>
          <label className="h5p-editor-checkbox">
            <input
              type="checkbox"
              checked={question.params.behaviour.enableSolutionsButton ?? true}
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
              checked={question.params.behaviour.randomAnswers ?? true}
              onChange={(e) =>
                handleBehaviourChange('randomAnswers', e.target.checked)
              }
              disabled={readOnly}
            />
            Randomize answer order
          </label>
        </div>
      </div>
    </div>
  );
}
