package utils

import (
	"context"
	"fmt"
	"log"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func InputToText(s string) (string, error) {
	return "Every minute", nil
}

type NaturalLanguageProcessor struct {
	llm    *ollama.LLM
	ctx    context.Context
	cancel context.CancelFunc
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

func (n NaturalLanguageProcessor) TextToCronjobExpression(text string) string {
	completion, _ := n.llm.Call(n.ctx, fmt.Sprintf("You are given this phrase \"%s\". If the phrase does not specify a date and time then say \"error\". If there is a date and time, write the cronjob expression using the phrase without any user or command. Do not provide any explanation either.", text),
		llms.WithTemperature(0.1),
	)

	return completion
}
