package glu

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

func OpenAI(model string) ClientOption {
	return func(c *Client) {
		c.config = ClientConfig{
			APIKey:     os.Getenv("OPENAI_API_KEY"),
			RequestURL: "https://api.openai.com/v1/chat/completions",
			Model:      model,
		}
		c.MessageDecoder = &OpenAIChoiceResponseBody{}
		c.StreamingMessageDecoder = &OpenAIStreamResponseBody{}
		c.BuildRequest = OpenAIRequestBuilder()
	}
}

func OpenAIRequestBuilder() RequestBuilder {
	return func(config ClientConfig, thread *Thread) (*http.Request, error) {
		body := map[string]any{
			"model":       config.Model,
			"stream":      config.Stream,
			"temperature": config.Temperature,
		}
		if config.ResponseFormat == "json" {
			body["response_format"] = map[string]any{
				"type": "json_object",
			}
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

type OpenAIChoiceResponseBody struct {
	Choices []MessageResponseBody `json:"choices"`
}

func (c *OpenAIChoiceResponseBody) Decode(bytes []byte) Message {
	json.Unmarshal(bytes, c)
	return c.Choices[0].Message
}

type OpenAIStreamChoiceResponseBody struct {
	Message Message `json:"delta"`
}

type OpenAIStreamResponseBody struct {
	Choices []OpenAIStreamChoiceResponseBody `json:"choices"`
}

func (c *OpenAIStreamResponseBody) Decode(bytes []byte) Message {
	json.Unmarshal(bytes, c)
	return c.Choices[0].Message
}
