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
		c.MessageDecoder = &MessageResponseBody{}
		c.StreamingMessageDecoder = &MessageResponseBody{}
		c.BuildRequest = OllamaRequestBuilder()
	}
}

func OllamaRequestBuilder() RequestBuilder {
	return func(config ClientConfig, thread *Thread) (*http.Request, error) {
		body := map[string]any{
			"model":  config.Model,
			"stream": config.Stream,
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

type MessageResponseBody struct {
	Message Message `json:"message"`
}

func (m *MessageResponseBody) Decode(bytes []byte) Message {
	json.Unmarshal(bytes, m)
	return m.Message
}
