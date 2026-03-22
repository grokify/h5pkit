import { useEditor } from './EditorContext';

interface EditorToolbarProps {
  onSave?: () => void;
  readOnly?: boolean;
}

export function EditorToolbar({ onSave, readOnly }: EditorToolbarProps) {
  const { state, undo, redo, canUndo, canRedo } = useEditor();

  return (
    <div className="h5p-editor-toolbar">
      <div className="h5p-editor-toolbar-left">
        <h2 className="h5p-editor-title">
          {state.quiz.metadata.title || 'Untitled Quiz'}
        </h2>
        {state.isDirty && <span className="h5p-editor-dirty-indicator">*</span>}
      </div>

      <div className="h5p-editor-toolbar-right">
        {!readOnly && (
          <>
            <button
              type="button"
              onClick={undo}
              disabled={!canUndo}
              className="h5p-editor-btn h5p-editor-btn-icon"
              title="Undo (Ctrl+Z)"
            >
              <UndoIcon />
            </button>
            <button
              type="button"
              onClick={redo}
              disabled={!canRedo}
              className="h5p-editor-btn h5p-editor-btn-icon"
              title="Redo (Ctrl+Shift+Z)"
            >
              <RedoIcon />
            </button>
          </>
        )}

        {!state.validationResult.valid && (
          <span className="h5p-editor-validation-badge h5p-editor-validation-error">
            {state.validationResult.errors.length} error(s)
          </span>
        )}

        {state.validationResult.valid && state.quiz.params.questions.length > 0 && (
          <span className="h5p-editor-validation-badge h5p-editor-validation-valid">
            Valid
          </span>
        )}

        {onSave && !readOnly && (
          <button
            type="button"
            onClick={onSave}
            disabled={!state.validationResult.valid}
            className="h5p-editor-btn h5p-editor-btn-primary"
          >
            <SaveIcon />
            Save Quiz
          </button>
        )}
      </div>
    </div>
  );
}

// Simple SVG icons
function UndoIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M3 10h10a5 5 0 0 1 5 5v2" />
      <polyline points="3 10 8 5 8 15" transform="rotate(90 5.5 10)" />
      <path d="M3 10l5-5M3 10l5 5" />
    </svg>
  );
}

function RedoIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M21 10H11a5 5 0 0 0-5 5v2" />
      <path d="M21 10l-5-5M21 10l-5 5" />
    </svg>
  );
}

function SaveIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
      <path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z" />
      <polyline points="17 21 17 13 7 13 7 21" />
      <polyline points="7 3 7 8 15 8" />
    </svg>
  );
}
