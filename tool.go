package talkative

type Tools []Tool

// func (t Tools) String() string {
// 	json, _ := json.Marshal(t)

// 	return string(json)
// }

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// func (t Tool) String() string {
// 	json, _ := json.Marshal(t)

// 	return string(json)
// }

type ToolCall struct {
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Index     int                       `json:"index,omitempty"`
	Name      string                    `json:"name"`
	Arguments ToolCallFunctionArguments `json:"arguments"`
}

type ToolCallFunctionArguments map[string]any

// func (t *ToolCallFunctionArguments) String() string {
// 	json, _ := json.Marshal(t)

// 	return string(json)
// }

type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  ToolFunctionParameters `json:"parameters"`
}

type ToolFunctionParameters struct {
	Type       string                            `json:"type"`
	Required   []string                          `json:"required"`
	Properties map[string]ToolFunctionProperties `json:"properties"`
}

type ToolFunctionProperties struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Enum        []string `json:"enum,omitempty"`
}

// func (t *ToolFunction) String() string {
// 	json, _ := json.Marshal(t)

// 	return string(json)
// }
