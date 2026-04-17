//go:build !voice

package asr

import (
	"context"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
)

// TranscriptionResponse is a stub response type when voice support is disabled.
type TranscriptionResponse struct {
	Text     string
	Language string
	Duration float64
}

// Transcriber is a stub interface when voice support is disabled.
type Transcriber interface {
	Name() string
	Transcribe(ctx context.Context, audioFilePath string) (*TranscriptionResponse, error)
}

// stubTranscriber is a no-op transcriber returned when voice is disabled.
type stubTranscriber struct{}

func (s *stubTranscriber) Name() string {
	return "stub"
}

func (s *stubTranscriber) Transcribe(ctx context.Context, audioFilePath string) (*TranscriptionResponse, error) {
	return nil, ctx.Err()
}

// Agent is a stub voice agent when voice support is disabled.
type Agent struct{}

// NewAgent returns a stub agent when voice support is disabled.
func NewAgent(mb *bus.MessageBus, t Transcriber) *Agent {
	return &Agent{}
}

// Start does nothing when voice support is disabled.
func (a *Agent) Start(ctx context.Context) {}

// DetectTranscriber returns nil when voice support is disabled.
func DetectTranscriber(cfg *config.Config) Transcriber {
	return nil
}