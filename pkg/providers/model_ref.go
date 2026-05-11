package providers

import "strings"

// ModelRef represents a parsed model reference with provider and model name.
type ModelRef struct {
	Provider string
	Model    string
}

// ParseModelRef parses "anthropic/claude-opus" into {Provider: "anthropic", Model: "claude-opus"}.
// If no slash present, uses defaultProvider.
// Returns nil for empty input.
func ParseModelRef(raw string, defaultProvider string) *ModelRef {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	provider, model := SplitModelProviderAndID(raw, defaultProvider)
	if model == "" {
		return nil
	}
	return &ModelRef{
		Provider: provider,
		Model:    model,
	}
}

// NormalizeProvider normalizes provider identifiers to canonical form.
func NormalizeProvider(provider string) string {
	p := strings.ToLower(strings.TrimSpace(provider))

	switch p {
	case "z.ai", "z-ai":
		return "zai"
	case "opencode-zen":
		return "opencode"
	case "qwen":
		return "qwen-portal"
	case "gpt":
		return "openai"
	case "glm":
		return "zhipu"
	case "alibaba-coding", "qwen-coding":
		return "coding-plan"
	case "qwen-international", "dashscope-intl":
		return "qwen-intl"
	case "dashscope-us":
		return "qwen-us"
	}

	return p
}

// ModelKey returns a canonical "provider/model" key for deduplication.
func ModelKey(provider, model string) string {
	return NormalizeProvider(provider) + "/" + strings.ToLower(strings.TrimSpace(model))
}
