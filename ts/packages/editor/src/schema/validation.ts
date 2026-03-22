/**
 * H5P Quiz Validation
 *
 * Validates H5P quiz content using JSON Schema and semantic rules.
 */

import type {
  H5PQuiz,
  Question,
  MultiChoiceQuestion,
  ValidationResult,
  ValidationError,
} from './types';

/**
 * Validates an H5P quiz structure.
 * Performs both structural and semantic validation.
 */
export function validateQuiz(quiz: unknown): ValidationResult {
  const errors: ValidationError[] = [];

  // Type guard
  if (!isH5PQuiz(quiz)) {
    return {
      valid: false,
      errors: [{ path: '/', message: 'Invalid quiz structure' }],
    };
  }

  // Validate metadata
  if (!quiz.metadata?.title?.trim()) {
    errors.push({
      path: '/metadata/title',
      message: 'Quiz title is required',
    });
  }

  // Validate params
  const params = quiz.params;

  // Pass percentage must be 0-100
  if (params.passPercentage < 0 || params.passPercentage > 100) {
    errors.push({
      path: '/params/passPercentage',
      message: 'Pass percentage must be between 0 and 100',
    });
  }

  // Must have at least one question
  if (!params.questions || params.questions.length === 0) {
    errors.push({
      path: '/params/questions',
      message: 'Quiz must have at least one question',
    });
  }

  // Validate each question
  params.questions?.forEach((question, index) => {
    const questionErrors = validateQuestion(question, index);
    errors.push(...questionErrors);
  });

  return {
    valid: errors.length === 0,
    errors,
  };
}

/**
 * Validates a single question.
 */
function validateQuestion(question: Question, index: number): ValidationError[] {
  const errors: ValidationError[] = [];
  const basePath = `/params/questions/${index}`;

  switch (question.library) {
    case 'H5P.MultiChoice 1.16':
      errors.push(...validateMultiChoice(question, basePath));
      break;
    case 'H5P.TrueFalse 1.8':
      errors.push(...validateTrueFalse(question, basePath));
      break;
    case 'H5P.Blanks 1.14':
      errors.push(...validateFillInBlanks(question, basePath));
      break;
    default:
      errors.push({
        path: `${basePath}/library`,
        message: `Unknown question type: ${(question as { library: string }).library}`,
      });
  }

  return errors;
}

/**
 * Validates a multiple choice question.
 */
function validateMultiChoice(
  question: MultiChoiceQuestion,
  basePath: string
): ValidationError[] {
  const errors: ValidationError[] = [];
  const params = question.params;

  // Question text required
  if (!params.question?.trim()) {
    errors.push({
      path: `${basePath}/params/question`,
      message: 'Question text is required',
    });
  }

  // Must have at least 2 answers
  if (!params.answers || params.answers.length < 2) {
    errors.push({
      path: `${basePath}/params/answers`,
      message: 'Multiple choice must have at least 2 answers',
    });
  }

  // Must have at least one correct answer
  const hasCorrect = params.answers?.some((a) => a.correct);
  if (!hasCorrect) {
    errors.push({
      path: `${basePath}/params/answers`,
      message: 'Multiple choice must have at least one correct answer',
    });
  }

  // Each answer must have text
  params.answers?.forEach((answer, i) => {
    if (!answer.text?.trim()) {
      errors.push({
        path: `${basePath}/params/answers/${i}/text`,
        message: 'Answer text is required',
      });
    }
  });

  return errors;
}

/**
 * Validates a true/false question.
 */
function validateTrueFalse(
  question: { library: 'H5P.TrueFalse 1.8'; params: { question: string; correct: string } },
  basePath: string
): ValidationError[] {
  const errors: ValidationError[] = [];
  const params = question.params;

  // Question text required
  if (!params.question?.trim()) {
    errors.push({
      path: `${basePath}/params/question`,
      message: 'Question text is required',
    });
  }

  // Correct must be 'true' or 'false'
  if (params.correct !== 'true' && params.correct !== 'false') {
    errors.push({
      path: `${basePath}/params/correct`,
      message: 'Correct answer must be "true" or "false"',
    });
  }

  return errors;
}

/**
 * Validates a fill in the blanks question.
 */
function validateFillInBlanks(
  question: { library: 'H5P.Blanks 1.14'; params: { text: string; questions: unknown[] } },
  basePath: string
): ValidationError[] {
  const errors: ValidationError[] = [];
  const params = question.params;

  // Must have text with blanks
  if (!params.text?.trim()) {
    errors.push({
      path: `${basePath}/params/text`,
      message: 'Fill in blanks text is required',
    });
  }

  // Text should contain at least one blank marker (asterisks)
  if (params.text && !params.text.includes('*')) {
    errors.push({
      path: `${basePath}/params/text`,
      message: 'Text must contain at least one blank (marked with *asterisks*)',
    });
  }

  return errors;
}

/**
 * Type guard for H5PQuiz.
 */
function isH5PQuiz(obj: unknown): obj is H5PQuiz {
  if (!obj || typeof obj !== 'object') return false;
  const quiz = obj as Record<string, unknown>;
  return (
    quiz.library === 'H5P.QuestionSet 1.20' &&
    typeof quiz.metadata === 'object' &&
    typeof quiz.params === 'object'
  );
}

/**
 * Checks if a question type is supported.
 */
export function isSupportedQuestionType(library: string): boolean {
  return [
    'H5P.MultiChoice 1.16',
    'H5P.TrueFalse 1.8',
    'H5P.Blanks 1.14',
  ].includes(library);
}
