package providers

import (
	"sort"
	"strings"
)

type ModelProviderOption struct {
	ID                  string `json:"id"`
	DefaultAPIBase      string `json:"default_api_base"`
	EmptyAPIKeyAllowed  bool   `json:"empty_api_key_allowed"`
	CreateAllowed       bool   `json:"create_allowed"`
	DefaultModelAllowed bool   `json:"default_model_allowed"`
	DefaultAuthMethod   string `json:"default_auth_method,omitempty"`
	AuthMethodLocked    bool   `json:"auth_method_locked,omitempty"`
}

func ModelProviderOptions() []ModelProviderOption {
	optionsByID := make(map[string]ModelProviderOption, len(protocolMetaByName))
	for provider := range protocolMetaByName {
		if NormalizeProvider(provider) != provider {
			continue
		}
		optionsByID[provider] = ModelProviderOption{
			ID:                  provider,
			DefaultAPIBase:      DefaultAPIBaseForProtocol(provider),
			EmptyAPIKeyAllowed:  IsEmptyAPIKeyAllowedForProtocol(provider),
			CreateAllowed:       true,
			DefaultModelAllowed: true,
		}
	}

	options := make([]ModelProviderOption, 0, len(optionsByID))
	for _, option := range optionsByID {
		options = append(options, option)
	}
	sort.Slice(options, func(i, j int) bool {
		return options[i].ID < options[j].ID
	})
	return options
}

func IsSupportedModelProvider(provider string) bool {
	normalized := NormalizeProvider(provider)
	if normalized == "" {
		return false
	}
	_, ok := protocolMetaByName[normalized]
	return ok
}

func IsCreatableModelProvider(provider string) bool {
	normalized := NormalizeProvider(provider)
	if normalized == "" {
		return false
	}
	_, ok := protocolMetaByName[normalized]
	return ok
}

func IsDefaultModelProvider(provider string) bool {
	normalized := NormalizeProvider(provider)
	if normalized == "" {
		return false
	}
	_, ok := protocolMetaByName[normalized]
	return ok
}

func SplitModelProviderAndID(model, defaultProvider string) (provider, modelID string) {
	model = strings.TrimSpace(model)
	if model == "" {
		return "", ""
	}

	provider, modelID = splitKnownProviderModel(model)
	if provider != "" || modelID != "" {
		return provider, modelID
	}

	return NormalizeProvider(defaultProvider), model
}

func splitKnownProviderModel(model string) (provider, modelID string) {
	provider, modelID, found := strings.Cut(strings.TrimSpace(model), "/")
	if !found {
		return "", ""
	}
	provider = strings.TrimSpace(provider)
	modelID = strings.TrimSpace(modelID)
	if provider == "" {
		return "", modelID
	}
	if !IsSupportedModelProvider(provider) {
		return "", ""
	}
	return NormalizeProvider(provider), modelID
}
