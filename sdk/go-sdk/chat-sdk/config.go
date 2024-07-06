package chat

import (
	"net/http"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	AuthToken  string
	BaseURL    string
	AppId      string
	UserId     string
	HTTPClient *http.Client
}

func DefaultConfig(baseURL, authToken, appId, userId string) ClientConfig {
	return ClientConfig{
		AuthToken:  authToken,
		BaseURL:    baseURL,
		AppId:      appId,
		UserId:     userId,
		HTTPClient: &http.Client{},
	}
}

func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}
