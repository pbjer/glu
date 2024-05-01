# glu
Go LLM Utility

Super simple chat completion support for Ollama, Groq, and OpenAI.

```go
package main

import (
	"fmt"
	l "github.com/pbjer/glu"
)

func main() {
	thread := l.NewThread(
		l.SystemMessage("When the user says 'hello', you say 'world'."),
		l.UserMessage("hello"))

	client := l.NewClient(
		l.OpenAI("gpt-4-turbo"),
		l.WithAPIKey("YOUR-API-KEY"))

	err := client.Generate(thread)
	if err != nil {
		panic(err)
	}
	fmt.Println(thread.LastMessage().Content)
	// world
	
	thread.AddMessages(l.UserMessage("hello"))
	
	err = client.GenerateStream(thread, func(response l.StreamMessage) error {
		fmt.Println(response.Message.Content)
		return nil
	})
	// wor
	// ld
}
```