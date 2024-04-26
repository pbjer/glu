package glu

import "os"

func Groq(model string) ClientOption {
	return func(c *Client) {
		c.config = ClientConfig{
			APIKey:     os.Getenv("GROQ_API_KEY"),
			RequestURL: "https://api.groq.com/openai/v1/chat/completions",
			Model:      model,
		}
		c.MessageDecoder = &OpenAIChoiceResponseBody{}
		c.StreamingMessageDecoder = &OpenAIStreamResponseBody{}
		c.BuildRequest = OpenAIRequestBuilder()
	}
}
