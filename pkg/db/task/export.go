package task

import (
	"bufio"
	"errors"
	"fmt"
	"gitee.com/teamide/zorm"
	"github.com/dop251/goja"
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	excelPath        string
	exportDataIndex  int
	UseTime          int64             `json:"useTime"`
	IsStop           bool              `json:"isStop"`
	GenerateParam    *db.GenerateParam `json:"-"`
	Service          *db.Service       `json:"-"`
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
	path := ExportTask.excelPath
	if path == "" {
		err = errors.New("任务导出文件丢失")
		return
	}

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

	err = this_.toSelectDataList()
	if err != nil {
		panic(err)
	}
}

func (this_ *ExportTask) doExport(dataList []map[string]interface{}) (err error) {
	if this_.ExportType == "excel" {
		err = this_.doExportExcel(dataList)
		if err != nil {
			return
		}
	} else if this_.ExportType == "sql" {
		err = this_.doExportSql(dataList)
		if err != nil {
			return
		}
	}
	return

}

func (this_ *ExportTask) doExportExcel(dataList []map[string]interface{}) (err error) {
	var excelPath = this_.excelPath
	if excelPath == "" {
		excelPath, err = util.GetTempDir()
		if err != nil {
			return
		}

		var isExists bool
		isExists, err = util.PathExists(excelPath)
		if err != nil {
			return
		}
		if !isExists {
			err = os.MkdirAll(excelPath, 0777)
			if err != nil {
				return
			}
		}

		excelPath += "/" + "导出库" + this_.Database + "-表" + this_.Table + "数据-" + time.Now().Format("20060102150405000") + ".xlsx"

		xlsxF := xlsx.NewFile()
		var sheet *xlsx.Sheet
		sheet, err = xlsxF.AddSheet("Sheet")
		if err != nil {
			return
		}

		row := sheet.AddRow()
		for _, exportColumn := range this_.ExportColumnList {

			cell := row.AddCell()
			var exportName = exportColumn["exportName"]
			if exportColumn != nil {
				cell.Value = exportName.(string)
			}
		}
		err = xlsxF.Save(excelPath)
		if err != nil {
			return
		}

		this_.excelPath = excelPath
	}

	xlsxF, err := xlsx.OpenFile(excelPath)
	if err != nil {
		return
	}
	sheet := xlsxF.Sheets[0]

	scriptContext := javascript.GetContext()

	vm := goja.New()

	for key, value := range scriptContext {
		err = vm.Set(key, value)
		if err != nil {
			return
		}
	}

	for _, data := range dataList {
		err = vm.Set("_$index", this_.exportDataIndex)
		if err != nil {
			return
		}

		row := sheet.AddRow()

		for _, exportColumn := range this_.ExportColumnList {
			var exportName string
			if exportColumn["exportName"] != nil {
				exportName = exportColumn["exportName"].(string)
			}
			var value = exportColumn["value"]
			if value == nil || value == "" {
				if exportColumn["column"] != nil {
					columnName := exportColumn["column"].(string)
					value = data[columnName]
				}
			} else {
				valueScript := value.(string)
				var scriptValue goja.Value
				scriptValue, err = vm.RunString(valueScript)
				if err != nil {
					util.Logger.Error("表达式执行异常", zap.Any("script", valueScript), zap.Error(err))
					return
				}
				value = scriptValue.Export()
			}

			cell := row.AddCell()
			if value != nil {
				stringValue := db.GetStringValue(value)
				cell.Value = stringValue
			}

			err = vm.Set(exportName, value)
			if err != nil {
				return
			}
		}
		this_.exportDataIndex++
		this_.SuccessCount++
	}
	err = xlsxF.Save(excelPath)
	if err != nil {
		return
	}
	return

}
func GetColumnFromList(columnList []*db.TableColumnModel, name string) *db.TableColumnModel {
	if len(columnList) == 0 {
		return nil
	}

	for _, column := range columnList {
		if column.Name == name {
			return column
		}
	}
	return nil

}
func (this_ *ExportTask) doExportSql(dataList []map[string]interface{}) (err error) {

	var sqlF *os.File
	var excelPath = this_.excelPath
	if excelPath == "" {
		excelPath, err = util.GetTempDir()
		if err != nil {
			return
		}

		var isExists bool
		isExists, err = util.PathExists(excelPath)
		if err != nil {
			return
		}
		if !isExists {
			err = os.MkdirAll(excelPath, 0777)
			if err != nil {
				return
			}
		}

		excelPath += "/" + "导出库" + this_.Database + "-表" + this_.Table + "数据-" + time.Now().Format("20060102150405000") + ".sql"

		sqlF, err = os.Create(excelPath)
		if err != nil {
			return
		}

		this_.excelPath = excelPath
	} else {
		sqlF, err = os.OpenFile(excelPath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return
		}
	}

	defer func() {
		_ = sqlF.Close()
	}()

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(sqlF)
	scriptContext := javascript.GetContext()

	vm := goja.New()

	for key, value := range scriptContext {
		err = vm.Set(key, value)
		if err != nil {
			return
		}
	}

	for _, data := range dataList {
		err = vm.Set("_$index", this_.exportDataIndex)
		if err != nil {
			return
		}

		insertColumns := ""
		insertValues := ""
		for _, exportColumn := range this_.ExportColumnList {
			var exportName string
			if exportColumn["exportName"] != nil {
				exportName = exportColumn["exportName"].(string)
			}
			var value = exportColumn["value"]
			var column *db.TableColumnModel
			if exportColumn["column"] != nil {
				columnName := exportColumn["column"].(string)
				column = GetColumnFromList(this_.ColumnList, columnName)
			}
			if column == nil {
				column = &db.TableColumnModel{
					Type: "VARCHAR",
				}
			}
			if value == nil || value == "" {
				if column != nil {
					value = data[column.Name]
				}
			} else {
				valueScript := value.(string)
				var scriptValue goja.Value
				scriptValue, err = vm.RunString(valueScript)
				if err != nil {
					util.Logger.Error("表达式执行异常", zap.Any("script", valueScript), zap.Error(err))
					return
				}
				value = scriptValue.Export()
			}

			insertColumns += this_.GenerateParam.PackingCharacterColumn(exportName) + ", "
			insertValues += this_.GenerateParam.PackingCharacterColumnStringValue(column, value) + ", "

			err = vm.Set(exportName, value)
			if err != nil {
				return
			}
		}

		insertColumns = strings.TrimSuffix(insertColumns, ", ")
		insertValues = strings.TrimSuffix(insertValues, ", ")

		sql := "INSERT INTO "

		if this_.GenerateParam.AppendDatabase && this_.ExportDatabase != "" {
			sql += this_.GenerateParam.PackingCharacterDatabase(this_.ExportDatabase) + "."
		}
		sql += this_.GenerateParam.PackingCharacterTable(this_.ExportTable)
		if insertColumns != "" {
			sql += "(" + insertColumns + ")"
		}
		if insertValues != "" {
			sql += " VALUES (" + insertValues + ")"
		}

		_, err = write.WriteString(sql + ";\n")
		if err != nil {
			return
		}
		this_.exportDataIndex++
		this_.SuccessCount++
	}
	err = write.Flush()
	return

}

func (this_ *ExportTask) toSelectDataList() (err error) {

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
