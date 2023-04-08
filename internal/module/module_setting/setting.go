package module_setting

import (
	"github.com/team-ide/go-tool/util"
	"teamide/internal/context"
	"teamide/internal/module/module_id"
)

// NewSettingService 根据库配置创建 SettingService
func NewSettingService(ServerContext *context.ServerContext) (res *SettingService) {

	idService := module_id.NewIDService(ServerContext)
	res = &SettingService{
		ServerContext: ServerContext,
		idService:     idService,
	}
	return
}

// SettingService 设置服务
type SettingService struct {
	*context.ServerContext
	idService *module_id.IDService
}

// Save 保存
func (this_ *SettingService) Save(data map[string]interface{}) (err error) {
	if data == nil {
		return
	}
	s := context.NewSetting()
	setData := map[string]string{}
	for name, value := range data {
		setData[name] = util.GetStringValue(value)
		var find bool
		find, err = s.Set(name, setData[name])
		if err != nil {
			return
		}
		if !find {
			delete(setData, name)
		}
	}

	var find *SettingModel
	for name, value := range setData {
		find, err = this_.Get(name)
		if err != nil {
			return
		}
		setting := &SettingModel{
			Name:  name,
			Value: value,
		}
		if find != nil {
			err = this_.Update(setting)
		} else {
			err = this_.Insert(setting)
		}
		if err != nil {
			return
		}
	}

	for name, value := range setData {
		_, err = this_.Setting.Set(name, value)
		if err != nil {
			return
		}
	}
	return
}

// Insert 新增
func (this_ *SettingService) Insert(setting *SettingModel) (err error) {

	if setting.CreateTime == 0 {
		setting.CreateTime = util.GetNowTime()
	}

	sql := `INSERT INTO ` + TableSetting + `(name, value, createTime) VALUES (?, ?, ?) `

	_, err = this_.DatabaseWorker.Exec(sql, []interface{}{setting.Name, setting.Value, setting.CreateTime})
	if err != nil {
		return
	}
	return
}

func (this_ *SettingService) Update(setting *SettingModel) (err error) {

	setting.UpdateTime = util.GetNowTime()

	_, err = this_.DatabaseWorker.Exec(`UPDATE `+TableSetting+` SET value=? , updateTime=? WHERE name=? `, []interface{}{setting.Value, setting.UpdateTime, setting.Name})
	if err != nil {
		return
	}
	return
}

// Query 查询
func (this_ *SettingService) Query() (res []*SettingModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableSetting + ` `

	err = this_.DatabaseWorker.Query(sql, values, &res)
	if err != nil {
		return
	}

	return
}

// Get 查询
func (this_ *SettingService) Get(name string) (res *SettingModel, err error) {

	var values []interface{}
	sql := `SELECT * FROM ` + TableSetting + ` WHERE name=? `
	values = append(values, name)

	var list []*SettingModel
	err = this_.DatabaseWorker.Query(sql, values, &list)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		res = nil
	}

	return
}

func (this_ *SettingService) delete(name string) (err error) {
	var sql string
	var values []interface{}

	sql += "DELETE FROM " + TableSetting + " WHERE name=?"
	values = append(values, name)
	_, err = this_.DatabaseWorker.Exec(sql, values)
	if err != nil {
		return
	}
	return
}
