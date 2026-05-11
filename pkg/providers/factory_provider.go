package providers

import (
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/config"
)

type protocolMeta struct {
	defaultAPIBase     string
	emptyAPIKeyAllowed bool
}

var protocolMetaByName = map[string]protocolMeta{
	"openai":                   {defaultAPIBase: "https://api.openai.com/v1"},
	"venice":                   {defaultAPIBase: "https://api.venice.ai/api/v1"},
	"openrouter":               {defaultAPIBase: "https://openrouter.ai/api/v1"},
	"litellm":                  {defaultAPIBase: "http://localhost:4000/v1"},
	"lmstudio":                 {defaultAPIBase: "http://localhost:1234/v1", emptyAPIKeyAllowed: true},
	"novita":                   {defaultAPIBase: "https://api.novita.ai/openai"},
	"groq":                     {defaultAPIBase: "https://api.groq.com/openai/v1"},
	"zhipu":                    {defaultAPIBase: "https://open.bigmodel.cn/api/paas/v4"},
	"nvidia":                   {defaultAPIBase: "https://integrate.api.nvidia.com/v1"},
	"ollama":                   {defaultAPIBase: "http://localhost:11434/v1", emptyAPIKeyAllowed: true},
	"moonshot":                 {defaultAPIBase: "https://api.moonshot.cn/v1"},
	"shengsuanyun":             {defaultAPIBase: "https://router.shengsuanyun.com/api/v1"},
	"deepseek":                 {defaultAPIBase: "https://api.deepseek.com/v1"},
	"cerebras":                 {defaultAPIBase: "https://api.cerebras.ai/v1"},
	"vivgrid":                  {defaultAPIBase: "https://api.vivgrid.com/v1"},
	"volcengine":               {defaultAPIBase: "https://ark.cn-beijing.volces.com/api/v3"},
	"qwen":                     {defaultAPIBase: "https://dashscope.aliyuncs.com/compatible-mode/v1"},
	"qwen-portal":              {defaultAPIBase: "https://dashscope.aliyuncs.com/compatible-mode/v1"},
	"qwen-intl":                {defaultAPIBase: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"},
	"qwen-international":       {defaultAPIBase: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"},
	"dashscope-intl":           {defaultAPIBase: "https://dashscope-intl.aliyuncs.com/compatible-mode/v1"},
	"qwen-us":                  {defaultAPIBase: "https://dashscope-us.aliyuncs.com/compatible-mode/v1"},
	"dashscope-us":             {defaultAPIBase: "https://dashscope-us.aliyuncs.com/compatible-mode/v1"},
	"coding-plan":              {defaultAPIBase: "https://coding-intl.dashscope.aliyuncs.com/v1"},
	"alibaba-coding":           {defaultAPIBase: "https://coding-intl.dashscope.aliyuncs.com/v1"},
	"qwen-coding":              {defaultAPIBase: "https://coding-intl.dashscope.aliyuncs.com/v1"},
	"zai":                      {defaultAPIBase: "https://api.z.ai/api/coding/paas/v4"},
	"vllm":                     {defaultAPIBase: "http://localhost:8000/v1", emptyAPIKeyAllowed: true},
	"mistral":                  {defaultAPIBase: "https://api.mistral.ai/v1"},
	"avian":                    {defaultAPIBase: "https://api.avian.io/v1"},
	"minimax":                  {defaultAPIBase: "https://api.minimaxi.com/v1"},
	"longcat":                  {defaultAPIBase: "https://api.longcat.chat/openai"},
	"modelscope":               {defaultAPIBase: "https://api-inference.modelscope.cn/v1"},
	"mimo":                     {defaultAPIBase: "https://api.xiaomimimo.com/v1"},
}

func ExtractProtocol(cfg *config.ModelConfig) (protocol, modelID string) {
	if cfg == nil {
		return "", ""
	}

	model := strings.TrimSpace(cfg.Model)
	if provider := strings.TrimSpace(cfg.Provider); provider != "" {
		return NormalizeProvider(provider), model
	}
	return SplitModelProviderAndID(model, "openai")
}

func ResolveAPIBase(cfg *config.ModelConfig) string {
	if cfg == nil {
		return ""
	}
	if apiBase := strings.TrimSpace(cfg.APIBase); apiBase != "" {
		return strings.TrimRight(apiBase, "/")
	}
	protocol, _ := ExtractProtocol(cfg)
	return strings.TrimRight(getDefaultAPIBase(protocol), "/")
}

func CreateProviderFromConfig(cfg *config.ModelConfig) (LLMProvider, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("config is nil")
	}

	if cfg.Model == "" {
		return nil, "", fmt.Errorf("model is required")
	}

	protocol, modelID := ExtractProtocol(cfg)

	userAgent := cfg.UserAgent
	if userAgent == "" {
		userAgent = fmt.Sprintf("PicoClaw/%s", config.Version)
	}

	if !IsSupportedModelProvider(protocol) {
		return nil, "", fmt.Errorf("unknown protocol %q in model %q", protocol, cfg.Model)
	}

	if cfg.APIKey() == "" && cfg.APIBase == "" && !isEmptyAPIKeyAllowed(protocol) {
		return nil, "", fmt.Errorf("api_key or api_base is required for HTTP-based protocol %q", protocol)
	}
	apiBase := cfg.APIBase
	if apiBase == "" {
		apiBase = getDefaultAPIBase(protocol)
	}

	extraBody := cfg.ExtraBody
	if protocol == "minimax" {
		if extraBody == nil {
			extraBody = make(map[string]any)
		}
		if _, ok := extraBody["reasoning_split"]; !ok {
			extraBody["reasoning_split"] = true
		}
	}

	provider := NewHTTPProviderWithMaxTokensFieldAndRequestTimeout(
		cfg.APIKey(),
		apiBase,
		cfg.Proxy,
		cfg.MaxTokensField,
		userAgent,
		cfg.RequestTimeout,
		extraBody,
		cfg.CustomHeaders,
	)
	provider.SetProviderName(protocol)
	return finalizeProviderFromConfig(provider, modelID, cfg)
}

func finalizeProviderFromConfig(
	provider LLMProvider,
	modelID string,
	cfg *config.ModelConfig,
) (LLMProvider, string, error) {
	wrapped, err := wrapProviderWithToolSchemaTransform(provider, cfg.ToolSchemaTransform)
	if err != nil {
		return nil, "", err
	}
	return wrapped, modelID, nil
}

func isEmptyAPIKeyAllowed(protocol string) bool {
	meta, ok := protocolMetaForName(protocol)
	return ok && meta.emptyAPIKeyAllowed
}

func IsEmptyAPIKeyAllowedForProtocol(protocol string) bool {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	return isEmptyAPIKeyAllowed(protocol)
}

func DefaultAPIBaseForProtocol(protocol string) string {
	protocol = strings.ToLower(strings.TrimSpace(protocol))
	return getDefaultAPIBase(protocol)
}

func getDefaultAPIBase(protocol string) string {
	meta, ok := protocolMetaForName(protocol)
	if !ok {
		return ""
	}
	return meta.defaultAPIBase
}

func protocolMetaForName(protocol string) (protocolMeta, bool) {
	if meta, ok := protocolMetaByName[protocol]; ok {
		return meta, true
	}
	return protocolMeta{}, false
}
