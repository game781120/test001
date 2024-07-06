package openai

import (
	"context"
	"net/http"
)

type ModelInfo struct {
	// Id    uint64 `json:"id"`
	Model string `json:"nickName"`
}

type ModelResponse struct {
	Content []ModelInfo `json:"context"`
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	httpHeader
}

func (c *Client) ListModels(ctx context.Context) (response ModelResponse, err error) {
	urlSuffix := "/brain/billing/v1/models/online"

	req, err := c.newRequest(ctx, http.MethodGet, c.fullURL(urlSuffix))
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
