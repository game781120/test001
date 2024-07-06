package openai

import (
	"net/http"
)

// ClientConfig is a configuration of a client.
type ClientConfig struct {
	authToken  string
	BaseURL    string
	HTTPClient *http.Client
}

func DefaultConfig(baseURL, authToken string) ClientConfig {
	return ClientConfig{
		authToken: authToken,
		BaseURL:   baseURL,

		HTTPClient: &http.Client{},
	}
}

func (ClientConfig) String() string {
	return "<OpenAI API ClientConfig>"
}
