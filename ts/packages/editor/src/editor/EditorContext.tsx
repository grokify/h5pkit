import {
  createContext,
  useContext,
  useReducer,
  useCallback,
  type ReactNode,
  type Dispatch,
} from 'react';
import type {
  H5PQuiz,
  Question,
  ValidationResult,
} from '../schema/types';
import { validateQuiz } from '../schema/validation';
import { DEFAULT_QUIZ } from '../schema/types';

// =============================================================================
// State
// =============================================================================

interface EditorState {
  quiz: H5PQuiz;
  selectedQuestionIndex: number | null;
  isDirty: boolean;
  validationResult: ValidationResult;
  undoStack: H5PQuiz[];
  redoStack: H5PQuiz[];
}

const initialState: EditorState = {
  quiz: DEFAULT_QUIZ,
  selectedQuestionIndex: null,
  isDirty: false,
  validationResult: { valid: true, errors: [] },
  undoStack: [],
  redoStack: [],
};

// =============================================================================
// Actions
// =============================================================================

type EditorAction =
  | { type: 'SET_QUIZ'; payload: H5PQuiz }
  | { type: 'UPDATE_METADATA'; payload: Partial<H5PQuiz['metadata']> }
  | { type: 'UPDATE_PARAMS'; payload: Partial<H5PQuiz['params']> }
  | { type: 'ADD_QUESTION'; payload: Question }
  | { type: 'UPDATE_QUESTION'; payload: { index: number; question: Question } }
  | { type: 'REMOVE_QUESTION'; payload: number }
  | { type: 'MOVE_QUESTION'; payload: { from: number; to: number } }
  | { type: 'SELECT_QUESTION'; payload: number | null }
  | { type: 'UNDO' }
  | { type: 'REDO' }
  | { type: 'MARK_CLEAN' };

function editorReducer(state: EditorState, action: EditorAction): EditorState {
  switch (action.type) {
    case 'SET_QUIZ': {
      const quiz = action.payload;
      return {
        ...state,
        quiz,
        isDirty: false,
        validationResult: validateQuiz(quiz),
        undoStack: [],
        redoStack: [],
        selectedQuestionIndex: null,
      };
    }

    case 'UPDATE_METADATA': {
      const quiz: H5PQuiz = {
        ...state.quiz,
        metadata: { ...state.quiz.metadata, ...action.payload },
      };
      return {
        ...state,
        quiz,
        isDirty: true,
        validationResult: validateQuiz(quiz),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: [],
      };
    }

    case 'UPDATE_PARAMS': {
      const quiz: H5PQuiz = {
        ...state.quiz,
        params: { ...state.quiz.params, ...action.payload },
      };
      return {
        ...state,
        quiz,
        isDirty: true,
        validationResult: validateQuiz(quiz),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: [],
      };
    }

    case 'ADD_QUESTION': {
      const questions = [...state.quiz.params.questions, action.payload];
      const quiz: H5PQuiz = {
        ...state.quiz,
        params: { ...state.quiz.params, questions },
      };
      return {
        ...state,
        quiz,
        isDirty: true,
        validationResult: validateQuiz(quiz),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: [],
        selectedQuestionIndex: questions.length - 1,
      };
    }

    case 'UPDATE_QUESTION': {
      const questions = [...state.quiz.params.questions];
      questions[action.payload.index] = action.payload.question;
      const quiz: H5PQuiz = {
        ...state.quiz,
        params: { ...state.quiz.params, questions },
      };
      return {
        ...state,
        quiz,
        isDirty: true,
        validationResult: validateQuiz(quiz),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: [],
      };
    }

    case 'REMOVE_QUESTION': {
      const questions = state.quiz.params.questions.filter(
        (_, i) => i !== action.payload
      );
      const quiz: H5PQuiz = {
        ...state.quiz,
        params: { ...state.quiz.params, questions },
      };
      return {
        ...state,
        quiz,
        isDirty: true,
        validationResult: validateQuiz(quiz),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: [],
        selectedQuestionIndex:
          state.selectedQuestionIndex === action.payload
            ? null
            : state.selectedQuestionIndex !== null &&
              state.selectedQuestionIndex > action.payload
            ? state.selectedQuestionIndex - 1
            : state.selectedQuestionIndex,
      };
    }

    case 'MOVE_QUESTION': {
      const questions = [...state.quiz.params.questions];
      const [moved] = questions.splice(action.payload.from, 1);
      questions.splice(action.payload.to, 0, moved);
      const quiz: H5PQuiz = {
        ...state.quiz,
        params: { ...state.quiz.params, questions },
      };
      return {
        ...state,
        quiz,
        isDirty: true,
        validationResult: validateQuiz(quiz),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: [],
        selectedQuestionIndex: action.payload.to,
      };
    }

    case 'SELECT_QUESTION':
      return {
        ...state,
        selectedQuestionIndex: action.payload,
      };

    case 'UNDO': {
      if (state.undoStack.length === 0) return state;
      const previous = state.undoStack[state.undoStack.length - 1];
      return {
        ...state,
        quiz: previous,
        isDirty: true,
        validationResult: validateQuiz(previous),
        undoStack: state.undoStack.slice(0, -1),
        redoStack: [...state.redoStack, state.quiz],
      };
    }

    case 'REDO': {
      if (state.redoStack.length === 0) return state;
      const next = state.redoStack[state.redoStack.length - 1];
      return {
        ...state,
        quiz: next,
        isDirty: true,
        validationResult: validateQuiz(next),
        undoStack: [...state.undoStack, state.quiz],
        redoStack: state.redoStack.slice(0, -1),
      };
    }

    case 'MARK_CLEAN':
      return {
        ...state,
        isDirty: false,
      };

    default:
      return state;
  }
}

