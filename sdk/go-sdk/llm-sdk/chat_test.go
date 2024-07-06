package openai_test

import (
	"context"
	"errors"
	"io"
	"testing"

	"thundersoft.com/brainos/openai"
	"thundersoft.com/brainos/openai/internal/test/checks"
)

func TestChatCompletionsWithStream(t *testing.T) {
	client := openai.NewClient("http://10.0.36.13:8888", "KJInf01E1p5Q1zvn65704c7501Ef4e83B85aB44a0128E5Dd")
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     "rubik6-chat",
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Lorem ipsum",
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	checks.NoError(t, err, "unexpected error")
	defer stream.Close()

	t.Log("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			t.Log("\nStream finished")
			return
		}
		checks.NoError(t, err, "unexpected error")

		t.Log(response.Choices[0].Delta.Content)
	}
}

// TestCompletions Tests the completions endpoint of the API using the mocked server.
func TestChatCompletions(t *testing.T) {
	client := openai.NewClient("http://10.0.36.13:8888", "KJInf01E1p5Q1zvn65704c7501Ef4e83B85aB44a0128E5Dd")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "rubik6-chat",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)
	t.Log(resp.Choices[0].Message.Content)

	checks.NoError(t, err, "CreateChatCompletion error")
}
