package chat

import (
	"context"
	"errors"
	"io"
	"testing"

	"thundersoft.com/brainos/chat/internal/test/checks"
)

func TestAddChat(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.AddChat(context.Background())

	t.Log(resp.Data)

	checks.NoError(t, err, "TestAddChat error")
}

func TestCreateChatCompletion(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.CreateChatCompletion(context.Background(),
		ChatCompletionRequest{
			ChatId:  1765630418937131008,
			Message: "hello",
		})

	for _, v := range resp.Data.MessageList {
		t.Log(v.Message)
	}

	checks.NoError(t, err, "TestCreateChatCompletion error")
}

func TestChatCompletionsWithStream(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	stream, err := client.CreateChatCompletionStream(context.Background(), ChatCompletionRequest{
		ChatId:  1765630418937131008,
		Message: "请帮我写一个有意思的关于穿越的g短篇小说，主题和和情节以及主要角色你发挥你的想象，字数控制在200字以上",
	})
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

		t.Log(response.Message)
	}
}

func TestCreateQaChatCompletion(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.CreateQaChatCompletion(context.Background(),
		QaChatCompletionRequest{
			ChatId:  1765630418937131008,
			Message: "hello",
			Knowledges: []QaChatRequestKnowledges{
				{
					Id:    "1761989896572108802",
					Value: "hello",
				},
			},
		})

	for _, v := range resp.Data.MessageList {
		t.Log(v.Message)
	}

	checks.NoError(t, err, "TestCreateChatCompletion error")
}

func TestQaChatCompletionsWithStream(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	stream, err := client.CreateQaChatCompletionStream(context.Background(), QaChatCompletionRequest{
		ChatId:  1765630418937131008,
		Message: "hello",
		Knowledges: []QaChatRequestKnowledges{
			{
				Id:    1761989896572108802,
				Value: "hello",
			},
		},
	})

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

		t.Log(response.Message)
	}
}

func TestChatList(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.ChatList(context.Background())

	for _, v := range resp.Data {
		t.Log(v.Id, v.Name)
	}

	checks.NoError(t, err, "TestChatList error")
}

func TestDeleteChat(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	err := client.DeleteChat(context.Background(), 1765630418937131008)
	checks.NoError(t, err, "TestDeleteChat error")
}

func TestChatHistory(t *testing.T) {
	client := NewClient("http://10.0.36.13:8888/",
		"IHDKqCSoT2oMlYq51f5c9b9274B64c7eA66f8dC2Fe5a86Fc", "1762455766604193792", "test-11111")

	resp, err := client.ChatHistory(context.Background(), ChatHistoryRequest{
		ChatId:   1765630418937131008,
		PageSize: 10,
		Id:       0,
	})

	for _, v := range resp.Data.HistoryList {
		t.Log("Q:", v.Message)
		for _, v := range v.Children {
			t.Log("A:", v.Message)
		}
	}

	checks.NoError(t, err, "TestChatHistory error")
}
