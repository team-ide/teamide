package toolbox

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"teamide/pkg/data_engine"
	"teamide/pkg/db"
	"teamide/pkg/util"
	"time"
)

var (
	OtherExportTaskCache = map[string]*OtherExportTask{}
)

func StartOtherExportTask(task *OtherExportTask) {
	OtherExportTaskCache[task.Key] = task
	go task.Start()
}

func GetOtherExportTask(taskKey string) *OtherExportTask {
	task := OtherExportTaskCache[taskKey]
	return task
}

func StopOtherExportTask(taskKey string) *OtherExportTask {
	task := OtherExportTaskCache[taskKey]
	if task != nil {
		task.Start()
	}
	return task
}

func CleanOtherExportTask(taskKey string) *OtherExportTask {
	task := OtherExportTaskCache[taskKey]
	if task != nil {
		delete(OtherExportTaskCache, taskKey)
	}
	return task
}

type OtherExportTask struct {
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

	OtherExportTask := OtherExportTaskCache[taskKey]
	if OtherExportTask == nil {
		err = errors.New("任务不存在")
		return
	}
	if OtherExportTask.exportTask == nil || OtherExportTask.exportTask.ExcelPath == "" {
		err = errors.New("任务导出文件丢失")
		return
	}
	path := OtherExportTask.exportTask.ExcelPath

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

func (this_ *OtherExportTask) Stop() {
	this_.IsStop = true
}
func (this_ *OtherExportTask) Start() {
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
		var exportValueKey = util.GenerateUUID()
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
		Name:             "生成数据-" + time.Now().Format("20060102150405") + "." + fileSuffix,
		CellSeparator:    this_.CellSeparator,
		ExportColumnList: this_.ExportColumnList,
	}
	err = this_.toStrategyData()
	if err != nil {
		panic(err)
	}
}

func (this_ *OtherExportTask) doExport(dataList []map[string]interface{}) (err error) {
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
func (this_ *OtherExportTask) success() {
	this_.SuccessCount++
	return

}

func (this_ *OtherExportTask) toStrategyData() (err error) {

	return

}
