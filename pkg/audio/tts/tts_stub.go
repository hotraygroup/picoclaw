//go:build !voice

package tts

import (
	"context"

	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/media"
)

// TTSProvider is a stub interface when voice support is disabled.
type TTSProvider interface {
	Name() string
}

// stubTTSProvider is a no-op TTS provider returned when voice is disabled.
type stubTTSProvider struct{}

func (s *stubTTSProvider) Name() string {
	return "stub"
}

// DetectTTS returns nil when voice support is disabled.
func DetectTTS(cfg *config.Config) TTSProvider {
	return nil
}

// SynthesizeAndStore returns an error when voice support is disabled.
func SynthesizeAndStore(
	ctx context.Context,
	provider TTSProvider,
	store media.MediaStore,
	text string,
	filename string,
	channel string,
	chatID string,
) (string, error) {
	return "", context.Canceled
}