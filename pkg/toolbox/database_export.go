package toolbox

import (
	"fmt"
	"gitee.com/teamide/zorm"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

var (
	databaseExportTaskCache = map[string]*databaseExportTask{}
)

func addDatabaseExportTask(task *databaseExportTask) {
	databaseExportTaskCache[task.Key] = task
	go task.Start()
}

type databaseExportTask struct {
	Key              string                   `json:"key,omitempty"`
	Database         string                   `json:"database,omitempty"`
	Table            string                   `json:"table,omitempty"`
	ExportType       string                   `json:"exportType,omitempty"`
	ColumnList       []*db.TableColumnModel   `json:"columnList,omitempty"`
	Wheres           []*db.Where              `json:"wheres"`
	Orders           []*db.Order              `json:"orders"`
	ExportColumnList []map[string]interface{} `json:"exportColumnList,omitempty"`
	BatchNumber      int                      `json:"batchNumber,omitempty"`
	DataCount        int                      `json:"dataCount"`
	ReadyDataCount   int                      `json:"readyDataCount"`
	SuccessCount     int                      `json:"successCount"`
	ErrorCount       int                      `json:"errorCount"`
	IsEnd            bool                     `json:"isEnd,omitempty"`
	StartTime        time.Time                `json:"startTime,omitempty"`
	EndTime          time.Time                `json:"endTime,omitempty"`
	Error            string                   `json:"error,omitempty"`
	excelPath        string
	UseTime          int64 `json:"useTime"`
	IsStop           bool  `json:"isStop"`
	service          DatabaseService
	generateParam    *db.GenerateParam
}

func (this_ *databaseExportTask) Stop() {
	this_.IsStop = true
}
func (this_ *databaseExportTask) Start() {
	this_.StartTime = time.Now()
	defer func() {
		if err := recover(); err != nil {
			util.Logger.Error("根据策略导入数据异常", zap.Any("error", err))
			this_.Error = fmt.Sprint(err)
		}
		this_.EndTime = time.Now()
		this_.IsEnd = true
		this_.UseTime = util.GetTimeTime(this_.EndTime) - util.GetTimeTime(this_.StartTime)
	}()
	var err error

	err = this_.toSelectDataList()
	if err != nil {
		panic(err)
	}
}

func (this_ *databaseExportTask) doExport(dataList []map[string]interface{}) (err error) {
	if this_.ExportType == "excel" {
		err = this_.doExcel(dataList)
		if err != nil {
			return
		}
	}
	return

}
func (this_ *databaseExportTask) doExcel(dataList []map[string]interface{}) (err error) {
	var excelPath = this_.excelPath
	if excelPath == "" {
		excelPath, err = GetTempDir()
		if err != nil {
			return
		}
		excelPath += "/" + util.GenerateUUID() + ".xlsx"

		xlsxF := xlsx.NewFile()
		_, err = xlsxF.AddSheet("Sheet")
		if err != nil {
			return
		}
		err = xlsxF.Save(excelPath)
		if err != nil {
			return
		}
	}

	xlsxF, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return
	}
	sheet := xlsxF.Sheets[0]

	for _, data := range dataList {

		row := sheet.AddRow()

		for _, exportColumn := range this_.ExportColumnList {
			var columnName string
			columnName = exportColumn["column"].(string)

			cell := row.AddCell()
			value := data[columnName]
			if value != nil {
				stringValue := db.GetStringValue(value)
				cell.Value = stringValue
			}

		}
	}
	err = xlsxF.Save(excelPath)
	if err != nil {
		return
	}
	return

}

func (this_ *databaseExportTask) toSelectDataList() (err error) {

	var pageSize = this_.BatchNumber
	if pageSize <= 0 {
		pageSize = 100
	}

	sql, values, err := db.DataListSelectSql(this_.generateParam, this_.Database, this_.Table, this_.ColumnList, this_.Wheres, this_.Orders)
	if err != nil {
		return
	}
	var dataList []map[string]interface{}
	var pageIndex = 1
	for {
		if this_.IsStop {
			break
		}
		dataList, err = this_.selectDataList(sql, values, pageSize, pageIndex)
		if err != nil {
			return
		}

		var dataCount = len(dataList)
		this_.ReadyDataCount += dataCount
		err = this_.doExport(dataList)
		if err != nil {
			return
		}
		if len(dataList) < pageSize {
			break
		}
		pageIndex++
	}
	return

}

func (this_ *databaseExportTask) selectDataList(sql string, values []interface{}, pageSize int, pageIndex int) (dataList []map[string]interface{}, err error) {

	finder := zorm.NewFinder()
	finder.InjectionCheck = false

	finder.Append(sql, values...)

	page := zorm.NewPage()
	page.PageSize = pageSize
	page.PageNo = pageIndex
	dataList, err = this_.service.GetDatabaseWorker().FinderQueryMapPage(finder, page)
	if err != nil {
		return
	}
	this_.DataCount = page.TotalCount
	for _, one := range dataList {
		for k, v := range one {
			t, tOk := v.(time.Time)
			if tOk {
				if t.IsZero() {
					one[k] = nil
				} else {
					one[k] = t.Format("2006-01-02 15:04:05")
				}
			}
		}
	}
	return

}