// =============================================================================
// Context
// =============================================================================

interface EditorContextValue {
  state: EditorState;
  dispatch: Dispatch<EditorAction>;

  // Convenience methods
  setQuiz: (quiz: H5PQuiz) => void;
  getQuiz: () => H5PQuiz;
  updateTitle: (title: string) => void;
  updatePassPercentage: (percentage: number) => void;
  addQuestion: (question: Question) => void;
  updateQuestion: (index: number, question: Question) => void;
  removeQuestion: (index: number) => void;
  moveQuestion: (from: number, to: number) => void;
  selectQuestion: (index: number | null) => void;
  undo: () => void;
  redo: () => void;
  canUndo: boolean;
  canRedo: boolean;
  validate: () => ValidationResult;
}

const EditorContext = createContext<EditorContextValue | null>(null);

// =============================================================================
// Provider
// =============================================================================

interface EditorProviderProps {
  children: ReactNode;
  initialQuiz?: H5PQuiz;
  onChange?: (quiz: H5PQuiz) => void;
}

export function EditorProvider({
  children,
  initialQuiz,
  onChange,
}: EditorProviderProps) {
  const [state, dispatch] = useReducer(editorReducer, {
    ...initialState,
    quiz: initialQuiz || DEFAULT_QUIZ,
    validationResult: validateQuiz(initialQuiz || DEFAULT_QUIZ),
  });

  // Notify parent of changes
  const wrappedDispatch: Dispatch<EditorAction> = useCallback(
    (action) => {
      dispatch(action);
      // Note: We can't get the new state here, so parent should use getQuiz()
      // or we could use useEffect to watch state.quiz changes
    },
    []
  );

  const contextValue: EditorContextValue = {
    state,
    dispatch: wrappedDispatch,

    setQuiz: useCallback((quiz: H5PQuiz) => {
      dispatch({ type: 'SET_QUIZ', payload: quiz });
      onChange?.(quiz);
    }, [onChange]),

    getQuiz: useCallback(() => state.quiz, [state.quiz]),

    updateTitle: useCallback((title: string) => {
      dispatch({ type: 'UPDATE_METADATA', payload: { title } });
    }, []),

    updatePassPercentage: useCallback((percentage: number) => {
      dispatch({ type: 'UPDATE_PARAMS', payload: { passPercentage: percentage } });
    }, []),

    addQuestion: useCallback((question: Question) => {
      dispatch({ type: 'ADD_QUESTION', payload: question });
    }, []),

    updateQuestion: useCallback((index: number, question: Question) => {
      dispatch({ type: 'UPDATE_QUESTION', payload: { index, question } });
    }, []),

    removeQuestion: useCallback((index: number) => {
      dispatch({ type: 'REMOVE_QUESTION', payload: index });
    }, []),

    moveQuestion: useCallback((from: number, to: number) => {
      dispatch({ type: 'MOVE_QUESTION', payload: { from, to } });
    }, []),

    selectQuestion: useCallback((index: number | null) => {
      dispatch({ type: 'SELECT_QUESTION', payload: index });
    }, []),

    undo: useCallback(() => {
      dispatch({ type: 'UNDO' });
    }, []),

    redo: useCallback(() => {
      dispatch({ type: 'REDO' });
    }, []),

    canUndo: state.undoStack.length > 0,
    canRedo: state.redoStack.length > 0,

    validate: useCallback(() => validateQuiz(state.quiz), [state.quiz]),
  };

  return (
    <EditorContext.Provider value={contextValue}>
      {children}
    </EditorContext.Provider>
  );
}

// =============================================================================
// Hook
// =============================================================================

export function useEditor(): EditorContextValue {
  const context = useContext(EditorContext);
  if (!context) {
    throw new Error('useEditor must be used within an EditorProvider');
  }
  return context;
}
