package toolbox

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
	databaseExportTaskCache = map[string]*databaseExportTask{}
)

func addDatabaseExportTask(task *databaseExportTask) {
	databaseExportTaskCache[task.Key] = task
	go task.Start()
}

type databaseExportTask struct {
	request          *DatabaseBaseRequest
	generateParam    *db.GenerateParam
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
	UseTime          int64 `json:"useTime"`
	IsStop           bool  `json:"isStop"`
	service          DatabaseService
}

func DatabaseExportDownload(data map[string]string, c *gin.Context) (err error) {

	taskKey := data["taskKey"]
	if taskKey == "" {
		err = errors.New("taskKey获取失败")
		return
	}

	databaseExportTask := databaseExportTaskCache[taskKey]
	if databaseExportTask == nil {
		err = errors.New("任务不存在")
		return
	}
	path := databaseExportTask.excelPath
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
	defer closeFile(fileInfo)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprint(fileSize))
	c.Header("download-file-name", fileName)

	err = CopyBytes(c.Writer, fileInfo, func(readSize int64, writeSize int64) {
	})
	if err != nil {
		return
	}

	c.Status(http.StatusOK)
	return
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

func (this_ *databaseExportTask) doExportExcel(dataList []map[string]interface{}) (err error) {
	var excelPath = this_.excelPath
	if excelPath == "" {
		excelPath, err = GetTempDir()
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

		excelPath += "/" + "导出库" + this_.request.Database + "-表" + this_.request.Table + "数据-" + time.Now().Format("20060102150405000") + ".xlsx"

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
func (this_ *databaseExportTask) doExportSql(dataList []map[string]interface{}) (err error) {

	var sqlF *os.File
	var excelPath = this_.excelPath
	if excelPath == "" {
		excelPath, err = GetTempDir()
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

		excelPath += "/" + "导出库" + this_.request.Database + "-表" + this_.request.Table + "数据-" + time.Now().Format("20060102150405000") + ".sql"

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
				column = GetColumnFromList(this_.request.ColumnList, columnName)
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

			insertColumns += this_.generateParam.PackingCharacterColumn(exportName) + ", "
			insertValues += this_.generateParam.PackingCharacterColumnStringValue(column, value) + ", "

			err = vm.Set(exportName, value)
			if err != nil {
				return
			}
		}

		insertColumns = strings.TrimSuffix(insertColumns, ", ")
		insertValues = strings.TrimSuffix(insertValues, ", ")

		sql := "INSERT INTO "

		if this_.generateParam.AppendDatabase && this_.ExportDatabase != "" {
			sql += this_.generateParam.PackingCharacterDatabase(this_.ExportDatabase) + "."
		}
		sql += this_.generateParam.PackingCharacterTable(this_.ExportTable)
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

func (this_ *databaseExportTask) toSelectDataList() (err error) {

	var pageSize = this_.BatchNumber
	if pageSize <= 0 {
		pageSize = 100
	}

	sql, values, err := db.DataListSelectSql(this_.generateParam, this_.request.Database, this_.request.Table, this_.request.ColumnList, this_.request.Wheres, this_.request.Orders)
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
