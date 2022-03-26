package module_toolbox

import "errors"

// Work 执行
func (this_ *ToolboxService) Work(toolboxId int64, data map[string]interface{}) (res interface{}, err error) {

	toolbox, err := this_.Get(toolboxId)
	if err != nil {
		return
	}
	if toolbox == nil {
		err = errors.New("工具配置丢失")
		return
	}
	return
}
