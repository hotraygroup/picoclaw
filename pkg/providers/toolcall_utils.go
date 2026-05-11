package providers

import (
	"encoding/json"
)

func NormalizeToolCall(tc ToolCall) ToolCall {
	normalized := tc

	if normalized.ThoughtSignature == "" &&
		normalized.ExtraContent != nil &&
		normalized.ExtraContent.Google != nil {
		normalized.ThoughtSignature = normalized.ExtraContent.Google.ThoughtSignature
	}

	if normalized.Name == "" && normalized.Function != nil {
		normalized.Name = normalized.Function.Name
	}

	if normalized.Arguments == nil {
		normalized.Arguments = map[string]any{}
	}

	if len(normalized.Arguments) == 0 && normalized.Function != nil && normalized.Function.Arguments != "" {
		var parsed map[string]any
		if err := json.Unmarshal([]byte(normalized.Function.Arguments), &parsed); err == nil && parsed != nil {
			normalized.Arguments = parsed
		}
	}

	argsJSON, _ := json.Marshal(normalized.Arguments)
	if normalized.Function == nil {
		normalized.Function = &FunctionCall{
			Name:             normalized.Name,
			Arguments:        string(argsJSON),
			ThoughtSignature: normalized.ThoughtSignature,
		}
	} else {
		if normalized.Function.Name == "" {
			normalized.Function.Name = normalized.Name
		}
		if normalized.Name == "" {
			normalized.Name = normalized.Function.Name
		}
		if normalized.Function.Arguments == "" {
			normalized.Function.Arguments = string(argsJSON)
		}
		if normalized.Function.ThoughtSignature == "" {
			normalized.Function.ThoughtSignature = normalized.ThoughtSignature
		}
		if normalized.ThoughtSignature == "" {
			normalized.ThoughtSignature = normalized.Function.ThoughtSignature
		}
	}

	return normalized
}
