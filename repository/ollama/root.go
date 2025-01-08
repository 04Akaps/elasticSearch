package ollama

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
	"log"
)

type Ollama struct {
	llm *ollama.LLM
}

func NewOllama(model, url string) Ollama {
	log.Println("Start to connect ollama")

	opts := []ollama.Option{
		ollama.WithModel(model), // example : "llama3.2"
	}

	// url 이 빈값이면, local 환경 연결
	if url != "" {
		opts = append(opts, ollama.WithServerURL(url))
	}

	llm, err := ollama.New(opts...)

	if err != nil {
		log.Panic("Failed to connect to model", "err", err)
	}

	log.Println("success to connect ollama")

	return Ollama{llm}
}

func (o Ollama) Call(text string) (result string, err error) {
	ctx := context.Background()

	result, err = o.llm.Call(
		ctx,
		text,
		llms.WithTemperature(0.3), // 랜덤성 설정 -> 높을수록 상상력이 풍부
	)

	if err != nil {
		return "", err
	}

	return result, nil
}
