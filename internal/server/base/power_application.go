package base

var (
	// 应用 权限

	// 应用基本 权限
	PowerApplication        = addPower(&PowerAction{Action: "application", Text: "应用", ShouldLogin: true, AllowNative: true})
	PowerApplicationPage    = addPower(&PowerAction{Action: "application_page", Text: "应用页面", Parent: PowerApplication, ShouldLogin: true, AllowNative: true})
	PowerApplicationList    = addPower(&PowerAction{Action: "application_list", Text: "应用列表", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationInsert  = addPower(&PowerAction{Action: "application_insert", Text: "应用新增", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationUpdate  = addPower(&PowerAction{Action: "application_update", Text: "应用修改", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationRename  = addPower(&PowerAction{Action: "application_rename", Text: "应用重命名", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationDelete  = addPower(&PowerAction{Action: "application_delete", Text: "应用删除", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationStart   = addPower(&PowerAction{Action: "application_start", Text: "应用启动", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationStop    = addPower(&PowerAction{Action: "application_stop", Text: "应用停止", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})
	PowerApplicationRestart = addPower(&PowerAction{Action: "application_restart", Text: "应用重启", Parent: PowerApplicationPage, ShouldLogin: true, AllowNative: true})

	// 应用Context 权限
	PowerApplicationContext     = addPower(&PowerAction{Action: "application_context", Text: "应用Context", ShouldLogin: true, AllowNative: true})
	PowerApplicationContextLoad = addPower(&PowerAction{Action: "application_context_load", Text: "应用Context加载", Parent: PowerApplicationContext, ShouldLogin: true, AllowNative: true})
	PowerApplicationContextSave = addPower(&PowerAction{Action: "application_context_save", Text: "应用Context保存", Parent: PowerApplicationContext, ShouldLogin: true, AllowNative: true})
)
