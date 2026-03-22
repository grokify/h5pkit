package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	h5p "github.com/grokify/h5pkit"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("must supply JSON filename")
		os.Exit(1)
	}
	// Read JSON file
	jsonData, err := os.ReadFile(filepath.Clean(os.Args[1]))
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON into QuestionSet
	questionSet, err := h5p.FromJSON(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// Validate the loaded question set
	err = questionSet.Validate()
	if err != nil {
		log.Fatal(err)
	}

	// Use the question set
	fmt.Printf("Loaded question set: %s\n", questionSet.Title)
	fmt.Printf("Number of questions: %d\n", len(questionSet.Questions))
	fmt.Printf("Pass percentage: %d%%\n", questionSet.PassPercentage)

	// Access individual questions
	for i, question := range questionSet.Questions {
		fmt.Printf("Question %d Library: %s\n", i+1, question.Library)

		// For MultiChoice questions, you can access params
		if paramsMap, ok := question.Params.(map[string]interface{}); ok {
			if questionText, exists := paramsMap["question"]; exists {
				fmt.Printf("Question %d Text: %s\n", i+1, questionText)
			}

			if answers, exists := paramsMap["answers"]; exists {
				if answersList, ok := answers.([]interface{}); ok {
					fmt.Printf("Question %d has %d answers\n", i+1, len(answersList))
				}
			}
		}
	}
}
