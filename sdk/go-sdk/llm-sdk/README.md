# Go Sdk of BrainOS LLM

## Installation

### GO111MODULE=off

下载 llm-sdk.zip 包，解压到 `$GOPATH/src/thundersoft.com/brainos/openai` 目录下即可

### GO111MODULE=on

下载 llm-sdk.zip 包，解压到任意目录，添加以下信息到 `go.mod`

```go
require (
	thundersoft.com/brainos/openai v0.0.0-incompatible
)

replace (
	thundersoft.com/brainos/openai => {llm-sdk 绝对或相对目录}/llm-sdk
)

```

## Usage

### 获取模型列表

```go
package main

import (
	"context"
	"fmt"
	openai "thundersoft.com/brainos/openai"
)

func main() {
	client := openai.NewClient("baseURL", "your token")
	resp, err := client.ListModels(context.Background())
	if err != nil {
		fmt.Printf("ListModels error: %v\n", err)
		return
	}

	for _, model := range resp.Content {
		fmt.Println(model.Model)
	}
}

```

### 非流式对话

```go
package main

import (
	"context"
	"fmt"
	openai "thundersoft.com/brainos/openai"
)

func main() {
	client := openai.NewClient("baseURL", "your token")
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

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
```

### 流式对话

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	openai "thundersoft.com/brainos/openai"
)

func main() {
	c := openai.NewClient("baseURL", "your token")
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
	stream, err := c.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}

		if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			return
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
}
```