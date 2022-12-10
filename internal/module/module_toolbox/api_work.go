package module_toolbox

type WorkRequest struct {
	ToolboxId   int64                  `json:"toolboxId,omitempty"`
	ToolboxType string                 `json:"toolboxType,omitempty"`
	Work        string                 `json:"work,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}
