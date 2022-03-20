package base

var (
	// 工具 权限

	// 工具基本 权限
	PowerToolbox       = addPower(&PowerAction{Action: "toolbox", Text: "工具", ShouldLogin: true, AllowNative: true})
	PowerToolboxPage   = addPower(&PowerAction{Action: "toolbox_page", Text: "工具页面", Parent: PowerToolbox, ShouldLogin: true, AllowNative: true})
	PowerToolboxList   = addPower(&PowerAction{Action: "toolbox_list", Text: "工具列表", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerToolboxInsert = addPower(&PowerAction{Action: "toolbox_insert", Text: "工具新增", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerToolboxUpdate = addPower(&PowerAction{Action: "toolbox_update", Text: "工具修改", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerToolboxRename = addPower(&PowerAction{Action: "toolbox_rename", Text: "工具重命名", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerToolboxDelete = addPower(&PowerAction{Action: "toolbox_delete", Text: "工具删除", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
)
