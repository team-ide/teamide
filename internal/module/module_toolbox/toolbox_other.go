package module_toolbox

var (
	OtherToolboxId int64 = 9999999910000
	Others         []*ToolboxModel
	OtherMap       = map[int64]*ToolboxModel{}
)

func init() {
	Others = append(Others, &ToolboxModel{
		ToolboxId: OtherToolboxId + 1,
		Name:      "格式转换",
		Option:    `{"type":"format"}`,
	})

	Others = append(Others, &ToolboxModel{
		ToolboxId: OtherToolboxId + 2,
		Name:      "生成数据文件",
		Option:    `{"type":"generateDataFile"}`,
	})

	for _, one := range Others {
		if one.ToolboxId <= OtherToolboxId {
			panic("其它工具中，内置工具ID异常")
		}
		if OtherMap[one.ToolboxId] != nil {
			panic("其它工具中，存在相同工具ID")
		}
		one.ToolboxType = "other"
		OtherMap[one.ToolboxId] = one
	}
}

// GetOtherToolbox 查询单个
func (this_ *ToolboxService) GetOtherToolbox(toolboxId int64) (res *ToolboxModel, err error) {
	res = OtherMap[toolboxId]
	return
}
