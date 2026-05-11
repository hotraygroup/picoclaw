package providers

import "testing"

func TestParseModelRef_WithSlash(t *testing.T) {
	ref := ParseModelRef("openrouter/auto", "openai")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "openrouter" {
		t.Errorf("provider = %q, want openrouter", ref.Provider)
	}
	if ref.Model != "auto" {
		t.Errorf("model = %q, want auto", ref.Model)
	}
}

func TestParseModelRef_WithoutSlash(t *testing.T) {
	ref := ParseModelRef("gpt-4", "openai")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "openai" {
		t.Errorf("provider = %q, want openai", ref.Provider)
	}
	if ref.Model != "gpt-4" {
		t.Errorf("model = %q, want gpt-4", ref.Model)
	}
}

func TestParseModelRef_Empty(t *testing.T) {
	ref := ParseModelRef("", "openai")
	if ref != nil {
		t.Errorf("expected nil for empty string, got %+v", ref)
	}
}

func TestParseModelRef_EmptyModelAfterSlash(t *testing.T) {
	ref := ParseModelRef("openai/", "default")
	if ref != nil {
		t.Errorf("expected nil for empty model, got %+v", ref)
	}
}

func TestParseModelRef_WhitespaceHandling(t *testing.T) {
	ref := ParseModelRef("  openrouter / auto  ", "openai")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "openrouter" {
		t.Errorf("provider = %q, want openrouter", ref.Provider)
	}
	if ref.Model != "auto" {
		t.Errorf("model = %q, want auto", ref.Model)
	}
}

func TestNormalizeProvider(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"OpenAI", "openai"},
		{"z.ai", "zai"},
		{"z-ai", "zai"},
		{"Z.AI", "zai"},
		{"opencode-zen", "opencode"},
		{"qwen", "qwen-portal"},
		{"gpt", "openai"},
		{"glm", "zhipu"},
		{"groq", "groq"},
		{"alibaba-coding", "coding-plan"},
		{"qwen-coding", "coding-plan"},
		{"qwen-international", "qwen-intl"},
		{"dashscope-intl", "qwen-intl"},
		{"dashscope-us", "qwen-us"},
		{"", ""},
	}

	for _, tt := range tests {
		got := NormalizeProvider(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeProvider(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestModelKey(t *testing.T) {
	tests := []struct {
		provider string
		model    string
		want     string
	}{
		{"openai", "gpt-4", "openai/gpt-4"},
		{"z.ai", "Model-X", "zai/model-x"},
	}

	for _, tt := range tests {
		got := ModelKey(tt.provider, tt.model)
		if got != tt.want {
			t.Errorf("ModelKey(%q, %q) = %q, want %q", tt.provider, tt.model, got, tt.want)
		}
	}
}

func TestParseModelRef_ProviderNormalization(t *testing.T) {
	ref := ParseModelRef("Z.AI/model-x", "default")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "zai" {
		t.Errorf("provider = %q, want zai", ref.Provider)
	}
}

func TestParseModelRef_DefaultProviderNormalization(t *testing.T) {
	ref := ParseModelRef("gpt-4o", "GPT")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "openai" {
		t.Errorf("provider = %q, want openai (normalized from GPT)", ref.Provider)
	}
}

func TestParseModelRef_UnknownPrefixFallsBackToDefaultProvider(t *testing.T) {
	ref := ParseModelRef("meta-llama/Llama-3.1-8B-Instruct", "openai")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "openai" {
		t.Fatalf("provider = %q, want openai", ref.Provider)
	}
	if ref.Model != "meta-llama/Llama-3.1-8B-Instruct" {
		t.Fatalf("model = %q, want full original model ID", ref.Model)
	}
}

func TestParseModelRef_UnknownPrefixPreservesEmptyDefaultProvider(t *testing.T) {
	ref := ParseModelRef("meta-llama/Llama-3.1-8B-Instruct", "")
	if ref == nil {
		t.Fatal("expected non-nil ref")
	}
	if ref.Provider != "" {
		t.Fatalf("provider = %q, want empty", ref.Provider)
	}
	if ref.Model != "meta-llama/Llama-3.1-8B-Instruct" {
		t.Fatalf("model = %q, want full original model ID", ref.Model)
	}
}
