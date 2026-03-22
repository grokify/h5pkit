import { useEditor } from './EditorContext';

interface QuizSettingsProps {
  readOnly?: boolean;
}

export function QuizSettings({ readOnly }: QuizSettingsProps) {
  const { state, updateTitle, updatePassPercentage, dispatch } = useEditor();
  const { quiz } = state;

  const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    updateTitle(e.target.value);
  };

  const handlePassPercentageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = parseInt(e.target.value, 10);
    if (!isNaN(value) && value >= 0 && value <= 100) {
      updatePassPercentage(value);
    }
  };

  const handleRandomQuestionsChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    dispatch({
      type: 'UPDATE_PARAMS',
      payload: { randomQuestions: e.target.checked },
    });
  };

  const handleProgressTypeChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    dispatch({
      type: 'UPDATE_PARAMS',
      payload: { progressType: e.target.value as 'dots' | 'textual' | 'none' },
    });
  };

  return (
    <div className="h5p-quiz-settings">
      <h3 className="h5p-quiz-settings-header">Quiz Settings</h3>

      <div className="h5p-quiz-settings-field">
        <label htmlFor="quiz-title">Title</label>
        <input
          id="quiz-title"
          type="text"
          value={quiz.metadata.title}
          onChange={handleTitleChange}
          disabled={readOnly}
          className="h5p-editor-input"
          placeholder="Enter quiz title"
        />
      </div>

      <div className="h5p-quiz-settings-field">
        <label htmlFor="pass-percentage">Pass Percentage</label>
        <div className="h5p-quiz-settings-field-inline">
          <input
            id="pass-percentage"
            type="number"
            min={0}
            max={100}
            value={quiz.params.passPercentage}
            onChange={handlePassPercentageChange}
            disabled={readOnly}
            className="h5p-editor-input h5p-editor-input-small"
          />
          <span>%</span>
        </div>
      </div>

      <div className="h5p-quiz-settings-field">
        <label htmlFor="progress-type">Progress Display</label>
        <select
          id="progress-type"
          value={quiz.params.progressType}
          onChange={handleProgressTypeChange}
          disabled={readOnly}
          className="h5p-editor-select"
        >
          <option value="dots">Dots</option>
          <option value="textual">Text (1/5)</option>
          <option value="none">None</option>
        </select>
      </div>

      <div className="h5p-quiz-settings-field h5p-quiz-settings-field-checkbox">
        <input
          id="random-questions"
          type="checkbox"
          checked={quiz.params.randomQuestions}
          onChange={handleRandomQuestionsChange}
          disabled={readOnly}
        />
        <label htmlFor="random-questions">Randomize question order</label>
      </div>
    </div>
  );
}
