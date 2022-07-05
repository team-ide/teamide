package module_toolbox

import (
	"encoding/json"
	"errors"
	"teamide/pkg/toolbox"
)

// Work 执行
func (this_ *ToolboxService) Work(toolboxId int64, work string, data map[string]interface{}) (res interface{}, err error) {

	toolboxData, err := this_.Get(toolboxId)
	if err != nil {
		return
	}
	if toolboxData == nil {
		toolboxData = this_.GetOtherToolbox(toolboxId)
	}
	if toolboxData == nil {
		err = errors.New("工具配置丢失")
		return
	}

	option := map[string]interface{}{}
	if toolboxData.Option != "" {
		err = json.Unmarshal([]byte(toolboxData.Option), &option)
		if err != nil {
			return
		}
	}

	toolboxWorker := toolbox.GetWorker(toolboxData.ToolboxType)
	if toolboxWorker == nil {
		err = errors.New("不支持的工具类型[" + toolboxData.ToolboxType + "]")
		return
	}

	publicKey, publicKeyOk := option["publicKey"]
	if publicKeyOk && publicKey != nil && publicKey != "" {
		option["publicKey"] = this_.GetFilesFile(publicKey.(string))
	}

	res, err = toolboxWorker.Work(work, option, data)
	if err != nil {
		return
	}

	return
}
