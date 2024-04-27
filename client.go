package glu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type MessageDecoder interface {
	Decode([]byte) Message
}

type ClientOption func(c *Client)

func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.config.APIKey = apiKey
	}
}

func WithJSONResponse() ClientOption {
	return func(c *Client) {
		c.config.ResponseFormat = "json"
	}
}

func WithTemperature(tmp float32) ClientOption {
	return func(c *Client) {
		c.config.Temperature = tmp
	}
}

type ClientConfig struct {
	RequestURL     string
	APIKey         string
	Model          string
	ResponseFormat string
	Stream         bool
	Temperature    float32
}

type RequestBuilder func(config ClientConfig, thread *Thread) (*http.Request, error)

type Client struct {
	config                  ClientConfig
	BuildRequest            RequestBuilder
	MessageDecoder          MessageDecoder
	StreamingMessageDecoder MessageDecoder
}

func NewClient(options ...ClientOption) *Client {
	c := &Client{}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Client) Generate(thread *Thread) error {
	req, err := c.BuildRequest(c.config, thread)
	if err != nil {
		return err
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	thread.AddMessages(c.MessageDecoder.Decode(respBody))
	return nil
}

func (c *Client) GenerateStream(thread *Thread, handler StreamResponseHandler) error {
	c.config.Stream = true
	req, err := c.BuildRequest(c.config, thread)
	if err != nil {
		return err
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimPrefix(line, []byte("data: "))
		// Attempt to parse the line as JSON
		var rawMessage json.RawMessage
		err := json.Unmarshal(line, &rawMessage)
		if err != nil {
			continue // Skip lines that cannot be parsed as JSON
		}
		message := c.StreamingMessageDecoder.Decode(rawMessage)
		if err := handler(StreamResponse{Message: message}); err != nil {
			return err // Handle errors from the handler
		}
	}
	if err := scanner.Err(); err != nil {
		return err // Handle scanner errors
	}

	return nil
}

type StreamResponse struct {
	Message Message `json:"message"`
}

type StreamResponseHandler func(StreamResponse) error
