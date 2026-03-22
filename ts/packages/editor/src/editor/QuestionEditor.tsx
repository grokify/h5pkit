import type { Question } from '../schema/types';
import { MultiChoiceEditor } from '../content-types/MultiChoiceEditor';
import { TrueFalseEditor } from '../content-types/TrueFalseEditor';
import { FillInBlanksEditor } from '../content-types/FillInBlanksEditor';

interface QuestionEditorProps {
  question: Question;
  index: number;
  readOnly?: boolean;
}

export function QuestionEditor({ question, index, readOnly }: QuestionEditorProps) {
  switch (question.library) {
    case 'H5P.MultiChoice 1.16':
      return (
        <MultiChoiceEditor
          question={question}
          index={index}
          readOnly={readOnly}
        />
      );

    case 'H5P.TrueFalse 1.8':
      return (
        <TrueFalseEditor
          question={question}
          index={index}
          readOnly={readOnly}
        />
      );

    case 'H5P.Blanks 1.14':
      return (
        <FillInBlanksEditor
          question={question}
          index={index}
          readOnly={readOnly}
        />
      );

    default:
      return (
        <div className="h5p-question-editor-unknown">
          <p>Unknown question type: {(question as { library: string }).library}</p>
        </div>
      );
  }
}
