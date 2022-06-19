package task

import (
	"errors"
	"fmt"
	"gitee.com/teamide/zorm"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"teamide/pkg/data_engine"
	"teamide/pkg/db"
	"teamide/pkg/javascript"
	"teamide/pkg/util"
	"time"
)

var (
	ExportTaskCache = map[string]*ExportTask{}
)

func StartExportTask(task *ExportTask) {
	ExportTaskCache[task.Key] = task
	go task.Start()
}

func GetExportTask(taskKey string) *ExportTask {
	task := ExportTaskCache[taskKey]
	return task
}

func StopExportTask(taskKey string) *ExportTask {
	task := ExportTaskCache[taskKey]
	if task != nil {
		task.Start()
	}
	return task
}

func CleanExportTask(taskKey string) *ExportTask {
	task := ExportTaskCache[taskKey]
	if task != nil {
		delete(ExportTaskCache, taskKey)
	}
	return task
}

type ExportTask struct {
	CellSeparator    string                   `json:"cellSeparator,omitempty"`
	FileSuffix       string                   `json:"fileSuffix,omitempty"`
	Database         string                   `json:"database,omitempty"`
	Table            string                   `json:"table,omitempty"`
	ColumnList       []*db.TableColumnModel   `json:"columnList,omitempty"`
	Wheres           []*db.Where              `json:"wheres,omitempty"`
	Orders           []*db.Order              `json:"orders,omitempty"`
	Key              string                   `json:"key,omitempty"`
	ExportType       string                   `json:"exportType,omitempty"`
	ExportDatabase   string                   `json:"exportDatabase,omitempty"`
	ExportTable      string                   `json:"exportTable,omitempty"`
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
	exportDataIndex  int
	UseTime          int64             `json:"useTime"`
	IsStop           bool              `json:"isStop"`
	GenerateParam    *db.GenerateParam `json:"-"`
	Service          *db.Service       `json:"-"`
	exportTask       *data_engine.DataExportTask
}

func DatabaseExportDownload(data map[string]string, c *gin.Context) (err error) {

	taskKey := data["taskKey"]
	if taskKey == "" {
		err = errors.New("taskKey获取失败")
		return
	}

	ExportTask := ExportTaskCache[taskKey]
	if ExportTask == nil {
		err = errors.New("任务不存在")
		return
	}
	if ExportTask.exportTask == nil || ExportTask.exportTask.ExcelPath == "" {
		err = errors.New("任务导出文件丢失")
		return
	}
	path := ExportTask.exportTask.ExcelPath

	var fileName string
	var fileSize int64
	ff, err := os.Lstat(path)
	if err != nil {
		return
	}
	fileName = ff.Name()
	fileSize = ff.Size()

	var fileInfo *os.File
	fileInfo, err = os.Open(path)
	if err != nil {
		return
	}
	defer func() {
		_ = fileInfo.Close()
	}()

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	err = util.CopyBytes(c.Writer, fileInfo, func(readSize int64, writeSize int64) {
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
}

func (this_ *ExportTask) Stop() {
	this_.IsStop = true
}
func (this_ *ExportTask) Start() {
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

	for _, exportColumn := range this_.ExportColumnList {
		var exportValueKey = util.UUID()
		exportColumn["exportValueKey"] = exportValueKey
	}
	var fileSuffix = ""
	if this_.ExportType == "excel" {
		fileSuffix = "xlsx"
	} else if this_.ExportType == "sql" {
		fileSuffix = "sql"
	} else if this_.ExportType == "csv" {
		fileSuffix = this_.FileSuffix
		if fileSuffix == "" {
			fileSuffix = "csv"
		}
	}

	this_.exportTask = &data_engine.DataExportTask{
		Name:             "导出库" + this_.Database + "-表" + this_.Table + "数据-" + time.Now().Format("20060102150405") + "." + fileSuffix,
		CellSeparator:    this_.CellSeparator,
		ExportColumnList: this_.ExportColumnList,
	}
	err = this_.toSelectDataList()
	if err != nil {
		panic(err)
	}
}

func (this_ *ExportTask) doExport(dataList []map[string]interface{}) (err error) {
	if this_.ExportType == "excel" {
		err = this_.exportTask.ExportExcel(dataList, this_.success)
		if err != nil {
			return
		}
	} else if this_.ExportType == "sql" {
		err = this_.exportTask.ExportSql(this_.GenerateParam, this_.ExportDatabase, this_.ExportTable, this_.ColumnList, dataList, this_.success)
		if err != nil {
			return
		}
	} else if this_.ExportType == "csv" {
		err = this_.exportTask.ExportCsv(dataList, this_.success)
		if err != nil {
			return
		}
	}
	return

}
func (this_ *ExportTask) success() {
	this_.SuccessCount++
	return

}

func (this_ *ExportTask) toSelectDataList() (err error) {

	script, err := javascript.NewScript()
	if err != nil {
		return
	}
	var pageSize = this_.BatchNumber
	if pageSize <= 0 {
		pageSize = 100
	}

	sql, values, err := db.DataListSelectSql(this_.GenerateParam, this_.Database, this_.Table, this_.ColumnList, this_.Wheres, this_.Orders)
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

		// 格式化

		for _, data := range dataList {
			err = script.Set("_$index", this_.exportDataIndex)
			if err != nil {
				return
			}
			for k, v := range data {
				err = script.Set(k, v)
				if err != nil {
					return
				}
			}

			for _, exportColumn := range this_.ExportColumnList {
				var exportValueKey = exportColumn["exportValueKey"].(string)
				var value = exportColumn["value"]
				if value == nil || value == "" {
					var column string
					if exportColumn["column"] != nil {
						column = exportColumn["column"].(string)
					}
					if column != "" {
						value = data[column]
					}

				} else {
					value, err = script.GetScriptValue(value.(string))
					if err != nil {
						return
					}
				}
				data[exportValueKey] = value
			}

			this_.exportDataIndex++
		}

		this_.ReadyDataCount += dataCount
		err = this_.doExport(dataList)
		if err != nil {
			return
		}

		if dataCount < pageSize {
			break
		}
		pageIndex++
	}
	return

}

func (this_ *ExportTask) selectDataList(sql string, values []interface{}, pageSize int, pageIndex int) (dataList []map[string]interface{}, err error) {

	finder := zorm.NewFinder()
	finder.InjectionCheck = false

	finder.Append(sql, values...)

	page := zorm.NewPage()
	page.PageSize = pageSize
	page.PageNo = pageIndex
	dataList, err = this_.Service.GetDatabaseWorker().FinderQueryMapPage(finder, page)
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
