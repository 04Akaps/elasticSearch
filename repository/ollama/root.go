package ollama

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

type Ollama struct {
	llm *ollama.LLM
}

func NewOllama(model string) Ollama {
	log.Println("Start to connect ollama")

	llm, err := ollama.New(ollama.WithModel(model)) // example : "llama3.2"

	if err != nil {
		log.Panic("Failed to connect to model", "err", err)
	}

	return Ollama{llm}
}

func (o Ollama) Call(text string) (result string, err error) {
	ctx := context.Background()

	result, err = o.llm.Call(ctx, text,
		llms.WithTemperature(0.8),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)

	if err != nil {
		return "", err
	}

	return result, nil
}
