package chat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

var ErrChatCompletionStreamNotSupported = errors.New("streaming is not supported with this method, please use CreateChatCompletionStream") //nolint:lll

const (
	chatCompletionsSuffix   = "/brain/brain/api-public/v1/completions"
	qaChatCompletionsSuffix = "/brain/brain/api-public/v1/qa/completions"
)

type ChatHistoryListResult struct {
	HistoryList []ChatHistoryResult `json:"list"`
}

type ChatHistoryResult struct {
	Id         string              `json:"id"`
	ParentId   string              `json:"parentId"`
	ChildrenId []string            `json:"childrenId"`
	Children   []ChatHistoryResult `json:"children"`
	Role       string              `json:"role"`
	Message    string              `json:"message"`
	Evaluate   uint32              `json:"evaluate"`
	CreateTime int64               `json:"createTime"`
}

type ChatCompletionRequest struct {
	ChatId        uint64 `json:"chatId"`
	ReGenerate    uint32 `json:"reGenerate"`
	MessagesId    uint64 `json:"messagesId"`
	Message       string `json:"message"`
	IgnoreHistory bool   `json:"ignoreHistory"`
	Stream        bool   `json:"stream"`
}
type ChatHistoryRequest struct {
	ChatId   uint64 `json:"chatId"`
	PageSize uint32 `json:"pageSize"`
	Id       uint64 `json:"id", default:0`
}

type QaChatCompletionRequest struct {
	ChatId        uint64                    `json:"chatId"`
	ReGenerate    uint32                    `json:"reGenerate"`
	MessagesId    uint64                    `json:"messagesId"`
	Message       string                    `json:"message"`
	IgnoreHistory bool                      `json:"ignoreHistory"`
	Stream        bool                      `json:"stream"`
	Knowledges    []QaChatRequestKnowledges `json:"knowledges"`
}

type QaChatRequestKnowledges struct {
	Id    string `json:"id"`
	Value string `json:"value"`
}

type ChatCompletionMessage struct {
	MessageId string `json:"messageId"`
	Message   string `json:"message"`
}

type ChatCompletionResult struct {
	ChatId           string                  `json:"chatId"`
	ParentMessagesId string                  `json:"parentMessagesId"`
	MessageList      []ChatCompletionMessage `json:"messageList"`
}

type ChatInfo struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AppInfo struct {
	Id             string                    `json:"id"`
	Name           string                    `json:"name"`
	Stream         bool                      `json:"stream"`
	Model          string                    `json:"model"`
	Role           string                    `json:"role"`
	HaveKnowledges bool                      `json:"haveKnowledges"`
	Knowledges     []QaChatRequestKnowledges `json:"knowledges"`
	OverstepType   int                       `json:"overstepType"`
}

func (c *Client) AddChat(ctx context.Context) (resp CommonResponse[ChatInfo, uint32], err error) {
	urlSuffix := "/brain/brain/api-public/v1/chat"

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &resp)
	return
}

func (c *Client) CreateChatCompletion(ctx context.Context,
	request ChatCompletionRequest,
) (response CommonResponse[ChatCompletionResult, uint32], err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(chatCompletionsSuffix),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) CreateQaChatCompletion(ctx context.Context,
	request QaChatCompletionRequest,
) (response CommonResponse[ChatCompletionResult, uint32], err error) {
	if request.Stream {
		err = ErrChatCompletionStreamNotSupported
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(qaChatCompletionsSuffix),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

func (c *Client) ChatList(ctx context.Context) (resp CommonResponse[[]ChatInfo, uint32], err error) {
	urlSuffix := "/brain/brain/api-public/v1/chat/list"

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		c.fullURL(urlSuffix),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &resp)
	return
}

func (c *Client) DeleteChat(ctx context.Context, chatId uint64) (err error) {
	urlSuffix := "/brain/brain/api-public/v1/chat/%d"

	req, err := c.newRequest(
		ctx,
		http.MethodDelete,
		c.fullURL(fmt.Sprintf(urlSuffix, chatId)),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, nil)
	return
}

func (c *Client) ChatHistory(ctx context.Context, request ChatHistoryRequest) (resp CommonResponse[ChatHistoryListResult, uint32], err error) {
	urlSuffix := "/brain/brain/api-public/v1/history/%d?pageSize=%d&id=%d"

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		c.fullURL(fmt.Sprintf(urlSuffix, request.ChatId, request.PageSize, request.Id)),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &resp)
	return
}

func (c *Client) GetAppInfo(ctx context.Context) (resp CommonResponse[AppInfo, uint32], err error) {
	urlSuffix := "/brain/brain/api-public/v1/app/%s"

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		c.fullURL(fmt.Sprintf(urlSuffix, c.config.AppId)),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &resp)
	return
}
