/**
 * H5P Quiz Types
 *
 * TypeScript types for H5P QuestionSet and related content types.
 * Based on H5P.QuestionSet 1.20 specification.
 */

// =============================================================================
// Metadata
// =============================================================================

export interface H5PMetadata {
  title: string;
  license?: string;
  licenseVersion?: string;
  licenseExtras?: string;
  authors?: Author[];
  changes?: Change[];
  source?: string;
  yearFrom?: number;
  yearTo?: number;
}

export interface Author {
  name: string;
  role?: string;
}

export interface Change {
  date: string;
  author: string;
  log: string;
}

// =============================================================================
// Question Set (Root Container)
// =============================================================================

export interface H5PQuiz {
  library: 'H5P.QuestionSet 1.20';
  metadata: H5PMetadata;
  params: QuestionSetParams;
}

export interface QuestionSetParams {
  introPage?: IntroPage;
  progressType: 'dots' | 'textual' | 'none';
  passPercentage: number;
  showResultPage: boolean;
  noResultMessage?: string;
  randomQuestions: boolean;
  poolSize?: number;
  questions: Question[];
  endGame: EndGame;
  texts?: QuestionSetTexts;
  override?: BehaviourOverride;
}

export interface IntroPage {
  showIntroPage: boolean;
  title?: string;
  introduction?: string;
  startButtonText?: string;
  backgroundImage?: H5PImage;
}

export interface EndGame {
  showResultPage: boolean;
  showSolutionButton: boolean;
  showRetryButton: boolean;
  noResultMessage?: string;
  message?: string;
  successGreeting?: string;
  successComment?: string;
  failGreeting?: string;
  failComment?: string;
  solutionButtonText?: string;
  retryButtonText?: string;
  finishButtonText?: string;
  showAnimations?: boolean;
  skippable?: boolean;
  skipButtonText?: string;
}

export interface QuestionSetTexts {
  prevButton?: string;
  nextButton?: string;
  finishButton?: string;
  textualProgress?: string;
  jumpToQuestion?: string;
  questionLabel?: string;
  readSpeakerProgress?: string;
  unansweredText?: string;
  answeredText?: string;
  currentQuestionText?: string;
}

export interface BehaviourOverride {
  enableRetry?: boolean;
  enableSolutionsButton?: boolean;
  showSolutionsRequiresInput?: boolean;
}

// =============================================================================
// Question Types
// =============================================================================

export type Question =
  | MultiChoiceQuestion
  | TrueFalseQuestion
  | FillInBlanksQuestion;

export type QuestionType = Question['library'];

// -----------------------------------------------------------------------------
// Multiple Choice (H5P.MultiChoice 1.16)
// -----------------------------------------------------------------------------

export interface MultiChoiceQuestion {
  library: 'H5P.MultiChoice 1.16';
  subContentId?: string;
  params: MultiChoiceParams;
}

export interface MultiChoiceParams {
  question: string;
  answers: MultiChoiceAnswer[];
  behaviour: MultiChoiceBehaviour;
  UI?: MultiChoiceUI;
  confirmCheck?: ConfirmationDialog;
  confirmRetry?: ConfirmationDialog;
  media?: H5PMedia;
}

export interface MultiChoiceAnswer {
  text: string;
  correct: boolean;
  tipsAndFeedback?: TipsAndFeedback;
}

export interface TipsAndFeedback {
  tip?: string;
  chosenFeedback?: string;
  notChosenFeedback?: string;
}

export interface MultiChoiceBehaviour {
  enableRetry?: boolean;
  enableSolutionsButton?: boolean;
  enableCheckButton?: boolean;
  type?: 'auto' | 'single' | 'multi';
  singlePoint?: boolean;
  randomAnswers?: boolean;
  showSolutionsRequiresInput?: boolean;
  autoCheck?: boolean;
  passPercentage?: number;
  showScorePoints?: boolean;
}

export interface MultiChoiceUI {
  checkAnswerButton?: string;
  submitAnswerButton?: string;
  showSolutionButton?: string;
  tryAgainButton?: string;
  tipsLabel?: string;
  scoreBarLabel?: string;
  tipAvailable?: string;
  feedbackAvailable?: string;
  readFeedback?: string;
  wrongAnswer?: string;
  correctAnswer?: string;
  shouldCheck?: string;
  shouldNotCheck?: string;
  noInput?: string;
  a11yCheck?: string;
  a11yShowSolution?: string;
  a11yRetry?: string;
}

export interface ConfirmationDialog {
  header?: string;
  body?: string;
  cancelLabel?: string;
  confirmLabel?: string;
}

// -----------------------------------------------------------------------------
// True/False (H5P.TrueFalse 1.8)
// -----------------------------------------------------------------------------

export interface TrueFalseQuestion {
  library: 'H5P.TrueFalse 1.8';
  subContentId?: string;
  params: TrueFalseParams;
}

