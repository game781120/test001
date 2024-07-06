# Go Sdk of BrainOS Chat

## Installation

### GO111MODULE=off

下载 chat-sdk.zip 包，解压到 `$GOPATH/src/thundersoft.com/brainos/chat` 目录下即可

### GO111MODULE=on

下载 chat-sdk.zip 包，解压到任意目录，添加以下信息到 `go.mod`

```go
require (
	thundersoft.com/brainos/chat v0.0.0-incompatible
)

replace (
	thundersoft.com/brainos/chat => {chat-sdk 绝对或相对目录}/chat-sdk
)

```

## Usage

### 添加对话
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.AddChat(context.Background())
	if err != nil {
		fmt.Printf("AddChat error: %v\n", err)
		return
	}

	fmt.Println(resp.Data)
}

```

### 删除对话
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	err := client.DeleteChat(context.Background(), 1765630418937131008)
	if err != nil {
		fmt.Printf("DeleteChat error: %v\n", err)
	}
}
```

### 获取对话列表
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.ChatList(context.Background())
	if err != nil {
		fmt.Printf("ChatList error: %v\n", err)
		return
	}

	for _, v := range resp.Data {
		fmt.Println(v.Id, v.Name)
	}
}
```

### 非流式对话
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.CreateChatCompletion(context.Background(),
		chat.ChatCompletionRequest{
			ChatId:  1765630418937131008,
			Message: "hello",
		})
	if err != nil {
		fmt.Printf("CreateChatCompletion error: %v\n", err)
		return
	}

	for _, v := range resp.Data.MessageList {
		fmt.Println(v.Message)
	}
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

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	stream, err := client.CreateChatCompletionStream(context.Background(), chat.ChatCompletionRequest{
		ChatId:  1765630418937131008,
		Message: "hello",
	})
	if err != nil {
		fmt.Printf("ListModels error: %v\n", err)
		return
	}

	defer stream.Close()

	fmt.Println("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}
		fmt.Println(response.Message)
	}
}
```

### 知识库非流式对话
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.CreateQaChatCompletion(context.Background(),
		chat.QaChatCompletionRequest{
			ChatId:  1765630418937131008,
			Message: "hello",
			Knowledges: []chat.QaChatRequestKnowledges{
				{
					Id:    1761989896572108802,
					Value: "hello",
				},
			},
		})
	if err != nil {
		fmt.Printf("CreateQaChatCompletion error: %v\n", err)
		return
	}

	for _, v := range resp.Data.MessageList {
		fmt.Println(v.Message)
	}
}
```

### 知识库流式对话
```go
package main

import (
	"context"
	"errors"
	"fmt"
	"io"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	stream, err := client.CreateQaChatCompletionStream(context.Background(), chat.QaChatCompletionRequest{
		ChatId:  1765630418937131008,
		Message: "hello",
		Knowledges: []chat.QaChatRequestKnowledges{
			{
				Id:    1761989896572108802,
				Value: "hello",
			},
		},
	})
	if err != nil {
		fmt.Printf("CreateQaChatCompletionStream error: %v\n", err)
		return
	}

	defer stream.Close()

	fmt.Println("Stream response: ")
	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\nStream finished")
			return
		}
		fmt.Println(response.Message)
	}
}
```

### 获取对话历史记录
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.ChatHistory(context.Background(), chat.ChatHistoryRequest{
		ChatId:   1765630418937131008,
		PageSize: 10,
		PageNum:  1,
	})
	if err != nil {
		fmt.Printf("ChatHistory error: %v\n", err)
		return
	}

	for _, v := range resp.Data.HistoryList {
		fmt.Println("Q:", v.Message)
		for _, v := range v.Children {
			fmt.Println("A:", v.Message)
		}
	}
}
```

### 创建知识空间
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.CreateKnowledge(context.Background(), chat.KnowledgeAddRequest{
		KnowledgeName: "test",
		VisibleState:  1,
	})
	if err != nil {
		fmt.Printf("CreateKnowledge error: %v\n", err)
		return
	}

	fmt.Println(resp.KnowledgeId, resp.KnowledgeName, resp.VisibleState)
}
```

### 获取空间列表
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.KnowledgeList(context.Background(), chat.KnowledgeListRequest{})
	if err != nil {
		fmt.Printf("KnowledgeList error: %v\n", err)
		return
	}

	for _, v := range resp.MyList {
		fmt.Println(v.KnowledgeId, v.KnowledgeName, v.VisibleState)
	}
}
```

### 删除知识空间
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	err := client.DeleteKnowledge(context.Background(), "1767799159370416129")
	if err != nil {
		fmt.Printf("DeleteKnowledge error: %v\n", err)
	}
}
```

### 知识库空间上传文件
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.UploadFile(context.Background(), "/home/ehlxr/README.md")
	if err != nil {
		fmt.Printf("UploadFile error: %v\n", err)
		return
	}

	for _, file := range resp {
		fmt.Println(file.OriginalFilename, file.PreviewUrl, file.Size, file.CustomFileType)
	}
}
```

### 知识库文件学习
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	err := client.Learning(context.Background(), chat.KnowledgeFileLearningRequest{
		KnowledgeId: "1765698335274041345",
		Files: []chat.KnowledgeFileLearningInfo{{
			OriginUrl:      "http://oss.xxxxxx.com/test/20240308/d41d8cd98f00b204e9800998ecf8427e.md",
			Size:           "0",
			CustomFileType: "NORMAL",
			Name:           "README.md",
		}},
	})
	if err != nil {
		fmt.Printf("Learning error: %v\n", err)
	}
}
```

### 获取知识库文件列表

```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	resp, err := client.Files(context.Background(), "test-11111")

	for _, v := range resp {
		fmt.Println(v.KnowledgeId, v.FileId, v.Name, v.Size, v.CustomFileType)
	}
	if err != nil {
		fmt.Printf("Files error: %v\n", err)
	}
}
```

### 删除知识库文件
```go
package main

import (
	"context"
	"fmt"

	chat "thundersoft.com/brainos/chat"
)

func main() {
	client := chat.NewClient("baseURL", "your token", "appId", "userId")

	err := client.DeleteFiles(context.Background(), chat.KnowledgeFileDeleteRequest{
		KnowledgeId: "1765698335274041345",
		FileIds:     []string{"1765929183516688385"},
	})
	if err != nil {
		fmt.Printf("DeleteFiles error: %v\n", err)
	}
}
```