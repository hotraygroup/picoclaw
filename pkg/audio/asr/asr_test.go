package asr

import (
	"testing"

	"github.com/sipeed/picoclaw/pkg/config"
)

func TestDetectTranscriber(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.Config
		wantNil  bool
		wantName string
	}{
		{
			name:    "no config",
			cfg:     &config.Config{},
			wantNil: true,
		},
		{
			name: "voice model name selects audio model transcriber",
			cfg: &config.Config{
				Voice: config.VoiceConfig{ModelName: "voice-test"},
				ModelList: []*config.ModelConfig{
					{
						ModelName: "voice-test",
						Model:     "openai/gpt-4o-audio-preview",
						APIKeys:   config.SimpleSecureStrings("sk-openai-model"),
					},
				},
			},
			wantName: "audio-model",
		},
		{
			name: "voice model name alias selects whisper transcriber for groq",
			cfg: &config.Config{
				Voice: config.VoiceConfig{ModelName: "my-asr-model"},
				ModelList: []*config.ModelConfig{
					{
						ModelName: "my-asr-model",
						Model:     "groq/whisper-large-v3",
						APIKeys:   config.SimpleSecureStrings("sk-groq-model"),
					},
				},
			},
			wantName: "whisper",
		},
		{
			name: "openai whisper alias selects whisper transcriber",
			cfg: &config.Config{
				Voice: config.VoiceConfig{ModelName: "my-asr-model"},
				ModelList: []*config.ModelConfig{
					{
						ModelName: "my-asr-model",
						Model:     "openai/whisper-1",
						APIKeys:   config.SimpleSecureStrings("sk-openai-model"),
					},
				},
			},
			wantName: "whisper",
		},
		{
			name: "whisper via model list fallback",
			cfg: &config.Config{
				ModelList: []*config.ModelConfig{
					{ModelName: "openai", Model: "openai/gpt-4o", APIKeys: config.SimpleSecureStrings("sk-openai")},
					{
						ModelName: "groq",
						Model:     "groq/whisper-large-v3-turbo",
						APIKeys:   config.SimpleSecureStrings("sk-groq-model"),
					},
				},
			},
			wantName: "whisper",
		},
		{
			name: "groq model list entry without key is skipped",
			cfg: &config.Config{
				ModelList: []*config.ModelConfig{
					{Model: "groq/whisper-large-v3"},
				},
			},
			wantNil: true,
		},
		{
			name: "provider key takes priority over model list",
			cfg: &config.Config{
				ModelList: []*config.ModelConfig{
					{
						ModelName: "groq",
						Model:     "groq/whisper-large-v3",
						APIKeys:   config.SimpleSecureStrings("sk-groq-model"),
					},
				},
			},
			wantName: "whisper",
		},
		{
			name: "missing voice model name config returns nil",
			cfg: &config.Config{
				Voice: config.VoiceConfig{ModelName: "missing"},
				ModelList: []*config.ModelConfig{
					{
						ModelName: "other",
						Model:     "openai/gpt-4o",
						APIKeys:   config.SimpleSecureStrings("sk-other-model"),
					},
				},
			},
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tr := DetectTranscriber(tc.cfg)
			if tc.wantNil {
				if tr != nil {
					t.Errorf("DetectTranscriber() = %v, want nil", tr)
				}
				return
			}
			if tr == nil {
				t.Fatal("DetectTranscriber() = nil, want non-nil")
			}
			if got := tr.Name(); got != tc.wantName {
				t.Errorf("Name() = %q, want %q", got, tc.wantName)
			}
		})
	}
}
