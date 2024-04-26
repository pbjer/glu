# glu
Go LLM Utility

Super simple chat completion support for Ollama, Groq, and OpenAI.

```go
package main

import (
	"fmt"
	"github.com/pbjer/glu"
)

func main() {
	thread := glu.NewThread(
		glu.SystemMessage("When the user says 'hello', you say 'world'."),
		glu.UserMessage("hello"))

	client := glu.NewClient(
		glu.OpenAI("gpt-4-turbo"),
		glu.WithAPIKey("YOUR-API-KEY"))

	err := client.Generate(thread)
	if err != nil {
		panic(err)
	}
	fmt.Println(thread.LastMessage().Content)
	// world
	
	thread.AddMessages(glu.UserMessage("hello"))
	
	err := client.GenerateStream(thread, func(response glu.StreamResponse) error {
		fmt.Println(response.Message.Content)
	})
	// wor
	// ld
}
```