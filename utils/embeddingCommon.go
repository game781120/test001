package utils

type KnowledgeData struct {
	DocumentContent string `json:"document_content"`
}
type KnowledgeResponse struct {
	Code int `json:"code"`
	Data struct {
		DataList []KnowledgeData `json:"data_list"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

type EmbeddingData struct {
	Value [][]float32 `json:"value"`
}

type EmbeddingResponseEx struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data EmbeddingData `json:"data"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}
