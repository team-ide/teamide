package data_engine

//
//import (
//	"bufio"
//	"github.com/tealeg/xlsx"
//	"github.com/team-ide/go-dialect/dialect"
//	"os"
//	"strings"
//	"github.com/team-ide/go-tool/util"
//)
//
//type DataExportTask struct {
//	CellSeparator    string                   `json:"cellSeparator,omitempty"`
//	ExcelPath        string                   `json:"excelPath,omitempty"`
//	Name             string                   `json:"name,omitempty"`
//	ExportColumnList []map[string]interface{} `json:"exportColumnList,omitempty"`
//}
//
//func (this_ *DataExportTask) ExportExcel(dataList []map[string]interface{}, success func()) (err error) {
//	var excelPath = this_.ExcelPath
//	if excelPath == "" {
//		excelPath, err = util.GetTempDir()
//		if err != nil {
//			return
//		}
//
//		var isExists bool
//		isExists, err = util.PathExists(excelPath)
//		if err != nil {
//			return
//		}
//		if !isExists {
//			err = os.MkdirAll(excelPath, 0777)
//			if err != nil {
//				return
//			}
//		}
//
//		excelPath += "/" + this_.Name
//
//		xlsxF := xlsx.NewFile()
//		var sheet *xlsx.Sheet
//		sheet, err = xlsxF.AddSheet("Sheet")
//		if err != nil {
//			return
//		}
//
//		row := sheet.AddRow()
//		for _, exportColumn := range this_.ExportColumnList {
//
//			cell := row.AddCell()
//			var exportName = exportColumn["exportName"]
//			if exportColumn != nil {
//				cell.Value = exportName.(string)
//			}
//		}
//		err = xlsxF.Save(excelPath)
//		if err != nil {
//			return
//		}
//
//		this_.ExcelPath = excelPath
//	}
//
//	xlsxF, err := xlsx.OpenFile(excelPath)
//	if err != nil {
//		return
//	}
//	sheet := xlsxF.Sheets[0]
//
//	for _, data := range dataList {
//
//		row := sheet.AddRow()
//
//		for _, exportColumn := range this_.ExportColumnList {
//
//			cell := row.AddCell()
//			exportValueKey := exportColumn["exportValueKey"].(string)
//			var value = data[exportValueKey]
//			var stringValue string
//			if value != nil {
//				stringValue, err = util.GetStringValue(value)
//				if err != nil {
//					return
//				}
//				cell.Value = stringValue
//			}
//
//		}
//		if success != nil {
//			success()
//		}
//	}
//	err = xlsxF.Save(excelPath)
//	if err != nil {
//		return
//	}
//	return
//
//}
//
//func (this_ *DataExportTask) ExportCsv(dataList []map[string]interface{}, success func()) (err error) {
//	var cellSeparator = this_.CellSeparator
//	if cellSeparator == "" {
//		cellSeparator = ","
//	}
//	var csvF *os.File
//	var excelPath = this_.ExcelPath
//	if excelPath == "" {
//		excelPath, err = util.GetTempDir()
//		if err != nil {
//			return
//		}
//
//		var isExists bool
//		isExists, err = util.PathExists(excelPath)
//		if err != nil {
//			return
//		}
//		if !isExists {
//			err = os.MkdirAll(excelPath, 0777)
//			if err != nil {
//				return
//			}
//		}
//
//		excelPath += "/" + this_.Name
//
//		csvF, err = os.Create(excelPath)
//		if err != nil {
//			return
//		}
//
//		var values []string
//		for _, exportColumn := range this_.ExportColumnList {
//			var exportName = exportColumn["exportName"]
//			if exportColumn != nil {
//				values = append(values, exportName.(string))
//			} else {
//				values = append(values, "")
//			}
//		}
//		_, err = csvF.Write([]byte(strings.Join(values, cellSeparator) + "\n"))
//		if err != nil {
//			return
//		}
//
//		this_.ExcelPath = excelPath
//	} else {
//		csvF, err = os.OpenFile(excelPath, os.O_WRONLY|os.O_APPEND, 0666)
//		if err != nil {
//			return
//		}
//	}
//
//	defer func() {
//		_ = csvF.Close()
//	}()
//
//	//写入文件时，使用带缓存的 *Writer
//	write := bufio.NewWriter(csvF)
//
//	if err != nil {
//		return
//	}
//
//	for _, data := range dataList {
//
//		var values []string
//		for _, exportColumn := range this_.ExportColumnList {
//
//			exportValueKey := exportColumn["exportValueKey"].(string)
//			var value = data[exportValueKey]
//			var stringValue string
//			if value != nil {
//				stringValue, err = util.GetStringValue(value)
//				if err != nil {
//					return
//				}
//			}
//			values = append(values, stringValue)
//		}
//
//		_, err = write.WriteString(strings.Join(values, cellSeparator) + "\n")
//		if err != nil {
//			return
//		}
//
//		if success != nil {
//			success()
//		}
//	}
//	err = write.Flush()
//	return
//}
//
//func GetColumnFromList(columnList []*dialect.ColumnModel, name string) *dialect.ColumnModel {
//	if len(columnList) == 0 {
//		return nil
//	}
//
//	for _, column := range columnList {
//		if column.ColumnName == name {
//			return column
//		}
//	}
//	return nil
//
//}
//func (this_ *DataExportTask) ExportSql(param *dialect.ParamModel, database string, table string, columnList []*dialect.ColumnModel, dataList []map[string]interface{}, success func()) (err error) {
//
//	var sqlF *os.File
//	var excelPath = this_.ExcelPath
//	if excelPath == "" {
//		excelPath, err = util.GetTempDir()
//		if err != nil {
//			return
//		}
//
//		var isExists bool
//		isExists, err = util.PathExists(excelPath)
//		if err != nil {
//			return
//		}
//		if !isExists {
//			err = os.MkdirAll(excelPath, 0777)
//			if err != nil {
//				return
//			}
//		}
//
//		excelPath += "/" + this_.Name
//
//		sqlF, err = os.Create(excelPath)
//		if err != nil {
//			return
//		}
//
//		this_.ExcelPath = excelPath
//	} else {
//		sqlF, err = os.OpenFile(excelPath, os.O_WRONLY|os.O_APPEND, 0666)
//		if err != nil {
//			return
//		}
//	}
//
//	defer func() {
//		_ = sqlF.Close()
//	}()
//
//	//写入文件时，使用带缓存的 *Writer
//	write := bufio.NewWriter(sqlF)
//
//	if err != nil {
//		return
//	}
//
//	for _, data := range dataList {
//
//		insertColumns := ""
//		insertValues := ""
//		for _, exportColumn := range this_.ExportColumnList {
//			var exportName string
//			if exportColumn["exportName"] != nil {
//				exportName = exportColumn["exportName"].(string)
//			}
//
//			exportValueKey := exportColumn["exportValueKey"].(string)
//			var value = data[exportValueKey]
//			var column *dialect.ColumnModel
//			if exportColumn["column"] != nil {
//				columnName := exportColumn["column"].(string)
//				column = GetColumnFromList(columnList, columnName)
//			}
//			if column == nil {
//				column = &dialect.ColumnModel{
//					Type: "varchar",
//				}
//			}
//
//			insertColumns += param.PackingCharacterColumn(exportName) + ", "
//			insertValues += param.PackingCharacterColumnStringValue(column, value) + ", "
//
//		}
//
//		insertColumns = strings.TrimSuffix(insertColumns, ", ")
//		insertValues = strings.TrimSuffix(insertValues, ", ")
//
//		sql := "INSERT INTO "
//
//		if param.AppendDatabase && database != "" {
//			sql += param.PackingCharacterDatabase(database) + "."
//		}
//		sql += param.PackingCharacterTable(table)
//		if insertColumns != "" {
//			sql += "(" + insertColumns + ")"
//		}
//		if insertValues != "" {
//			sql += " VALUES (" + insertValues + ")"
//		}
//
//		_, err = write.WriteString(sql + ";\n")
//		if err != nil {
//			return
//		}
//		if success != nil {
//			success()
//		}
//	}
//	err = write.Flush()
//	return
//
//}
