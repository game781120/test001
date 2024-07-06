package service

import (
	"encoding/json"
	"fmt"
	"thundersoft.com/brainos/openai"
)

const HeaderData = "data: "

func GetCompletionStreamChoice(content string) string {
	jsonData, _ := json.Marshal(openai.ChatCompletionStreamResponse{
		Choices: []openai.ChatCompletionStreamChoice{
			{
				Delta: openai.ChatCompletionStreamChoiceDelta{
					Role:    openai.ChatMessageRoleAssistant,
					Content: content,
				},
			},
		},
	})

	return fmt.Sprintf("%s%s\n\n", HeaderData, jsonData)
}

func GetCompletionStreamChoiceEx(content string) string {
	jsonData, _ := json.Marshal(openai.ChatCompletionStreamResponse{
		PromptAnnotations: []openai.PromptAnnotation{
			{
				ContentFilterResults: openai.ContentFilterResults{
					Hate: openai.Hate{
						Filtered: true,
						Severity: content,
					},
				},
			},
		},
	})

	return fmt.Sprintf("%s%s\n\n", HeaderData, jsonData)
}
