import { useEditor } from './EditorContext';
import type { Question } from '../schema/types';
import {
  DEFAULT_MULTICHOICE,
  DEFAULT_TRUEFALSE,
  DEFAULT_FILLBLANKS,
} from '../schema/types';
import clsx from 'clsx';

type QuestionTypeKey = 'multipleChoice' | 'trueFalse' | 'fillInBlanks';

interface QuestionListProps {
  allowedTypes?: QuestionTypeKey[];
  readOnly?: boolean;
}

const QUESTION_TYPE_INFO: Record<
  QuestionTypeKey,
  { label: string; icon: string; default: Question }
> = {
  multipleChoice: {
    label: 'Multiple Choice',
    icon: '☑',
    default: DEFAULT_MULTICHOICE,
  },
  trueFalse: {
    label: 'True/False',
    icon: '⚖',
    default: DEFAULT_TRUEFALSE,
  },
  fillInBlanks: {
    label: 'Fill in Blanks',
    icon: '___',
    default: DEFAULT_FILLBLANKS,
  },
};

const LIBRARY_TO_TYPE: Record<string, QuestionTypeKey> = {
  'H5P.MultiChoice 1.16': 'multipleChoice',
  'H5P.TrueFalse 1.8': 'trueFalse',
  'H5P.Blanks 1.14': 'fillInBlanks',
};

export function QuestionList({ allowedTypes, readOnly }: QuestionListProps) {
  const { state, addQuestion, removeQuestion, selectQuestion, moveQuestion } =
    useEditor();
  const { quiz, selectedQuestionIndex } = state;

  const types = allowedTypes || (['multipleChoice', 'trueFalse', 'fillInBlanks'] as QuestionTypeKey[]);

  const handleAddQuestion = (type: QuestionTypeKey) => {
    const info = QUESTION_TYPE_INFO[type];
    addQuestion(structuredClone(info.default));
  };

  const handleMoveUp = (index: number) => {
    if (index > 0) {
      moveQuestion(index, index - 1);
    }
  };

  const handleMoveDown = (index: number) => {
    if (index < quiz.params.questions.length - 1) {
      moveQuestion(index, index + 1);
    }
  };

  const getQuestionTitle = (question: Question, index: number): string => {
    let title = '';
    switch (question.library) {
      case 'H5P.MultiChoice 1.16':
        title = question.params.question;
        break;
      case 'H5P.TrueFalse 1.8':
        title = question.params.question;
        break;
      case 'H5P.Blanks 1.14':
        title = question.params.text;
        break;
    }
    return title?.slice(0, 50) || `Question ${index + 1}`;
  };

  return (
    <div className="h5p-question-list">
      <h3 className="h5p-question-list-header">Questions</h3>

      {/* Question list */}
      <ul className="h5p-question-list-items">
        {quiz.params.questions.map((question, index) => {
          const typeKey = LIBRARY_TO_TYPE[question.library];
          const typeInfo = typeKey ? QUESTION_TYPE_INFO[typeKey] : null;

          return (
            <li
              key={index}
              className={clsx(
                'h5p-question-list-item',
                selectedQuestionIndex === index && 'h5p-question-list-item-selected'
              )}
            >
              <button
                type="button"
                onClick={() => selectQuestion(index)}
                className="h5p-question-list-item-main"
              >
                <span className="h5p-question-list-item-icon">
                  {typeInfo?.icon || '?'}
                </span>
                <span className="h5p-question-list-item-title">
                  {getQuestionTitle(question, index)}
                </span>
              </button>

              {!readOnly && (
                <div className="h5p-question-list-item-actions">
                  <button
                    type="button"
                    onClick={() => handleMoveUp(index)}
                    disabled={index === 0}
                    className="h5p-question-list-item-action"
                    title="Move up"
                  >
                    ↑
                  </button>
                  <button
                    type="button"
                    onClick={() => handleMoveDown(index)}
                    disabled={index === quiz.params.questions.length - 1}
                    className="h5p-question-list-item-action"
                    title="Move down"
                  >
                    ↓
                  </button>
                  <button
                    type="button"
                    onClick={() => removeQuestion(index)}
                    className="h5p-question-list-item-action h5p-question-list-item-action-danger"
                    title="Delete"
                  >
                    ×
                  </button>
                </div>
              )}
            </li>
          );
        })}
      </ul>

      {quiz.params.questions.length === 0 && (
        <p className="h5p-question-list-empty">No questions yet. Add one below.</p>
      )}

      {/* Add question buttons */}
      {!readOnly && (
        <div className="h5p-question-list-add">
          <p className="h5p-question-list-add-label">Add Question:</p>
          <div className="h5p-question-list-add-buttons">
            {types.map((type) => {
              const info = QUESTION_TYPE_INFO[type];
              return (
                <button
                  key={type}
                  type="button"
                  onClick={() => handleAddQuestion(type)}
                  className="h5p-editor-btn h5p-editor-btn-secondary"
                >
                  <span>{info.icon}</span> {info.label}
                </button>
              );
            })}
          </div>
        </div>
      )}
    </div>
  );
}
