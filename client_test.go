package glu

import (
	"fmt"
	"testing"
)

func TestEmbeddings(t *testing.T) {
	e, err := NewClient(Ollama("nomic-embed-text")).GenerateEmbedding("hello")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(e)
}
