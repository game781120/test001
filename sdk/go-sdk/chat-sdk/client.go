package chat

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	utils "thundersoft.com/brainos/chat/internal"
)

// Client is OpenAI API client.
type Client struct {
	config ClientConfig

	requestBuilder    utils.RequestBuilder
	createFormBuilder func(io.Writer) utils.FormBuilder
}

type Response interface {
	SetHeader(http.Header)
}

type httpHeader http.Header

func (h *httpHeader) SetHeader(header http.Header) {
	*h = httpHeader(header)
}

// NewClient creates new OpenAI API client.
func NewClient(baseURL, authToken, appId, userId string) *Client {
	config := DefaultConfig(strings.TrimSuffix(baseURL, "/"), authToken, appId, userId)
	return NewClientWithConfig(config)
}

// NewClientWithConfig creates new OpenAI API client for specified config.
func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config:         config,
		requestBuilder: utils.NewRequestBuilder(),
		createFormBuilder: func(body io.Writer) utils.FormBuilder {
			return utils.NewFormBuilder(body)
		},
	}
}

type requestOptions struct {
	body   any
	header http.Header
}

type requestOption func(*requestOptions)

func withBody(body any) requestOption {
	return func(args *requestOptions) {
		args.body = body
	}
}

type Void struct{}

type responseInfo interface {
	[]File | KnowledgeInfoResult | KnowledgeListResult | []KnowledgeFileListResult | ChatInfo |
		ChatCompletionResult | []ChatInfo | Void | ChatHistoryListResult | AppInfo | []KnowledgeFileLearningResult | GetKnowledgeInfoResult
}

type responseCode interface {
	uint32 | string
}

type CommonResponse[T responseInfo, K responseCode] struct {
	Data    T      `json:"data"`
	Code    K      `json:"code"`
	Message string `json:"message"`
	httpHeader
}

func (c *Client) newRequest(ctx context.Context, method, url string, setters ...requestOption) (*http.Request, error) {
	// Default Options
	args := &requestOptions{
		body:   nil,
		header: make(http.Header),
	}
	for _, setter := range setters {
		setter(args)
	}
	req, err := c.requestBuilder.Build(ctx, method, url, args.body, args.header)
	if err != nil {
		return nil, err
	}
	c.setCommonHeaders(req)
	return req, nil
}

func (c *Client) sendRequest(req *http.Request, v Response) error {
	req.Header.Set("Accept", "application/json")

	// Check whether Content-Type is already set, Upload Files API requires
	// Content-Type == multipart/form-data
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	if req.Header.Get("Accept-Language") == "" {
		req.Header.Set("Accept-Language", "zh-CN")
	}

	res, err := c.config.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if isFailureStatusCode(res) {
		return c.handleErrorResp(res)
	}

	if v != nil {
		v.SetHeader(res.Header)
	}

	return decodeResponse(res.Body, v)
}

func sendRequestStream(client *Client, req *http.Request) (*ChatCompletionStream, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.config.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if isFailureStatusCode(resp) {
		return nil, client.handleErrorResp(resp)
	}
	return &ChatCompletionStream{
		reader:         bufio.NewReader(resp.Body),
		response:       resp,
		errAccumulator: utils.NewErrorAccumulator(),
		unmarshaler:    &utils.JSONUnmarshaler{},
		httpHeader:     httpHeader(resp.Header),
	}, nil
}

func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("Authorization", c.config.AuthToken)
	req.Header.Set("appId", c.config.AppId)
	req.Header.Set("userId", c.config.UserId)
}

func isFailureStatusCode(resp *http.Response) bool {
	return resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest
}

func decodeResponse(body io.Reader, v any) error {
	if v == nil {
		return nil
	}

	switch o := v.(type) {
	case *string:
		return decodeString(body, o)
	default:
		// b, _ := io.ReadAll(body)
		// body = bytes.NewBuffer(b)
		// log.Println(string(b))
		return json.NewDecoder(body).Decode(v)
	}
}

func decodeString(body io.Reader, output *string) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	*output = string(b)
	return nil
}

// fullURL returns full URL for request.
func (c *Client) fullURL(suffix string) string {
	return fmt.Sprintf("%s%s", c.config.BaseURL, suffix)
}

func (c *Client) handleErrorResp(resp *http.Response) error {
	var errRes ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errRes)
	if err != nil || errRes.Error == nil {
		reqErr := &RequestError{
			HTTPStatusCode: resp.StatusCode,
			Err:            err,
		}
		if errRes.Error != nil {
			reqErr.Err = errRes.Error
		}
		return reqErr
	}

	errRes.Error.HTTPStatusCode = resp.StatusCode
	return errRes.Error
}

func withContentType(contentType string) requestOption {
	return func(args *requestOptions) {
		args.header.Set("Content-Type", contentType)
	}
}
