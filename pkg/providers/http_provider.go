package providers

import (
	"context"
	"time"

	"github.com/sipeed/picoclaw/pkg/providers/openai_compat"
)

type HTTPProvider struct {
	delegate *openai_compat.Provider
}

func NewHTTPProvider(apiKey, apiBase, proxy string) *HTTPProvider {
	return &HTTPProvider{
		delegate: openai_compat.NewProvider(apiKey, apiBase, proxy),
	}
}

func NewHTTPProviderWithMaxTokensField(apiKey, apiBase, proxy, maxTokensField string) *HTTPProvider {
	return NewHTTPProviderWithMaxTokensFieldAndRequestTimeout(apiKey, apiBase, proxy, maxTokensField, "", 0, nil, nil)
}

func NewHTTPProviderWithMaxTokensFieldAndRequestTimeout(
	apiKey, apiBase, proxy, maxTokensField, userAgent string,
	requestTimeoutSeconds int,
	extraBody map[string]any,
	customHeaders map[string]string,
) *HTTPProvider {
	return &HTTPProvider{
		delegate: openai_compat.NewProvider(
			apiKey,
			apiBase,
			proxy,
			openai_compat.WithMaxTokensField(maxTokensField),
			openai_compat.WithRequestTimeout(time.Duration(requestTimeoutSeconds)*time.Second),
			openai_compat.WithExtraBody(extraBody),
			openai_compat.WithCustomHeaders(customHeaders),
			openai_compat.WithUserAgent(userAgent),
		),
	}
}

func (p *HTTPProvider) Chat(
	ctx context.Context,
	messages []Message,
	tools []ToolDefinition,
	model string,
	options map[string]any,
) (*LLMResponse, error) {
	return p.delegate.Chat(ctx, messages, tools, model, options)
}

func (p *HTTPProvider) ChatStream(
	ctx context.Context,
	messages []Message,
	tools []ToolDefinition,
	model string,
	options map[string]any,
	onChunk func(accumulated string),
) (*LLMResponse, error) {
	return p.delegate.ChatStream(ctx, messages, tools, model, options, onChunk)
}

func (p *HTTPProvider) GetDefaultModel() string {
	return ""
}

func (p *HTTPProvider) SupportsNativeSearch() bool {
	return p.delegate.SupportsNativeSearch()
}

func (p *HTTPProvider) SetProviderName(providerName string) {
	if p == nil || p.delegate == nil {
		return
	}
	p.delegate.SetProviderName(providerName)
}