export interface TrueFalseParams {
  question: string;
  correct: 'true' | 'false';
  behaviour?: TrueFalseBehaviour;
  l10n?: TrueFalseL10n;
  media?: H5PMedia;
  feedbackOnCorrect?: string;
  feedbackOnWrong?: string;
}

export interface TrueFalseBehaviour {
  enableRetry?: boolean;
  enableSolutionsButton?: boolean;
  enableCheckButton?: boolean;
  confirmCheckDialog?: boolean;
  confirmRetryDialog?: boolean;
  autoCheck?: boolean;
}

export interface TrueFalseL10n {
  trueText?: string;
  falseText?: string;
  checkAnswer?: string;
  showSolutionButton?: string;
  tryAgain?: string;
  wrongAnswerMessage?: string;
  correctAnswerMessage?: string;
  scoreBarLabel?: string;
  a11yCheck?: string;
  a11yShowSolution?: string;
  a11yRetry?: string;
}

// -----------------------------------------------------------------------------
// Fill in the Blanks (H5P.Blanks 1.14)
// -----------------------------------------------------------------------------

export interface FillInBlanksQuestion {
  library: 'H5P.Blanks 1.14';
  subContentId?: string;
  params: FillInBlanksParams;
}

export interface FillInBlanksParams {
  text: string; // Contains *blanks* marked with asterisks
  questions: BlankQuestion[];
  behaviour?: FillInBlanksBehaviour;
  overallFeedback?: OverallFeedback[];
  showSolutions?: string;
  tryAgain?: string;
  checkAnswer?: string;
  submitAnswer?: string;
  notFilledOut?: string;
  answerIsCorrect?: string;
  answerIsWrong?: string;
  answeredCorrectly?: string;
  answeredIncorrectly?: string;
  solutionLabel?: string;
  inputLabel?: string;
  inputHasTipLabel?: string;
  tipLabel?: string;
  media?: H5PMedia;
}

export interface BlankQuestion {
  text: string; // The sentence with *blank* markers
  tip?: string;
}

export interface FillInBlanksBehaviour {
  enableRetry?: boolean;
  enableSolutionsButton?: boolean;
  enableCheckButton?: boolean;
  autoCheck?: boolean;
  caseSensitive?: boolean;
  showSolutionsRequiresInput?: boolean;
  separateLines?: boolean;
  confirmCheckDialog?: boolean;
  confirmRetryDialog?: boolean;
  acceptSpellingErrors?: boolean;
}

export interface OverallFeedback {
  from: number;
  to: number;
  feedback: string;
}

// =============================================================================
// Shared Types
// =============================================================================

export interface H5PImage {
  path?: string;
  mime?: string;
  copyright?: Copyright;
  width?: number;
  height?: number;
  alt?: string;
}

export interface H5PMedia {
  type?: H5PImage | H5PVideo;
  disableImageZooming?: boolean;
}

export interface H5PVideo {
  path?: string;
  mime?: string;
  copyright?: Copyright;
}

export interface Copyright {
  license?: string;
  author?: string;
  year?: string;
  source?: string;
  version?: string;
}

// =============================================================================
// Editor Types
// =============================================================================

export interface ValidationError {
  path: string;
  message: string;
  keyword?: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
}

// =============================================================================
// Default Values
// =============================================================================

export const DEFAULT_QUIZ: H5PQuiz = {
  library: 'H5P.QuestionSet 1.20',
  metadata: {
    title: 'New Quiz',
  },
  params: {
    progressType: 'dots',
    passPercentage: 70,
    showResultPage: true,
    randomQuestions: false,
    questions: [],
    endGame: {
      showResultPage: true,
      showSolutionButton: true,
      showRetryButton: true,
      message: 'Your score: @score / @total',
    },
  },
};

export const DEFAULT_MULTICHOICE: MultiChoiceQuestion = {
  library: 'H5P.MultiChoice 1.16',
  params: {
    question: '',
    answers: [
      { text: '', correct: true },
      { text: '', correct: false },
    ],
    behaviour: {
      enableRetry: true,
      enableSolutionsButton: true,
      type: 'auto',
      randomAnswers: true,
    },
  },
};

export const DEFAULT_TRUEFALSE: TrueFalseQuestion = {
  library: 'H5P.TrueFalse 1.8',
  params: {
    question: '',
    correct: 'true',
    behaviour: {
      enableRetry: true,
      enableSolutionsButton: true,
    },
  },
};

export const DEFAULT_FILLBLANKS: FillInBlanksQuestion = {
  library: 'H5P.Blanks 1.14',
  params: {
    text: '',
    questions: [],
    behaviour: {
      enableRetry: true,
      enableSolutionsButton: true,
      caseSensitive: false,
    },
  },
};
