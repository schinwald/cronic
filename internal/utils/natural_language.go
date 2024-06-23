package utils

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms/ollama"
)

func InputToText(s string) (string, error) {
	return "Every minute", nil
}

type NaturalLanguageProcessor struct {
	llm               *ollama.LLM
	ctx               context.Context
	cancel            context.CancelFunc
	generatedCallback func()
}

func MakeNaturalLanguageProcessor() *NaturalLanguageProcessor {
	llm, err := ollama.New(ollama.WithModel("llama3"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &NaturalLanguageProcessor{
		llm:    llm,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (n *NaturalLanguageProcessor) OnGenerated(callback func()) {
	n.generatedCallback = callback
}

func (n NaturalLanguageProcessor) Done() <-chan struct{} {
	return n.ctx.Done()
}

func (n *NaturalLanguageProcessor) Cancel() {
	n.cancel()

	// Get a new context to use
	ctx, cancel := context.WithCancel(context.Background())
	n.ctx = ctx
	n.cancel = cancel
}

func (n NaturalLanguageProcessor) TextToCronjobExpression(text string) (string, error) {
	channel := make(chan string)

	go func() {
		completion, _ := n.llm.Call(n.ctx, fmt.Sprintf(`Create a cronjob expression from the input "%s" without any explanation.`, text))
		channel <- completion
	}()

	select {
	case <-n.ctx.Done():
		return "", errors.New("cancelled generating cronjob expression from text")
	case result := <-channel:
		n.generatedCallback()
		return result, nil
	}
}
