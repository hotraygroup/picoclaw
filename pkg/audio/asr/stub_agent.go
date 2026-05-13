//go:build !audio

package asr

import (
	"context"

	"github.com/sipeed/picoclaw/pkg/bus"
)

type speechAccumulator struct{}

func (a *speechAccumulator) Push(chunk bus.AudioChunk) {}
func (a *speechAccumulator) Close()                     {}

type Agent struct{}

func NewAgent(mb *bus.MessageBus, t Transcriber) *Agent {
	return nil
}

func (a *Agent) Start(ctx context.Context) error {
	return nil
}

func (a *Agent) listenChunks(ctx context.Context)   {}
func (a *Agent) handleChunk(chunk bus.AudioChunk)   {}
func (a *Agent) vadTick(ctx context.Context)        {}
func (a *Agent) checkSilence(ctx context.Context)   {}
func (a *Agent) processUtterance(ctx context.Context, acc *speechAccumulator) {
}