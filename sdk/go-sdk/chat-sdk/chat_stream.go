package chat

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	utils "thundersoft.com/brainos/chat/internal"
)

var headerData = []byte("data: ")

type ChatCompletionStreamResponse struct {
	Status string `json:"status"`
	ChatId int64  `json:"chatId"`

	MsgId       int64  `json:"msgId"`
	ParentMsgId int64  `json:"parentMsgId"`
	Message     string `json:"chunkMessage"`
}

type ChatCompletionStream struct {
	isFinished bool

	reader         *bufio.Reader
	response       *http.Response
	errAccumulator utils.ErrorAccumulator
	unmarshaler    utils.Unmarshaler

	httpHeader
}

func (stream *ChatCompletionStream) Recv() (response ChatCompletionStreamResponse, err error) {
	if stream.isFinished {
		err = io.EOF
		return
	}

	response, err = stream.processLines()
	return
}

func (stream *ChatCompletionStream) processLines() (ChatCompletionStreamResponse, error) {
	var emptyMessagesCount uint

	for {
		rawLine, readErr := stream.reader.ReadBytes('\n')
		// log.Println(string(rawLine))
		if readErr != nil {
			respErr := stream.unmarshalError()
			if respErr != nil {
				return *new(ChatCompletionStreamResponse), fmt.Errorf("error, %w", respErr.Error)
			}
			return *new(ChatCompletionStreamResponse), readErr
		}

		noSpaceLine := bytes.TrimSpace(rawLine)
		if !bytes.HasPrefix(noSpaceLine, headerData) {
			writeErr := stream.errAccumulator.Write(noSpaceLine)
			if writeErr != nil {
				return *new(ChatCompletionStreamResponse), writeErr
			}
			emptyMessagesCount++

			continue
		}
		noPrefixLine := bytes.TrimPrefix(noSpaceLine, headerData)

		var response ChatCompletionStreamResponse
		unmarshalErr := stream.unmarshaler.Unmarshal(noPrefixLine, &response)
		if unmarshalErr != nil {
			return *new(ChatCompletionStreamResponse), unmarshalErr
		}

		if response.Status == "done" {
			stream.isFinished = true
			return *new(ChatCompletionStreamResponse), io.EOF
		}

		return response, nil
	}
}

func (stream *ChatCompletionStream) unmarshalError() (errResp *ErrorResponse) {
	errBytes := stream.errAccumulator.Bytes()
	if len(errBytes) == 0 {
		return
	}

	err := stream.unmarshaler.Unmarshal(errBytes, &errResp)
	if err != nil {
		errResp = nil
	}

	return
}

func (stream *ChatCompletionStream) Close() {
	stream.response.Body.Close()
}

func (c *Client) CreateChatCompletionStream(
	ctx context.Context,
	request ChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	request.Stream = true
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(chatCompletionsSuffix), withBody(request))
	if err != nil {
		return nil, err
	}

	stream, err = sendRequestStream(c, req)
	if err != nil {
		return
	}
	return
}

func (c *Client) CreateQaChatCompletionStream(
	ctx context.Context,
	request QaChatCompletionRequest,
) (stream *ChatCompletionStream, err error) {
	request.Stream = true
	req, err := c.newRequest(ctx, http.MethodPost, c.fullURL(qaChatCompletionsSuffix), withBody(request))
	if err != nil {
		return nil, err
	}

	stream, err = sendRequestStream(c, req)
	if err != nil {
		return
	}
	return
}
