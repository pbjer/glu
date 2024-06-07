package glu

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func Ollama(model string) ClientOption {
	return func(c *Client) {
		c.config = ClientConfig{
			RequestURL: "http://localhost:11434/api/chat",
			Model:      model,
		}

		c.MessageDecoder = &OllamaMessageResponseBody{}
		c.StreamingMessageDecoder = &OllamaStreamMessageResponseBody{}
		c.EmbeddingDecoder = &OllamaEmbeddingResponseBody{}

		c.BuildRequest = OllamaRequestBuilder()
		c.BuildEmbeddingRequest = OllamaEmbeddingRequestBuilder()
	}
}

func OllamaEmbeddingRequestBuilder() EmbeddingRequestBuilder {
	return func(config ClientConfig, text string) (*http.Request, error) {
		body := map[string]any{
			"model":  config.Model,
			"prompt": text,
		}

		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, "http://localhost:11434/api/embeddings", bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+config.APIKey)

		return req, nil
	}
}

func OllamaRequestBuilder() RequestBuilder {
	return func(config ClientConfig, thread *Thread) (*http.Request, error) {
		body := map[string]any{
			"model":       config.Model,
			"stream":      config.Stream,
			"temperature": config.Temperature,
		}
		if config.ResponseFormat == "json" {
			body["format"] = "json"
			thread.LastMessage().Content = thread.LastMessage().Content + "\nRespond with JSON."
		}
		body["messages"] = thread.Messages
		reqBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest(http.MethodPost, config.RequestURL, bytes.NewBuffer(reqBody))
		if err != nil {
			return nil, err
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+config.APIKey)

		return req, nil
	}
}

type OllamaMessageResponseBody struct {
	Message Message `json:"message"`
}

func (m *OllamaMessageResponseBody) Decode(bytes []byte) Message {
	json.Unmarshal(bytes, m)
	return m.Message
}

type OllamaStreamMessageResponseBody struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func (m *OllamaStreamMessageResponseBody) Decode(bytes []byte) StreamMessage {
	json.Unmarshal(bytes, m)
	return StreamMessage{
		Message: m.Message,
		Done:    m.Done,
	}
}

type OllamaEmbeddingResponseBody struct {
	Embedding []float32 `json:"embedding"`
}

func (m *OllamaEmbeddingResponseBody) Decode(bytes []byte) Embedding {
	json.Unmarshal(bytes, m)
	return m.Embedding
}
