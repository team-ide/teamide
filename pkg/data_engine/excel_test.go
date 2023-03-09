package data_engine

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"testing"
)

func TestExcel(t *testing.T) {

	defer func() {
		if rec := recover(); rec != nil {
			util.Logger.Error("数据读取异常", zap.Any("error", rec))
		}
	}()
	task := &ExcelTask{}
	task.Path = `C:\Users\ZhuLiang\Downloads\导出库TEST_DB-表TB_USER1数据-20220522155318000.xlsx`
	task.OnData = func(data map[string]interface{}) (err error) {
		util.Logger.Info("on data_engine", zap.Any("data_engine", data))
		return
	}
	task.OnError = func(err error) {
		util.Logger.Error("on error", zap.Error(err))
		return
	}
	task.OnEnd = func() {
		util.Logger.Info("on end")
		return
	}
	task.SheetList = []*ExcelSheet{
		{
			NameList: []string{"c1", "c2", "c3", "c4"},
		},
		{
			NameList: []string{"c1", "c2", "c3", "c4"},
		},
	}
	task.Start()
}
