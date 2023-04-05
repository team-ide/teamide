package module_toolbox

import (
	"teamide/internal/context"
	"teamide/pkg/base"
)

type ToolboxApi struct {
	*context.ServerContext
	ToolboxService *ToolboxService
}

func NewToolboxApi(ToolboxService *ToolboxService) *ToolboxApi {
	return &ToolboxApi{
		ServerContext:  ToolboxService.ServerContext,
		ToolboxService: ToolboxService,
	}
}

var (
	// 工具 权限

	// Power 工具基本 权限
	Power          = base.AppendPower(&base.PowerAction{Action: "toolbox", Text: "工具箱", ShouldLogin: true, StandAlone: true})
	PowerList      = base.AppendPower(&base.PowerAction{Action: "list", Text: "工具箱列表", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerGet       = base.AppendPower(&base.PowerAction{Action: "get", Text: "工具箱查询", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerCount     = base.AppendPower(&base.PowerAction{Action: "count", Text: "工具箱统计", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerInsert    = base.AppendPower(&base.PowerAction{Action: "insert", Text: "工具箱新增", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerUpdate    = base.AppendPower(&base.PowerAction{Action: "update", Text: "工具箱修改", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerRename    = base.AppendPower(&base.PowerAction{Action: "rename", Text: "工具箱重命名", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerDelete    = base.AppendPower(&base.PowerAction{Action: "delete", Text: "工具箱删除", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerMoveGroup = base.AppendPower(&base.PowerAction{Action: "moveGroup", Text: "工具箱分组", Parent: Power, ShouldLogin: true, StandAlone: true})

	PowerGroup       = base.AppendPower(&base.PowerAction{Action: "group", Text: "工具箱分组列表", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerGroupList   = base.AppendPower(&base.PowerAction{Action: "list", Text: "工具箱分组列表", Parent: PowerGroup, ShouldLogin: true, StandAlone: true})
	PowerGroupInsert = base.AppendPower(&base.PowerAction{Action: "insert", Text: "工具箱分组新增", Parent: PowerGroup, ShouldLogin: true, StandAlone: true})
	PowerGroupUpdate = base.AppendPower(&base.PowerAction{Action: "update", Text: "工具箱分组修改", Parent: PowerGroup, ShouldLogin: true, StandAlone: true})
	PowerGroupDelete = base.AppendPower(&base.PowerAction{Action: "delete", Text: "工具箱分组删除", Parent: PowerGroup, ShouldLogin: true, StandAlone: true})

	PowerOpen                = base.AppendPower(&base.PowerAction{Action: "open", Text: "打开工具箱工具", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerUpdateOpenExtend    = base.AppendPower(&base.PowerAction{Action: "updateOpenExtend", Text: "更新工具箱工具扩展", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerUpdateOpenSequence  = base.AppendPower(&base.PowerAction{Action: "updateOpenSequence", Text: "更新工具箱排序", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerQueryOpens          = base.AppendPower(&base.PowerAction{Action: "queryOpens", Text: "查询打开的工具箱工具", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerGetOpen             = base.AppendPower(&base.PowerAction{Action: "getOpen", Text: "查询工具箱工具打开信息", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerClose               = base.AppendPower(&base.PowerAction{Action: "close", Text: "工具箱工具关闭", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerOpenTab             = base.AppendPower(&base.PowerAction{Action: "openTab", Text: "工具箱工具打开Tab", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerQueryOpenTabs       = base.AppendPower(&base.PowerAction{Action: "queryOpenTabs", Text: "查询工具箱工具打开的Tab", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerCloseTab            = base.AppendPower(&base.PowerAction{Action: "closeTab", Text: "工具箱工具关闭Tab", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerUpdateOpenTabExtend = base.AppendPower(&base.PowerAction{Action: "updateOpenTabExtend", Text: "更新工具箱工具Tab扩展", Parent: Power, ShouldLogin: true, StandAlone: true})

	PowerQuickCommand       = base.AppendPower(&base.PowerAction{Action: "quickCommand", Text: "工具快速指令", Parent: Power, ShouldLogin: true, StandAlone: true})
	PowerQuickCommandQuery  = base.AppendPower(&base.PowerAction{Action: "query", Text: "工具快速指令查询", Parent: PowerQuickCommand, ShouldLogin: true, StandAlone: true})
	PowerQuickCommandInsert = base.AppendPower(&base.PowerAction{Action: "insert", Text: "工具快速指令新增", Parent: PowerQuickCommand, ShouldLogin: true, StandAlone: true})
	PowerQuickCommandUpdate = base.AppendPower(&base.PowerAction{Action: "update", Text: "工具快速指令修改", Parent: PowerQuickCommand, ShouldLogin: true, StandAlone: true})
	PowerQuickCommandDelete = base.AppendPower(&base.PowerAction{Action: "delete", Text: "工具快速指令删除", Parent: PowerQuickCommand, ShouldLogin: true, StandAlone: true})
)

func (this_ *ToolboxApi) GetApis() (apis []*base.ApiWorker) {
	apis = append(apis, &base.ApiWorker{Power: PowerList, Do: this_.list})
	apis = append(apis, &base.ApiWorker{Power: PowerGet, Do: this_.get})
	apis = append(apis, &base.ApiWorker{Power: PowerCount, Do: this_.count})
	apis = append(apis, &base.ApiWorker{Power: PowerInsert, Do: this_.insert})
	apis = append(apis, &base.ApiWorker{Power: PowerUpdate, Do: this_.update})
	apis = append(apis, &base.ApiWorker{Power: PowerRename, Do: this_.rename})
	apis = append(apis, &base.ApiWorker{Power: PowerDelete, Do: this_.delete})
	apis = append(apis, &base.ApiWorker{Power: PowerMoveGroup, Do: this_.moveGroup})

	apis = append(apis, &base.ApiWorker{Power: PowerGroupList, Do: this_.listGroup})
	apis = append(apis, &base.ApiWorker{Power: PowerGroupInsert, Do: this_.insertGroup})
	apis = append(apis, &base.ApiWorker{Power: PowerGroupUpdate, Do: this_.updateGroup})
	apis = append(apis, &base.ApiWorker{Power: PowerGroupDelete, Do: this_.deleteGroup})

	apis = append(apis, &base.ApiWorker{Power: PowerOpen, Do: this_.open})
	apis = append(apis, &base.ApiWorker{Power: PowerGetOpen, Do: this_.getOpen})
	apis = append(apis, &base.ApiWorker{Power: PowerQueryOpens, Do: this_.queryOpens})
	apis = append(apis, &base.ApiWorker{Power: PowerClose, Do: this_.close})
	apis = append(apis, &base.ApiWorker{Power: PowerUpdateOpenExtend, Do: this_.updateOpenExtend})
	apis = append(apis, &base.ApiWorker{Power: PowerUpdateOpenSequence, Do: this_.UpdateOpenSequence})
	apis = append(apis, &base.ApiWorker{Power: PowerQueryOpenTabs, Do: this_.queryOpenTabs})
	apis = append(apis, &base.ApiWorker{Power: PowerOpenTab, Do: this_.openTab})
	apis = append(apis, &base.ApiWorker{Power: PowerCloseTab, Do: this_.closeTab})
	apis = append(apis, &base.ApiWorker{Power: PowerUpdateOpenTabExtend, Do: this_.updateOpenTabExtend})

	apis = append(apis, &base.ApiWorker{Power: PowerQuickCommandQuery, Do: this_.queryQuickCommand})
	apis = append(apis, &base.ApiWorker{Power: PowerQuickCommandInsert, Do: this_.insertQuickCommand})
	apis = append(apis, &base.ApiWorker{Power: PowerQuickCommandUpdate, Do: this_.updateQuickCommand})
	apis = append(apis, &base.ApiWorker{Power: PowerQuickCommandDelete, Do: this_.deleteQuickCommand})

	return
}

type QuickCommandType struct {
	Name  string `json:"name,omitempty"`
	Text  string `json:"text,omitempty"`
	Value int    `json:"value,omitempty"`
}

var (
	QuickCommandTypes []*QuickCommandType
)

func init() {
	QuickCommandTypes = append(QuickCommandTypes, &QuickCommandType{Name: "SSH Command", Text: "", Value: 1})
}

func GetQuickCommandTypes() []*QuickCommandType {
	return QuickCommandTypes
}
